package api

import (
	"log"
	"net/http"

	"github.com/JKolios/FieldWorkClassifier/Common/config"
	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v5"
	"golang.org/x/net/context"
	"time"
	"code.google.com/p/go-uuid/uuid"
)

type deviceDataDoc struct {
	CompanyId int       `json:"company_id"`
	DriverId  int       `json:"driver_id"`
	Timestamp time.Time `json:"timestamp"`
	Latitude  float32   `json:"latitude"`
	Longitude float32   `json:"longitude"`
	Accuracy  float32   `json:"accuracy"`
	Speed     float32   `json:"speed"`

}

func handleIncomingDeviceData(ginContext *gin.Context) {
	var incoming deviceDataDoc
	var err error
	var resp *elastic.IndexResponse

	client := ginContext.MustGet("ESClient").(*elastic.Client)
	defaultIndex := ginContext.MustGet("Conf").(*config.Settings).DefaultIndex

	err = ginContext.Bind(&incoming)
	if err != nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	resp, err = indexDeviceDataDoc(client, defaultIndex, "device_data", incoming)

	ginContext.JSON(http.StatusOK, resp)
	return

}


func indexDeviceDataDoc(client *elastic.Client, indexName string, docType string, doc deviceDataDoc) (*elastic.IndexResponse, error) {

	resp, err := client.Index().
		Index(indexName).
		Type(docType).
		Id(uuid.New()).
		BodyJson(doc).
		Do(context.TODO())

	return resp, err
}
