package main

import(
	"fmt"
	"github.com/gorilla/mux"
)

const StudentAPIbaseURL =  "http://localhost:0000/api/v1"

type studentInfo struct{
	studentID string
}

// 3.15.1: give all students 20 ETI credits
func addAll(w http.ResponseWriter, r *http.Request){
	//get list of students
	response,err := http.Get(StudentAPIbaseURL)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else{
		if response.StatusCode == http.StatusOK{
			data,_ := ioutil.ReadAll(response.Body)
			var students studentInfo[]
			json.Unmarshal([]byte(data), &students)
		} else{
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte(
				"422 - failed to retrieve all students from student API"))
			return
		}
	}
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