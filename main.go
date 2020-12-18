package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var layout = "2006-01-02T15:04:05.000Z"
var currentTimeUUID string
var firstLogs = false
var byTID = false
var byService = false
var jwt string
var follow = false
var url string
var service string
var transactionId string
var startTime string
var endTime string
var showTimeUUID = false


func main() {
	args := os.Args[1:]
	getArgs(args)

	if url == "" {
		log.Fatal("No url provided")
	}

	client := &http.Client{}

	if byTID {
		initGetLogs(url, transactionId, *client)
	} else if byService {
		initGetLogs(url, service, *client)
	} else {
		initGetLogs(url, "", *client)
	}

	if follow {
		for {
			endTime = time.Now().UTC().Format(layout)
			time.Sleep(10*time.Second)
			if byTID {
				getLogs(url, transactionId, *client)
			} else if byService {
				getLogs(url, service, *client)
			} else {
				getLogs(url, "", *client)
			}
		}
	}
}

func initGetLogs(url string, optional string, client http.Client) {
	ts := timespan{
		StartTime: time.Now().UTC().Add(-2*time.Minute).Format(layout),
	}
	if startTime != ""{
		ts.StartTime = startTime
	}
	if endTime != "" {
		ts.EndTime = endTime
	}


	if byTID {
		ts.TransactionId = optional
	}
	if byService {
		ts.Service = optional
	}

	bytesTS, _ := json.Marshal(ts)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(bytesTS))
	if err != nil {
		log.Panicf("error while creating request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer " + jwt)

	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("error while performing request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 401 {
		fmt.Println("401 error")
		fmt.Println("please provide a new jwt")
		fmt.Scanln(&jwt)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println("non 200 status")
		log.Fatal(body)
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	// read open bracket
	_, err = dec.Token()
	if err != nil {
		log.Println("failed looking for opening token")
		log.Println(err)
	}

	// while the array contains values
	for dec.More() {
		var m logMsg
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			log.Println("failed decoding")
			log.Println(err)
		}

		if showTimeUUID {
			fmt.Printf("%v: %v\n", m.TimeUUID, m.JsonLog.Log)
		} else {
			fmt.Printf("%v\n", m.JsonLog.Log)
		}
		currentTimeUUID = m.TimeUUID
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		log.Println("failed looking for closing token")
		log.Println(err)
	}
	if currentTimeUUID != "" {
		firstLogs = true
	}
}

func getLogs(url string, optional string, client http.Client) {
	if firstLogs == false {
		initGetLogs(url, optional, client)
		return
	}
	ts := timespan{
		StartTime: currentTimeUUID,
		EndTime:   endTime,
	}

	if byTID {
		ts.TransactionId = optional
	}
	if byService {
		ts.Service = optional
	}

	bytesTS, _ := json.Marshal(ts)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(bytesTS))
	if err != nil {
		log.Panicf("error while creating request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer " + jwt)

	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("error while performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		fmt.Println("401 error")
		fmt.Println("please provide a new jwt")
		fmt.Scanln(&jwt)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println("non 200 status")
		log.Fatal(body)
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	// read open bracket
	_, err = dec.Token()
	if err != nil {
		log.Println("failed looking for opening token in loop")
		log.Println(err)
	}

	// while the array contains values
	for dec.More() {
		var m logMsg
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			log.Println("failed decoding in loop")
			log.Println(err)
		}

		if showTimeUUID {
			fmt.Printf("%v: %v\n", m.TimeUUID, m.JsonLog.Log)
		} else {
			fmt.Printf("%v\n", m.JsonLog.Log)
		}
		currentTimeUUID = m.TimeUUID
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		log.Println("failed looking for closing token in loop")
		log.Println(err)
	}
}



type timespan struct {
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	TransactionId string `json:"transactionId"`
	Service string `json:"service"`
}

type logMsg struct {
	TimeUUID string `json:"timeUUID"`
	JsonLog logJson `json:"jsonLog"`
}

type logJson struct {
	Container_id string `json:"container_id"`
	Container_name string `json:"container_name"`
	Source string `json:"source"`
	Log string `json:"log"`
	Ts_uuid string `json:"ts_uuid"`
	Ts string `json:"ts"`
	Year string `json:"year"`
	Month string `json:"month"`
	Day string `json:"day"`
}

func getArgs(args []string){
	for _, v := range args {
		if v == "--help" {
			println("Use this app to retrieve and tail logs from the logger service")
			println("required flags are: ")
			println("--url <url path>: url for accessing the logger via REST. Include the endpoint you want to use")
			println("")
			println("optional flags are: ")
			println("--follow: set this to tail logs")
			println("--jwt <jwt>: the jwt for access the logger")
			println("--service <service name>: the name of the service to get logs for the service endpoint")
			println("--transactionid <transaction ID>: the transaction ID to get logs for the transaction endpoint")
			println("--starttime <timestamp>: the timestamp for start time given in the form: " + layout)
			println("--endtime <timestamp>: the timestamp for end time given in the form: " + layout)
			println("--showtimeuuid: set this to show timeUUID for each log")
			log.Fatal()
		}
	}

	for k, v := range args {
		if v == "--jwt" {
			jwt = args[k+1]
		}
	}

	for _, v := range args {
		if v == "--follow" {
			follow = true
		}
	}

	for k, v := range args {
		if v == "--url" {
			url = args[k+1]
		}
	}

	for k, v := range args {
		if v == "--service" {
			byService = true
			service = args[k+1]
		}
	}

	for k, v := range args {
		if v == "--transactionid" {
			byTID = true
			transactionId = args[k+1]
		}
	}

	for k, v := range args {
		if v == "--starttime" {
			startTime = args[k+1]
		}
	}

	for k, v := range args {
		if v == "--endtime" {
			endTime = args[k+1]
		}
	}

	for _, v := range args {
		if v == "--showtimeuuid" {
			showTimeUUID = true
		}
	}
}