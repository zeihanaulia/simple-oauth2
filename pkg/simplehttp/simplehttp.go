package simplehttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

func Post(url string, headers map[string]string, bodyRequest []byte) (body []byte, err error) {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyRequest))
	if err != nil {
		return
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}

func HTMLRender(w http.ResponseWriter, path string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func JSONRender(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func Middleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}
