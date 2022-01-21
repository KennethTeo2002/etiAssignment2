package main

import(
	"fmt"
	"time"
	"math/rand"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Class struct {
    ClassCode string
    Schedule string
    Tutor    string
    Capacity int32
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

func getSemStart(currentDate time.Time)(time.Time){
	daysUntilMon = (1 - int(currentDate.Weekday())+7) % 7
	semStartDate = currentDate.AddDate(0,0,daysuntilMon).Format("02 Jan 2006")
	return semStartDate
}

const ClassAPIbaseURL =  "http://localhost:0000/api/v1/classes"

func timeTable(w http.ResponseWriter, r *http.Request) {
	// connect to mongoDB cluster
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic(err.Error())
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)
	timeTableDatabase := client.Database("TimeTable")

	params := mux.Vars(r)

	if r.Method == "GET" {
		v := r.URL.Query()
		if semester,ok := v["semester"]; ok {
			semesterCollection = timeTableDatabase.semesterCollection(semester)
		} else{
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte(
				"422 - Missing semester value"))
			return
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
		newSem := getSemStart(time.Now())
		response,err := http.Get(ClassAPIbaseURL+"?semester=" + newSem)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else{
			if response.StatusCode == http.StatusOK{
				data,_ := ioutil.ReadAll(response.Body)
				var sem Semester 
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
			"Friday 09:00 - 11:00","Friday 11:00 - 13:00","Friday 14:00 - 16:00","Friday 16:00 - 18:00"
		}
	
		for _,module := range sem.SemesterModules{
			randomNumber :=  rand.Intn(len(availableTimeSchedule))
			assignedTimeSlot := availableTimeSchedule[randomNumber]
			availableTimeSchedule = append(availableTimeSchedule[:randomNumber],availableTimeSchedule[randomNumber+1:]...)
	
			// for each class
			for _,class := range module.ModuleClasses{
				class.Schedule = assignedTimeSlot
				// add schedule to db
				InsertSchedule(timeTableDatabase, ctx ,class.ClassCode,assignedTimeSlot)
				
				// send put request to set schedule datetime
				classToUpdate,_ := json.Marshal(class)
				request, _ := http.NewRequest(http.MethodPut,
					ClassAPIbaseURL+"/"+sem.SemesterStartDate + "?moduleCode=" + module.ModuleCode + "&classCode=" + class.classCode,
					bytes.NewBuffer(classToUpdate))
				
			}
				
		}
	}
}


func main(){
	router := mux.NewRouter()
	router.HandleFunc("api/timetable",timeTable).Methods(
		"GET", "POST")
	fmt.Println("Listening at port 8072")
	log.Fatal(http.ListenAndServe(":8072", router))
}