# rail [![GoDoc](https://godoc.org/github.com/go-india/rail?status.svg)](https://godoc.org/github.com/go-india/rail) [![Build Status](https://travis-ci.org/go-india/rail.svg?branch=master)](https://travis-ci.org/go-india/rail) [![Coverage Status](https://coveralls.io/repos/github/go-india/rail/badge.svg?branch=master)](https://coveralls.io/github/go-india/rail?branch=master) [![Report card](https://goreportcard.com/badge/github.com/go-india/rail)](https://goreportcard.com/report/github.com/go-india/rail)

rail is a [Go](http://golang.org/) client library for accessing the [railwayapi.com API](https://railwayapi.com/api).

> <img src="https://railwayapi.com/api/images/logo.png" width="196">  

> RailwayAPI.com provides API for Indian Railways data of Trains and Stations, along with features like Train Live status, PNR status, Arrivals on Station, Trains Between Stations etc

### Installation

```bash
$ go get -u github.com/go-india/rail
```

### Usage

Construct a new rail client, then use the various methods on the client to access different parts of the RailwayAPI.

For demonstration:

```go
package main

import (
  "context"
  "github.com/go-india/rail"
)

var ctx = context.Background()

func main() {
  client := rail.NewClient(API_KEY)

  // Gets PNR status details.
  resp, err := client.PNRStatus(ctx, 2124289856)

  // Gets Live running status of Train.
  resp, err := client.LiveTrainStatus(ctx, 14311, time.Now())

  // Gets fares of train.
  resp, err := c.TrainFare(ctx, 14311, "BE", "ADI", 24, "SL", "GN", time.Now())
}
```

`Notes`
* Using the [context](https://godoc.org/context) package for passing context.  
* Make sure you have a valid API Key. If not, you can get a new one by registering at railwayapi.com [registration page](https://railwayapi.com/register).

For complete usage of rail, see the full [package docs](https://godoc.org/github.com/go-india/rail).

#### Authentication

If you are using concrete Client, then you need to assign `client.Auth` field to make the client methods use authenticator for requests.

```go
client := rail.Client{
  Auth: rail.NewAuth(API_KEY),
}
```

This will add API Key to each request made by client methods.

#### Integration Tests

You can run integration tests from the directory.

```bash
$ go test -v
```

`Note`: Use `-update` flag to update the testdata. When using update flag, you will need to define `RAILWAYAPI_TEST_API_KEY` in your environment for tests to use the API Key for testing.

### Contributing

We welcome pull requests, bug fixes and issue reports. Before proposing a change, please discuss your change by raising an issue.

### License

This library is distributed under the MIT license found in the [LICENSE](./LICENSE) file.

### Author

[Yash Raj Singh](http://yashrajsingh.net/)