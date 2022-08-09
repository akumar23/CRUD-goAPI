package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"golang.org/x/net/html"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	ISBN     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connect-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
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
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
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
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func getRatings() []string {
	var titleList []string

	for i := 0; i < len(movies); i++ {
		titleList = append(titleList, movies[i].Title)
	}

	//for j := 0; j < len(titleList); j++ {
	response, err := http.Get("https://letterboxd.com/film/the-dark-knight/")

	if err != nil {
		fmt.Printf("%s", err)
	} else {
		doc, err := html.Parse(response.Body)

		parse_html(doc)
		fmt.Printf("%s", err)
	}

	defer response.Body.Close()
	//}

	return titleList
}

func parse_html(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, element := range n.Attr {
			if element.Key == "href" {
				fmt.Printf("%s\n", element.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse_html(c)
	}

}

func main() {

	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", ISBN: "438", Title: "The Dark Knight", Director: &Director{Firstname: "Christopher", Lastname: "Nolan"}})
	movies = append(movies, Movie{ID: "2", ISBN: "439", Title: "Her", Director: &Director{Firstname: "Spike", Lastname: "Jonez"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at post 8081\n")
	fmt.Print(getRatings())
	log.Fatal(http.ListenAndServe(":8081", r))

}
