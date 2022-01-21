package main

import(
	"fmt"
	"github.com/gorilla/mux"
)

type Bids struct{
	StudentID string
	BidAmount int
}

type ClassBids struct{
	ClassCode string
	ClassBids []Bids
}

type ModuleBids struct {
    ModuleCode string
    ModuleClasses []ClassBids
}

type SemesterBids struct {
    SemesterStartDate string
    SemesterModules []ModuleBids
}

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

type BidInfo struct{
	StudentID string
	BidAmount int
	ClassCode string
	ModuleCode string
}

const BidAPIbaseURL  = "http://localhost:0000/api/v1"
const ClassAPIbaseURL =  "http://localhost:0000/api/v1/classes"

func getSemStart(currentDate time.Time){
	daysUntilMon = (1 - int(currentDate.Weekday())+7) % 7
	semStartDate = currentDate.AddDate(0,0,daysuntilMon).Format("02 Jan 2006")
	return semStartDate
}

// 3.15.2: allocate classes by bids
func allocateBid(w http.ResponseWriter, r *http.Request){
	newSem = getSemStart(time.Now())
	// get all bids from 3.14
	resBids,errBids := http.Get(BidAPIbaseURL+"?semester=" + newSem)
	
	if errBids != nil {
		fmt.Printf("The HTTP request failed with error %s\n", errBids)
	} else{
		if resBids.StatusCode == http.StatusOK{
			data,_ := ioutil.ReadAll(resBids.Body)
			var semBids SemesterBids 
			json.Unmarshal([]byte(data), &semBids)
		} else{
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte(
				"422 - failed to retrieve all bids from bids API"))
			return
		}
	}
	// get all classes from 3.8
	resClass,errClass := http.Get(ClassAPIbaseURL+"?semester=" + newSem)
	
	if errClass != nil {
		fmt.Printf("The HTTP request failed with error %s\n", errClass)
	} else{
		if resClass.StatusCode == http.StatusOK{
			data,_ := ioutil.ReadAll(resClass.Body)
			var semClasses Semester 
			json.Unmarshal([]byte(data), &semClasses)
		} else{
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte(
				"422 - failed to retrieve all classes from classes API"))
			return
		}
	}

	allBids := []BidInfo{}
	for _,module := range semBids.SemesterModules{
		for _,class :=  range module.ModuleClasses{
			for _,bids := range class.ClassBids{
				studentBid = BidInfo{
					StudentID: bids.StudentID,
					BidAmount: bids.BidAmount,
					ClassCode: class.ClassCode,
					ModuleCode module.ModuleCode
				}

				allBids = append(allBids,studentBid)
			}
		}
	}

	// sort bids by descending
	sort.Slice(allBids, func(i, j int) bool {
		return allBids[i].BidAmount > allBids[j].BidAmount
	  })


	// Allocation algo
	for _,bid := range allBids{
		classApplying := getClass(semClasses,bid)
		if len(classApplying.Students) < classApplying.Capacity{
			classApplying.Students = append(classApplying.Students,bid.StudentID)
		} else{
			// todo: refund
		}
		
	}
	// todo: check if class students < 3, delete class and refund

	// api call 3.8 to set student array

}

func getClass(semClasses Semester, bid BidInfo) Class{
	for _,module := range semClasses.SemesterModules{
		if module.ModuleCode == bid.ModuleCode{
			for _,class := range module.ModuleClasses{
				if class.ClassCode == bid.ClassCode{
					return class
				}
			}
		}
	}
	return Class{}
}


func main(){
	router := mux.NewRouter()
	router.HandleFunc("/allocateBid",allocateBid)
	fmt.Println("Listening at port 8073")
	log.Fatal(http.ListenAndServe(":8073", router))
}