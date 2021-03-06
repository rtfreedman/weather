package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/rtfreedman/weather"
)

var zipCode = flag.Int("zip", 20175, "supply zip code")

func main() {
	// create a client from our api key
	client := weather.NewClient(os.Getenv("EXAMPLE_API_KEY"))

	// grab the weather based on the zip code we want
	w, err := client.GetFromZip(*zipCode)
	if err != nil {
		panic(err.Error())
	}

	// print it out in a nice way
	b, err := json.MarshalIndent(w, "", "\t")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(b))
}
