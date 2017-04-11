package es

import (
	"text/template"
)

const latestDataForDriverQuery = `
{
		"bool": {
			"must": [{
				"range": {
					"timestamp": {
						"gte": "{{.Timestamp}}||/d",
						"lte": "{{.Timestamp}}"
					}
				}
			}, {
				"term": {
					"driver_id": {{.DriverId}}
				}
			}, {
				"term": {
					"company_id": {{.CompanyId}}
				}

			}]
		}
}
`

type LatestDataforDriverParams struct {
	DriverId, CompanyId int
	Timestamp string
}


var LatestDataforDriverTemplate = template.Must(
	template.New("LatestDataForDriver").
		Parse(latestDataForDriverQuery))

