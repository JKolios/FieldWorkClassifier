package api

import (
	"github.com/JKolios/FieldWorkClassifier/Common/config"
	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v5"
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

	//API v0 HTTP endpoints
	v0 := router.Group("/v0")
	{
		v0.GET("/status", status)
	}

	return router
}
