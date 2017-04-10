package es


import (
	"time"
	"gopkg.in/olivere/elastic.v5"
)

type DriverTimetable struct {
	DriverId, CompanyId int64
	Day time.Time
	Activities []Activity
}

type Activity struct {
	From, To time.Time
	ActivityName string
	Duration time.Duration
}

func getDriverTimetable (client *elastic.Client, driverId, companyId int64, timestamp time.Time ) (string, error) {

	////Get the latest document with the same driver and company id
	////For the same day
	//queryParams := LatestDataforDriverParams{
	//	DriverId:  doc.DriverId,
	//	CompanyId: doc.CompanyId,
	//	Timestamp: doc.Timestamp.Format(time.RFC3339),
	//}
	//
	//queryBody := new(bytes.Buffer)
	//LatestDataforDriverTemplate.Execute(queryBody, queryParams)
	//
	//searchResult, err := client.Search().
	//	Query(elastic.NewRawStringQuery(queryBody.String())).
	//	Sort("timestamp", false).
	//	Size(1).
	//	Do(context.TODO())
	//
	//if err != nil {
	//	return "", err
	//}
	//
	//if searchResult.TotalHits() == 0 {
	//	return uuid.New(), nil
	//}
	//
	//var latestDoc AdaptedDataDoc
	//
	//// Iterate through results
	//for _, hit := range searchResult.Hits.Hits {
	//
	//	err := json.Unmarshal(*hit.Source, &latestDoc)
	//	if err != nil {
	//		return "", err
	//	}
	//}
}