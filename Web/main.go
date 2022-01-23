package main

import(
	"fmt"
	"html/template"
	"log"
	"net/http"
	
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("Pages/homepage.html"))
		tmpl.Execute(w, nil)
	}
}

func addTokens(w http.ResponseWriter, r *http.Request){
	fmt.Println("add tokens")
	_, err := http.Post("http://localhost:8071/addCredits",
	"application/json",nil)

	if err != nil{
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else{
		fmt.Println("add tokens")
	}

	http.Redirect(w, r, "/",http.StatusFound)
	
}

func allocateSchedule(w http.ResponseWriter, r *http.Request){
	fmt.Println("test")
	_, err := http.Post("http://localhost:8072/api/timetable","application/json",nil)
	if err != nil{
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else{
		fmt.Println("allocate class timings")
	}
	http.Redirect(w, r, "/",http.StatusFound)
}

func allocateBids(w http.ResponseWriter, r *http.Request){
	_, err := http.Post("http://localhost:8073/allocateBid","application/json",nil)
	if err != nil{
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else{
		fmt.Println("allocate classes to students")
	}
	http.Redirect(w, r, "/",http.StatusFound)
}

func getTimeTable(w http.ResponseWriter, r *http.Request){
	_,err := http.Get("http://localhost:8072/api/timetable" + "?studentID=")
	if err != nil{
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else{
		fmt.Println("allocate classes to students")
	}
}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/", homePage)
	router.HandleFunc("/addTokens", addTokens)
	router.HandleFunc("/allocateSchedule", allocateSchedule)
	router.HandleFunc("/allocateBids", allocateBids)

	fmt.Println("Listening at port 8070")
	log.Fatal(http.ListenAndServe(":8070", router))
}