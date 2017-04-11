package es

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"github.com/JKolios/FieldWorkClassifier/Common/geojson"
	"context"
	"gopkg.in/olivere/elastic.v5"
	"log"
	"time"
)

/* Field Locations are stored as an array of polygons in an ES document.*/
const FIELD_DOC_ID = "field_locations"

/*DeviceDataDoc is a representation of the JSON object
in the form it's received from the mobile devices*/
type DeviceDataDoc struct {
	CompanyId int       `json:"company_id"`
	DriverId  int       `json:"driver_id"`
	Timestamp time.Time `json:"timestamp"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Accuracy  float64   `json:"accuracy"`
	Speed     float64   `json:"speed"`
}

/*AdaptedDataDoc is a representation of the JSON object
after it is adapted for querying with Elasticsearch*/
type AdaptedDataDoc struct {
	CompanyId         int           `json:"company_id"`
	DriverId          int           `json:"driver_id"`
	Timestamp         time.Time     `json:"timestamp"`
	Location          geojson.Point `json:"location"`
	Accuracy          float64       `json:"accuracy"`
	Speed             float64       `json:"speed"`
	Activity          string        `json:"activity"`
	ActivitySessionId string        `json:"activity_session_id"`
}

/*FieldDoc is a representation an elasticsearch datatype
containing a GeoJSON polygon field*/
type FieldDoc struct {
	FieldPolygons geojson.Multipolygon `json:"field_polygons"`
}

func CreateDeviceDataDoc(client *elastic.Client, doc DeviceDataDoc) (elastic.IndexResponse, error) {

	adaptedDoc := AdaptedDataDoc{
		CompanyId: doc.CompanyId,
		DriverId:  doc.DriverId,
		Timestamp: doc.Timestamp,
		Location:  geojson.NewPoint(geojson.Coordinate{doc.Longitude, doc.Latitude}),
		Accuracy:  doc.Accuracy,
		Speed:     doc.Speed,
	}

	activity, err := GetActivityFromDeviceData(client, adaptedDoc)

	if err != nil {
		return elastic.IndexResponse{}, err
	}

	adaptedDoc.Activity = activity

	activitySessionId, err := getActivitySessionId(client, adaptedDoc)

	if err != nil {
		return elastic.IndexResponse{}, err
	}

	adaptedDoc.ActivitySessionId = activitySessionId

	response, err := IndexDeviceDataDoc(client, adaptedDoc)

	if err != nil {
		return elastic.IndexResponse{}, err
	}

	return *response, nil
}

func GetActivityFromDeviceData(client *elastic.Client, doc AdaptedDataDoc) (string, error) {

	percolationQuery := elastic.NewPercolatorQuery().
		DocumentType("device_data").
		Field("query").
		Document(doc)

	percolationResult, err := client.Search("device_data").
		Query(percolationQuery).
		Do(context.TODO())

	if err != nil {
		log.Printf("Error encountered while percolating: %v", err.Error())
		return "", err
	}

	// 0 matches means the activity cannot be classified into the given categories
	// >1 matches should not be possible given the problem description
	if percolationResult.TotalHits() != 1 {
		return "other", nil
	}

	return percolationResult.Hits.Hits[0].Id, err

}

func getActivitySessionId(client *elastic.Client, incomingDoc AdaptedDataDoc) (string, error) {

	//Get the latest document with the same driver and company id
	//For the same day
	queryParams := LatestDataforDriverParams{
		DriverId:  incomingDoc.DriverId,
		CompanyId: incomingDoc.CompanyId,
		Timestamp: incomingDoc.Timestamp.Format(time.RFC3339),
	}

	queryBody := new(bytes.Buffer)
	LatestDataforDriverTemplate.Execute(queryBody, queryParams)

	searchResult, err := client.Search().
		Query(elastic.NewRawStringQuery(queryBody.String())).
		Sort("timestamp", false).
		Size(1).
		Do(context.TODO())

	if err != nil {
		return "", err
	}

	if searchResult.TotalHits() == 0 {
		return uuid.New(), nil
	}

	var latestDoc AdaptedDataDoc

	// Iterate through results
	for _, hit := range searchResult.Hits.Hits {

		err := json.Unmarshal(*hit.Source, &latestDoc)
		if err != nil {
			return "", err
		}
	}

	if latestDoc.Activity == incomingDoc.Activity {
		return latestDoc.ActivitySessionId, nil
	} else {
		return uuid.New(), nil
	}
}

func IndexDeviceDataDoc(client *elastic.Client, doc AdaptedDataDoc) (*elastic.IndexResponse, error) {

	indexResp, err := client.Index().
		Index("device_data").
		Type("device_data").
		Id(uuid.New()).
		BodyJson(doc).
		Do(context.TODO())

	return indexResp, err
}
