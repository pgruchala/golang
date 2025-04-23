package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func fetchStopsData(url string) ([]Stop, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var stopsResp StopsRes
	err = json.Unmarshal(body, &stopsResp)
	if err != nil {
		return nil, err
	}
	return stopsResp.Stops, nil
}

func fetchEstimatedDepartures(stopID string) ([]EstimatedDeparture, error) {
	url := fmt.Sprintf("https://ckan2.multimediagdansk.pl/departures?stopId=%s", stopID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var departuresResp DeparturesRes
	err = json.Unmarshal(body, &departuresResp)
	if err != nil {
		return nil, err
	}

	return departuresResp.EstimatedDepartures, nil
}
