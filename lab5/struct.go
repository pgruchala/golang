package main

import "time"

type StopsRes struct {
	Stops []Stop `json:"stops"`
}


type Stop struct {
	StopID        int    `json:"stopId"`
	StopCode      string `json:"stopCode"`
	StopName      string `json:"stopName"`
	StopShortName string `json:"stopShortName"`
	StopDesc      string `json:"stopDesc"`
	SubName       string `json:"subName"`
	StopLat       float64 `json:"stopLat"`
	StopLon       float64 `json:"stopLon"`
	ZoneID        int    `json:"zoneId"`
	StopURL       string `json:"stopUrl"`
	LocationType  int    `json:"locationType"`
	ParentStation string `json:"parentStation"`
	StopTimezone  string `json:"stopTimezone"`
	WheelchairBoarding int `json:"wheelchairBoarding"`
	Virtual       int   `json:"virtual"`
	NonPassenger  int   `json:"nonPassenger"`
}

type DeparturesRes struct {
	EstimatedDepartures []EstimatedDeparture `json:"departures"`
}

type EstimatedDeparture struct {
	RouteID         int `json:"routeId"`
	Headsign        string `json:"headsign"`
	EstimatedTime   string `json:"estimatedTime"`
	DelayInSeconds  int    `json:"delayInSeconds"`
	ActualTime      string `json:"actualTime"`
}

func formatTime(timeStr string) string {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return timeStr 
	}
	return t.Format("15:04:05")
}