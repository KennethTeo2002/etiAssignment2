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

const StudentAPIbaseURL =  "http://localhost:0000/api/v1"
const TransactionAPIbaseURL = "http://localhost:8053/Transaction/new"

type StudentInfo struct{
	StudentID string
}

type TransactionInfo struct{
	ttype string
	sid string
	rid string
	ts string
	tsym string
	ta int
	stat string
}

var students []StudentInfo

// 3.15.1: give all students 20 ETI credits
func addAll(w http.ResponseWriter, r *http.Request){
	//get list of students
	resStudent,errStudent := http.Get(StudentAPIbaseURL)

	if errStudent != nil {
		fmt.Printf("The HTTP request failed with error %s\n", errStudent)
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

	for _,student := range students{

		currentDateTime := time.Now()
		formattedDT := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
        currentDateTime.Year(), currentDateTime.Month(), currentDateTime.Day(),
        currentDateTime.Hour(), currentDateTime.Minute(), currentDateTime.Second())

		transactionDetails := TransactionInfo{
			ttype: "Timetable",
			sid: "Timetable",
			rid: student.StudentID,
			ts: formattedDT,
			tsym: "ETI",
			ta: 20,
			stat: "ping"}
		
		transactionToAdd, _ := json.Marshal(transactionDetails)

		response, err := http.Post(TransactionAPIbaseURL,
		"application/json", bytes.NewBuffer(transactionToAdd))

		if err != nil{
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else{
			if response.StatusCode == http.StatusOK{
				fmt.Println("add tokens to " + student.StudentID)
			}
		}
	}

	
}


func main(){
	router := mux.NewRouter()
	router.HandleFunc("/addCredits",addAll) 
	fmt.Println("Listening at port 8071")
	log.Fatal(http.ListenAndServe(":8071", router))
	
}