package main

import (
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
	"github.com/JKolios/FieldWorkClassifier/Common/geojson"
)

const FIELD_ENDPOINT  = "http://localhost:8090/v0/field"

func main() {

	fileContents, err := ioutil.ReadFile("polygons.json")
	if err != nil {
		log.Fatalf("Cannot read polygons.json: %v", err.Error())
	}
	featureCollection := geojson.FeatureCollection{}
	err = json.Unmarshal(fileContents, &featureCollection)
	if err != nil {
		log.Fatalf("Cannot parse polygons.json as valid JSON: %v", err.Error())
	}

	allPolygonCoordinates := [][][]geojson.Coordinate{}

	for _, feature := range featureCollection.Features {

		allPolygonCoordinates = append(allPolygonCoordinates, feature.Geometry.Coordinates)
	}

		buffer := new(bytes.Buffer)
		jsonEncoder := json.NewEncoder(buffer)
		err = jsonEncoder.Encode(allPolygonCoordinates)
		if err != nil {
			log.Printf("Error encountered when JSON encoding payload: %v", err)
			return
		}

		response, err := http.Post(FIELD_ENDPOINT, "application/json", buffer)

		if err != nil {
			log.Printf("Error encountered when posting payload: %v", err)
			return
		}

		if response.StatusCode != http.StatusOK {
			log.Printf("HTTP Error encountered: %v", response.StatusCode)
			return
		}


		response.Body.Close()
		log.Println("Polygons Imported successfully")
}
