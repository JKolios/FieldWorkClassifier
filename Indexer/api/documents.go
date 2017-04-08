package api

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v5"
	"golang.org/x/net/context"
	"time"
	"code.google.com/p/go-uuid/uuid"
)

/*deviceDataDoc is a representation of the JSON object
 in the form it's received from the mobile devices*/
type deviceDataDoc struct {
	CompanyId int       `json:"company_id"`
	DriverId  int       `json:"driver_id"`
	Timestamp time.Time `json:"timestamp"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Accuracy  float64   `json:"accuracy"`
	Speed     float64   `json:"speed"`

}

/*deviceDataDoc is a representation of the JSON object
 after it is adapted for Elasticsearch's schema*/
type adaptedDataDoc struct {
	CompanyId int       `json:"company_id"`
	DriverId  int       `json:"driver_id"`
	Timestamp time.Time `json:"timestamp"`
	Location  [2]float64   `json:"location"`
	Accuracy  float64   `json:"accuracy"`
	Speed     float64   `json:"speed"`

}

/*fieldDoc is a representation of a JSON object
representing an elasticsearch document containing
the location of a single (agricultural) field*/
type fieldDoc struct {
	Location Polygon		`json:"field_location"`
}

/*Polygon is a representation of a JSON object
containing the coordinates of each edge of a field
formatted for Elasticsearch's geo_shape type*/
type Polygon struct {
	Type string		       `json:"type"`
	Coordinates [][]Coordinate       `json:"coordinates"`
}

func NewPolygon(coordinates [][]Coordinate) Polygon {
	return Polygon{
		Type:"polygon",
		Coordinates:coordinates,
	}
}

/*Coordinate is a pair of floats representing the latitude and longitude
  of a single geographic point*/

type Coordinate [2]float64

func handleIncomingDeviceData(ginContext *gin.Context) {
	var incoming deviceDataDoc
	var err error
	var resp *elastic.IndexResponse

	client := ginContext.MustGet("ESClient").(*elastic.Client)

	err = ginContext.Bind(&incoming)
	if err != nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	resp, err = indexDeviceDataDoc(client, "device_data", "device_data", incoming)

	if err!=nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	ginContext.JSON(http.StatusOK, resp)
	return

}

func handleIncomingFieldDoc(ginContext *gin.Context) {
	//Assuming that the incoming coordinate data is formatted
	//as a GeoJSON polygon
	var incoming [][]Coordinate
	var err error
	var resp *elastic.IndexResponse

	client := ginContext.MustGet("ESClient").(*elastic.Client)

	err = ginContext.Bind(&incoming)
	if err != nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("%+v", incoming)
	newDoc := fieldDoc{Location: NewPolygon(incoming)}
	log.Printf("%+v", newDoc)

	resp, err = client.Index().
		Index("fields").
		Type("field").
		Id(uuid.New()).
		BodyJson(newDoc).
		Do(context.TODO())

	if err!=nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	ginContext.JSON(http.StatusOK, resp)
	return

}


func indexDeviceDataDoc(client *elastic.Client, indexName string, docType string, doc deviceDataDoc) (*elastic.IndexResponse, error) {

	adaptedDoc := adaptedDataDoc{
		CompanyId:doc.CompanyId,
		DriverId:doc.DriverId,
		Timestamp:doc.Timestamp,
		Location: [2]float64{doc.Latitude, doc.Longitude},
		Accuracy:doc.Accuracy,
		Speed:doc.Speed,
	}

	resp, err := client.Index().
		Index(indexName).
		Type(docType).
		Id(uuid.New()).
		BodyJson(adaptedDoc).
		Do(context.TODO())

	return resp, err
}
