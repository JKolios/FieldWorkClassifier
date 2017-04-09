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
 after it is adapted for Elasticsearch's schema*/
type adaptedDataDoc struct {
	CompanyId int       `json:"company_id"`
	DriverId  int       `json:"driver_id"`
	Timestamp time.Time `json:"timestamp"`
	Location  geojson.Point   `json:"location"`
	Accuracy  float64   `json:"accuracy"`
	Speed     float64   `json:"speed"`

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

	//exists, err := client.Exists().
	//	Index("fields").
	//	Type("field_locations").
	//	Id(FIELD_DOC_ID).
	//	Do(context.TODO())
	//
	//if err != nil {
	//	log.Println("Failed: Checking for field doc existence")
	//	ginContext.JSON(http.StatusInternalServerError, exists)
	//	return
	//}
	//
	//if !exists {
	//	defaultDoc := fieldDoc{FieldPolygons:NewMultipolygon([][][]Coordinate{})}
	//
	//	exists, err := client.Index().
	//		Index("fields").
	//		Type("field_locations").
	//		Id(FIELD_DOC_ID).
	//		BodyJson(defaultDoc).
	//		Do(context.TODO())
	//
	//	if err != nil {
	//		log.Println("Failed: Indexing default field doc ")
	//		ginContext.JSON(http.StatusInternalServerError, exists)
	//		return
	//	}
	//
	//}


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


func indexDeviceDataDoc(client *elastic.Client, indexName string, docType string, doc deviceDataDoc) (*elastic.IndexResponse, error) {

	adaptedDoc := adaptedDataDoc{
		CompanyId:doc.CompanyId,
		DriverId:doc.DriverId,
		Timestamp:doc.Timestamp,
		Location: geojson.NewPoint(geojson.Coordinate{doc.Latitude, doc.Longitude}),
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
