# CountryInfoService

> **CountryInfoService** is a simple REST API that provides country-related information  by integrating external APIs.

The code is deployed to Render and can be accessed at [https://countryinfoservice.onrender.com/](https://countryinfoservice.onrender.com/).

## Features

The service provides the following endpoints:

- **Country Information**: Get the name, capital, population, cities of a given country.
    - Optional: Limit the amount of cities returned.
- **Population Information**: Get the population history of a given country.
    - Optional: Get the population history for a given time interval.
- **Status Information**: Get the status of the service and external APIs.

For more information, see the [API Endpoints](#api-endpoints) section.

### External APIs
The project integrates the following external APIs to provide country-related information:

- **CountriesNow API**: Used to fetch details about countries, including historical population data and city listings.
- **RestCountries API**: Supplies additional country information, such as capitals, bordering nations, and national flags.

### Error Handling

All API responses follow a structured `JSON` format. If an error occurs, the response will include an error flag and a descriptive message.

**Example:**
```json
{
  "error": true,
  "message": "error fetching cities",
  "data": null
}
```

### Known issues

- The `/info` endpoint does not return data for South Sudan (`SS`/`SSD`). This is due to the absence of South Sudan in the `CountriesNow API`, preventing the retrieval of its city data.
- The "default" handler does _not_ function correctly on Render, as it attempts to serve the `index.html` file from the `static` directory, which is not currently working.

## Requirements

- Go 1.23
- External APIs: CountriesNow API and RestCountries API
- Go modules (use `go mod` for managing dependencies)

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/SigurdRiseth/CountryInfoService.git
   cd CountryInfoService
    ```

2.	Install dependencies (if you haven’t already):

 ```bash
 go mod tidy
 ```

3.	(optional) Create a .env file for environment variables. Here’s an example of what you might want to include:

 ```bash
 PORT=8080
 ```

4.	Run the service:

 ```bash
 go run main.go
 ```

5.	The service should now be running at http://localhost:8080/.

## API Endpoints

### GET /countryinfo/v1/info/

Returns the name, capital, population, and area of a country. The country name should be passed as a query parameter.

Example: http://localhost:8080/countryinfo/v1/info/no?limit=3

Response:
```json
{
  "error": false,
  "message": "Country information retrieved successfully",
  "data": {
    "name": "Norway",
    "continents": [
      "Europe"
    ],
    "population": 5379475,
    "languages": {
      "nno": "Norwegian Nynorsk",
      "nob": "Norwegian Bokmål",
      "smi": "Sami"
    },
    "borders": [
      "FIN",
      "SWE",
      "RUS"
    ],
    "flag": "🇳🇴",
    "capital": "Oslo",
    "cities": [
      "Abelvaer",
      "Adalsbruk",
      "Adland"
    ]
  }
}
```

### GET /countryinfo/v1/population/

Returns the population of a country. The country name should be passed as a query parameter.

Example: http://localhost:8080/countryinfo/v1/population/no?limit=2002-2004

Response:
```json
{
  "error": false,
  "message": "Population data retrieved successfully",
  "data": {
    "mean": 4564974,
    "values": [
      {
        "year": 2002,
        "value": 4538159
      },
      {
        "year": 2003,
        "value": 4564855
      },
      {
        "year": 2004,
        "value": 4591910
      }
    ]
  }
}
```

### GET /countryinfo/v1/status/

Returns the uptime of the service, API version, and status of the external APIs.

Example: http://localhost:8080/countryinfo/v1/status

Response:
```json
{
  "error": false,
  "message": "Service status retrieved successfully",
  "data": {
    "countriesnowapi": "Online (HTTP 200 OK)",
    "restcountriesapi": "Online (HTTP 200 OK)",
    "version": "1.0",
    "uptime": "104 seconds"
  }
}
```