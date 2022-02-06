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
	BidAmount int `json:"BidAmt"`
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


const BidAPIbaseURL  = "http://10.31.11.11:8221/api/v1/bid"
const ClassAPIbaseURL =  "http://10.31.11.11:8041/api/v1/classes"
const TransactionAPIbaseURL = "http://10.31.11.11:8053/Transaction/new"

func getSemStart(currentDate time.Time)string{
	daysUntilMon := (1 - int(currentDate.Weekday())+7) % 7
	semStartDate := currentDate.AddDate(0,0,daysUntilMon).Format("02-01-2006")
	return semStartDate
}

// 3.15.2: allocate classes by bids
func allocateBid(w http.ResponseWriter, r *http.Request){
	var semBids SemesterBids 
	var semClasses Semester 

	// get all bids from 3.14
	
	newSem := getSemStart(time.Now().In(time.FixedZone("UTC+8", 8*60*60)))
	resBids,errBids := http.Get(BidAPIbaseURL+ "/" + newSem)
	if errBids != nil {
		fmt.Printf("The HTTP request failed with error %s\n", errBids)
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte(
			"422 - failed to retrieve all bids from bids API"))
		return
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
	
	resClass,errClass := http.Get(ClassAPIbaseURL + "/"+ newSem)
	if errClass != nil {
		fmt.Printf("The HTTP request failed with error %s\n", errClass)
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte(
			"422 - failed to retrieve all classes from classes API"))
		return
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

	// unpack bidding struct into 1d array
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

	// sort bids by descending bid amount
	sort.SliceStable(allBids, func(i, j int) bool {
		return allBids[i].BidAmount > allBids[j].BidAmount
	  })

	// loop through each bid 
	for _,bid := range allBids{
		bidStatus := true 
		var classApplying *Class
		for moduleIndex,module := range semClasses.SemesterModules{
			if module.ModuleCode == bid.ModuleCode{
				for classIndex,class := range module.ModuleClasses{
					if class.ClassCode == bid.ClassCode{
						classApplying = &semClasses.SemesterModules[moduleIndex].ModuleClasses[classIndex]
					}
				}
			}
		}
		
		// Allocation algorithm checks
		// check if class is already full
		if len(classApplying.Students) >= classApplying.Capacity{
			bidStatus = false
		} 
		// check if student already assigned to another class in same module
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
		// check if class bids less than 3
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
			(*classApplying).Students = append(classApplying.Students, bid.StudentID)
			
			// update bidding status
			err := updateBid(bid,semBids,"Successful")
			if err{
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - failed to update bid api"))
				return
			}

		}else{
			// refund
			transactionError := refundTransaction(bid)
			if transactionError{
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - failed to send transaction api"))
				return
			}

			// update bidding status
			bidError := updateBid(bid,semBids,"Failed")
			if bidError{
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - failed to update bid api"))
				return
			}
		}
		
	}

	classUpdateErr := false

	// api call 3.8 to set student array
	for _,module := range semClasses.SemesterModules{
		for _,class := range module.ModuleClasses{
			fmt.Println(class)
			if len(class.Students) >= 0{
				// class object has students assigned
				classToUpdate,_ := json.Marshal(class)

				request, _ := http.NewRequest(http.MethodPut,
					ClassAPIbaseURL + "/" + semClasses.SemesterStartDate + "?moduleCode=" + module.ModuleCode + "&classCode=" + class.ClassCode,
					bytes.NewBuffer(classToUpdate))
		
				request.Header.Set("Content-Type", "application/json")
		
				client := &http.Client{}
				response, err := client.Do(request)
				if err != nil {
					fmt.Printf("The HTTP request failed with error %s\n", err)
					classUpdateErr = true
				}else{
					if response.StatusCode != http.StatusOK{
						classUpdateErr = true
					}
				}
			}else{
				// if no one in class object, delete
				request, _ := http.NewRequest(http.MethodDelete,
					ClassAPIbaseURL+"/"+semClasses.SemesterStartDate+"?moduleCode=" + module.ModuleCode + "&classCode=" + class.ClassCode, nil)
			
				client := &http.Client{}
				response, err := client.Do(request)
			
				if err != nil {
					fmt.Printf("The HTTP request failed with error %s\n", err)
					classUpdateErr = true
				}else{
					if response.StatusCode != http.StatusOK{
						classUpdateErr = true
					}
				}
			}
			
		}
	}

	if classUpdateErr{
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - failed to update class api"))
		return
	}else{
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 - allocated all students to classes"))
		return
	}
}

func refundTransaction(bid BidInfo)bool{
	apiErr := true
	currentDateTime := time.Now().In(time.FixedZone("UTC+8", 8*60*60))
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
			apiErr = false
		}
	}

	return apiErr
}

func updateBid(bid BidInfo, semBids SemesterBids, bidStatus string)bool{
	apiErr := true
	updatedBid := Bids{
		StudentID : bid.StudentID,
		BidAmount : bid.BidAmount,
		BidStatus : bidStatus,
	}

	bidToUpdate, _ := json.Marshal(updatedBid)

	request, _ := http.NewRequest(http.MethodPut,
		BidAPIbaseURL + "/" + semBids.SemesterStartDate + "?classCode=" + bid.ClassCode,
		bytes.NewBuffer(bidToUpdate))

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resBid, errBid := client.Do(request)

	if errBid != nil {
		fmt.Printf("The HTTP request failed with error %s\n", errBid)
	}else{
		if resBid.StatusCode == http.StatusOK{
			apiErr = false
		}
	}
	return apiErr
}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/allocateBid",allocateBid).Methods("GET","POST")
	fmt.Println("Listening at port 8072")
	log.Fatal(http.ListenAndServe(":8072", router))
}