package main

import (
	"github.com/gorilla/pat"
	"io/ioutil"
	"net/http"
)

var (
	indexPage []byte
)

func init() {
	var err error
	indexPage, err = ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
}

type command func()

func (cmd command) Handle(w http.ResponseWriter, req *http.Request) {
	cmd()
}

func main() {
	camera, err := NewCamera()
	if err != nil {
		panic(err)
	}
	r := pat.New()
	r.Post("/panLeft", command(func() { camera.PanLeft() }).Handle)
	r.Post("/panRight", command(func() { camera.PanRight() }).Handle)
	r.Post("/tiltUp", command(func() { camera.TiltUp() }).Handle)
	r.Post("/tiltDown", command(func() { camera.TiltDown() }).Handle)
	r.Get("/", DefaultHandler)

	http.Handle("/", r)
	http.ListenAndServe(":8090", nil)
}

func DefaultHandler(w http.ResponseWriter, req *http.Request) {
	w.Write(indexPage)
}
