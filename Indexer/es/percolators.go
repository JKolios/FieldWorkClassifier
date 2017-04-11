package es

import (
	"github.com/JKolios/FieldWorkClassifier/Common/utils"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

const drivingPercolator = `
{
"query": {
        "bool": {
            "filter": [
                {
                    "geo_shape": {
                        "location": {
                            "relation": "disjoint",
                            "indexed_shape": {
                                "id": "field_locations",
                                "type": "field_locations",
                                "index": "fields",
                                "path": "field_polygons"
                            }
                        }
                    }
                },
                {
                    "range": {
                        "speed": {
                            "gt": 5
                        }
                    }
                }
            ]
        }
    }
}
`

const cultivatingPercolator = `
{
"query": {
        "bool": {
            "filter": [
                {
                    "geo_shape": {
                        "location": {
                            "relation": "within",
                            "indexed_shape": {
                                "id": "field_locations",
                                "type": "field_locations",
                                "index": "fields",
                                "path": "field_polygons"
                            }
                        }
                    }
                },
                {
                    "range": {
                        "speed": {
                            "gt": 1
                        }
                    }
                }
            ]
        }
    }
}
`

const repairingPercolator = `
{
"query": {
        "bool": {
            "filter": [
                {
                    "geo_shape": {
                        "location": {
                            "relation": "within",
                            "indexed_shape": {
                                "id": "field_locations",
                                "type": "field_locations",
                                "index": "fields",
                                "path": "field_polygons"
                            }
                        }
                    }
                },
                {
                    "range": {
                        "speed": {
                            "lt": 1
                        }
                    }
                }
            ]
        }
    }
}

`

var percolators = map[string]string{
	"Driving":     drivingPercolator,
	"Cultivating": cultivatingPercolator,
	"Repairing":   repairingPercolator,
}

func InitPercolators(elasticClient *elastic.Client) {
	for percolatorId, query := range percolators {

		log.Printf("Initializing Percolator: %s", percolatorId)
		_, err := elasticClient.Index().
			Index("device_data").
			BodyJson(query).
			Id(percolatorId).
			Type("queries").
			Refresh("true").
			Do(context.TODO())

		utils.CheckFatalError(err)

		log.Printf("Percolator Initialized: %s", percolatorId)
	}
}
