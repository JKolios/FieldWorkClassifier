package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"gopkg.in/olivere/elastic.v5"
	"github.com/JKolios/FieldWorkClassifier/QueryRunner/es"
)

func handleTimeTableQuery(ginContext *gin.Context) {
	var incoming *es.DriverTimetableQueryParams
	var err error
	var response es.DriverDailyTimetable

	client := ginContext.MustGet("ESClient").(*elastic.Client)

	err = ginContext.Bind(&incoming)
	if err != nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	response, err = es.GetDriverDailyTimetable(client, incoming)

	if err!=nil {
		log.Println(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	ginContext.JSON(http.StatusOK, response)
	return

}