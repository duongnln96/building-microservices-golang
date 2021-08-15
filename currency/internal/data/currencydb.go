package data

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type CurrencyRates interface {
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

func NewCurrencyData(deps CurrencyRatesDeps) (CurrencyRates, error) {
	cr := &currencyRates{
		log:   deps.Log,
		ctx:   deps.Ctx,
		rates: map[string]float64{},
	}

	err := cr.getRates()

	return cr, err
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
