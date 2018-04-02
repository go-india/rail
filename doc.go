/*
Package rail provides a client for using RailwayAPI.com's API.

You can read the API server documentation at https://railwayapi.com/api

Usage

Construct a new rail client, then use the various methods on the client to
access different parts of the RailwayAPI.

For demonstration:

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

NOTE: Using the https://godoc.org/context package for passing context.

Authentication

If you are using concrete Client, then you need to assign client.Auth field to
make the client methods use authenticator for requests.

  client := rail.Client{
    Auth: rail.NewAuth(API_KEY),
  }

This will add API Key to each request made by client methods.
*/
package rail
