# Countries Dashboard Service

> **CountriesDashboardService** is a...

## Table of Contents

- [Features](#features)
  - [External APIs](#external-apis)
  - [Error Handling](#error-handling)
  - [Known issues](#known-issues)
- [Requirements](#requirements)
- [Setup](#setup)
- [API Endpoints](#api-endpoints)
  - [Endpoint 'Registrations': Registering dashboard configuration](#endpoint-registrations-registering-dashboard-configuration)
  - [/dashboard/v1/dashboards/](#dashboardv1dashboards)
  - [/dashboard/v1/notifications/](#dashboardv1notifications)
  - [/dashboard/v1/status/](#dashboardv1status)
  
## Features

### External APIs

The services used in this application are:

#### REST Countries API

- **Endpoint:** [http://129.241.150.113:8080/v3.1](http://129.241.150.113:8080/v3.1)
- **Documentation:** [http://129.241.150.113:8080/](http://129.241.150.113:8080/)

#### Open-Meteo APIs

- **Documentation:** [https://open-meteo.com/en/features#available-apis](https://open-meteo.com/en/features#available-apis)

#### Currency API

- **Endpoint:** [http://129.241.150.113:9090/currency/](http://129.241.150.113:9090/currency/)
- **Documentation:** [http://129.241.150.113:9090/](http://129.241.150.113:9090/)

### Error Handling

### Known issues

## Requirements

- Go 1.23
- External APIs: CountriesNow API and RestCountries API
- Go modules (use `go mod` for managing dependencies)

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/SigurdRiseth/CountriesDashboardService.git
   cd CountriesDashboardService
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

### Endpoint 'Registrations': Registering dashboard configuration

The initial endpoint focuses on the management of dashboard configurations that can later be used via the `dashboards` endpoint.

#### Register new dashboard configuration:

Manages the registration of new dashboard configurations that indicate which information is to be shown for registered dashboards (via the `dashboards` endpoint – see later). This includes weather, country and currency exchange information.

##### - Request (POST)

```
Method: POST
Path: /dashboard/v1/registrations/
Content type: application/json
```

Body (exemplary code):

```json
{
   "country": "Norway",                                     // Indicates country name (alternatively to ISO code, i.e., country name can be empty if ISO code field is filled and vice versa)
   "isoCode": "NO",                                         // Indicates two-letter ISO code for country (alternatively to country name)
   "features": {
                  "temperature": true,                      // Indicates whether temperature in degree Celsius is shown
                  "precipitation": true,                    // Indicates whether precipitation (rain, showers and snow) is shown
                  "capital": true,                          // Indicates whether the name of the capital is shown
                  "coordinates": true,                      // Indicates whether country coordinates are shown
                  "population": true,                       // Indicates whether population is shown
                  "area": true,                             // Indicates whether land area size is shown
                  "targetCurrencies": ["EUR", "USD", "SEK"] // Indicates which exchange rates (to target currencies) relative to the base currency of the registered country (in this case NOK for Norway) are shown
               }
}
```

##### - Response

The response to the POST request on the endpoint stores the configuration on the server and returns the associated ID. In the example below, it is the ID `1`. Responses show be encoded in the above-mentioned JSON format, with the `lastChange` field highlighting the last change to the configuration (including updates via `PUT` – see later)

- Content type: `application/json`
- Status code: Appropriate error code. Ensure to deal with errors gracefully.

Body (exemplary code for registered configuration):

```json
{
    "id": 1
    "lastChange": "20240229 12:31"
}
```

#### View a specific registered dashboard configuration

Enables retrieval of a specific registered dashboard configuration.

##### - Request (GET)

The following shows a request for an individual configuration identified by its ID.

```
Method: GET
Path: /dashboard/v1/registrations/{id}
```

- id is the ID associated with the specific configuration.

Example request: `/dashboard/v1/registrations/1`

##### - Response

- Content type: application/json
- Status code: Appropriate error code. Ensure to deal with errors gracefully.

Body (exemplary code):

```json
{
   "id": 1,
   "country": "Norway",
   "isoCode": "NO",
   "features": {
                  "temperature": true,
                  "precipitation": true,
                  "capital": true,
                  "coordinates": true,
                  "population": true,
                  "area": false,
                  "targetCurrencies": ["EUR", "USD", "SEK"]
               },
    "lastChange": "20240229 14:07"
}
```

#### View all registered dashboard configurations

Enables retrieval of all registered dashboard configurations.

##### - Request (GET)

A GET request to the endpoint should return all registered configurations including IDs and timestamps of last change.

```
Method: GET
Path: /dashboard/v1/registrations/
```

##### - Response

- Content type: application/json
- Status code: Appropriate error code. Ensure to deal with errors gracefully.

Body (exemplary code):

```json
[
   {
      "id": 1,
      "country": "Norway",
      "isoCode": "NO",
      "features": {
                     "temperature": true,
                     "precipitation": true,
                     "capital": true,
                     "coordinates": true,
                     "population": true,
                     "area": false,
                     "targetCurrencies": ["EUR", "USD", "SEK"]
                  }, 
      "lastChange": "20240229 14:07"
   },
   {
      "id": 2,
      "country": "Denmark",
      "isoCode": "DK",
      "features": {
                     "temperature": false,
                     "precipitation": true,
                     "capital": true,
                     "coordinates": true,
                     "population": false,
                     "area": true,
                     "targetCurrencies": ["NOK", "MYR", "JPY", "EUR"]
                  },
       "lastChange": "20240224 08:27"
   },
   ...
]
```

The response should return a collection of return all stored configurations.

**Advanced Task:** Implement the HEAD method functionality (only return the header, not the body).

#### Replace a specific registered dashboard configuration

Enables the replacing of specific registered dashboard configurations.

##### - Request (PUT)

The following shows a request for an updated of individual configuration identified by its ID. This update should lead to an update of the configuration and an update of the associated timestamp (lastChange).

```
Method: PUT
Path: /dashboard/v1/registrations/{id}
```

- `id` is the ID associated with the specific configuration.

Example request: `/dashboard/v1/registrations/1`

Body (exemplary code):

```json
{
   "country": "Norway",
   "isoCode": "NO",
   "features": {
                  "temperature": false, // this value is to be changed
                  "precipitation": true,
                  "capital": true,
                  "coordinates": true, 
                  "population": true,
                  "area": false,
                  "targetCurrencies": ["EUR", "SEK"] // this value is to be changed
               }
}
```

Note that the request neither contains ID in the body (only in the URL), and neither contains the timestamp. The timestamp should be generated on the server side.

**Advanced Task:** Implement the PATCH method functionality.

##### - Response

This is the response to the change request.

- Status code: Appropriate error code. Ensure to deal with errors gracefully.
- Body: empty

#### Delete a specific registered dashboard configuration

Enabling the deletion of a specific registered dashboard configuration.

##### - Request (DELETE)

The following shows a request for deletion of an individual configuration identified by its ID. This update should lead to a deletion of the configuration on the server.

```
Method: DELETE
Path: /dashboard/v1/registrations/{id}
```

- `id` is the ID associated with the specific configuration.

Example request: `/dashboard/v1/registrations/1`

##### - Response

This is the response to the delete request.

- Status code: Appropriate error code. Ensure to deal with errors gracefully.
- Body: empty

### /dashboard/v1/dashboards/

### /dashboard/v1/notifications/

### /dashboard/v1/status/
