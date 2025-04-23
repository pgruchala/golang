package main

import("time")

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
	Virtual       bool   `json:"virtual"`
	NonPassenger  bool   `json:"nonPassenger"`
}

type Prediction struct {
	TripID      string    `json:"tripId"`
	StopID      int       `json:"stopId"`
	StopSequence int       `json:"stopSequence"`
	ArrivalTime time.Time `json:"arrivalTime"`
	DepartureTime time.Time `json:"departureTime"`
	Timepoint     bool      `json:"timepoint"`
	StopHeadsign string    `json:"stopHeadsign"`
	Realtime      bool      `json:"realtime"`
   
}
