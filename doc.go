/*
Package rail provides a client for using RailwayAPI.com's API.

You can read the API server documentation at https://railwayapi.com/api

Usage

Construct a new Rail client, then use the various methods on the client to
access different parts of the RailwayAPI. For demonstration:

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

NOTE: Using the https://godoc.org/context package, one can easily
pass cancelation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then context.Background()
can be used as a starting point.

Authentication

If you are using concrete Client, then you need to assign client.Auth field to
make the client methods use authenticator for requests.

  client := rail.Client{
    Auth: rail.NewAuth(API_KEY),
  }

This will add API Key to each request made by client methods.
*/
package rail
