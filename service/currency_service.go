package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const layoutDateOnly = "2006-01-02"

type ApiResponse []struct {
	Date       time.Time `json:"date"`
	Currencies []struct {
		Code          string    `json:"code"`
		Quantity      int       `json:"quantity"`
		RateFormated  string    `json:"rateFormated"`
		DiffFormated  string    `json:"diffFormated"`
		Rate          float64   `json:"rate"`
		Name          string    `json:"name"`
		Diff          float64   `json:"diff"`
		Date          time.Time `json:"date"`
		ValidFromDate time.Time `json:"validFromDate"`
	} `json:"currencies"`
}

func GetCurrencyAtDate(date time.Time, currency string) (float64, error) {
	const BaseUrl = "https://nbg.gov.ge/gw/api/ct/monetarypolicy/currencies/en/json/"
	const RequestCurrencyKey = "currencies"
	const RequestDateKey = "date"

	urlParsed, errUrl := url.Parse(BaseUrl)
	if errUrl != nil {
		log.Fatal(errUrl)
		return 1, errUrl
	}
	query := urlParsed.Query()
	query.Add(RequestCurrencyKey, currency)
	query.Add(RequestDateKey, date.Format(layoutDateOnly))
	urlParsed.RawQuery = query.Encode()

	response, errFetching := http.Get(urlParsed.String())
	if errFetching != nil {
		log.Fatal(errUrl)
		return 1, errUrl
	}
	defer response.Body.Close()

	rawBody, _ := io.ReadAll(response.Body)
	var result ApiResponse
	if err := json.Unmarshal(rawBody, &result); err != nil {
		log.Fatal(err)
		return 1, err
	}
	return result[0].Currencies[0].Rate, nil

}
