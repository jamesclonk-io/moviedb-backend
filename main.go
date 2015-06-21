package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jamesclonk-io/moviedb-backend/modules/database"
	"github.com/jamesclonk-io/moviedb-backend/modules/database/migration"
	"github.com/jamesclonk-io/moviedb-backend/modules/moviedb"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
	"github.com/jamesclonk-io/stdlib/web/negroni"
)

var (
	log *logrus.Logger
	mdb moviedb.MovieDB
)

func init() {
	log = logger.GetLogger()
}

func setup() *negroni.Negroni {
	// setup movie database
	adapter := database.NewAdapter()
	migration.RunMigrations("./migrations", adapter)
	mdb = moviedb.NewMovieDB(adapter)

	// create backend service
	backend := web.NewBackend()

	// setup API routes on backend
	backend.NewRoute("/movie/{id}", getMovie).Methods("GET")
	backend.NewSecuredRoute("/movie", postMovie).Methods("POST")
	backend.NewSecuredRoute("/movie/{id}", putMovie).Methods("PUT")
	backend.NewSecuredRoute("/movie/{id}", deleteMovie).Methods("DELETE")

	backend.NewRoute("/movies", getMovies)
	backend.NewRoute("/languages", getLanguages)
	backend.NewRoute("/genres", getGenres)
	backend.NewRoute("/person/{id}", getPerson)
	backend.NewRoute("/actors", getActors)
	backend.NewRoute("/directors", getDirectors)
	backend.NewRoute("/datecount", getDateCount)
	backend.NewRoute("/statistics", getStatistics)

	backend.NewRoute("/500", createError)

	n := negroni.Sbagliato()
	n.UseHandler(backend.Router)

	return n
}

func main() {
	// setup http handler
	n := setup()

	// start backend server
	server := web.NewServer()
	server.Start(n)
}

func postMovie(w http.ResponseWriter, req *http.Request) *web.Page {
	decoder := json.NewDecoder(req.Body)
	var movie moviedb.Movie
	if err := decoder.Decode(&movie); err != nil {
		return web.Error("Error", http.StatusInternalServerError, err)
	}
	if err := mdb.AddMovie(&movie); err != nil {
		return web.Error("Error", http.StatusInternalServerError, err)
	}
	return &web.Page{
		Content: map[string]string{"Result": "OK"},
	}
}

func putMovie(w http.ResponseWriter, req *http.Request) *web.Page {
	return web.Error("Error", http.StatusNotImplemented, errors.New("Not implemented"))
}

func deleteMovie(w http.ResponseWriter, req *http.Request) *web.Page {
	id := mux.Vars(req)["id"]
	rows, err := mdb.DeleteMovie(id)
	if err != nil {
		log.Error(err)
		return web.Error("Error", http.StatusInternalServerError, err)
	}
	return &web.Page{
		Content: map[string]interface{}{"RowsDeleted": rows},
	}
}

func getMovie(w http.ResponseWriter, req *http.Request) *web.Page {
	id := mux.Vars(req)["id"]
	data, err := mdb.GetMovie(id)
	return getData(data, err)
}

func getMovies(w http.ResponseWriter, req *http.Request) *web.Page {
	data, err := mdb.GetMovieListings(moviedb.ParseMovieListingOptions(req))
	return getData(data, err)
}

func getLanguages(w http.ResponseWriter, req *http.Request) *web.Page {
	data, err := mdb.GetLanguages()
	return getData(data, err)
}

func getGenres(w http.ResponseWriter, req *http.Request) *web.Page {
	data, err := mdb.GetGenres()
	return getData(data, err)
}

func getPerson(w http.ResponseWriter, req *http.Request) *web.Page {
	id := mux.Vars(req)["id"]
	data, err := mdb.GetPerson(id)
	return getData(data, err)
}

func getActors(w http.ResponseWriter, req *http.Request) *web.Page {
	data, err := mdb.GetActors()
	return getData(data, err)
}

func getDirectors(w http.ResponseWriter, req *http.Request) *web.Page {
	data, err := mdb.GetDirectors()
	return getData(data, err)
}

func getDateCount(w http.ResponseWriter, req *http.Request) *web.Page {
	data, err := mdb.GetDateCount()
	return getData(data, err)
}

func getStatistics(w http.ResponseWriter, req *http.Request) *web.Page {
	data, err := mdb.GetStatistics()
	return getData(data, err)
}

func getData(data interface{}, err error) *web.Page {
	if err != nil {
		log.Error(err)
		return web.Error("Error", http.StatusInternalServerError, err)
	}
	return &web.Page{
		Content: data,
	}
}

func createError(w http.ResponseWriter, req *http.Request) *web.Page {
	return web.Error("Error", http.StatusInternalServerError, fmt.Errorf("Error!"))
}
