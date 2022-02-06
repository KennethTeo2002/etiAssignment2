package main

import(
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"log"
	"bytes"

	"github.com/gorilla/mux"
)

const StudentAPIbaseURL =  "http://10.31.11.11:8103/api/v1/students"
const TransactionAPIbaseURL = "http://10.31.11.11:8053/Transaction/new"


type TransactionInfo struct{
	Ttype string
	Sid string
	Rid string
	Ts string
	Tsym string
	Ta int
	Stat string
}

// 3.15.1: give all students 20 ETI credits
func addAll(w http.ResponseWriter, r *http.Request){
	var students []string
	//get list of all student id from 3.5
	resStudent,errStudent := http.Get(StudentAPIbaseURL)

	if errStudent != nil {
		fmt.Printf("The HTTP request failed with error %s\n", errStudent)
		// student API fail safe
		students = append(students,"S0001")
		students = append(students,"S0002")

	} else{
		if resStudent.StatusCode == http.StatusOK{
			data,_ := ioutil.ReadAll(resStudent.Body)
			json.Unmarshal([]byte(data), &students)
		} else{
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte(
				"422 - failed to retrieve all students from student API"))
			return
		}
	}

	allTransactionPassed := true
	// loop through all students
	for _,student := range students{
		currentDateTime := time.Now().In(time.FixedZone("UTC+8", 8*60*60))
		formattedDT := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
        currentDateTime.Year(), currentDateTime.Month(), currentDateTime.Day(),
        currentDateTime.Hour(), currentDateTime.Minute(), currentDateTime.Second())

		transactionDetails := TransactionInfo{
			Ttype: "Timetable",
			Sid: "Timetable",
			Rid: student,
			Ts: formattedDT,
			Tsym: "ETI",
			Ta: 20,
			Stat: "ping"}
		
		transactionToAdd, _ := json.Marshal(transactionDetails)

		// send an ETI +20 transaction to 3.12
		_, err := http.Post(TransactionAPIbaseURL,
		"application/json", bytes.NewBuffer(transactionToAdd))
		
		if err != nil{
			fmt.Printf("The HTTP request failed with error %s\n", err)
			allTransactionPassed = false
		} else{
			fmt.Println("Succesfully added tokens to " + student)
			
		}
	}
	// return status code
	if !allTransactionPassed{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - One or more transaction(s) failed to work"))
		return
	} else{
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("All tokens successfully added"))
		return
	}
}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/addCredits",addAll).Methods("GET","POST")
	fmt.Println("Listening at port 8071")
	log.Fatal(http.ListenAndServe(":8071", router))
	
}