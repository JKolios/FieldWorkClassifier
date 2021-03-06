package api

import (
	"github.com/JKolios/FieldWorkClassifier/Common/geojson"
	"github.com/JKolios/FieldWorkClassifier/Indexer/es"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
	"log"
	"net/http"
)

func handleIncomingDeviceData(ginContext *gin.Context) {

	var incoming es.DeviceDataDoc
	var resp elastic.IndexResponse
	var err error

	client := ginContext.MustGet("ESClient").(*elastic.Client)

	err = ginContext.Bind(&incoming)
	if err != nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	resp, err = es.CreateDeviceDataDoc(client, incoming)

	if err != nil {
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

	newDoc := es.FieldDoc{FieldPolygons: geojson.NewMultipolygon(incoming)}

	update, err = client.Update().
		Index("fields").
		Type("field_locations").
		Id(es.FIELD_DOC_ID).
		Doc(newDoc).
		DocAsUpsert(true).
		Do(context.TODO())

	if err != nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	ginContext.JSON(http.StatusOK, update)
	return

}
