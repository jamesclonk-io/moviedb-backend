package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
	"github.com/jamesclonk-io/stdlib/web/negroni"
)

var (
	log *logrus.Logger
)

func init() {
	log = logger.GetLogger()
}

func main() {
	backend := web.NewBackend()
	backend.NewRoute("/", index)

	n := negroni.Sbagliato()
	n.UseHandler(backend.Router)

	server := web.NewServer()
	server.Start(n)
}

func index(w http.ResponseWriter, req *http.Request) *web.Page {
	return &web.Page{
		Content: []string{"foo", "bar"},
	}
}
