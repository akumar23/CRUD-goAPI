package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func ping(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://api.hatchways.io/assessment/blog/posts")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	response.Body.Close()
	json.NewEncoder(w).Encode(1)

}

type Response struct {
	Tags      string `json:"tags"`
	Sort      string `json:sortBy`
	Direction string `json:direction`
}

func posts(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	tag := params["tags"]
	response, err := http.Get("https://api.hatchways.io/assessment/blog/posts?tag=" + tag)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	//var object Response
	//json.Unmarshal(data, &object)

	var prettyData bytes.Buffer
	error := json.Indent(&prettyData, data, "", "\t")
	if error != nil {
		log.Println(error)
	}

	w.Write(prettyData.Bytes())

}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/api/ping", ping).Methods("GET")
	r.HandleFunc("/api/posts/{tags}", posts).Methods("GET")

	fmt.Printf("Starting server at post 8081\n")
	log.Fatal(http.ListenAndServe(":8081", r))

}
