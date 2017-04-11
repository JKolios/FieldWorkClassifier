package es

import (
	"context"
	"encoding/json"
	"gopkg.in/olivere/elastic.v5"
	"log"
	"time"
)

type DriverTimetableQueryParams struct {
	DriverId  int64     `json:"driverId"`
	CompanyId int64     `json:"companyId"`
	Timestamp time.Time `json:"day"`
}

type DriverDailyTimetable struct {
	DriverId   int64      `json:"driverId"`
	CompanyId  int64      `json:"companyId"`
	Day        string     `json:"day"`
	Activities []Activity `json:"activities"`
}

type Activity struct {
	From         time.Time `json:"from"`
	To           time.Time `json:"to"`
	ActivityName string    `json:"activity"`
	Duration     string    `json:"duration"`
}

type ActivityDoc struct {
	Activity string `json:"activity"`
}

//GetDriverDailyTimetable returns the daily timetable (listing of tasks with their durations)
//for one particular driver of one particular company on a given day.
func GetDriverDailyTimetable(client *elastic.Client, params *DriverTimetableQueryParams) (DriverDailyTimetable, error) {

	searchResult, err := client.Search("device_data").
		Query(constructDriverTimetableQuery(params)).
		Aggregation("sessions", constructDriverTimetableAggregation()).
		Size(0).
		Do(context.Background())

	if err != nil {
		return DriverDailyTimetable{}, err
	}

	day := params.Timestamp.Format("02/01/2006")

	timetable := DriverDailyTimetable{
		DriverId:  params.DriverId,
		CompanyId: params.CompanyId,
		Day:       day,
	}

	aggResult, found := searchResult.Aggregations.Terms("sessions")

	if !found {
		return timetable, nil
	}

	timetable.Activities = populateDailyActivityList(aggResult.Buckets)

	return timetable, nil

}

func constructDriverTimetableQuery(params *DriverTimetableQueryParams) elastic.Query {

	// The top-level filtering query
	baseQuery := elastic.NewBoolQuery()
	filters := []elastic.Query{}

	// Adding filters

	activityTypeFilter := elastic.NewTermsQuery("activity",
		"Driving", "Cultivating", "Repairing")
	filters = append(filters, activityTypeFilter)

	companyIdFilter := elastic.NewTermQuery("company_id", params.CompanyId)
	filters = append(filters, companyIdFilter)

	driverIdFilter := elastic.NewTermQuery("driver_id", params.DriverId)
	filters = append(filters, driverIdFilter)

	//Adapting the timestamp format to ElasticSearch's date_optional_time
	minimumTs := params.Timestamp.Format("2006-01-02T15:04:05-07:00") + "||/d"
	maximumTs := params.Timestamp.Format("2006-01-02T15:04:05-07:00") + "||+1d/d"

	timestampFilter := elastic.NewRangeQuery("timestamp").
		Gte(minimumTs).
		Lte(maximumTs)

	filters = append(filters, timestampFilter)

	baseQuery.Filter(filters...)

	return baseQuery
}

func constructDriverTimetableAggregation() *elastic.TermsAggregation {

	//Creating lower level aggregations
	firstPointAggregation := elastic.NewMinAggregation().
		Field("timestamp")

	lastPointAggregation := elastic.NewMaxAggregation().
		Field("timestamp")

	sessionActivityAggregation := elastic.NewTopHitsAggregation().
		Size(1).
		DocvalueField("activity")

		//The top-level aggregation
	baseAggregation := elastic.NewTermsAggregation().
		Field("activity_session_id").
		Order("session_first_point", true).
		MinDocCount(2).
		SubAggregation("session_first_point", firstPointAggregation).
		SubAggregation("session_last_point", lastPointAggregation).
		SubAggregation("session_activity", sessionActivityAggregation)

	baseAggregation = baseAggregation.SubAggregation("session_first_point", firstPointAggregation)

	return baseAggregation
}

func populateDailyActivityList(sessions []*elastic.AggregationBucketKeyItem) []Activity {

	activities := []Activity{}
	for _, session := range sessions {

		sessionStart, found := session.Min("session_first_point")
		if !found {
			log.Println("Malformed result pulled from timetable query.")
			continue
		}

		sessionEnd, found := session.Min("session_last_point")
		if !found {
			log.Println("Malformed result pulled from timetable query.")
			continue
		}

		//Elasticsearch returns Unix-like timestamps in milliseconds since the Epoch
		sessionStartTime := time.Unix(int64(*sessionStart.Value)/1000, 0)
		sessionEndTime := time.Unix(int64(*sessionEnd.Value)/1000, 0)

		sessionDuration := sessionEndTime.Sub(sessionStartTime).String()

		sessionActivityDocAggregation, found := session.TopHits("session_activity")
		if !found || sessionActivityDocAggregation.Hits.TotalHits == 0 {
			log.Println("Malformed result pulled from timetable query.")
			continue
		}

		sessionActivityDoc := ActivityDoc{}

		//This bucket will only ever include one result, by design

		err := json.Unmarshal(*sessionActivityDocAggregation.Hits.Hits[0].Source, &sessionActivityDoc)

		if err != nil {
			log.Printf("Malformed result pulled from timetable query: %v", err.Error())
			continue
		}

		newActivity := Activity{
			From:         sessionStartTime,
			To:           sessionEndTime,
			ActivityName: sessionActivityDoc.Activity,
			Duration:     sessionDuration}

		activities = append(activities, newActivity)

	}

	return activities

}
