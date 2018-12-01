package server

import (
	"net/http"
	"testing"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"gopkg.in/guregu/null.v3"
)

func TestFetch_All(t *testing.T) {
	aggregationTypes := []string{
		"", "mean", "median", "mode", "nonexistent",
	}
	for _, at := range aggregationTypes {
		t.Run(at, func (t *testing.T) {
			rr := RunResult{
				Data: Params{
					API: []string{
						"https://www.bitstamp.net/api/v2/ticker/btcusd/",
						"https://api.pro.coinbase.com/products/btc-usd/ticker",
					},
					Paths: []string{
						"$.last",
						"$.price",
					},
					AggregationType: at,
				},
			}
			rec, err := serveFetch(&rr)
			assert.NoError(t, err)
			assert.Equal(t, rec.Code, http.StatusOK)

			rb, err := ioutil.ReadAll(rec.Body)
			err = json.Unmarshal(rb, &rr)
			assert.NoError(t, err)

			assert.NotEqual(t, rr.Data.AggregateValue, "")
		})
	}
}

func TestFetch_EmptyParam(t *testing.T) {
	rr := RunResult{}
	rec, err := serveFetch(&rr)
	assert.NoError(t, err)
	assert.Equal(t, rec.Code, http.StatusBadRequest)

	rb, err := ioutil.ReadAll(rec.Body)
	err = json.Unmarshal(rb, &rr)
	assert.NoError(t, err)

	assert.Equal(t, rr.Error, null.StringFrom("invalid api and path array"))
}

func TestFetch_InvalidArray(t *testing.T) {
	rr := RunResult{
		Data: Params{
			API: []string{
				"https://www.bitstamp.net/api/v2/ticker/btcusd/",
				"https://api.pro.coinbase.com/products/btc-usd/ticker",
			},
			Paths: []string{
				"$.last",
			},
		},
	}
	rec, err := serveFetch(&rr)
	assert.NoError(t, err)
	assert.Equal(t, rec.Code, http.StatusBadRequest)

	rb, err := ioutil.ReadAll(rec.Body)
	err = json.Unmarshal(rb, &rr)
	assert.NoError(t, err)

	assert.Equal(t, rr.Error, null.StringFrom("invalid api and path array"))
}

func serveFetch(runResult *RunResult) (*httptest.ResponseRecorder, error) {
	rec := httptest.NewRecorder()
	b, err := json.Marshal(runResult)
	req, err := http.NewRequest("POST", "/fetch", bytes.NewReader(b))
	if err != nil {
		return rec, err
	}
	handler := http.HandlerFunc(fetch)
	handler.ServeHTTP(rec, req)
	return rec, nil
}