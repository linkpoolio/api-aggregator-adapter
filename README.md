# API Aggregator External Adapter 
Adapter that can generically aggregate numerical values for any given amount of APIs.

This adaptor is built using the [bridges](https://github.com/linkpoolio/bridges) framework.

**Supported Aggregation Methods:**
- Mode
- Median
- Mean

### Contract Usage

Example: https://github.com/linkpoolio/example-chainlinks/blob/master/contracts/APIAggregatorConsumer.sol

### Setup Instructions
#### Local Install
Make sure [Golang](https://golang.org/pkg/) is installed.

Build (in the root of the bridges repository):
```
GO111MODULE=on go build examples/apiaggregator/main -o apiaggregator
```

Then run the bridge:
```
./apiaggregator
```

#### Docker
To run the container:
```
docker run -it -p 8080:8080 linkpool/api-aggregator-adapter:latest
```

#### AWS Lambda

```bash
zip api_aggregator.zip ./apiaggregator
```

Upload the the zip file into AWS and then use `apiaggregator` as the
handler.

**Important:** Set the `LAMBDA` environment variable to `true` in AWS for
the adaptor to be compatible with Lambda.

### Testing

To call the API, you need to send a POST request to `http://localhost:<port>/` with the request body being of the ChainLink `RunResult` type.

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
		"type": "median"
	}
}
EOF
```
Should return something similar to:
```json
{
    "jobRunId": "1234",
    "status": "completed",
    "error": null,
    "pending": false,
    "data": {
        "EUR": 141.65,
        "JPY": 17864.71,
        "USD": 160.11,
        "type": "median",
        "api": [
            "https://www.bitstamp.net/api/v2/ticker/btcusd/",
            "https://api.pro.coinbase.com/products/btc-usd/ticker"
        ],
        "paths": [
            "$.last",
            "$.price"
        ]
    }
}
```
