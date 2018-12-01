package server

import (
	"encoding/json"
	"flag"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/aws/aws-lambda-go/lambda"
	"fmt"
)

type Config struct {
	Port int
	Level string
}

func Start() {
	if len(os.Getenv("LAMBDA")) > 0 {
		lambda.Start(getAggregate)
	} else {
		c := getConfig()
		http.HandleFunc("/fetch", fetch)
		log.Print("API Aggregator Listening on Port ", c.Port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil))
	}
}

func fetch(w http.ResponseWriter, r *http.Request) {
	var rr RunResult
	sc := http.StatusOK
	if b, err := ioutil.ReadAll(r.Body); err != nil {
		rr.Error = null.StringFrom(err.Error())
		sc = http.StatusInternalServerError
	} else if err = json.Unmarshal(b, &rr); err != nil {
		rr.Error = null.StringFrom(err.Error())
		sc = http.StatusBadRequest
	} else if getAggregate(&rr); !rr.Error.IsZero() {
		sc = http.StatusBadRequest
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sc)
	json.NewEncoder(w).Encode(rr)
	logRequest(sc, &rr)
}

func getConfig() *Config {
	var c Config
	flag.IntVar(&c.Port, "p", 8080, "Port number to serve")
	flag.StringVar(&c.Level, "level", "release", "Logging level")
	flag.Parse()
	return &c
}

func logRequest(statusCode int, runResult *RunResult) {
	if statusCode == 200 {
		log.Print("Request Completed, Status Code ", statusCode, ", API Count ", len(runResult.Data.API))
	} else {
		log.Print("Request Error, Status Code ", statusCode, ", Error ", runResult.Error.String)
	}
}