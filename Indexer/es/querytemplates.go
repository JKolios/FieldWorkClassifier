package es

import (
	"text/template"
	"time"
)

const latestDataForDriverQuery = `
{
	"query": {
		"bool": {
			"must": [{
				"range": {
					"timestamp": {
						"gte": "{.Timestamp}||/d",
						"lte": "{.Timestamp}||/d+1d"
					}
				}
			}, {
				"term": {
					"driver_id": {.DriverId}
				}
			}, {
				"term": {
					"company_id": {.CompanyId}
				}

			}]
		}
	},
	"size": 1,
	"sort": [{
		"timestamp": {
			"order": "desc"
		}
	}]
}
`

type LatestDataforDriverParams struct {
	DriverId, CompanyId int
	Timestamp time.Time
}

var LatestDataforDriverTemplate = template.Must(
	template.New("LatestDataForDriver").
		Parse(latestDataForDriverQuery))

