# Weather Library
This is a weather library for golang that takes advantage of OpenweatherMap to retrieve the weather. It also will cache the 
## Usage
As can be seen in `cmd/example/main.go` the usage is pretty straightforward
```go
package main

import (
    ...

	"github.com/rtfreedman/weather"
)

func main() {
	// set the api key
	weather.APIKey = os.Getenv("EXAMPLE_API_KEY") // or however you want to do it

	// grab the weather based on the zip code we want
	w, err := weather.GetFromZip(*zipCode)
	if err != nil {
		panic(err.Error())
	}
    ...
}
```
