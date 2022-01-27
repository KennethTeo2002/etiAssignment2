package main

import(
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"log"
	"bytes"

	"github.com/gorilla/mux"
)

const StudentAPIbaseURL =  "http://localhost:8103/api/v1/students"
const TransactionAPIbaseURL = "http://localhost:8053/Transaction/new"

type StudentInfo struct{
	StudentID string
}

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
	var students []StudentInfo
	//get list of all student id from 3.5
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
	// ------------------------------------- remove --------------------------------------------
	jsonString := 
	`[
		{
			"StudentID":"S001"
		},
		{
			"StudentID":"S002"
		},
		{
			"StudentID":"S004"
		}
	]`
	json.Unmarshal([]byte(jsonString), &students)
	fmt.Println(students)
	// ------------------------------------- remove --------------------------------------------

	allTransactionPassed := true
	// loop through all students and send an ETI +20 transaction to 3.12
	for _,student := range students{
		currentDateTime := time.Now()
		formattedDT := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
        currentDateTime.Year(), currentDateTime.Month(), currentDateTime.Day(),
        currentDateTime.Hour(), currentDateTime.Minute(), currentDateTime.Second())

		transactionDetails := TransactionInfo{
			Ttype: "Timetable",
			Sid: "Timetable",
			Rid: student.StudentID,
			Ts: formattedDT,
			Tsym: "ETI",
			Ta: 20,
			Stat: "ping"}
		
		transactionToAdd, _ := json.Marshal(transactionDetails)

		response, err := http.Post(TransactionAPIbaseURL,
		"application/json", bytes.NewBuffer(transactionToAdd))
		
		if err != nil{
			fmt.Printf("The HTTP request failed with error %s\n", err)
			allTransactionPassed = false
		} else{
			if response.StatusCode == http.StatusOK{
				fmt.Println("Succesfully added tokens to " + student.StudentID)
			}else{
				fmt.Println("Failed to add tokens to " + student.StudentID)
				allTransactionPassed = false
			}
		}
	}
	if !allTransactionPassed{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - One or more transaction(s) failed to work"))
		return
	}
}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/addCredits",addAll) 
	fmt.Println("Listening at port 8071")
	log.Fatal(http.ListenAndServe(":8071", router))
	
}