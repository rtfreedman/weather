# Weather Library
This is a weather library for golang that takes advantage of OpenweatherMap to retrieve the weather. It also will cache the 
## Usage
### Importing
As can be seen in `cmd/example/main.go`:
```go
package main

import (
    ...

	"github.com/rtfreedman/weather"
)

func main() {
	// set the api key
	weather.APIKey = os.Getenv("EXAMPLE_API_KEY")

	// grab the weather based on the zip code we want
	w, err := weather.GetFromZip(*zipCode)
	if err != nil {
		panic(err.Error())
	}
    ...
}
```
### Running the example
Run the following command: `EXAMPLE_API_KEY=YOUR_API_KEY go run cmd/example/main.go --zip ZIP` replacing `ZIP` with the zip code you're interested in and `YOUR_API_KEY` with your api key.