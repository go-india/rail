# rail

rail is a Go client library for accessing the [RailwayAPI.com API](https://railwayapi.com/api).

## Usage

Construct a new Rail client, then use the various methods on the client to
access different parts of the RailwayAPI. For demonstration:

```go
package main

import (
  "context"
  "github.com/go-india/rail"
)

func main() {
  client := rail.NewClient(API_KEY)

  // Gets PNR status details.
  resp, err := client.PNRStatus(context.Background(), 2124289856)

  // Gets Live running status of Train.
  resp, err := client.LiveTrainStatus(context.Background(), 14311, time.Now())
}
```

`Note`: Using the [context](https://godoc.org/context) package for passing context.  
For complete usage of rail, see the full [api docs](https://railwayapi.com/api).

## Authentication

If you are using concrete Client, then you need to assign client.Auth field to
make the client methods use authenticator for requests.

```go
client := rail.Client{
  Auth: rail.NewAuth(API_KEY),
}
```

This will add API Key to each request made by client methods.

## Integration Tests

You can run integration tests from the directory.

```bash
$ go test -v
```

Use `-update` flag to run tests and update the testdata as well.  

`Note`: Define `RAILWAYAPI_TEST_API_KEY` in your environment for tests to use the API Key for testing.

## License

This library is distributed under the BSD-style license found in the [LICENSE](./LICENSE)
file.