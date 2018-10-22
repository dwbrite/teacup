package teacup

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type teacup struct {
	dbInfo  string
	TlsPair *TlsKeyPair

	FileWhitelist regexp.Regexp
	DirBlacklist  []string//regexp.Regexp

	Log    *log.Logger
	errors tcErrors

	tables   map[string]bool
	webpages map[string]func(http.Request, string) (*template.Template, interface{})

	Server *http.Server
	Mux    *http.ServeMux
}

type TlsKeyPair struct {
	Cert string
	Key  string
}

type PageContents struct {
	Uid      uint32
	Summary  sql.NullString
	Title    string
	Body     template.HTML
	PostDate time.Time
}

type TemplateContent struct {
	Template *template.Template
	Content  interface{}
}

type ContentFn func(r http.Request, dbInfo string) (*TemplateContent, error)

type teacupFnPair struct {
	t *teacup
	fn ContentFn
}

func NewTeacup(port uint16, dbInfo string, pair *TlsKeyPair, fileWhitelist []string, dirBlacklist []string, log *log.Logger) *teacup {
	// create whitelist regex. The final result will look something like:
	// "^/.*\\.(html|css|scss|map|js|png|jpg|gif|webm|ico|md|mp3|mp4|ttf|woff|woff2|eot)$"
	whitelistStr := "^/.*\\.("
	for i, str := range fileWhitelist {
		whitelistStr += str
		if i < len(fileWhitelist) -1 {
			whitelistStr += "|"
		}
	}
	whitelistStr += ")$"

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:           ":"+strconv.Itoa(int(port)),
		Handler:        mux,
	}

	t := teacup{
		dbInfo,
		pair,
		*regexp.MustCompile(whitelistStr),
		dirBlacklist,
		log,
		newTcErrors(),
		make(map[string]bool),
		make(map[string]func(http.Request, string) (*template.Template, interface{})),
		server,
		mux,
	}

	return &t
}

func (t *teacup) StartServer() {
	for _, pattern := range t.DirBlacklist {
		t.Mux.HandleFunc(pattern, t.denyRequest)
	}

	if t.TlsPair == nil {
		t.Log.Fatal(t.Server.ListenAndServe())
	} else {
		t.Log.Fatal(t.Server.ListenAndServeTLS(t.TlsPair.Cert, t.TlsPair.Key))
	}
}

func (t *teacup) ServeFile(writer http.ResponseWriter, request *http.Request) {
	if t.FileWhitelist.MatchString(request.URL.Path) {
		http.ServeFile(writer, request, request.URL.Path[1:])
	} else {
		t.serveError(writer, http.StatusNotFound)
	}
}

func (tFn teacupFnPair) serve(writer http.ResponseWriter, request *http.Request) {
	tc, err := tFn.fn(*request, tFn.t.dbInfo)
	if tc.Content == nil {
		http.Redirect(writer, request, request.URL.Path, http.StatusSeeOther)
		return
	}

	err = tc.Template.Execute(writer, tc.Content)
	if err != nil {
		tFn.t.ServeFile(writer, request)
	}
}

func (t *teacup) HandleFunc(pattern string, fn ContentFn) {
	tFn := teacupFnPair{t, fn}
	t.Mux.HandleFunc(pattern, tFn.serve)
}

func (t *teacup) denyRequest(writer http.ResponseWriter, request *http.Request) {
	t.serveError(writer, http.StatusNotFound)
}

