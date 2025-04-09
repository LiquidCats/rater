# Rater

This microservice is used to get exchange rates for crypto 

## How to use 

### Configure Env 

```dotenv
RATER_APP_PORT=8080
RATER_APP_PAIRS="BTC_USD,BTC_EUR,ETH_USD,ETH_EUR"

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
GET http://localhost:8080/v1/rate/{pair}

{
    "data": {
        "price": "29295.92969",
        "pair": "BTC_USD",
    },
    "status": "success"
}
```

