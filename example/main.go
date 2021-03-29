package main

import (
	"encoding/json"
	"fmt"

	"github.com/rtfreedman/weather"
)

func main() {
	w, err := weather.GetFromZip(20175)
	if err != nil {
		panic(err.Error())
	}
	b, err := json.MarshalIndent(w, "", "\t")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(b)
}
