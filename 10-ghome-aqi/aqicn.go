package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var token = "c3bfc1119947754409a5b92bfc9eb1e404ae953b"

// AQICNSearchResponse struct
type AQICNSearchResponse struct {
	Status string `json:"status"`
	Data   []struct {
		UID int `json:"uid"`
	} `json:"data"`
}

// AQICNFeedResponse struct
type AQICNFeedResponse struct {
	Status string `json:"status"`
	Data   struct {
		AQI int `json:"aqi"`
	} `json:"data"`
}

func getAirQualityByCoordinates(r *http.Request, lat float32, long float32) (int, string, error) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	resp, err := client.Get(fmt.Sprintf("http://api.waqi.info/feed/geo:%.2f;%.2f/?token=%s", lat, long, token))
	if err != nil {
		return 0, "", err
	}

	defer resp.Body.Close()

	aqi := AQICNFeedResponse{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&aqi)
	if decodeErr != nil {
		return 0, "", decodeErr
	}

	if aqi.Status != "ok" {
		return 0, "", fmt.Errorf("unable to get info for %.2f %.2f", lat, long)
	}

	return aqi.Data.AQI, getAirQualityDescription(aqi.Data.AQI), nil
}

// https://airnow.gov/index.cfm?action=aqibasics.aqi
func getAirQualityDescription(aqi int) string {
	if aqi <= 50 {
		return "Air quality is considered satisfactory, and air pollution poses little or no risk."
	} else if aqi <= 100 {
		return "Air quality is acceptable; however, for some pollutants there may be a moderate health concern for a very small number of people who are unusually sensitive to air pollution."
	} else if aqi <= 150 {
		return "Members of sensitive groups may experience health effects. The general public is not likely to be affected."
	} else if aqi <= 200 {
		return "Everyone may begin to experience health effects; members of sensitive groups may experience more serious health effects."
	} else if aqi <= 250 {
		return "Health alert: everyone may experience more serious health effects."
	}

	return "Health warnings of emergency conditions. The entire population is more likely to be affected."
}
