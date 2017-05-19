## Background

Language: Golang

Database: Redis 

## Flow

___Create___

1. Create a record in database with given path
2. Return a token (ID) to the client
3. Send an async request to Google to get the duration and distance
4. Update the record with the duration and distance when Google's calculation is completed


## Setup

1. Install Docker and Docker Compose
2. Create a Google API key and enable ___Google Maps Distance Matrix API___
3. Replace the value of ___GOOGLE_API_KEY___ in ___ENV___ in Dockerfile with your Google API key
4. Run `docker-compose build` 
4. Run `docker-compose up` 
5. Send a request to `localhost:3000`

## Remove 

Run `docker-compose down`
