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
	"strings"
	"os"

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
	ModuleName string
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

const ClassAPIbaseURL =  "http://localhost:8041/api/v1/classes"

var sem Semester 

func timeTable(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		v := r.URL.Query()
		/*
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
		*/
		jsonClass, err := os.Open("sampleClasses.json")
		if err != nil {
			fmt.Println(err)
		}
		byteValue, _ := ioutil.ReadAll(jsonClass)
		err = json.Unmarshal(byteValue, &sem)
		if err != nil {
			fmt.Println(err)
		}

		var timetable []Class
		if studentID, ok := v["studentID"]; ok {
			for _,module := range sem.SemesterModules{
				for _,class := range module.ModuleClasses{
					for _,student := range class.Students{
						if student == studentID[0]{
							classDetails := Class{
								ClassCode: class.ClassCode,
								Schedule: class.Schedule,
								Tutor: class.Tutor,
							}
							timetable = append(timetable,classDetails)
						}
					}
				}
			}	
		} else if tutorID, ok := v["tutorID"]; ok {
			for _,module := range sem.SemesterModules{
				for _,class := range module.ModuleClasses{
					if class.Tutor == tutorID[0]{
						classDetails := Class{
							ClassCode: class.ClassCode,
							Schedule: class.Schedule,
							Tutor: class.Tutor,
						}
						timetable = append(timetable,classDetails)
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
		fmt.Println("Lessons: ",timetable)
		daysOfWeek := []string{
			"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
		}
		// return timetable schedule
		timetableHTML := `<table>
		<tr>
		  <th>Time</th>
		  <th>Monday</th>
		  <th>Tuesday</th>
		  <th>Wednesday</th>
		  <th>Thursday</th>
		  <th>Friday</th>
		</tr>`
		for hour := 8;hour < 18;hour++{
			var suffix string
			if hour >=12{
				suffix = "pm"
			}else{
				suffix = "am"
			}
			var lessonTime int
			if hour >12{
				lessonTime = hour - 12
			}else{
				lessonTime = hour
			}

			timetableHTML += "<tr>"
			timetableHTML += fmt.Sprintf("<th>%d %s</th>", lessonTime,suffix)
			for _,day := range daysOfWeek{
				filled := false
				for _, lesson := range timetable{
					lessonDetails := strings.Split(lesson.Schedule, " ")
					if lessonDetails[0] == day{
						timing,_ := time.Parse("15:04",lessonDetails[1])
						if timing.Hour() == hour{
							timetableHTML += fmt.Sprintf("<td rowspan='2'>%s <br>%s</td>",lesson.ClassCode,lesson.Tutor)
							filled = true
						}else if timing.Hour() == hour +1{
							// skip table data
							filled = true
						}
					}
				}
				if !filled{
					timetableHTML += "<td></td>"
				}
			}

			timetableHTML += "</tr>"
			
		}
		

		timetableHTML += "</table>"

		json.NewEncoder(w).Encode(timetableHTML)

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