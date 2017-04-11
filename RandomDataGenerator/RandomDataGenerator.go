package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"io/ioutil"
	"log"
)

const MIN_TIMESTAMP = 1493510400
const MAX_TIMESTAMP = 1496188800

const MIN_LAT = 52.10
const MAX_LAT = 52.20

const MIN_LONG = 11.28
const MAX_LONG = 11.48

const INDEXER_ENDPOINT = "http://localhost:8090/v0/deviceData"

type DevicePayload struct {
	CompanyId int       `json:"company_id"`
	DriverId  int       `json:"driver_id"`
	Timestamp time.Time `json:"timestamp"`
	Latitude  float32   `json:"latitude"`
	Longitude float32   `json:"longitude"`
	Accuracy  float32   `json:"accuracy"`
	Speed     float32   `json:"speed"`
}

func generateRandomPayload() DevicePayload {
	companyId := rand.Intn(10)
	driverId := rand.Intn(10)
	timestamp := randomTimestampBetween(MIN_TIMESTAMP, MAX_TIMESTAMP)
	latitude := randomFloatBetween(MIN_LAT, MAX_LAT)
	longitude := randomFloatBetween(MIN_LONG, MAX_LONG)
	accuracy := randomFloatBetween(MIN_LONG, MAX_LONG)
	speed := randomFloatBetween(0.1, 15.0)

	return DevicePayload{CompanyId: companyId, DriverId: driverId, Timestamp: timestamp,
		Latitude: latitude, Longitude: longitude, Accuracy: accuracy, Speed: speed}

}

func randomTimestampBetween(min, max int64) time.Time {
	randomTime := rand.Int63n(max-min) + min
	randomNow := time.Unix(randomTime, 0)
	return randomNow
}

func randomFloatBetween(min, max float32) float32 {
	return rand.Float32()*(max-min) + min
}

func main() {

	rand.Seed(time.Now().Unix())
	log.Println("Generating random data, use Ctrl-C to stop...")

	for {
		randomPayload := generateRandomPayload()
		buffer := new(bytes.Buffer)
		jsonEncoder := json.NewEncoder(buffer)
		err := jsonEncoder.Encode(randomPayload)
		log.Println(buffer)
		if err != nil {
			fmt.Printf("Error encountered when JSON encoding payload: %v", err)
			return
		}

		client := &http.Client{}
		request, err := http.NewRequest("POST", INDEXER_ENDPOINT, buffer)
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)

		if err != nil {
			fmt.Printf("Error encountered when posting payload: %v", err)
			return
		}

		if response.StatusCode != http.StatusOK {
			fmt.Printf("HTTP Error encountered: %v", response.StatusCode)
			responseBody, _ := ioutil.ReadAll(response.Body)
			response.Body.Close()
			fmt.Printf("Response: %v", string(responseBody))
			return
		}


		response.Body.Close()

	}

}
