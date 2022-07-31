package main
import (
    "fmt"
    "html"
    "log"
    "net/http"
    "strconv"
    "sync"	
)

var counter int
var mutex = &sync.Mutex{}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
    mutex.Lock()
    counter++
    fmt.Fprintf(w, strconv.Itoa(counter))
    mutex.Unlock()
}

func main() {

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    })

	http.HandleFunc("/increment", incrementCounter)

    http.HandleFunc("/database", func(w http.ResponseWriter, r *http.Request){
        fmt.Fprintf(w, "we'll store the database here")
    })

    log.Fatal(http.ListenAndServe(":8081", nil))

}