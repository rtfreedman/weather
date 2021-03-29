package weather

import (
	"errors"
	"time"
)

// Weather contains the information about weather at a certain zipcode
type Weather struct {
	CreatedAt   time.Time
	Expiry      time.Time
	Humidity    int
	Temperature float32
	WindSpeed   float32
	ZipCode     int
}

func (client *WeatherClient) GetFromZip(zip int) (w Weather, err error) {
	if client.apiKey == "" {
		err = errors.New("set APIKey")
		return
	}
	// get the weather from the cache if it's there
	w, present, expired := cachedWeather.retrieve(zip)
	if present && !expired {
		return
	}
	// otherwise fetch the weather
	if w, err = client.fetch(zip); err != nil {
		return
	}
	cachedWeather.add(w)
	return
}
