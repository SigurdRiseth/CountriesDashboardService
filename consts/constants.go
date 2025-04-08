package consts

// Firestore Collections
const (
	RegistrationsCollection = "registrations"
)

// External API URLs
const (
	RestCountriesAPI = "http://129.241.150.113:8080/v3.1"
	OpenMeteoAPI     = "https://api.open-meteo.com/v1/forecast"
	CurrencyAPI      = "http://129.241.150.113:9090/currency"

	WeatherURL = "%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m,precipitation"
)

// API Endpoint Base Paths
const (
	BaseDashboardPath    = "/dashboard/v1/dashboards/"
	BaseRegistrationPath = "/dashboard/v1/registrations/"
	BaseStatusPath       = "/dashboard/v1/status/"
	BaseNotificationPath = "/dashboard/v1/notifications/"

	PortInUse = "8080"
)

// Mock/Test Endpoints for Dashboards
var (
	MockDashboardEndpointWithTestID    = BaseDashboardPath + "?id=test-id"
	MockDashboardEndpointWithoutID     = BaseDashboardPath
	MockDashboardEndpointWithInvalidID = BaseDashboardPath + "?id=nonexistent"
)

// HTTP header values
const (
	ContentTypeHeader   = "Content-Type"
	ApplicationJSON     = "application/json"
	ContentLengthHeader = "Content-Length"
)

// Query parameter keys
const (
	Dash              = "/"
	QueryParamID      = "id"
	QueryParamName    = "/name/"
	CurrencyParam     = "/currency/"
	QueryNameUnknown  = "%s/name/unknown"
	QueryMeteoLatLong = "?latitude=0&longitude=0&hourly=temperature_2m"
)

// Logging message templates (optional, can be used in `log.Printf`)
const (
	LogAddSuccess    = "Document added to %s. Identifier of returned document: %s "
	LogUpdateSuccess = "Successfully updated configuration with ID: %s "
	LogDeleteSuccess = "Successfully deleted configuration with ID: %s "
	LogPatchSuccess  = "Successfully patched configuration with ID: %s "
	LogGetSingle     = "Retrieving single dashboard configuration with ID: %s "
	LogGetAll        = "Retrieving all dashboard configurations "
	LogPutProcessing = "Processing PUT request for dashboard configuration "
	LogGETForService = "Processing GET request for service status"

	LogServerStarted  = "Server started on port"
	LogStartingServer = "Starting server..."
	LogColon          = ":"
)

// Errors
const (
	MissingIDParamInURL         = "Missing configuration ID in URL path"
	NotFound                    = "not found"
	ClosingResponseBody         = "Error closing response body: "
	FailedEncodeJSON            = "Failed to encode response as JSON: "
	ErrorEncodingResponse       = "Error encoding response: "
	MethodNotAllowed            = "Method not allowed"
	QueryingWebhooks            = "error querying webhooks: %v"
	EncodingErrorResponse       = "Error encoding error response: %v "
	CheckingServiceAvailability = "Error checking service availability for %s: %v"
	FailedCountRegisteredWH     = "Failed to count registered webhooks: "
)

const (
	InvalidJSONPayload = "Invalid JSON payload "
)

// firebase
const (
	FBNotInitialized = "firebase client is not initialized"
	FTNotInitialized = "firestore client is not initialized"
)
