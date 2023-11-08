# coleteonline

Go (golang) client library for accessing the [Colete Online API](https://docs.api.colete-online.ro/).

[![Go Reference](https://pkg.go.dev/badge/github.com/radulucut/coleteonline.svg)](https://pkg.go.dev/github.com/radulucut/coleteonline)
![Test](https://github.com/radulucut/coleteonline/actions/workflows/test.yml/badge.svg)

## Install

`go get github.com/radulucut/coleteonline`

## Endpoints

- [] /search/country/{needle}
- [] /search/location/{countryCode}/{needle}
- [] /search/city/{countryCode}/{county}/{needle}
- [] /search/street/{countryCode}/{city}/{county}/{needle}
- [] /search/postal-code/{countryCode}/{city}/{county}/{street}
- [] /search/validate-postal-code/{countryCode}/{city}/{county}/{street}/{postalCode}
- [] /search/postal-code-reverse/{countryCode}/{postalCode}
- [x] /address
- [x] /service/list
- [x] /order
- [x] /order/price
- [x] /order/status/{uniqueId}
- [] /order/awb/{uniqueId}
- [x] /user/balance

## Usage

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/radulucut/coleteonline"
)

func main() {
	client := coleteonline.NewClient(coleteonline.Config{
		ClientId:      "<ClientId>",
		ClientSecret:  "<ClientSecret>",
		UseProduction: true,
		Timeout:       10 * time.Second,
	})
	order := &coleteonline.Order{
		// ...
	}
	res, err := client.CreateOrder(order)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", res)
}
```
