package es

import (
	"text/template"
)

const DriverTimetableQuery = `
{
	"bool": {
		"filter": [{
			"term": {
				"company_id": "{{.CompanyId}}"
			}
		}, {
			"term": {
				"driver_id": "{{.DriverId}}"
			}
		}, {
			"range": {
				"timestamp": {
					"gte": "{{.Timestamp}}||/d",
					"lte": "{{.Timestamp}}||+1d/d"
				}
			}
		}]
	}
}`

type DriverTimetableQueryParams struct {
	DriverId, CompanyId int
	Timestamp string
}


var DriverTimetableQueryTemplate = template.Must(
	template.New("DriverTimetableQuery").
		Parse(DriverTimetableQuery))

const DriverTimetableAggregation = `
{
	"sessions": {
		"terms": {
			"field": "activity_session_id",
			"order": {
				"session_first_point": "asc"
			},
			"min_doc_count": 2
		},
		"aggs": {
			"session_first_point": {
				"min": {
					"field": "timestamp"
				}
			},
			"session_last_point": {
				"max": {
					"field": "timestamp"
				}
			},
			"session_activity": {
				"top_hits": {
					"size": 1,
					"_source": "activity"
				}
			}
		}
	}
}
`

