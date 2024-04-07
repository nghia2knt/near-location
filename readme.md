# FIND NEAR LOCATION README

This repository support for query near location by lat+lon input.

## Request

Imagine you are assigned to develop a feature similar to Zalo's "Search around here" feature. Let's design the DB schema, and processing model, for this feature. Input is the lat+long of that user, the DB contains 10 million datapoints (lat+lon) of 10 million other users, find an algorithm to search to effectively list users from near to far (API response has pagination)

## Database design

- Collection: users_locations

+----------------------+-----------------------------------------+
| Field                | Type                                    |
+----------------------+-----------------------------------------+
| _id                  | ObjectID  (uniqueIndex - primary key)   |
| created_at           | DATETIME                                |
| updated_at           | DATETIME                                |
| deleted_at           | DATETIME  (default not exits)           |
| user_id              | String    (uuid format) (uniqueIndex)   |
| location             | GeoJSON   (2dsphere index)              |
+----------------------+-----------------------------------------+


- GeoJSON format: 

<field>: { type: <GeoJSON type> , coordinates: <coordinates> }
<GeoJSON type>: currently support type "Point"
<coordinates>: [ longitude, latitude ]

## Prerequisites

Before running the program or tests, ensure that you have the following installed:

- Docker
- Docker Compose
- GoLang

## Running the Program

1. Run Unit Test:

    ```bash
   go test -coverprofile='coverage.out' ./...
   go tool cover -html='coverage.out'
   ```

2. Build Program:

    ```bash
    docker-compose build
    ```

3. Run Docker Compose:

    ```bash
    docker-compose up -d
    ```

4. Test:

    - Using `near-user.postman_collection.json` for run get location api.