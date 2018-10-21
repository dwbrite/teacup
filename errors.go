package teacup

import (
	"html/template"
	"log"
	"net/http"
)

type errorPage struct {
	Status  int
	Message string
}

type tcErrors struct {
	template      *template.Template
	errorMessages map[int]string
}

func newTcErrors() tcErrors {
	return tcErrors{
		template.Must(template.New("teacup_error").Parse(`
{{ define "teacup_error" }}
<html>
<head>
    <style>
body {
    text-align: center;
    margin-top: 240px;
}
    </style>
</head>
<body>
<h1>Error {{ .Status }}</h1>
<h2>{{ .Message }}</h2>
<footer>
    Developed by <a href="https://github.com/dwbrite/teacup">Devin Brite</a>.
</footer>
</body>
</html>
{{ end }}
`)),
		map[int]string{
			http.StatusBadRequest:                   "Bad Request",
			http.StatusUnauthorized:                 "Unauthorized",
			http.StatusPaymentRequired:              "Payment Required",
			http.StatusForbidden:                    "Forbidden",
			http.StatusNotFound:                     "Not Found",
			http.StatusMethodNotAllowed:             "Method Not Allowed",
			http.StatusNotAcceptable:                "Not Acceptable",
			http.StatusProxyAuthRequired:            "Proxy Authentication Required",
			http.StatusRequestTimeout:               "Request Timeout",
			http.StatusConflict:                     "Conflict",
			http.StatusGone:                         "Gone",
			http.StatusLengthRequired:               "Length Required",
			http.StatusPreconditionFailed:           "Precondition Failed",
			http.StatusRequestEntityTooLarge:        "Request Entity Too Large",
			http.StatusRequestURITooLong:            "Request URI Too Long",
			http.StatusUnsupportedMediaType:         "Unsupported Media Type",
			http.StatusRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
			http.StatusExpectationFailed:            "Expectation Failed",
			http.StatusTeapot:                       "I'm a teapot",
			//http.StatusMisdirectedRequest:           "Misdirected Request",
			http.StatusUnprocessableEntity:          "Unprocessable Entity",
			http.StatusLocked:                       "Locked",
			http.StatusFailedDependency:             "Failed Dependency",
			http.StatusUpgradeRequired:              "Upgrade Required",
			http.StatusPreconditionRequired:         "Precondition Required",
			http.StatusTooManyRequests:              "Too Many Requests",
			http.StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
			http.StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

			http.StatusInternalServerError:           "Internal Server Error",
			http.StatusNotImplemented:                "Not Implemented",
			http.StatusBadGateway:                    "Bad Gateway",
			http.StatusServiceUnavailable:            "Service Unavailable",
			http.StatusGatewayTimeout:                "Gateway Timeout",
			http.StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
			http.StatusVariantAlsoNegotiates:         "Variant Also Negotiates",
			http.StatusInsufficientStorage:           "Insufficient Storage",
			http.StatusLoopDetected:                  "Loop Detected",
			http.StatusNotExtended:                   "Not Extended",
			http.StatusNetworkAuthenticationRequired: "Network Authentication Required",
		},
	}
}

func (t *teacup) SetErrorTemplate(tmpl *template.Template) {
	t.errors.template = tmpl
}

func (t *teacup) SetErrorText(code int, media string) {
	t.errors.errorMessages[code] = media
}

func (t *teacup) serveError(writer http.ResponseWriter, httpStatus int) {
	writer.WriteHeader(httpStatus)
	err := t.errors.template.Execute(writer,
		errorPage{
			httpStatus,
			t.errors.errorMessages[httpStatus],
		})

	t.checkAndLogError(err)
}

func (t teacup) checkAndLogError(err error) bool {
	isError := err != nil
	if isError {
		log.Println(err)
		t.Log.Println(err)
	}
	return isError
}
