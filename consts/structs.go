package consts

type Features struct {
	Temperature      *bool     `json:"temperature"`
	Precipitation    *bool     `json:"precipitation"`
	Capital          *bool     `json:"capital"`
	Coordinates      *bool     `json:"coordinates"`
	Population       *bool     `json:"population"`
	Area             *bool     `json:"area"`
	TargetCurrencies *[]string `json:"targetCurrencies"`
}

type UserUpdateRequest struct {
	Features Features `json:"features"`
}

type RegistrationRequestBody struct {
	Country     string   `json:"country"`
	IsoCode     string   `json:"isoCode"`
	Features    Features `json:"features"`
	TimeChanged string   `json:"timeChanged"`
}

type RegistrationRequestResponse struct {
	Id         string `json:"id"`
	LastChange string `json:"lastChange"`
}

type WebhookRegistration struct {
	Url         string `json:"url"`
	Country     string `json:"country"`
	Event       string `json:"event"`
	TimeChanged string `json:"timeChanged"`
}

// WebhookPayload represents the JSON structure of the request body
type WebhookPayload struct {
	Id      string `json:"id"`
	Country string `json:"country"`
	Event   string `json:"event"`
	Time    string `json:"time"`
}
