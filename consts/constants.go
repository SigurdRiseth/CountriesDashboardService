package consts

// Constants for external API URLs and endpoints
var (
	RestCountriesAPI = "http://129.241.150.113:8080/v3.1"
	OpenMeteoAPI     = "https://api.open-meteo.com/v1/forecast"
	CurrencyAPI      = "http://129.241.150.113:9090/currency"
)

// Endpoints for the dashboard service tests
var (
	MockDashboardEndpointWithTestID    = "/dashboard/v1/dashboards/?id=test-id"
	MockDashboardEndpointWithoutID     = "/dashboard/v1/dashboards/"
	MockDashboardEndpointWithInvalidID = "/dashboard/v1/dashboards/?id=nonexistent"
)
