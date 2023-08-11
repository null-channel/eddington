package notfound

import (
	"encoding/json"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	rb, _ := json.Marshal(NotFound{
		err:     "Not Found",
		path:    r.URL.Path,
		verb:    r.Method,
		rawPath: r.URL.RawPath,
		headers: r.Header,
		uri:     r.RequestURI,
	})

	w.Write(rb)

}

type NotFound struct {
	err     string `json:"error"`
	path    string
	verb    string
	rawPath string
	headers map[string][]string
	uri     string
}
