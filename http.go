package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

var APIKey string
var client *http.Client

// 60 calls per minute
var rateLimiter = make(chan *time.Timer, 60)

// 1M items per month
var requestsSent *int64 = new(int64)
var lastResetTime = time.Now()

type apiResponse struct {
	Main main `json:"main"`
	Wind wind `json:"wind"`
}

type main struct {
	Temp     float32 `json:"temp"`
	Humidity int     `json:"humidity"`
}

type wind struct {
	Speed float32 `json:"speed"`
}

func init() {
	// initialize the client
	dialer := net.Dialer{
		Timeout: 10 * time.Second,
	}
	client = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Dial:                dialer.Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
}

func checkMillion() bool {
	if time.Now().Month() != lastResetTime.Month() {
		requestsSent = new(int64)
		lastResetTime = time.Now()
	}
	return atomic.AddInt64(requestsSent, 1) > 10
}

func (w *Weather) fetch(zip int) (err error) {
	if checkMillion() {
		return errors.New("monthly limit reached")
	}
	select {
	case rateLimiter <- time.NewTimer(time.Minute):
	default:
		return errors.New("rate limiting to 60 calls per minute reached")
	}
	// set the stuff we know about
	w.CreatedAt = time.Now()
	w.Expiry = w.CreatedAt.Add(time.Hour)
	w.ZipCode = zip

	// create the request to the api from the zip provided and the API Key
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%d&appid=%s", zip, APIKey), nil)
	if err != nil {
		return
	}

	// perform the request and read the body
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// if we didn't get a 200 the body is likely different. Here we return the first 100 characters of the body.
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Bad Status Code from API : %d; " + string(b)[:100] + "...")
		return
	}

	// map the response object back to the weather object
	var weatherResponse apiResponse
	if err = json.Unmarshal(b, &weatherResponse); err != nil {
		return
	}
	w.Humidity = weatherResponse.Main.Humidity
	w.Temperature = weatherResponse.Main.Temp
	w.WindSpeed = weatherResponse.Wind.Speed
	return
}
