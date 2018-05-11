package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", handle)
}

var (
	defaultIndex = 49
	userMsg      = "The air quality index in your city is %d right now. %s"
)

func handle(w http.ResponseWriter, r *http.Request) {
	dfReq := DialogFlowRequest{}
	dfErr := json.NewDecoder(r.Body).Decode(&dfReq)

	if dfErr == nil && dfReq.Result.Action == "location" {
		handleLocationPermissionAction(w, r, dfReq)
		return
	}

	if dfErr == nil && dfReq.Result.Action == "get" {
		handleGetAction(w, r, dfReq)
		return
	}

	returnDefaultResults(w)
}

func handleLocationPermissionAction(w http.ResponseWriter, r *http.Request, dfReq DialogFlowRequest) {
	json.NewEncoder(w).Encode(DialogFlowLocationResponse{
		Speech: "PLACEHOLDER_FOR_PERMISSION",
		Data: DialogFlowResponseData{
			Google: DialogFlowResponseGoogle{
				ExpectUserResponse: true,
				IsSsml:             false,
				SystemIntent: DialogFlowResponseSystemIntent{
					Intent: "actions.intent.PERMISSION",
					Data: DialogFlowResponseSystemIntentData{
						Type:        "type.googleapis.com/google.actions.v2.PermissionValueSpec",
						OptContext:  "To get city for air quality check",
						Permissions: []string{"DEVICE_PRECISE_LOCATION"},
					},
				},
			},
		},
	})
}

func returnDefaultResults(w http.ResponseWriter) {
	levelDescription := getAirQualityDescription(defaultIndex)
	json.NewEncoder(w).Encode(DialogFlowResponse{
		Speech: fmt.Sprintf(userMsg, defaultIndex, levelDescription),
	})
}

func handleGetAction(w http.ResponseWriter, r *http.Request, dfReq DialogFlowRequest) {
	lat := dfReq.OriginalRequest.Data.Device.Location.Coordinates.Lat
	long := dfReq.OriginalRequest.Data.Device.Location.Coordinates.Long
	if lat == 0 || long == 0 {
		returnDefaultResults(w)
		return
	}

	index, levelDescription, aqiErr := getAirQualityByCoordinates(r, lat, long)
	if aqiErr != nil {
		returnDefaultResults(w)
		return
	}

	json.NewEncoder(w).Encode(DialogFlowResponse{
		Speech: fmt.Sprintf(userMsg, index, levelDescription),
	})
}
