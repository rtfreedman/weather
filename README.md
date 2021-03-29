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

...

func main() {
	// create a client from our api key
	client := weather.NewClient(os.Getenv("EXAMPLE_API_KEY"))

	// grab the weather based on the zip code we want
	w, err := client.GetFromZip(*zipCode)
	if err != nil {
		panic(err.Error())
	}
    ...
}
```
That's it: create a client (for each api key if you have multiple) and request the zip. The client will return an error if you're outside of your rate limit.

### Running the example
Run the following command: `EXAMPLE_API_KEY=YOUR_API_KEY go run cmd/example/main.go --zip ZIP` replacing `ZIP` with the zip code you're interested in and `YOUR_API_KEY` with your api key.

### Future Work
- Rate limiting is based on your tier, it'd be nice to have a client setting that let you configure further
- Testing would be nice outside of manual
- Better errors