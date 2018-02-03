package v2

import (
	"fmt"
	"log"
	"net/url"

	"github.com/thrasher-/gocryptotrader/common"
	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/exchanges"
	"github.com/thrasher-/gocryptotrader/exchanges/ticker"
)

const (
	bitfinexAPIURL     = "https://api.bitfinex.com/v2/"
	bitfinexAPIVersion = "2"
	bitfinexTrades     = "trades/"
	// bitfinexMaxRequests if exceeded IP address blocked 10-60 sec, JSON response
	// {"error": "ERR_RATE_LIMIT"}
	bitfinexMaxRequests = 90
)

// Bitfinex is the overarching type across the bitfinex package
// Notes: Bitfinex has added a rate limit to the number of REST requests.
// Rate limit policy can vary in a range of 10 to 90 requests per minute
// depending on some factors (e.g. servers load, endpoint, etc.).
type Bitfinex struct {
	exchange.Base
}

// SetDefaults sets the basic defaults for bitfinex
func (b *Bitfinex) SetDefaults() {
	b.Name = "Bitfinex"
	b.Enabled = false
	b.Verbose = false
	b.Websocket = false
	b.RESTPollingDelay = 10
	b.RequestCurrencyPairFormat.Delimiter = ""
	b.RequestCurrencyPairFormat.Uppercase = true
	b.ConfigCurrencyPairFormat.Delimiter = ""
	b.ConfigCurrencyPairFormat.Uppercase = true
	b.AssetTypes = []string{ticker.Spot}
}

// Setup takes in the supplied exchange configuration details and sets params
func (b *Bitfinex) Setup(exch config.ExchangeConfig) {
	if !exch.Enabled {
		b.SetEnabled(false)
	} else {
		b.Enabled = true
		b.AuthenticatedAPISupport = exch.AuthenticatedAPISupport
		b.SetAPIKeys(exch.APIKey, exch.APISecret, "", false)
		b.RESTPollingDelay = exch.RESTPollingDelay
		b.Verbose = exch.Verbose
		b.Websocket = exch.Websocket
		b.BaseCurrencies = common.SplitStrings(exch.BaseCurrencies, ",")
		b.AvailablePairs = common.SplitStrings(exch.AvailablePairs, ",")
		b.EnabledPairs = common.SplitStrings(exch.EnabledPairs, ",")
		err := b.SetCurrencyPairFormat()
		if err != nil {
			log.Fatal(err)
		}
		err = b.SetAssetTypes()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (b *Bitfinex) GetTrades(currencyPair string, values url.Values) ([]Trade, error) {
	var response [][]interface{}
	path := common.EncodeURLValues(
		bitfinexAPIURL+bitfinexTrades+currencyPair+"/hist",
		values,
	)
	err := common.SendHTTPGetRequest(path, true, b.Verbose, &response)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	trades, err := NewTradeListFromResponse(response)
	return trades, err
}

func NewTradeListFromResponse(response [][]interface{}) (ts []Trade, err error) {
	if len(response) == 0 {
		return
	}
	for _, v := range response {
		t, err := NewTradeFromResponse(v)
		if err != nil {
			return ts, err
		}
		ts = append(ts, t)
	}

	return
}

// public trade update just looks like a trade
func NewTradeFromResponse(response []interface{}) (o Trade, err error) {
	if len(response) == 4 {
		o = Trade{
			ID:     i64ValOrZero(response[0]),
			MTS:    i64ValOrZero(response[1]),
			Amount: f64ValOrZero(response[2]),
			Price:  f64ValOrZero(response[3])}
		return
	}
	return o, fmt.Errorf("data slice too short for trade update: %#v", response)
}
