package api

import (
	"github.com/JKolios/FieldWorkClassifier/Common/config"
	"github.com/JKolios/FieldWorkClassifier/Common/geojson"
	"github.com/JKolios/FieldWorkClassifier/Indexer/es"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"gopkg.in/olivere/elastic.v5"
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"golang.org/x/net/context"
)


func contextInjector(ESClient *elastic.Client, conf *config.Settings) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ESClient", ESClient)
		c.Set("Conf", conf)
		c.Next()
	}
}

func SetupAPI(ESClient *elastic.Client, conf *config.Settings) *gin.Engine {
	if !conf.GinDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(contextInjector(ESClient, conf))

	wsRouter := melody.New()
	// Allow all origins for websocket connections
	// The requirements imply that this is valid for the use case
	wsRouter.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	//API v0 HTTP endpoints
	v0 := router.Group("/v0")
	{
		v0.GET("/status", status)
		v0.POST("/deviceData", handleIncomingDeviceData)
		v0.POST("/field", handleIncomingFieldDoc)
		v0.GET("/wsDeviceData", func(ginContext *gin.Context) {

			//Upgrade the HTTP connection to WebSocket session and add all needed objects
			requestKeys := map[string] interface{} {"Client": ESClient}
			wsRouter.HandleRequestWithKeys(ginContext.Writer, ginContext.Request, requestKeys)
		})
	}

	//Websocket Connection Handlers

	wsRouter.HandleMessage(func(session *melody.Session, msg []byte) {

		var incomingDoc es.DeviceDataDoc
		err := json.Unmarshal(msg, &incomingDoc)
		if err != nil {
			errorMessage := fmt.Sprintf("{\"error\": \"%v\"}", err.Error())
			session.Write([]byte(errorMessage))
			return
		}

		err = es.CreateDeviceDataDoc(ESClient, incomingDoc)


		if err != nil {
			errorMessage := fmt.Sprintf("{\"error\": \"%v\"}", err.Error())
			session.Write([]byte(errorMessage))
			return
		}

		session.Write([]byte("{\"indexed\": true}"))

	})

	return router
}

func handleIncomingDeviceData(ginContext *gin.Context) {
	var incoming es.DeviceDataDoc
	var err error
	var resp *elastic.IndexResponse

	client := ginContext.MustGet("ESClient").(*elastic.Client)

	err = ginContext.Bind(&incoming)
	if err != nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	err = es.CreateDeviceDataDoc(client, incoming)

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
	newDoc := es.FieldDoc{FieldPolygons: geojson.NewMultipolygon(incoming)}
	log.Printf("%+v", newDoc)

	update, err = client.Update().
		Index("fields").
		Type("field_locations").
		Id(es.FIELD_DOC_ID).
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
