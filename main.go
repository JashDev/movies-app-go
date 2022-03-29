package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var movies []Movie

type ErrorResponse struct {
	Message string `json:"message"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type SuccessMoviesResponse struct {
	ResponseMessage
	Movies []Movie `json:"movies"`
}

type SuccessMovieResponse struct {
	ResponseMessage
	Movie Movie `json:"movie"`
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint Hit: getMovies")
	json.NewEncoder(w).Encode(SuccessMoviesResponse{
		ResponseMessage: ResponseMessage{Message: "All movies"},
		Movies:          movies,
	})
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(SuccessMovieResponse{
				ResponseMessage: ResponseMessage{Message: "Movie found"},
				Movie:           item,
			})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ErrorResponse{Message: "Movie not found"})
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(len(movies) + 1)
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(SuccessMovieResponse{
		ResponseMessage: ResponseMessage{Message: "Movie created"},
		Movie:           movie,
	})
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(SuccessMovieResponse{
				ResponseMessage: ResponseMessage{Message: "Movie updated"},
				Movie:           movie,
			})
			return
		}
	}
	json.NewEncoder(w).Encode(ErrorResponse{Message: "Movie not found"})
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(ResponseMessage{Message: "Movie deleted"})
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "123", Title: "The first movie", Director: &Director{FirstName: "John", LastName: "Smith"}})
	movies = append(movies, Movie{ID: "2", Isbn: "456", Title: "The second movie", Director: &Director{FirstName: "John", LastName: "Smith"}})
	movies = append(movies, Movie{ID: "3", Isbn: "789", Title: "The third movie", Director: &Director{FirstName: "John", LastName: "Smith"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server on port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
