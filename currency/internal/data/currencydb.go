package data

import (
	"context"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type CurrencyRatesI interface {
	MonitorRate(time.Duration) chan struct{}
	GetRate(string, string) (float64, error)
}

type CurrencyRatesDeps struct {
	Log *zap.SugaredLogger
	Ctx context.Context
}

type currencyRates struct {
	log   *zap.SugaredLogger
	ctx   context.Context
	rates map[string]float64
}

func NewCurrencyData(deps CurrencyRatesDeps) (CurrencyRatesI, error) {
	cr := &currencyRates{
		log:   deps.Log,
		ctx:   deps.Ctx,
		rates: map[string]float64{},
	}

	err := cr.getRates()

	return cr, err
}

// MonitorRates checks the rates in the ECB API every interval and sends a message to the
// returned channel when there are changes
//
// Note: the ECB API only returns data once a day, this function only simulates the changes
// in rates for demonstration purposes
func (cd *currencyRates) MonitorRate(interval time.Duration) chan struct{} {
	retChan := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				// just add a random difference to the rate and return it
				// this simulates the fluctuations in currency rates
				for k, v := range cd.rates {
					// change can be 10% of original value
					change := (rand.Float64() / 10)

					// is this a postive or negative change
					direction := rand.Intn(1)

					if direction == 0 {
						// new value with be min 90% of old
						change = 1 - change
					} else {
						// new value will be 110% of old
						change = 1 + change
					}
					// modify the rate
					cd.rates[k] = v * change
				}
				// notify updates, this will block unless there is a listener on the other go routines
				retChan <- struct{}{}
			}
		}
	}()

	return retChan
}

func (cd *currencyRates) GetRate(base string, destination string) (float64, error) {
	br, ok := cd.rates[base]
	if !ok {
		return 0, fmt.Errorf("Base Currency is not in dataset %+v", ok)
	}

	dr, ok := cd.rates[destination]
	if !ok {
		return 0, fmt.Errorf("Destination Currency is not in dataset %+v", ok)
	}

	return dr / br, nil
}

func (cd *currencyRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		cd.log.Errorf("[ERROR] While fetching data from http %+v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[ERROR] Fetching Status Code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	data := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&data)

	for _, cube := range data.CubeData {
		rate, err := strconv.ParseFloat(cube.Rate, 64)
		if err != nil {
			cd.log.Errorf("[ERROR] string to float64 %+v", err)
			return err
		}

		cd.rates[cube.Currency] = rate
	}

	cd.rates["EUR"] = 1

	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
