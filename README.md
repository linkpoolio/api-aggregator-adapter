# API Aggregator External Adaptor ![Travis-CI](https://travis-ci.org/linkpoolio/api-aggregator-cl-ea.svg?branch=master) [![codecov](https://codecov.io/gh/linkpoolio/asset-price-cl-ea/branch/master/graph/badge.svg)](https://codecov.io/gh/linkpoolio/api-aggregator-cl-ea)
External adaptor that can generically aggregate numerical values for any given amount of APIs.

**Supported Aggregation Methods:**
- Mode
- Median
- Mean

### Setup Instructions
#### Local Install
Make sure [Golang](https://golang.org/pkg/) is installed.

Build:
```
go get && go build -o api-aggregator-cl-ea main.go
```

Then run the adaptor:
```
./api-aggregator-cl-ea -p <port>
```

#### Docker
To run the container:
```
docker run -it -p 8080:8080 linkpoolio/api-aggregator-cl-ea
```

Container also supports passing in CLI arguments.

#### AWS Lambda

```bash
go get && go build -o lambda_handler main.go
zip api_aggregator.zip ./lambda_handler
```

Upload the the zip file into AWS and then use `lambda_handler` as the
handler.

**Important:** Set the `LAMBDA` environment variable to `true` in AWS for
the adaptor to be compatible with Lambda.

### Testing

To call the API, you need to send a POST request to `http://localhost:<port>/fetch` with the request body being of the ChainLink `RunResult` type.

For example:
```bash
curl -X POST http://localhost:8080/fetch \
-H 'Content-Type: application/json' \
-d @- << EOF
{
	"jobId": "1234",
	"data": {
		"api": ["https://www.bitstamp.net/api/v2/ticker/btcusd/", "https://api.pro.coinbase.com/products/btc-usd/ticker"],
		"paths": ["$.last", "$.price"],
		"aggregationType": "median"
	}
}
EOF
```
Should return something similar to:
```json
{
    "jobRunId": "",
    "status": "",
    "error": null,
    "pending": false,
    "data": {
        "api": [
            "https://www.bitstamp.net/api/v2/ticker/btcusd/",
            "https://api.pro.coinbase.com/products/btc-usd/ticker"
        ],
        "paths": [
            "$.last",
            "$.price"
        ],
        "aggregationType": "median",
        "aggregateValue": "4152.26",
        "failedApiCount": 0,
        "apiErrors": null
    }
}
```

### Contribution
We welcome any contributors. The more exchanges supported, the better. Feel free to raise any PR's or issues.