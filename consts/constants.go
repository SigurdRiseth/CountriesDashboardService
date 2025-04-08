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
)

// API Endpoint Base Paths
const (
	BaseDashboardPath    = "/dashboard/v1/dashboards/"
	BaseRegistrationPath = "/dashboard/v1/registrations/"
	BaseStatusPath       = "/dashboard/v1/status/"
	BaseNotificationPath = "/dashboard/v1/notifications/"
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
	QueryParamID = "id"
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
)

// Errors
const (
	MissingIDParamInURL = "Missing configuration ID in URL path"
	NotFound            = "not found"
)

const (
	InvalidJSONPayload = "Invalid JSON payload "
)

// firebase
const (
	FBNotInitialized = "firebase client is not initialized"
)
