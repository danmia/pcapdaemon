package main

import (
  "flag"
  "fmt"
)

var cshostPtr *string
var cstokenPtr *string
var csschemePtr *string
var upPtr *bool

func main() {

    cshostPtr = flag.String("cshost", "localhost", "cloushark host")
    cstokenPtr = flag.String("cstoken", "xxxxxxx", "cloushark api token")
    csschemePtr = flag.String("csscheme", "https", "cloushark scheme http|https")
    redisnode := flag.String("redisnode", "127.0.0.1", "Hostname|IP of redis server.  Default localhost")
    redisport := flag.Int("redisport", 6379, "Port of redis server. Default 6379")
    redischannel := flag.String("redischannel", "capture", "Redis channel to subscribe to.  Default capture")

    upPtr = flag.Bool("upload", false, "Upload pcap")
    flag.Parse()

    // Channel for thread sync
    done := make(chan bool)

    go func()  {
        fmt.Println("Starting Redis Thread")
        subToRedis(*redisnode, *redisport, *redischannel)
        done <- true
    }()

    fmt.Println("Exiting?")
    <- done
}
