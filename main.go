package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// DB struct stores the dabase session imformation.
// Need to be initialized once.
type DB struct {
	session    *mgo.Session
	collection *mgo.Collection
}

// Movie struct hols a movie data
type Movie struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	Year      string        `bson:"year"`
	Directors []string      `bson:"directors"`
	Writers   []string      `bson:"writers"`
	BoxOffice `bson:"boxOffice"`
}

// BoxOffice is nested in Movie struct
type BoxOffice struct {
	Budget uint64 `bson:"budget"`
	Gross  uint64 `bson:"gross"`
}

// GetMovie method fetches a movie with given ID
func (db *DB) GetMovie(w http.ResponseWriter, r *http.Request) {

}

// GetMovies method fetches movies
func (db *DB) GetMovies(w http.ResponseWriter, r *http.Request) {

}

// PostMovie method adds a new movie
func (db *DB) PostMovie(w http.ResponseWriter, r *http.Request) {

}

// UpdateMovie modifies the data of given movie
func (db *DB) UpdateMovie(w http.ResponseWriter, r *http.Request) {

}

// DeleteMovie removes the data of given movie
func (db *DB) DeleteMovie(w http.ResponseWriter, r *http.Request) {

}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Fatalf("Error connecting MongoDB: %s", err)
	}
	defer session.Close()
	c := session.DB("movies").C("movies")
	db := &DB{session: session, collection: c}

	r := mux.NewRouter()
	sub := r.PathPrefix("/api/v1").Subrouter()
	sub.HandleFunc("/movies/{id:[0-9]*}", db.GetMovie).Methods("GET")
	sub.HandleFunc("/movies/{id:[0-9]*}", db.UpdateMovie).Methods("PUT")
	sub.HandleFunc("/movies/{id:[0-9]*}", db.DeleteMovie).Methods("DELETE")
	sub.HandleFunc("/movies", db.GetMovies).Methods("GET")
	sub.HandleFunc("/movies", db.PostMovie).Methods("POST")

	logFile, _ := os.OpenFile("server.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        handlers.CombinedLoggingHandler(logFile, r),
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
