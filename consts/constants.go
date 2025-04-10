package consts

// Firestore Collections
const (
	RegistrationsCollection  = "registrations"
	DashboardCacheCollection = "dashboardCache"
)

// home page
const (
	StaticFilePath = "./static/index.html"
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
	Bunny             = ""
	Dash              = "/"
	QueryParamID      = "id"
	QueryParamName    = "/name/"
	CurrencyParam     = "/currency/"
	QueryNameUnknown  = "%s/name/unknown"
	QueryMeteoLatLong = "?latitude=0&longitude=0&hourly=temperature_2m"
	SS                = "%s/%s"
)

// Status
const (
	StatusCountriesAPI   = "countries_api"
	StatusMeteoAPI       = "meteo_api"
	StatusCurrencyAPI    = "currency_api"
	StatusNotificationDB = "notification_db"
	StatusWebhooks       = "webhooks"
	StatusVersion        = "version"
	StatusUptime         = "uptime"
	StatusNOK            = "%s/NOK"
	V1                   = "v1"
)

// structs
const (
	FeaturesString         = "features"
	TimeChangedString      = "TimeChanged"
	CountryString          = "country"
	ISOCodeString          = "isoCode"
	LastRetrievedString    = "lastRetrieval"
	CapitalString          = "capital"
	CoordinatesString      = "coordinates"
	PopulationString       = "population"
	AreaString             = "area"
	LatitudeString         = "latitude"
	LongitudeString        = "longitude"
	HourlyString           = "hourly"
	TemperatureString2     = "temperature_2m"
	TemperatureString      = "temperature"
	PrecipitationString    = "precipitation"
	RatesString            = "rates"
	TargetCurrenciesString = "targetCurrencies"
	CurrenciesString       = "currencies"
	LatlngString           = "latlng"
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

	LogServerStarted              = "Server started on port: "
	LogStartingServer             = "Starting server..."
	LogColon                      = ":"
	LogFailedFetchWeatherAPI      = " Failed to fetch weather from API"
	FailCacheWeatherData          = "Failed to cache weather data: "
	LogHTTPReqFAil                = "HTTP request failed:"
	LogFailedJSONDecode           = "JSON decode failed:"
	LogHourlyDataMissingWeather   = "Hourly data missing in weather map"
	LogNoCurrencyCodeForCountry   = "No valid currency code found for the country"
	LogCurrencyCacheMiss          = "Currency cache MISS:"
	LogInvalidCacheTypesCurrency  = "Invalid cache type for currency rates"
	LogCurrencyCacheHIT           = "Currency cache HIT:"
	LogErrorFetchCurrency         = "Error fetching currency data:"
	LogErrorDecodeCurrency        = "Error decoding currency data:"
	LogMissingRatesCurrencyAPI    = "Missing 'rates' in currency API response"
	LogCountryDataCacheHIT        = "Country data cache HIT for: "
	LogCountryDataCacheMISS       = "Country data cache MISS for: "
	LogErrorRetrievingConfig      = "Error retrieving configuration:"
	LogErrorRetrievingDBConfig    = "Error retrieving dashboard configurations: "
	LogErrorRetrievingCountryData = "Error retrieving country data"
	LogErrorRetrievingDCFT        = "Error retrieving document from Firestore: "
	LogDBResponseSent             = "Dashboard response sent successfully:"
	LogErrorCacheCountryData      = "Error caching country data:"

	LogRegistrationsEndpointRecieved  = "Registrastions endpoint received: "
	LogRequest                        = " request"
	LogDocAddedToRC                   = "Document added to registrationsCollection. Identifier of returned document: "
	LogErrorConvertingDCtoDataStruct  = "Error converting document data to struct."
	LogErrorConvertingDCtoDataStruct2 = "Error converting document %s data to struct: %v"
	LogErrorIteratingFTDocs           = "Error iterating Firestore documents: %v"
	LogErrorUpdatingConfig            = "Error updating configuration: "
)

// Errors
const (
	Error                       = "error "
	MissingIDParamInURL         = "Missing configuration ID in URL path"
	NotFound                    = "not found"
	ClosingResponseBody         = "error closing response body: "
	FailedEncodeJSON            = "Failed to encode response as JSON: "
	ErrorEncodingResponse       = "Error encoding response: "
	MethodNotAllowed            = "Method not allowed"
	QueryingWebhooks            = "error querying webhooks: %v"
	EncodingErrorResponse       = "Error encoding error response: %v "
	CheckingServiceAvailability = "Error checking service availability for %s: %v"
	FailedCountRegisteredWH     = "Failed to count registered webhooks: "
	FailedRetrieveDBConfig      = "Failed to retrieve dashboard configuration from database"
	FailedCatchAllDBConfig      = "Error fetching all dashboard configurations: "
	FailedToFetchCountryData    = "Failed to fetch country data"
	UnsopprtedReqMethod         = "Unsupported request method: "
	MissingFields               = "Missing required fields: Name or Description"
	ErrorAddingDocToFT          = "Error when adding document to Firestore: "
	ErrorEncodingJSONResp       = "Error encoding JSON response: %v"
	MissingFieldCountryIso      = "Missing required fields: country or isoCode"
	DBClientUnavailable         = "Database client unavailable"
	ConfigNotFound              = "Configuration not found"
	ErrorDeletingConfig         = "Error deleting configuration: "
	ErrorMarshalDataToJSON      = "Error marshalling data to JSON: "
	ErrorCheckingDCExistence    = "error checking document existence: %w"
)

const (
	InvalidJSONPayload = "Invalid JSON payload "
)

// firebase
const (
	FBNotInitialized               = "firebase client is not initialized"
	FTNotInitialized               = "firestore client is not initialized"
	NoEnvFileFound                 = "No .env file found"
	GOOGLE_APPLICATION_CREDENTIALS = "GOOGLE_APPLICATION_CREDENTIALS"
	EnvironmentVarGACNotSet        = "environment variable GOOGLE_APPLICATION_CREDENTIALS is not set"
	FailedToInitializeFB           = "failed to initialize Firebase app: %v"
	FailedToInitializeFT           = "failed to initialize Firestore client: %v"
	FailedToCloseFT                = "failed to close Firestore client: %v"
	ErrorSettingDocument           = "error setting document: %v"
	ErrorUpdatingDocument          = "error updating document: %v"
)

// cache
const (
	CacheCollection        = "cache"
	FieldTimestamp         = "timestamp"
	FieldData              = "data"
	CacheCountryInfoPrefix = "country_"
	CacheWeatherPrefix     = "weather_"
	CacheCurrencyPrefix    = "currency_"
	TimeStamp              = "timestamp"
	LessThanAlligator      = "<"
	CountryKeyStamp        = "country_%s"
	CurrencyKeyStamp       = "currency_%s"
	WeatherKeyStamp        = "weather_%.2f:%.2f"

	CacheMissFor             = "cache miss for %s: %w"
	CacheExpiredFor          = "cache expired for %s: "
	WeatherCacheHit          = "Weather cache HIT"
	WeatherCacheMiss         = "Weather cache MISS"
	FailedToQueryExpiredDocs = "failed to query expired documents: %w"
	FailedToDeleteDocs       = "Failed to delete doc %s: %v"
	PurgedExpiredDocs        = "Purged %d expired documents from %s"
	AutoPurgedTriggered      = "Auto-purge triggered"
	CountryCache             = "CountryCache"
	CurrencyCache            = "CurrencyCache"
	WeatherCache             = "WeatherCache"
	ErrorPurging             = "Error purging %s: %v"
	PurgedSuccessfully       = "Successfully purged %s"
	FailedDecodeCacheEntry   = "failed to decode cache entry for %s: %w"
)
