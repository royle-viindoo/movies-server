package main

import (
	"encoding/json"
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
	vars := mux.Vars(r)
	id := vars["id"]
	result := Movie{}
	//err := db.collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	err := db.collection.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetches movie"))
		log.Printf("Error fetches movie: %v", err)
		return
	}
	jsonResp, err := json.Marshal(&result)
	//jsonResp, err := json.MarshalIndent(&result, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetches movie"))
		log.Printf("Error marshall movie: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// GetMovies method fetches movies
func (db *DB) GetMovies(w http.ResponseWriter, r *http.Request) {
	results := []Movie{}
	err := db.collection.Find(bson.M{}).All(&results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetches movies"))
		log.Printf("Error fetches movies: %v", err)
		return
	}
	jsonResp, err := json.Marshal(&results)
	//jsonResp, err := json.MarshalIndent(&results, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetches movies"))
		log.Printf("Error marshall movies: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// PostMovie method adds a new movie
func (db *DB) PostMovie(w http.ResponseWriter, r *http.Request) {

}

// UpdateMovie modifies the data of given movie
func (db *DB) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	data := Movie{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error update movie"))
		log.Printf("Error unmarshall movie: %v", err)
		return
	}
	err = db.collection.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": &data})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error update movie"))
		log.Printf("Error update movie: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Updated succesfully!"))
}

// DeleteMovie removes the data of given movie
func (db *DB) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := db.collection.Remove(bson.M{"_id": bson.ObjectIdHex(vars["id"])})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error delete movie"))
		log.Printf("Error delete movie: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted succesfully!"))
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
	sub.HandleFunc("/movies/{id:[a-zA-Z0-9]*}", db.GetMovie).Methods("GET")
	sub.HandleFunc("/movies/{id:[a-zA-Z0-9]*}", db.UpdateMovie).Methods("PUT")
	sub.HandleFunc("/movies/{id:[a-zA-Z0-9]*}", db.DeleteMovie).Methods("DELETE")
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
