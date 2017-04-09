package api

import (
	"github.com/JKolios/FieldWorkClassifier/Common/config"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"gopkg.in/olivere/elastic.v5"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/JKolios/FieldWorkClassifier/Common/geojson"
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

		var incomingDoc deviceDataDoc
		err := json.Unmarshal(msg, &incomingDoc)
		if err != nil {
			errorMessage := fmt.Sprintf("{\"error\": \"%v\"}", err.Error())
			session.Write([]byte(errorMessage))
			return
		}

		adaptedDoc := adaptedDataDoc{
			CompanyId:incomingDoc.CompanyId,
			DriverId:incomingDoc.DriverId,
			Timestamp:incomingDoc.Timestamp,
			Location: geojson.NewPoint(geojson.Coordinate{incomingDoc.Longitude, incomingDoc.Latitude}),
			Accuracy:incomingDoc.Accuracy,
			Speed:incomingDoc.Speed,
		}

		activity, err := activityFromDeviceData(ESClient, adaptedDoc)

		if err != nil {
			errorMessage := fmt.Sprintf("{\"error\": \"%v\"}", err.Error())
			session.Write([]byte(errorMessage))
			return
		}

		resp, err := indexDeviceDataDoc(ESClient, adaptedDoc, activity)


		if err != nil {
			errorMessage := fmt.Sprintf("{\"error\": \"%v\"}", err.Error())
			session.Write([]byte(errorMessage))
			return
		}

		response := fmt.Sprintf("{\"indexed\": \"%v\"}", resp.Created)
		session.Write([]byte(response))

	})

	return router
}
