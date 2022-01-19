package main

import(
	"fmt"
	"github.com/gorilla/mux"
)

type classInfo struct{
	ID string
	schedule string
}

type classesInfo struct{
	classes []classInfo
}


const ClassAPIbaseURL =  "http://localhost:0000/api/v1"

func timeTable(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Method == "GET" {
		v := r.URL.Query()

		// if semester value is missing, default to latest
		if semester,ok := v["semester"]; !ok {
			semester = "latest"
		}

		if studentID, ok := v["studentID"]; ok {
			// run get student timetable func
			
		} else if tutorID, ok := v["tutorID"]; ok {
			// run get tutor timetable func

		} else{
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte(
				"422 - Missing studentID or tutorID"))
			return
		}
	} else if r.Method == "POST" {
		// allocate class schedule

		// get all classes
		response,err := http.Get(ClassAPIbaseURL+"?semester=latest")

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else{
			if response.StatusCode == http.StatusOK{
				data,_ := ioutil.ReadAll(response.Body)
				var classArr classesInfo 
				json.Unmarshal([]byte(data), &classArr)
			} else{
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte(
					"422 - failed to retrieve all classes from class API"))
				return
			}
		}

		allocateSchedule(classArr)

	}
}

const availableTimeSchedule = [
	"Monday 09:00 - 11:00","Monday 11:00 - 13:00","Monday 14:00 - 16:00","Monday 16:00 - 18:00", 
	"Tuesday 09:00 - 11:00","Tuesday 11:00 - 13:00","Tuesday 14:00 - 16:00","Tuesday 16:00 - 18:00", 
	"Wednesday 09:00 - 11:00","Wednesday 11:00 - 13:00","Wednesday 14:00 - 16:00","Wednesday 16:00 - 18:00", 
	"Thursday 09:00 - 11:00","Thursday 11:00 - 13:00","Thursday 14:00 - 16:00","Thursday 16:00 - 18:00", 
	"Friday 09:00 - 11:00","Friday 11:00 - 13:00","Friday 14:00 - 16:00","Friday 16:00 - 18:00", 

]

func allocateSchedule(classArr classesInfo){
	fmt.Println(classArr)
}


func main(){
	router := mux.NewRouter()
	router.HandleFunc("api/timetable",timeTable).Methods(
		"GET", "POST")
	fmt.Println("Listening at port 8072")
	log.Fatal(http.ListenAndServe(":8072", router))
}