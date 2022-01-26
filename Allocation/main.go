package main

import(
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"sort"
	"bytes"
	"log"

	"github.com/gorilla/mux"
)

type Bids struct{
	StudentID string
	BidAmount int
	BidStatus string
}
type ClassBids struct{
	ClassCode string
	ClassBids []Bids
}
type ModuleBids struct {
    ModuleCode string
	ModuleName string
    ModuleClasses []ClassBids
}
type SemesterBids struct {
    SemesterStartDate string
    SemesterModules []ModuleBids
}
type BidInfo struct{
	StudentID string
	BidAmount int
	BidStatus string
	ClassCode string
	ModuleCode string
}

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

type TransactionInfo struct{
	ttype string
	sid string
	rid string
	ts string
	tsym string
	ta int
	stat string
}


const BidAPIbaseURL  = "http://bidding_api:8221/api/v1/bids"
const ClassAPIbaseURL =  "http://class:8041/api/v1/classes"
const TransactionAPIbaseURL = "http://localhost:8053/Transaction/new"

func getSemStart(currentDate time.Time)string{
	daysUntilMon := (1 - int(currentDate.Weekday())+7) % 7
	semStartDate := currentDate.AddDate(0,0,daysUntilMon).Format("02-01-2006")
	return semStartDate
}

// 3.15.2: allocate classes by bids
func allocateBid(w http.ResponseWriter, r *http.Request){
	newSem := getSemStart(time.Now())
	var semBids SemesterBids 
	var semClasses Semester 

	// get all bids from 3.14
	resBids,errBids := http.Get(BidAPIbaseURL+"?semester=" + newSem)
	if errBids != nil {
		fmt.Printf("The HTTP request failed with error %s\n", errBids)
	} else{
		if resBids.StatusCode == http.StatusOK{
			data,_ := ioutil.ReadAll(resBids.Body)
			
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
				studentBid := BidInfo{
					StudentID: bids.StudentID,
					BidAmount: bids.BidAmount,
					BidStatus: bids.BidStatus,
					ClassCode: class.ClassCode,
					ModuleCode: module.ModuleCode}

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
		bidStatus := true 
		classApplying := getClass(semClasses,bid)
		if len(classApplying.Students) >= classApplying.Capacity{
			bidStatus = false
		} 
		for _,module := range semClasses.SemesterModules{
			if module.ModuleCode == bid.ModuleCode{
				for _,class := range module.ModuleClasses{
					for _,student := range class.Students{
						if student == bid.StudentID{
							bidStatus = false 
						}
					}
				}
			}
		}
		count := 0
		for _,bid := range allBids{
			
			if bid.ClassCode == classApplying.ClassCode{
				count += 1
			}
		}
		if count < 3 {
			bidStatus = false
		}
		
		if bidStatus{
			// if all allocation checks are successful
			classApplying.Students = append(classApplying.Students,bid.StudentID)
			
			// update bidding status
			updatedBid := Bids{
				StudentID : bid.StudentID,
				BidAmount : bid.BidAmount,
				BidStatus : "Successful",
			}

			bidToUpdate, _ := json.Marshal(updatedBid)

			request, _ := http.NewRequest(http.MethodPut,
				BidAPIbaseURL + "/" + semBids.SemesterStartDate + "?classCode=" + bid.ClassCode,
				bytes.NewBuffer(bidToUpdate))
	
			request.Header.Set("Content-Type", "application/json")
	
			client := &http.Client{}
			_, err := client.Do(request)

			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			}
		}else{
			// refund
			currentDateTime := time.Now()
			formattedDT := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
			currentDateTime.Year(), currentDateTime.Month(), currentDateTime.Day(),
			currentDateTime.Hour(), currentDateTime.Minute(), currentDateTime.Second())

			transactionDetails := TransactionInfo{
				ttype: "BidRefund",
				sid: "BidRefund",
				rid: bid.StudentID,
				ts: formattedDT,
				tsym: "ETI",
				ta: bid.BidAmount,
				stat: "ping",
			}

			transactionToAdd, _ := json.Marshal(transactionDetails)

			response, err := http.Post(TransactionAPIbaseURL,
			"application/json",bytes.NewBuffer(transactionToAdd))

			if err != nil{
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else{
				if response.StatusCode == http.StatusOK{
					fmt.Println("refund tokens to " + bid.StudentID)
				}
			}

			// update bidding status
			updatedBid := Bids{
				StudentID : bid.StudentID,
				BidAmount : bid.BidAmount,
				BidStatus : "Failed",
			}

			bidToUpdate, _ := json.Marshal(updatedBid)

			request, _ := http.NewRequest(http.MethodPut,
				BidAPIbaseURL + "/" + semBids.SemesterStartDate + "?classCode=" + bid.ClassCode,
				bytes.NewBuffer(bidToUpdate))

			request.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			_, errBid := client.Do(request)

			if errBid != nil {
				fmt.Printf("The HTTP request failed with error %s\n", errBid)
			}
		}
		
	}

	// todo: api call 3.8 to set student array
	for _,module := range semClasses.SemesterModules{
		for _,class := range module.ModuleClasses{
			if len(class.Students) > 0{
				classToUpdate,_ := json.Marshal(class)

				request, _ := http.NewRequest(http.MethodPut,
					ClassAPIbaseURL + "/" + semClasses.SemesterStartDate + "?moduleCode=" + module.ModuleCode + "&classCode=" + class.ClassCode,
					bytes.NewBuffer(classToUpdate))
		
				request.Header.Set("Content-Type", "application/json")
		
				client := &http.Client{}
				_, err := client.Do(request)
				if err != nil {
					fmt.Printf("The HTTP request failed with error %s\n", err)
				}
			}else{
				request, _ := http.NewRequest(http.MethodDelete,
					ClassAPIbaseURL+"/"+semClasses.SemesterStartDate+"?moduleCode=" + module.ModuleCode + "&classCode=" + class.ClassCode, nil)
			
				client := &http.Client{}
				_, err := client.Do(request)
			
				if err != nil {
					fmt.Printf("The HTTP request failed with error %s\n", err)
				}
			}
			
		}
	}
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