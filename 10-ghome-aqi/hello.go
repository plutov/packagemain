package hello

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/", handle)
}

// DialogFlowRequest struct
type DialogFlowRequest struct {
	Result struct {
		Parameters map[string]string `json:"parameters"`
	} `json:"result"`
}

// DialogFlowResponse struct
type DialogFlowResponse struct {
	Speech      string `json:"speech"`
	DisplayText string `json:"displayText"`
}

// AQICNResponse struct
type AQICNResponse struct {
	Status string `json:"status"`
	Data   struct {
		AQI int `json:"aqi"`
	} `json:"data"`
}

var errMsg = "Sorry, I was unable to get data from AQICN. Please try later."

func handle(w http.ResponseWriter, r *http.Request) {
	dfReq := DialogFlowRequest{}
	dfErr := json.NewDecoder(r.Body).Decode(&dfReq)
	_, withCity := dfReq.Result.Parameters["geo-city"]
	if dfErr != nil || !withCity {
		json.NewEncoder(w).Encode(DialogFlowResponse{
			Speech:      errMsg,
			DisplayText: errMsg,
		})
		return
	}

	city := dfReq.Result.Parameters["geo-city"]

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	resp, err := client.Get(fmt.Sprintf("http://api.waqi.info/feed/%s/?token=c3bfc1119947754409a5b92bfc9eb1e404ae953b", city))
	if err != nil {
		json.NewEncoder(w).Encode(DialogFlowResponse{
			Speech:      errMsg,
			DisplayText: errMsg,
		})
		return
	}

	defer resp.Body.Close()

	aqi := AQICNResponse{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&aqi)
	if decodeErr != nil || aqi.Status != "ok" {
		json.NewEncoder(w).Encode(DialogFlowResponse{
			Speech:      errMsg,
			DisplayText: errMsg,
		})
		return
	}

	msg := fmt.Sprintf("The air quality index in %s is %d right now. Air quality conditions are: %s", city, aqi.Data.AQI, getAirQualityLevel(aqi.Data.AQI))
	json.NewEncoder(w).Encode(DialogFlowResponse{
		Speech:      msg,
		DisplayText: msg,
	})
}

// TODO: https://airnow.gov/index.cfm?action=aqibasics.aqi
func getAirQualityLevel(aqi int) string {
	return "Moderate"
}
