package notfound

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("not found!!!")
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	rb, err := json.Marshal(NotFound{
		Err:     "Not Found",
		Path:    r.URL.Path,
		Verb:    r.Method,
		RawPath: r.URL.RawPath,
		Headers: r.Header,
		Uri:     r.RequestURI,
	})

	if err != nil {
		fmt.Println("error with json!")
		fmt.Println(err)
	}

	w.Write(rb)

}

type NotFound struct {
	Err     string              `json:"error"`
	Path    string              `json:"path"`
	Verb    string              `json:"verb"`
	RawPath string              `json:"rawPath"`
	Headers map[string][]string `json:"headers"`
	Uri     string              `json:"uri"`
}
