# Rater

This microservice is used to get exchange rates for crypto 

## How to use 

### Configure Env 

```dotenv
RATER_PORT=8080
RATER_BASE_CURRENCIES="BTC,ETH"
RATER_QUOTE_CURRENCIES="USD,EUR"

RATER_COINGATE_URL="https://api.coingate.com/v2/rates/merchant"

RATER_COINAPI_URL="https://rest.coinapi.io/v1/exchangerate"
RATER_COINAPI_SECRET=""

RATER_REDIS_ADDRESS="redis:6379"
RATER_REDIS_DB=0
```

### Build and Run

```
docker compose up --build app
```

### Response
After running application you can access to endpoint
```http request
GET http://localhost:8080/v1/rate/:base/:quote

{
    "data": {
        "Price": "29295.92969",
        "base": "BTC",
        "quote": "USD"
    },
    "status": "success"
}
```

