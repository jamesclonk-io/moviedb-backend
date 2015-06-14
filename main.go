package main

import (
	"net/http"

	"github.com/JamesClonk/vcap"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jamesclonk-io/moviedb-backend/modules/database"
	"github.com/jamesclonk-io/moviedb-backend/modules/moviedb"
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
	"github.com/jamesclonk-io/stdlib/web/negroni"
)

var (
	log *logrus.Logger
	db  *database.Adapter
	mdb moviedb.MovieDB
)

func init() {
	log = logger.GetLogger()
	databaseSetup()
}

func databaseSetup() {
	var databaseType, databaseUri string

	// get db type
	databaseType = env.Get("JCIO_DATABASE_TYPE", "postgres")

	// check for VCAP_SERVICES first
	data, err := vcap.New()
	if err != nil {
		panic(err)
	}
	if service := data.GetService("moviedb"); service != nil {
		if uri, ok := service.Credentials["uri"]; ok {
			databaseUri = uri.(string)
		}
	}

	// if JCIO_DATABASE_URL is not yet set then try to read it from ENV
	if len(databaseUri) == 0 {
		databaseUri = env.MustGet("JCIO_DATABASE_URI")
	}

	// setup database adapter
	switch databaseType {
	case "postgres":
		db = database.NewPostgresAdapter(databaseUri)
	case "sqlite":
		db = database.NewSQLiteAdapter(databaseUri)
	default:
		log.Fatalf("Invalid database type: %s\n", databaseType)
	}

	// panic if no database adapter was set up
	if db == nil {
		panic("Could not set up database adapter")
	}
	mdb = moviedb.NewMovieDB(db)
}

func main() {
	backend := web.NewBackend()

	// setup API routes
	backend.NewRoute("/movie/{id}", getMovie)
	backend.NewRoute("/movies", getMovies)
	backend.NewRoute("/languages", getLanguages)
	backend.NewRoute("/genres", getGenres)
	backend.NewRoute("/actors", getActors)
	backend.NewRoute("/directors", getDirectors)
	backend.NewRoute("/datecount", getDateCount)
	backend.NewRoute("/statistics", getStatistics)

	n := negroni.Sbagliato()
	n.UseHandler(backend.Router)

	server := web.NewServer()
	server.Start(n)
}

func getMovie(w http.ResponseWriter, req *http.Request) *web.Page {
	vars := mux.Vars(req)
	id := vars["id"]

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
