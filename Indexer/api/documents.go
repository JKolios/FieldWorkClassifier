package api

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v5"
	"golang.org/x/net/context"
	"time"
	"code.google.com/p/go-uuid/uuid"
	"github.com/JKolios/FieldWorkClassifier/Common/geojson"
)


/* Field Locations are stored in a nested array on one
preset elasticsearch document. This increases geosearch performance.*/
const FIELD_DOC_ID = "field_locations"

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
 after it is adapted for querying with Elasticsearch's*/
type adaptedDataDoc struct {
	CompanyId int       `json:"company_id"`
	DriverId  int       `json:"driver_id"`
	Timestamp time.Time `json:"timestamp"`
	Location  geojson.Point   `json:"location"`
	Accuracy  float64   `json:"accuracy"`
	Speed     float64   `json:"speed"`
	Activity string     `json:"activity"`

}


/*FieldDoc is a representation an elasticsearch datatype
 containing a GeoJSON polygon field*/
type fieldDoc struct {
	FieldPolygons geojson.Multipolygon		`json:"field_polygons"`
}


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

	adaptedDoc := adaptedDataDoc{
		CompanyId:incoming.CompanyId,
		DriverId:incoming.DriverId,
		Timestamp:incoming.Timestamp,
		Location: geojson.NewPoint(geojson.Coordinate{incoming.Longitude, incoming.Latitude}),
		Accuracy:incoming.Accuracy,
		Speed:incoming.Speed,
	}

	activity, err := activityFromDeviceData(client, adaptedDoc)

	if err!=nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	resp, err = indexDeviceDataDoc(client, adaptedDoc, activity)

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
	var incoming [][][]geojson.Coordinate
	var err error
	var update *elastic.UpdateResponse

	client := ginContext.MustGet("ESClient").(*elastic.Client)

	err = ginContext.Bind(&incoming)
	if err != nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("%+v", incoming)
	newDoc := fieldDoc{FieldPolygons: geojson.NewMultipolygon(incoming)}
	log.Printf("%+v", newDoc)

	update, err = client.Update().
		Index("fields").
		Type("field_locations").
		Id(FIELD_DOC_ID).
		Doc(newDoc).
		DocAsUpsert(true).
		Do(context.TODO())

	if err!=nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	ginContext.JSON(http.StatusOK, update)
	return

}

func activityFromDeviceData (client *elastic.Client, doc adaptedDataDoc) (string, error) {



	percolationQuery := elastic.NewPercolatorQuery().
		DocumentType("device_data").
		Field("percolation_query").
		Document(doc)


	percolationResult, err := client.Search("device_data").
		Query(percolationQuery).
		Do(context.TODO())

	log.Println(percolationResult)

	if err != nil || percolationResult.TotalHits() == 0  {
		return "", err
	}


	return percolationResult.Hits.Hits[0].Id, err

}


func indexDeviceDataDoc(client *elastic.Client,  doc adaptedDataDoc, activity string) (*elastic.IndexResponse, error) {



	doc.Activity = activity

	indexResp, err := client.Index().
		Index("device_data").
		Type("device_data").
		Id(uuid.New()).
		BodyJson(doc).
		Do(context.TODO())

	return indexResp, err
}
