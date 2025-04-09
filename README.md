# Countries Dashboard Service

This project is a **RESTful web application** developed in **Golang** that allows users to configure and retrieve dynamic information dashboards. These dashboards are populated using real-time data from external APIs and are stored persistently to survive service restarts. Additionally, the application includes a simple **notification service** that listens for specific events using webhooks.

The project is part of a group assignment where we will implement state management using databases, utilize webhooks for event notifications, and deploy the service using **Docker** and an **Infrastructure-as-a-Service (IaaS) system**.

## Table of Contents

- [Features](#features)
  - [Units Used in Dashboard Responses](#units-used-in-dashboard-responses)
  - [External APIs](#external-apis)
  - [Error Handling](#error-handling)
  - [Known issues](#known-issues)
- [Requirements](#requirements)
- [Setup](#setup)
- [Deployment](#deployment)
- [API Endpoints](#api-endpoints)
  - [Registrations: Registering dashboard configuration](#registrations-registering-dashboard-configuration)
  - [Dashboards: Retrieve populated dashboard](#dashboards-retrieve-populated-dashboard)
  - [Notifications: Managing webhooks for event notifications](#notifications-managing-webhooks-for-event-notifications)
  - [Status: Monitoring service availability](#status-monitoring-service-availability)
- [Contributors](#contributors)
  
## Features

### Units Used in Dashboard Responses

When retrieving dashboard data via `/dashboard/v1/dashboards/{id}`, the values for certain features are presented in specific units:

| Feature          | Unit                      |
|------------------|---------------------------|
| `temperature`     | Degrees Celsius (°C)      |
| `precipitation`   | Millimeters (mm)          |
| `area`            | Square kilometers (km²)   |
| `population`      | Number of people          |
| `targetCurrencies`| Exchange rate (1 base currency → target currency) |

All timestamps such as `lastRetrieval` follow the **RFC3339** format, e.g., `"2025-04-03T10:15:00+01:00"`.

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

To run and develop this project locally, make sure you have the following installed:

### General
- [Go (Golang)](https://golang.org/dl/) v1.23
- [Docker](https://www.docker.com/) (optional for containerized deployment)
- [Git](https://git-scm.com/) (for cloning the repository)
- Internet access (for dependency downloads and external services)

### Firebase / Google Cloud
- A service account key JSON file with Firestore access
  - Set the path using the `GOOGLE_APPLICATION_CREDENTIALS` environment variable in a `.env` file

⚠️ Do not commit your credentials to version control.

### Environment Configuration
Create a `.env` file in the project root with the following (example format):

```env
GOOGLE_APPLICATION_CREDENTIALS=path/to/serviceAccountKey.json
```

## Setup

### Local Development

1. **Clone the repository**:

 ```bash
 git clone https://github.com/SigurdRiseth/CountriesDashboardService.git
 cd CountriesDashboardService
 ```

2. **Install dependencies** (if you haven’t already):

```bash
go mod tidy
```

3. **Create a `.env` file for environment variables** (if not already present):

```bash
GOOGLE_APPLICATION_CREDENTIALS=path/to/serviceAccountKey.json
```

⚠️ Make sure the service account file exists at the specified path and has Firestore access.

4. **Run the service**:

```bash
go run main.go
```

5. **Access the API**:

The service will be available at `http://localhost:8080/`.

### Deploy with Docker Compose

1. **Ensure Docker and Docker Compose is installed** on your machine.

2. **Place your Firebase service account key** in the root directory (or update the volume path accordingly).

For example:

```
CountriesDashboardService/
├── docker-compose.yml
├── credentials.json        <-- Your service account key
└── .env
```

3. **Build and run the Docker container**:

```bash
docker-compose up --build
```

4. **Access the API**:

The service will be available at `http://localhost:8080/`.

5. **To stop the service**, run:

```bash
docker-compose down
```

## Deployment

The application has been successfully deployed on an OpenStack instance (SkyHigh) and is accessible at the following address:
- [http://10.212.174.14:8080](http://10.212.174.14:8080)

## API Endpoints

### Registrations: Registering dashboard configuration

This endpoint handles the creation, retrieval, updating, and deletion of dashboard configurations. These configurations define the data to be shown in the dashboards (see the `dashboards` endpoint).

---

#### Register a New Dashboard Configuration

Registers a new dashboard configuration with selected features such as weather data, country info, and currency exchange rates.

##### Request

- **Method**: `POST`
- **Path**: `/dashboard/v1/registrations/`
- **Content-Type**: `application/json`

###### Example Request Body

```json
{
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": true,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": true,
    "targetCurrencies": ["EUR", "USD", "SEK"]
  }
}
```

##### Response

- **Status Code**: `201 Created` or appropriate error code
- **Content-Type**: `application/json`

###### Response Body

```json
{
  "id": 1,
  "lastChange": "2024-02-29 12:31"
}
```

---

#### View a Specific Dashboard Configuration

Retrieves a single dashboard configuration by its ID.

##### Request

- **Method**: `GET`
- **Path**: `/dashboard/v1/registrations/{id}`

##### Response

- **Status Code**: `200 OK` or appropriate error code
- **Content-Type**: `application/json`

###### Example Response Body

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

---

#### View All Dashboard Configurations

Returns all registered dashboard configurations with IDs and timestamps.

##### Request

- **Method**: `GET`
- **Path**: `/dashboard/v1/registrations/`

##### Response

- **Status Code**: `200 OK` or appropriate error code
- **Content-Type**: `application/json`

###### Example Response Body

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
    "lastChange": "2024-02-29 14:07"
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
    "lastChange": "2024-02-24 08:27"
  }
]
```

---

#### HEAD Request

Supports the ``HEAD`` method to return only the headers (e.g., status and metadata) without the response body for both `GET` requests.

##### Request

- **Method**: `HEAD`
- **Paths**:
  - `/dashboard/v1/registrations/`
  - `/dashboard/v1/registrations/{id}`

---

#### Replace a Dashboard Configuration

Enables the replacing of specific registered dashboard configurations.

##### - Request (PUT)

Updates a dashboard configuration by replacing it entirely.

##### Request

- **Method**: `PUT`
- **Path**: `/dashboard/v1/registrations/{id}`
- **Content-Type**: `application/json`

###### Example Request Body

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

##### Response

- **Status Code**: `204 No Content` or appropriate error code
- **Body**: empty

---

#### Partial Update (PATCH)

Supports partial updates to specific fields within a configuration.

##### Request

- **Method**: `PATCH`
- **Path**: `/dashboard/v1/registrations/{id}`
- **Content-Type**: `application/json`

###### Example Request Body

```json
{
   "features": {
                  "temperature": false, // this value is to be changed
                  "targetCurrencies": ["EUR", "SEK"] // this value is to be changed
               }
}
```

##### Response

- **Status Code**: `204 No Content` or appropriate error code
- **Body**: empty

---

#### Delete a Dashboard Configuration

Deletes a specific configuration by ID.

##### Request

- **Method**: `DELETE`
- **Path**: `/dashboard/v1/registrations/{id}`

##### Response

- **Status Code**: `204 No Content` or appropriate error code
- **Body**: empty

--- 

### Dashboards: Retrieve populated dashboard

This endpoint retrieves the populated dashboard data based on the registered configuration.

#### Request

- **Method**: `GET`
- **Path**: `/dashboard/v1/dashboards/{id}`

#### Response

- **Status Code**: `200 OK` or appropriate error code
- **Content-Type**: `application/json`

###### Example Response Body

```json
{
    "country": "Poland",
    "features": {
        "area": 312679,
        "targetCurrencies": {
            "SEK": 2.567791
        }
    },
    "isoCode": "NO",
    "lastRetrieval": "2025-04-09T15:55:51+02:00"
}
```
---

### Notifications: Managing webhooks for event notifications

This endpoint allows users to register and manage webhooks to receive notifications about changes and events related to dashboards and their configurations.

---

#### Register a Webhook

This endpoint allows users to register a webhook URL to receive notifications about events related to a specific country.

Supported event types:
- `REGISTER`: Triggered when a new dashboard configuration is registered.
- `CHANGE`: Triggered when an existing dashboard configuration is updated.
- `DELETE`: Triggered when a dashboard configuration is deleted.
- `INVOKE`: Triggered when a dashboard is accessed.

##### Request

- **Method**: `POST`
- **Path**: `/dashboard/v1/notifications/`
- **Content-Type**: `application/json`

###### Example Request Body

```json
{
  "country": "NO",
  "url": "https://example.com/webhook",
  "event": "DELETE"
}
```

This registers a webhook for the country "NO" to receive notifications when a dashboard configuration for this country is deleted.

##### Response

- **Status Code**: `201 Created` or appropriate error code
- **Content-Type**: `application/json`

###### Example Response Body

```json
{
  "id": "OIdksUDwveiwe"
}
```

---

#### Deletion of a Webhook

This endpoint allows users to delete a registered webhook by its ID.

##### Request

- **Method**: `DELETE`
- **Path**: `/dashboard/v1/notifications/{id}`

##### Response

- **Status Code**: `204 No Content` or appropriate error code
- **Body**: empty

---

#### View a _Specific_ Webhook

This endpoint retrieves a specific webhook by its ID.

##### Request

- **Method**: `GET`
- **Path**: `/dashboard/v1/notifications/{id}`

##### Response

- **Status Code**: `200 OK` or appropriate error code
- **Content-Type**: `application/json`

###### Example Response Body

```json
{
  "id": "OIdksUDwveiwe",
  "country": "NO",
  "url": "https://example.com/webhook",
  "event": "DELETE"
}
```

---

#### View All Webhooks
This endpoint retrieves all registered webhooks.

##### Request

- **Method**: `GET`
- **Path**: `/dashboard/v1/notifications/`

##### Response

- **Status Code**: `200 OK` or appropriate error code
- **Content-Type**: `application/json`

###### Example Response Body

```json
[
  {
    "id": "OIdksUDwveiwe",
    "country": "NO",
    "url": "https://example.com/webhook",
    "event": "DELETE"
  },
  {
    "id": "OIdksUDwdsaiwe",
    "country": "SE",
    "url": "https://example.com/webhook",
    "event": "INVOKE"
  }
]
```

---

#### Webhook Invocation

When a webhook is triggered, the service sends a (POST) notification to the registered URL with the relevant event data.

##### Request

- **Method**: `POST`
- **Path**: `<url specified in the corresponding webhook registration>`
- **Content-Type**: `application/json`

###### Example Request Body

```json
{
  "id": "OIdksUDwveiwe",
  "country": "NO",
  "event": "INVOKE",
  "timestamp": "2024-02-29T14:07:00+01:00"
}
```

---

### Status: Monitoring service availability

This endpoint provides a comprehensive overview of the service’s current health, dependencies, and operational metrics.

#### Request

- **Method**: `GET`
- **Path**: `/dashboard/v1/status/`

#### Response

- **Status Code**: `200 OK` or appropriate error code
- **Content-Type**: `application/json`

###### Example Response Body

```json
{
  "countries_api": 404,
  "currency_api": 200,
  "meteo_api": 200,
  "notification_db": 200,
  "uptime": 26,
  "version": "v1",
  "webhooks": 6
}
```

## Contributors

This project was developed by the following contributors:

### [**Sigurd Riseth**](https://github.com/SigurdRiseth)
**Main Contributions:**
- Created the initial repository and project structure
- Implemented key methods for the `registrations` endpoint:
  - `POST /dashboard/v1/registrations/`
  - `GET /dashboard/v1/registrations/{id}`
  - `GET /dashboard/v1/registrations/`
  - `PATCH /dashboard/v1/registrations/{id}`
  - `HEAD /dashboard/v1/registrations/`
- Developed the `notifications` endpoint, including webhook functionality and testing
- Deployed the service using Docker on NTNU’s OpenStack instance *SkyHigh*

---

### [**Theodor Sjetnan Utvik**](https://github.com/TheodorUtvik)
**Main Contributions:**
- Implemented the `status` and `dashboards` endpoints
- Contributed to the `registrations` endpoint with the following methods:
  - `PATCH /dashboard/v1/registrations/{id}`
  - `PUT /dashboard/v1/registrations/{id}`
  - `DELETE /dashboard/v1/registrations/{id}`
- Wrote and executed tests for the `registrations`, `dashboards`, and `status` endpoints

---

For detailed activity logs, visit the [**GitHub Insights**](https://github.com/SigurdRiseth/CountriesDashboardService/pulse) or the [**Contributors Page**](https://github.com/SigurdRiseth/CountriesDashboardService/graphs/contributors).  
You can also explore GitHub Issues and Pull Requests for a history of discussions and changes.