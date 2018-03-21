package hello

import (
	"encoding/json"
	"fmt"
	"html/template"
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
		Action string `json:"action"`
	} `json:"result"`
	OriginalRequest DialogFlowOriginalRequest `json:"originalRequest"`
}

type DialogFlowOriginalRequest struct {
	Data DialogFlowOriginalRequestData `json:"data"`
}

type DialogFlowOriginalRequestData struct {
	Device DialogFlowOriginalRequestDevice `json:"device"`
}

type DialogFlowOriginalRequestDevice struct {
	Location DialogFlowOriginalRequestLocation `json:"location"`
}

type DialogFlowOriginalRequestLocation struct {
	City string `json:"city"`
}

// DialogFlowResponse struct
type DialogFlowResponse struct {
	Speech string `json:"speech"`
}

// DialogFlowLocationResponse struct
type DialogFlowLocationResponse struct {
	Speech string                 `json:"speech"`
	Data   DialogFlowResponseData `json:"data"`
}

type DialogFlowResponseData struct {
	Google DialogFlowResponseGoogle `json:"google"`
}

type DialogFlowResponseGoogle struct {
	ExpectUserResponse bool                           `json:"expectUserResponse"`
	IsSsml             bool                           `json:"isSsml"`
	SystemIntent       DialogFlowResponseSystemIntent `json:"systemIntent"`
}

type DialogFlowResponseSystemIntent struct {
	Intent string                             `json:"intent"`
	Data   DialogFlowResponseSystemIntentData `json:"data"`
}

type DialogFlowResponseSystemIntentData struct {
	Type        string   `json:"@type"`
	OptContext  string   `json:"optContext"`
	Permissions []string `json:"permissions"`
}

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

var (
	apiErrMsg     = "Sorry, I was unable to get data from AQICN. Please try later."
	unknownErrMsg = "Sorry, I can't help you with this right now. Please try later. ok?"
	token         = "c3bfc1119947754409a5b92bfc9eb1e404ae953b"
)

func handle(w http.ResponseWriter, r *http.Request) {
	dfReq := DialogFlowRequest{}
	dfErr := json.NewDecoder(r.Body).Decode(&dfReq)

	if dfErr == nil && dfReq.Result.Action == "location_permission" {
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
		return
	}

	if dfErr == nil && dfReq.Result.Action == "get" {
		handleYesIntent(w, r, dfReq)
		return
	}

	json.NewEncoder(w).Encode(DialogFlowResponse{
		Speech: unknownErrMsg,
	})
}

func handleYesIntent(w http.ResponseWriter, r *http.Request, dfReq DialogFlowRequest) {
	var feedID = "here"
	geoCity := dfReq.OriginalRequest.Data.Device.Location.City
	if len(geoCity) > 0 {
		foundFeedID, feedErr := getFeedID(r, geoCity)
		if feedErr == nil && len(foundFeedID) > 0 {
			feedID = foundFeedID
		}
	}

	if feedID == "here" {
		geoCity = "your city"
	}

	aqi, aqiErr := getFeedAQI(r, feedID)
	if aqiErr != nil {
		json.NewEncoder(w).Encode(DialogFlowResponse{
			Speech: apiErrMsg,
		})
		return
	}

	msg := fmt.Sprintf("The air quality index in %s is %d right now. %s", geoCity, aqi, getAirQualityDescription(aqi))
	json.NewEncoder(w).Encode(DialogFlowResponse{
		Speech: msg,
	})
}

func getFeedID(r *http.Request, city string) (string, error) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	resp, err := client.Get(fmt.Sprintf("http://api.waqi.info/search/?token=%s&keyword=%s", token, template.URLQueryEscaper(city)))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	search := AQICNSearchResponse{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&search)
	if decodeErr != nil {
		return "", decodeErr
	}

	if search.Status != "ok" || len(search.Data) == 0 {
		return "", fmt.Errorf("unable to find station in %s", city)
	}

	return fmt.Sprintf("@%d", search.Data[0].UID), nil
}

func getFeedAQI(r *http.Request, feedID string) (int, error) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	resp, err := client.Get(fmt.Sprintf("http://api.waqi.info/feed/%s/?token=%s", feedID, token))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	aqi := AQICNFeedResponse{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&aqi)
	if decodeErr != nil {
		return 0, decodeErr
	}

	if aqi.Status != "ok" {
		return 0, fmt.Errorf("unable to find feed %s", feedID)
	}

	return aqi.Data.AQI, nil
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
