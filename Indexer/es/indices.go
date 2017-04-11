package es

import (
	"github.com/JKolios/FieldWorkClassifier/Common/geojson"
	"github.com/JKolios/FieldWorkClassifier/Common/utils"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

const deviceDataMapping = `
{
	"mappings": {
		"device_data": {
			"properties": {
				"company_id": {
					"type": "long"
				},
				"driver_id": {
					"type": "long"
				},
				"timestamp": {
					"type": "date"
				},
				"accuracy": {
					"type": "double"
				},
				"speed": {
					"type": "double"
				},
				"location": {
					"type": "geo_shape",
					"tree": "quadtree",
					"precision": "5m",
					"strategy": "recursive"
				},
				"activity": {
					"type": "keyword"
				},
				"activity_session_id": {
					"type": "keyword"
				}
			}
		},
		"queries": {
			"properties": {
				"query": {
					"type": "percolator"
				}
			}
		}
	}
}`

const FieldMapping = `{
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
	"fields":      FieldMapping,
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

func InitFieldLocationDocument(elasticClient *elastic.Client) {

	defaultDoc := FieldDoc{
		FieldPolygons: geojson.NewMultipolygon([][][]geojson.Coordinate{}),
	}

	exists, err := elasticClient.Exists().
		Index("fields").
		Type("field_locations").
		Id(FIELD_DOC_ID).
		Do(context.TODO())

	if err != nil {
		log.Fatalf("Failed: checking for existence of the field location doc: %v", err.Error())
	}

	if !exists {

		//Inserts a default document if none already exists.
		_, err := elasticClient.Index().
			Index("fields").
			Type("field_locations").
			Id(FIELD_DOC_ID).
			BodyJson(defaultDoc).
			Do(context.TODO())

		if err != nil {
			log.Fatalf("Failed: initializing default field location doc: %v", err.Error())
		}

		log.Println("Initialized the default field location doc")
	}

}
