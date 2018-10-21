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
	Port    uint16
	DbInfo  string
	TlsPair *TlsKeyPair

	FileWhitelist regexp.Regexp
	DirBlacklist  regexp.Regexp

	Log    *log.Logger
	errors tcErrors

	tables   map[string]bool
	webpages map[string]func(http.Request, string) (*template.Template, interface{})
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

func NewTeacup(port uint16, dbInfo string, pair *TlsKeyPair, fileWhitelist regexp.Regexp, dirBlacklist regexp.Regexp, log *log.Logger) *teacup {
	t := teacup{
		port,
		dbInfo,
		pair,
		fileWhitelist,
		dirBlacklist,
		log,
		newTcErrors(),
		make(map[string]bool),
		make(map[string]func(http.Request, string) (*template.Template, interface{})),
	}
	return &t
}

func (t *teacup) StartServer() {
	http.HandleFunc("/", t.matchRequest)
	if t.TlsPair == nil {
		t.Log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(t.Port)), nil))
	} else {
		t.Log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(int(t.Port)), t.TlsPair.Cert, t.TlsPair.Key, nil))
	}
}

func (t *teacup) matchRequest(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path

	switch {
	case t.DirBlacklist.MatchString(path):
		t.serveError(writer, http.StatusForbidden)
	default:
		t.serveContent(writer, request)
	}
}

func (t *teacup) serveFile(writer http.ResponseWriter, request *http.Request) {
	if t.FileWhitelist.MatchString(request.URL.Path) {
		http.ServeFile(writer, request, request.URL.Path[1:])
	} else {
		t.serveError(writer, http.StatusNotFound)
	}
}

func (t *teacup) serveContent(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	for pathRegex := range t.webpages {
		if regexp.MustCompile(pathRegex).MatchString(path) {
			dynPage := t.webpages[pathRegex]
			tmpl, content := dynPage(*request, t.DbInfo)
			if content == nil {
				http.Redirect(writer, request, request.URL.Path, http.StatusSeeOther)
				return
			}

			err := tmpl.Execute(writer, content)
			if t.checkAndLogError(err) {
				t.serveError(writer, http.StatusInternalServerError)
			}
			return
		}
	}

	t.serveFile(writer, request)
}

func (t *teacup) AddTemplateContent(pathRegex string, fn func(http.Request, string) (*template.Template, interface{})) {
	regexp.MustCompile(pathRegex)
	t.webpages[pathRegex] = fn
}

