package app

// DialogFlowRequest struct
type DialogFlowRequest struct {
	Result struct {
		Action string `json:"action"`
	} `json:"result"`
	OriginalRequest DialogFlowOriginalRequest `json:"originalRequest"`
}

// DialogFlowOriginalRequest struct
type DialogFlowOriginalRequest struct {
	Data DialogFlowOriginalRequestData `json:"data"`
}

// DialogFlowOriginalRequestData struct
type DialogFlowOriginalRequestData struct {
	Device DialogFlowOriginalRequestDevice `json:"device"`
}

// DialogFlowOriginalRequestDevice struct
type DialogFlowOriginalRequestDevice struct {
	Location DialogFlowOriginalRequestLocation `json:"location"`
}

// DialogFlowOriginalRequestLocation struct
type DialogFlowOriginalRequestLocation struct {
	Coordinates DialogFlowOriginalRequestCoordinates `json:"coordinates"`
}

// DialogFlowOriginalRequestCoordinates struct
type DialogFlowOriginalRequestCoordinates struct {
	Lat  float32 `json:"latitude"`
	Long float32 `json:"longitude"`
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

// DialogFlowResponseData struct
type DialogFlowResponseData struct {
	Google DialogFlowResponseGoogle `json:"google"`
}

// DialogFlowResponseGoogle struct
type DialogFlowResponseGoogle struct {
	ExpectUserResponse bool                           `json:"expectUserResponse"`
	IsSsml             bool                           `json:"isSsml"`
	SystemIntent       DialogFlowResponseSystemIntent `json:"systemIntent"`
}

// DialogFlowResponseSystemIntent struct
type DialogFlowResponseSystemIntent struct {
	Intent string                             `json:"intent"`
	Data   DialogFlowResponseSystemIntentData `json:"data"`
}

// DialogFlowResponseSystemIntentData struct
type DialogFlowResponseSystemIntentData struct {
	Type        string   `json:"@type"`
	OptContext  string   `json:"optContext"`
	Permissions []string `json:"permissions"`
}
