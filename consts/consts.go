package consts

type Features struct {
	Temperature      bool     `json:"temperature"`
	Precipitation    bool     `json:"precipitation"`
	Capital          bool     `json:"capital"`
	Coordinates      bool     `json:"coordinates"`
	Population       bool     `json:"population"`
	Area             bool     `json:"area"`
	TargetCurrencies []string `json:"targetCurrencies"`
}

type RegistrationRequestBody struct {
	Country     string   `json:"country"`
	IsoCode     string   `json:"isoCode"`
	Features    Features `json:"features"`
	TimeChanged string   `json:"timeChanged"` // TODO: Surely not ideal since users can now provide this?
}
type RegistrationRequestResponse struct {
	Id         string `json:"id"`
	LastChange string `json:"lastChange"`
}
