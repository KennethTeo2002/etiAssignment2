package main

import(
	"fmt"
	"github.com/gorilla/mux"
)

// 3.15.1: give all students 20 ETI credits
func addAll(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the refund REST API!")
	//get list of students

	//for each students, api call 3.12 add 20 eti credit
}

// 3.15.4: refund failed bids
func addSingle(w http.ResponseWriter, r *http.Request){
	// url param: studentID, url query: creditAmount
	params := mux.Vars(r)
	v := r.URL.Query()

	studentID := params["studentID"]

	if creditAmount, ok := v["creditAmount"]; !ok {
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte(
			"422 - Missing credit value "))
		return
	}

	//api call 3.12 add creditAmount to studentID

}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/addCredits",addAll) 
	router.HandleFunc("/addCredits/{studentID}",addSingle) 
	fmt.Println("Listening at port 8071")
	log.Fatal(http.ListenAndServe(":8071", router))
	
}