package main

import(
	"fmt"
	"github.com/gorilla/mux"
)

type bidInfo struct{
	studentID string
	classCode string
	bidAmount int
}

type bidsInfo struct{
	bids []bidInfo
}

const BidAPIbaseURL  = "http://localhost:0000/api/v1"
const ClassAPIbaseURL =  "http://localhost:0000/api/v1"

// 3.15.2: allocate classes by bids
func allocateBid (w http.ResponseWriter, r *http.Request){
	// get all bids from 3.14
	response,err := http.Get(BidAPIbaseURL+"?semester=latest")
	
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else{
		if response.StatusCode == http.StatusOK{
			data,_ := ioutil.ReadAll(response.Body)
			var bidArr bidsInfo 
			json.Unmarshal([]byte(data), &bidArr)
		} else{
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte(
				"422 - failed to retrieve all bids from bids API"))
			return
		}
	}
	// sort bids by descending

	// todo: allocate algo

	// api call 3.8 to set student array
}



func main(){
	router := mux.NewRouter()
	router.HandleFunc("/allocateBid",allocateBid)
	fmt.Println("Listening at port 8072")
	log.Fatal(http.ListenAndServe(":8072", router))
}