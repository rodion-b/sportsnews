# Sports News Service

This is a simple Sports News Service that continously montiors latest ECB articles, persists them and exposes REST API that can be used for analytics or when the ECB server is down. The API also can be furthure extended to standardise news articles across different providers for upstream apps consumption.


## Functional requirements
- The system must periodically poll external endpoints to retrieve news articles from the following sources:
  - List of the latest *n* news articles:  
    `https://content-ecb.pulselive.com/content/ecb/text/EN/?pageSize=20`  
  - Details of a specific article by ID:  
    `https://content-ecb.pulselive.com/content/ecb/text/EN/{id}`  
- Implement a data management mechanism that:  
  - Handles updates to articles from the source gracefully without causing duplicate entries.
  - Limits the retention of historical data to ensure efficient database usage.
-   The system must provide two REST API endpoints:
    1. **Retrieve a list of all articles**  
       - Endpoint returns JSON with a list of articles.
    2. **Retrieve a single article by ID**  
       - Endpoint returns detailed JSON information for a specific article.
  

## Technical requirements
- **Language:** Golang  
- **Database:** MongoDB  
- The server needs to expose RESTful API with JSend specification

## How to run
- `make dc` runs docker-compose with the app container on port 8080 for you.
- `make test` runs the tests
- `make run` runs the app locally on port 8080 without docker.
- `make lint` runs the linter


## Solution notes
-   There are 2 endpoints available:
    - GET /articles?clientId=ecb&limit=1&offset=1 - Retrieve a list of articles for a clientId with limit and offset
    - GET /articles/b982f8ac-d74c-4666-91db-c917a2e19e19?clientId=ecb - Retrieve a single article by ID


- clean architecture (handler->service->repository)
- standard Go project layout
- docker compose + Makefile included
- simple server test is included
- The clean up is achived via setting up TTL index for articles collection that expire after 24 hours (configurable)