package main

import(
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gorilla/mux"
)

type BusStopInfo struct{
	BusStopCode string
	Description string
}

var busStops map[string]BusStopInfo

func busStop(w http.ResponseWriter, r *http.Request){

	if r.Method == "GET" {
		params := mux.Vars(r)
        busCode := params["busStopCode"]
        if busCode == "" {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("400 - No Bus Code was provided"))
            return
        }
        // check if code exists
        if _, ok := busStops[busCode]; ok {
			retrievedBusStop := busStops[busCode]
            json.NewEncoder(w).Encode(retrievedBusStop)
        } else {
            w.WriteHeader(http.StatusNotFound)
        }
    } 
	if r.Method == "DELETE" {
		params := mux.Vars(r)
		busCode := params["busStopCode"]
		if busCode == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("404 - No Bus Code was provided"))
			return
		}
		// check if code exists
		if _, ok := busStops[busCode]; ok {
			delete(busStops, busCode)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("200 - BusStop Deleted: " + 
			busCode))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
    }

	if r.Header.Get("Content-type") == "application/json" {
		var newBusStop BusStopInfo
		reqBody, err := ioutil.ReadAll(r.Body)

		if err == nil {
			json.Unmarshal(reqBody, &newBusStop)
			
			if newBusStop.BusStopCode == "" || newBusStop.Description == ""{
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte(
				   "422 - Please supply busStop details "))
				   return
			}
			if r.Method == "POST"{
				if _, ok := busStops[newBusStop.BusStopCode]; !ok{
					busStops[newBusStop.BusStopCode] = newBusStop
					w.WriteHeader(http.StatusCreated)
                    w.Write([]byte("201 - BusStop added: " + 
					newBusStop.BusStopCode))

				}else{
					w.WriteHeader(http.StatusConflict)
                    w.Write([]byte(
                        "409 - Duplicate busStop code"))
				}

			}
			if r.Method == "PUT"{
				if _,ok := busStops[newBusStop.BusStopCode]; !ok{
					busStops[newBusStop.BusStopCode] = newBusStop
					w.WriteHeader(http.StatusCreated)
                    w.Write([]byte("201 - Bus Stop added: " + 
					newBusStop.BusStopCode))
				} else{
					busStops[newBusStop.BusStopCode] = newBusStop
                    w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 -  Bus Stop updated: " +
					newBusStop.BusStopCode))

				}
			}
		}else{
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply " +
				"bus stop information " +
				"in JSON format"))

		}
	}
}

func main(){
	busStops = make(map[string]BusStopInfo)
	router := mux.NewRouter()
	router.HandleFunc("/v1/BusStops/{busStopCode}",busStop).Methods("GET","PUT", "POST","DELETE")
	fmt.Println("Listening at port 5040")
	log.Fatal(http.ListenAndServe(":5040", router))
}