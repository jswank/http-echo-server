package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("output").Parse(TEMPLATE))
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/{rest:.*}", handler)
	//	router.HandleFunc("/echo", handler)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Printf("starting http-echo-server on :7000")
	log.Fatal(http.ListenAndServe(":7000", loggedRouter))

}

func handler(w http.ResponseWriter, r *http.Request) {

	var b bytes.Buffer

	w.Header().Set("Content-Type", "text/plain")

	err := r.Header.Write(&b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
	}
	headers := b.String()

	b.Reset()

	buf, _ := ioutil.ReadAll(r.Body)
	// b = bytes.NewBuffer(buf)
	body := string(buf)

	b.Reset()
	t := struct {
		Req     *http.Request
		Headers string
		Body    string
	}{
		Req:     r,
		Headers: headers,
		Body:    body,
	}

	err = tmpl.Execute(&b, t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b.Bytes())
	return

}

const TEMPLATE = `{{ .Req.Method }} {{ .Req.URL.String }} {{ .Req.Proto }}
{{ .Headers }}

{{ .Body }}
`
