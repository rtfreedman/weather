package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var ApiKey = "809f761fd91b3990cdc45262b01aa174"
var client *http.Client

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

// Weather contains the information about weather at a certain zipcode
type Weather struct {
	CreatedAt   time.Time
	Expiry      time.Time
	Humidity    int
	Temperature float32
	WindSpeed   float32
	ZipCode     int
}

func GetFromZip(zip int) (w Weather, err error) {
	// get the weather from the cache if it's there
	w, present, expired := cachedWeather.retrieve(zip)
	if present && !expired {
		return
	}
	// otherwise fetch the weather
	err = w.fetch(zip)
	return
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

func (w *Weather) fetch(zip int) (err error) {
	// set the stuff we know about
	w.CreatedAt = time.Now()
	w.Expiry = w.CreatedAt.Add(time.Hour)
	w.ZipCode = zip

	// create the request to the api from the zip provided and the API Key
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%d&appid=%s", zip, ApiKey), nil)
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
