package main

import(
	"fmt"
	"html/template"
	"log"
	"net/http"
	
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("Pages/homepage.html"))
	tmpl.Execute(w, nil)
}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/", homePage)
	fmt.Println("Listening at port 8070")
	log.Fatal(http.ListenAndServe(":8070", router))
}