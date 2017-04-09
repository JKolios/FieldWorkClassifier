package es

import (
	"gopkg.in/olivere/elastic.v5"
	"log"
	"golang.org/x/net/context"
	"github.com/JKolios/FieldWorkClassifier/Common/utils"
)
var deviceDataMapping = `
{
	"mappings":{
		"device_data":{
			"properties":{
				"company_id":{
					"type":"long"
				},
				"driver_id":{
					"type":"long"
				},
				"timestamp":{
					"type":"date"
				},
				"accuracy":{
					"type":"double"
				},
				"speed":{
					"type":"double"
				},
				"location":{
					"type": "geo_shape",
					"tree": "quadtree",
					"precision": "5m",
					"strategy": "recursive"
				}
			}
		}
	}
}`

var FieldMapping = `{
    "mappings": {
        "field_locations": {
            "properties": {
                "field_polygons": {
                    "type": "geo_shape",
                    "tree": "quadtree",
                    "precision": "5m",
                    "strategy": "recursive"
                }
            }
        }
    }
}`


var indices = map[string]string{
	"device_data": deviceDataMapping,
	"fields": FieldMapping,
}


func InitIndices(elasticClient *elastic.Client) {
	for index, mapping := range indices {

		log.Printf("Initializing Index: %s", index)

		indexExists, err := elasticClient.IndexExists(index).Do(context.TODO())
		utils.CheckFatalError(err)
		if !indexExists {
			resp, err := elasticClient.CreateIndex(index).BodyString(mapping).Do(context.TODO())
			utils.CheckFatalError(err)
			if !resp.Acknowledged {
				log.Fatalf("Cannot create index: %s on ES", index)
			}
			log.Printf("Created index: %s on ES", index)

		} else {
			log.Printf("Index: %s already exists on ES", index)
		}

		_, err = elasticClient.OpenIndex(index).Do(context.TODO())
		utils.CheckFatalError(err)
	}
}
