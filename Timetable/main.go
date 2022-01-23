package main

import(
	"fmt"
	"bytes"
	"time"
	"net/http"
	"log"
	"math/rand"
	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/mux"
)

type Class struct {
    ClassCode string
    Schedule string
    Tutor    string
    Capacity int
    Students []string
}

type Module struct {
    ModuleCode string
    ModuleClasses []Class
}

type Semester struct {
    SemesterStartDate string
    SemesterModules []Module
}

func getSemStart(currentDate time.Time)string{
	daysUntilMon := (1 - int(currentDate.Weekday())+7) % 7
	semStartDate := currentDate.AddDate(0,0,daysUntilMon).Format("02-01-2006")
	return semStartDate
}

const ClassAPIbaseURL =  "http://localhost:0000/api/v1/classes"

var sem Semester 

func timeTable(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "GET" {
		
		v := r.URL.Query()
		if semester,ok := v["semester"]; ok {
			response,err := http.Get(ClassAPIbaseURL+"?semester=" + semester[0])
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else{
				if response.StatusCode == http.StatusOK{
					data,_ := ioutil.ReadAll(response.Body)
					
					json.Unmarshal([]byte(data), &sem)
				} else{
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - failed to retrieve all classes from class API"))
					return
				}
			}
		} else{
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte(
				"422 - Missing semester value"))
			return
		}

		if studentID, ok := v["studentID"]; ok {
			timetable := []string{}
			for _,module := range sem.SemesterModules{
				for _,class := range module.ModuleClasses{
					for _,student := range class.Students{
						if student == studentID[0]{
							timetable = append(timetable,class.Schedule)
						}
					}
				}
			}	
		} else if tutorID, ok := v["tutorID"]; ok {
			timetable := []string{}
			for _,module := range sem.SemesterModules{
				for _,class := range module.ModuleClasses{
					if class.Tutor == tutorID[0]{
						timetable = append(timetable,class.Schedule)
					}
				}
			}

		} else{
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte(
				"422 - Missing studentID or tutorID"))
			return
		}

		// todo: generate timetable
		



	} else if r.Method == "POST" {
	// allocate class schedule
		newSem := getSemStart(time.Now())
		// get all classes
		response,err := http.Get(ClassAPIbaseURL+"?semester=" + newSem)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else{
			if response.StatusCode == http.StatusOK{
				data,_ := ioutil.ReadAll(response.Body)
				 
				json.Unmarshal([]byte(data), &sem)
			} else{
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte(
					"422 - failed to retrieve all classes from class API"))
				return
			}
		}
	

		availableTimeSchedule := []string{
			"Monday 09:00 - 11:00","Monday 11:00 - 13:00","Monday 14:00 - 16:00","Monday 16:00 - 18:00", 
			"Tuesday 09:00 - 11:00","Tuesday 11:00 - 13:00","Tuesday 14:00 - 16:00","Tuesday 16:00 - 18:00", 
			"Wednesday 09:00 - 11:00","Wednesday 11:00 - 13:00","Wednesday 14:00 - 16:00","Wednesday 16:00 - 18:00", 
			"Thursday 09:00 - 11:00","Thursday 11:00 - 13:00","Thursday 14:00 - 16:00","Thursday 16:00 - 18:00", 
			"Friday 09:00 - 11:00","Friday 11:00 - 13:00","Friday 14:00 - 16:00","Friday 16:00 - 18:00"}
	
		for _,module := range sem.SemesterModules{
			randomNumber :=  rand.Intn(len(availableTimeSchedule))
			assignedTimeSlot := availableTimeSchedule[randomNumber]
			availableTimeSchedule = append(availableTimeSchedule[:randomNumber],availableTimeSchedule[randomNumber+1:]...)
	
			// for each class
			for _,class := range module.ModuleClasses{
				class.Schedule = assignedTimeSlot	
				// send put request to set schedule datetime
				
				classToUpdate,_ := json.Marshal(class)
				_, err := http.NewRequest(http.MethodPut,
					ClassAPIbaseURL+"/"+sem.SemesterStartDate + "?moduleCode=" + module.ModuleCode + "&classCode=" + class.ClassCode,
					bytes.NewBuffer(classToUpdate))
				
				if err != nil {
					fmt.Printf("The HTTP request failed with error %s\n", err)
				}			
			}
				
		}
	}
}


func main(){
	router := mux.NewRouter()
	router.HandleFunc("/api/timetable",timeTable).Methods(
		"GET", "POST")
	fmt.Println("Listening at port 8072")
	log.Fatal(http.ListenAndServe(":8072", router))
}