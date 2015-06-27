package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/moviedb-backend/modules/moviedb"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web/negroni"
	"github.com/stretchr/testify/assert"
)

var (
	m                   *negroni.Negroni
	movieTestDbFile     string = "_fixtures/test.db"
	movieTestDbFileCopy string = "_fixtures/test_copy.db"
	testUser            string = "test123"
	testPassword        string = "testpw999"
)

func init() {
	os.Setenv("PORT", "4008")
	logrus.SetOutput(ioutil.Discard)
	logger.GetLogger().Out = ioutil.Discard

	os.Setenv("JCIO_DATABASE_TYPE", "sqlite")
	os.Setenv("JCIO_DATABASE_URI", fmt.Sprintf("sqlite3://%s", movieTestDbFileCopy))
	os.Setenv("JCIO_HTTP_CERT_FILE", "_fixtures/test.cert")
	os.Setenv("JCIO_HTTP_KEY_FILE", "_fixtures/test.key")
	os.Setenv("JCIO_HTTP_AUTH_USER", testUser)
	os.Setenv("JCIO_HTTP_AUTH_PASSWORD", testPassword)

	copyFile(movieTestDbFile, movieTestDbFileCopy)
	m = setup()
}

func copyFile(from, to string) {
	in, err := os.Open(from)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out, err := os.Create(to)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		panic(err)
	}
	if err := out.Close(); err != nil {
		panic(err)
	}
}

func Test_Main_404(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "https://localhost:4008/something", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusNotFound, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `{"Code":"404","Status":"This is not the JSON you are looking for.."}`)
}

func Test_Main_500(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "https://localhost:4008/500", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusInternalServerError, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `"Error!"`)
}

func Test_Main_GetMovie(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "https://localhost:4008/movie/914", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `"id":914,"title":"Argo"`)
	assert.Contains(t, body, `"genres":[{"id":27,"name":"Biography"},{"id":6,"name":"Drama"},{"id":28,"name":"History"},{"id":4,"name":"Thriller"}]`)
	assert.Equal(t, `{"id":914,"title":"Argo","alttitle":{"String":"","Valid":true},"year":2012,"description":"Acting under the cover of a Hollywood producer scouting a location for a science fiction film, a CIA agent launches a dangerous operation to rescue six Americans in Tehran during the U.S. hostage crisis in Iran in 1980.","format":"16:9","length":129,"region":"B","rating":12,"disks":1,"score":5,"picture":"argo.jpg","type":"BluRay","languages":[{"id":1,"name":"Deutsch","country":"Schweiz","native_name":"Deutsch"},{"id":2,"name":"Englisch","country":"USA","native_name":"English"},{"id":3,"name":"Franz\u0026#246;sisch","country":"Frankreich","native_name":"Fran\u0026#231;ais"},{"id":4,"name":"Spanisch","country":"Spanien","native_name":"Espa\u0026#241;ol"}],"genres":[{"id":27,"name":"Biography"},{"id":6,"name":"Drama"},{"id":28,"name":"History"},{"id":4,"name":"Thriller"}],"actors":[{"id":5310,"name":"Alan Arkin"},{"id":331,"name":"Ben Affleck"},{"id":5321,"name":"Bill Tangradi"},{"id":3665,"name":"Bob Gunton"},{"id":3470,"name":"Bryan Cranston"},{"id":4139,"name":"Chris Messina"},{"id":2490,"name":"Christopher Denham"},{"id":5325,"name":"Christopher Stanley"},{"id":40,"name":"Clea DuVall"},{"id":5317,"name":"Farshad Farahat"},{"id":5322,"name":"Jamie McShane"},{"id":942,"name":"John Goodman"},{"id":5319,"name":"Karina Logue"},{"id":5313,"name":"Keith Szarabajka"},{"id":3232,"name":"Kyle Chandler"},{"id":5323,"name":"Matthew Glave"},{"id":5316,"name":"Omid Abtahi"},{"id":4776,"name":"Page Leong"},{"id":5315,"name":"Richard Dillane"},{"id":5314,"name":"Richard Kind"},{"id":5324,"name":"Roberto Garcia"},{"id":5312,"name":"Rory Cochrane"},{"id":5320,"name":"Ryan Ahern"},{"id":1859,"name":"Scoot McNairy"},{"id":5318,"name":"Sheila Vand"},{"id":5311,"name":"Tate Donovan"},{"id":1590,"name":"Titus Welliver"},{"id":3122,"name":"Victor Garber"},{"id":1326,"name":"Zeljko Ivanek"}],"directors":[{"id":331,"name":"Ben Affleck"}]}`, body)

	response = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "https://localhost:4008/movie", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.Contains(t, response.Body.String(), `"Unauthorized!"`)

	response = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "https://localhost:4008/movie", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth("wrong!", "wrong!")

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.Contains(t, response.Body.String(), `"Unauthorized!"`)

	response = httptest.NewRecorder()
	req, err = http.NewRequest("PUT", "https://localhost:4008/movie/914", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.Contains(t, response.Body.String(), `"Unauthorized!"`)

	response = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "https://localhost:4008/movie/1", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.Contains(t, response.Body.String(), `"Unauthorized!"`)
}

func Test_Main_DeleteMovie(t *testing.T) {
	copyFile(movieTestDbFile, movieTestDbFileCopy)
	defer copyFile(movieTestDbFile, movieTestDbFileCopy)

	// first with wrong auth
	response := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "https://localhost:4008/movie/7", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.Contains(t, response.Body.String(), `"Unauthorized!"`)

	response = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "https://localhost:4008/movie/7", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth("wrong!", "wrong!")

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.Contains(t, response.Body.String(), `"Unauthorized!"`)

	// now with correct auth
	response = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "https://localhost:4008/movie/7", nil)
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth(testUser, testPassword)

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `{"RowsDeleted":10}`)

	// is it gone?
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "https://localhost:4008/movie/7", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func Test_Main_AddMovie(t *testing.T) {
	copyFile(movieTestDbFile, movieTestDbFileCopy)
	defer copyFile(movieTestDbFile, movieTestDbFileCopy)

	// is it not there yet?
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "https://localhost:4008/movie/915", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusInternalServerError, response.Code)

	// first with missing auth
	response = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "https://localhost:4008/movie", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.Contains(t, response.Body.String(), `"Unauthorized!"`)

	// now with correct auth
	newMovie := &moviedb.Movie{
		Id:       0,
		Title:    "Super Testfilm",
		Alttitle: sql.NullString{"The ultimate test!", true},
		Year:     2039,
		Score:    3,
		Rating:   16,
		Region:   "1",
		Format:   "16:9",
		Disks:    3,
		Type:     "BluRay",
		Length:   234,
		Picture:  "super_testfilm.jpg",
		Languages: []*moviedb.Language{
			&moviedb.Language{Name: "Deutsch"},
			&moviedb.Language{Name: "1337"},
			&moviedb.Language{Name: "Serbokroatisch"},
		},
		Genres: []*moviedb.Genre{
			&moviedb.Genre{Name: "Deutsche Soap"},
			&moviedb.Genre{Name: "Thriller"},
		},
		Actors: []*moviedb.Person{
			&moviedb.Person{Name: "Brad Pitt"},
			&moviedb.Person{Name: "Edward Norton"},
			&moviedb.Person{Name: "Looize de Testador"},
		},
		Directors: []*moviedb.Person{
			&moviedb.Person{Name: "David Fincher"},
			&moviedb.Person{Name: "Senõr Spielbergo"},
		},
	}
	json, err := json.Marshal(newMovie)
	if err != nil {
		t.Fatal(err)
	}

	response = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "https://localhost:4008/movie", bytes.NewBuffer(json))
	if err != nil {
		t.Error(err)
	}
	req.SetBasicAuth(testUser, testPassword)

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Equal(t, `{"Result":"OK"}`, body)

	// is it there now?
	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "https://localhost:4008/movie/915", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Contains(t, body, `{"id":915,"title":"Super Testfilm"`)
	assert.Equal(t, `{"id":915,"title":"Super Testfilm","alttitle":{"String":"The ultimate test!","Valid":true},"year":2039,"description":"","format":"16:9","length":234,"region":"1","rating":16,"disks":3,"score":3,"picture":"super_testfilm.jpg","type":"BluRay","languages":[{"id":23,"name":"1337","country":"","native_name":""},{"id":1,"name":"Deutsch","country":"Schweiz","native_name":"Deutsch"},{"id":24,"name":"Serbokroatisch","country":"","native_name":""}],"genres":[{"id":34,"name":"Deutsche Soap"},{"id":4,"name":"Thriller"}],"actors":[{"id":7,"name":"Brad Pitt"},{"id":8,"name":"Edward Norton"},{"id":5326,"name":"Looize de Testador"}],"directors":[{"id":11,"name":"David Fincher"},{"id":5327,"name":"Senõr Spielbergo"}]}`, body)
}

func Test_Main_Movies(t *testing.T) {
	copyFile(movieTestDbFile, movieTestDbFileCopy)

	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "https://localhost:4008/movies", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `{"id":1,"title":"Face/Off","year":1997,"score":5,"rating":16}`)
	assert.Contains(t, body, `{"id":135,"title":"James Bond 007:\u003cbr/\u003eOn her Majesty's Secret Service","year":1969,"score":2,"rating":16}`)
	assert.Contains(t, body, `{"id":283,"title":"Sie nannten ihn M\u0026#252;cke","year":1978,"score":3,"rating":12}`)

	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "https://localhost:4008/movies?sort=title&by=desc", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.True(t, strings.HasPrefix(body, `[{"id":151,"title":"Zwei sind nicht zu bremsen","year":1978,"score":3,"rating":12}`))

	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "https://localhost:4008/movies?sort=year&by=desc&sort=title&by=desc", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.True(t, strings.HasPrefix(body, `[{"id":893,"title":"World War Z","year":2013,"score":3,"rating":16}`))

	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "https://localhost:4008/movies?sort=title&by=asc&query=year&value=2013&query=score&value=4", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.True(t, strings.HasPrefix(body, `[{"id":856,"title":"House of Cards (1)","year":2013,"score":4,"rating":16}`))

	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "https://localhost:4008/movies?query=year&value=2013&query=language&value=3&query=genre&value=9&sort=title&by=asc", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.True(t, strings.HasPrefix(body, `[{"id":867,"title":"Kick-Ass 2","year":2013,"score":3,"rating":18}`))

	response = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "https://localhost:4008/movies?query=director&value=331&query=actor&value=2145", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body = response.Body.String()
	assert.Equal(t, `[{"id":647,"title":"The Town","year":2010,"score":4,"rating":16}]`, body)
}

func Test_Main_Languages(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "https://localhost:4008/languages", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `{"id":17,"name":"Koreanisch","country":"Korea","native_name":"Hangul"}`)
	assert.Contains(t, body, `{"id":11,"name":"Thail\u0026#228;ndisch","country":"Thailand","native_name":"Phasa Thai"}`)
}

func Test_Main_Genres(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "https://localhost:4008/genres", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `{"id":28,"name":"History"},{"id":3,"name":"Horror"}`)
	assert.Contains(t, body, `{"id":15,"name":"War"},{"id":10,"name":"Western"}`)
}

func Test_Main_GetPerson(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "https://localhost:4008/person/470", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Equal(t, `{"id":470,"name":"Roger Moore"}`, body)
}

func Test_Main_Actors(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "https://localhost:4008/actors", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `{"id":141,"name":"Alan Rickman"}`)
	assert.Contains(t, body, `{"id":137,"name":"Bruce Willis"}`)
	assert.Contains(t, body, `{"id":254,"name":"Leonardo DiCaprio"}`)
}

func Test_Main_Directors(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "https://localhost:4008/directors", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `{"id":11,"name":"David Fincher"}`)
	assert.Contains(t, body, `{"id":331,"name":"Ben Affleck"}`)
	assert.Contains(t, body, `{"id":417,"name":"John Carpenter"}`)
}

func Test_Main_Statistics(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "https://localhost:4008/statistics", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `"ground_zero":"1999-08-01T00:13:37Z","last_update":"2014-01-01T17:11:36Z","count":912`)
	assert.Contains(t, body, `"movie_types":[{"type":"DVD","disks":1238,"length":154253,"count":602},{"type":"BluRay","disks":493,"length":61691,"count":310}]`)
	assert.Contains(t, body, `"actors":4696,"directors":578,"people_total":5244`)
	assert.Contains(t, body, `"top5_actors":[{"id":483,"name":"Bud Spencer","count":26},{"id":338,"name":"Matt Damon","count":22},{"id":369,"name":"Desmond Llewelyn","count":21},{"id":76,"name":"Cate Blanchett","count":16},{"id":651,"name":"Morgan Freeman","count":16}]`)
	assert.Contains(t, body, `"top5_directors":[{"id":493,"name":"Kenji Kamiyama","count":17},{"id":127,"name":"George Lucas","count":11},{"id":70,"name":"Peter Jackson","count":11},{"id":64,"name":"Steven Spielberg","count":11},{"id":121,"name":"Ridley Scott","count":9}]`)
	assert.Contains(t, body, `"top5_actors_and_directors":[{"id":396,"name":"Clint Eastwood","count":20},{"id":189,"name":"Mel Gibson","count":14},{"id":220,"name":"George Clooney","count":12},{"id":211,"name":"Quentin Tarantino","count":12},{"id":331,"name":"Ben Affleck","count":11}]`)
	assert.Contains(t, body, `"regions":[{"type":"0","count":12},{"type":"1","count":74},{"type":"2","count":514},{"type":"3","count":1},{"type":"B","count":311}]`)
	assert.Contains(t, body, `"ratings":[{"type":"21","count":12},{"type":"18","count":132},{"type":"16","count":397},{"type":"12","count":304},{"type":"6","count":67}]`)
	assert.Contains(t, body, `"scores":[{"type":"5","count":92},{"type":"4","count":214},{"type":"3","count":520},{"type":"2","count":81},{"type":"1","count":5}]`)
	assert.Contains(t, body, `"dvd_movies":602,"bluray_movies":310,"dvd_disks":1238,"bluray_disks":493,"total_length":215944,"avg_length_per_movie":236,"avg_length_per_disk":124`)
	assert.Contains(t, body, `avg_movies_per_day`)
	assert.Contains(t, body, `new_movies_estimate`)
}
