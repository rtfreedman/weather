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

var ErrRateLimitReached = errors.New("rate limit reached")
var ErrMonthlyLimitReached = errors.New("monthly limit reached")

type WeatherClient struct {
	rateLimiter   chan bool
	requestsSent  *int64
	lastResetTime time.Time
	apiKey        string
	http.Client
}

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

// NewClient constructs a client for the api key provided
func NewClient(apiKey string) *WeatherClient {
	// initialize the client
	dialer := net.Dialer{
		Timeout: 10 * time.Second,
	}
	w := &WeatherClient{
		// TODO: make the size of the rate limiters configurable
		rateLimiter:  make(chan bool, 60),
		requestsSent: new(int64),
		apiKey:       apiKey,
		Client: http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				Dial:                dialer.Dial,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		}}
	return w
}

func (client *WeatherClient) checkMillion() bool {
	if time.Now().Month() != client.lastResetTime.Month() {
		client.requestsSent = new(int64)
		client.lastResetTime = time.Now()
	}
	return atomic.AddInt64(client.requestsSent, 1) > 1000000
}

func (client *WeatherClient) fetch(zip int) (w Weather, err error) {
	if client.checkMillion() {
		err = ErrMonthlyLimitReached
		return
	}
	// check the rate limiter
	// if we can push a new timer in then there's still space this minute
	select {
	case client.rateLimiter <- true:
	default:
		err = ErrRateLimitReached
		return
	}
	// after a minute, pull off the items on the rateLimiter
	time.AfterFunc(time.Minute, func() {
		<-client.rateLimiter
	})

	// set the stuff we know about
	w.CreatedAt = time.Now()
	w.Expiry = w.CreatedAt.Add(time.Hour)
	w.ZipCode = zip

	// create the request to the api from the zip provided and the API Key
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%d&appid=%s", zip, client.apiKey), nil)
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
		err = fmt.Errorf("Bad Status Code from API : %d; " + string(b))
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
