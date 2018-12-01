package server

import (
	"encoding/json"
	"fmt"
	"github.com/montanaflynn/stats"
	"github.com/oliveagle/jsonpath"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"gopkg.in/guregu/null.v3"
)

func getAggregate(rr *RunResult) (*RunResult, error) {
	var wg sync.WaitGroup
	t := len(rr.Data.API)
	wg.Add(t)

	values := make(chan float64, t)
	errs := make(chan error, t)

	if (len(rr.Data.API) == 0 && len(rr.Data.Paths) == 0) ||
		len(rr.Data.API) != len(rr.Data.Paths) {
		rr.Error = null.StringFrom("invalid api and path array")
		return rr, nil
	}

	for i, a := range rr.Data.API {
		go performRequest(&wg, a, rr.Data.Paths[i], values, errs)
	}

	wg.Wait()
	close(values)
	close(errs)

	if len(errs) > 0 {
		rr.Data.FailedAPICount = len(errs)
		var errArr []string
		for err := range errs {
			errArr = append(errArr, err.Error())
		}
		rr.Data.APIErrors = errArr
	}

	if aggValue, err := aggregateValues(rr.Data.AggregationType, values); err != nil {
		rr.Error = null.StringFrom(fmt.Sprintf("error aggregating value: %s", err))
	} else {
		rr.Data.AggregateValue = aggValue
	}

	return rr, nil
}

func performRequest(
	wg *sync.WaitGroup,
	api string,
	path string,
	values chan<- float64,
	errs chan<- error,
) {
	defer wg.Done()
	hc := http.Client{}

	resp, err := hc.Get(api)
	if err != nil {
		errs <- err
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errs <- err
		return
	}

	var jd interface{}
	json.Unmarshal([]byte(body), &jd)

	val, err := jsonpath.JsonPathLookup(jd, path)
	if err != nil {
		errs <- err
		return
	}
	fv, err := strconv.ParseFloat(fmt.Sprint(val), 64)
	if err != nil {
		errs <- err
		return
	}

	values <- fv
}

func aggregateValues(aggType string, values chan float64) (string, error) {
	var av float64
	var err error
	var valArr []float64

	for v := range values {
		valArr = append(valArr, v)
	}

	switch aggType {
	case "mode":
		var modeArr []float64
		modeArr, err = stats.Mode(valArr)
		if len(modeArr) == 0 {
			av = valArr[0]
		} else {
			av = modeArr[0]
		}
		break
	case "median":
		av, err = stats.Median(valArr)
		break
	default:
		av, err = stats.Mean(valArr)
		break
	}

	return fmt.Sprint(av), err
}
