package main

import (
    "flag"
    "fmt"
    "os"
    "log" 
)

var cshostPtr *string
var cstokenPtr *string
var csschemePtr *string
var upPtr *bool
var wLocal *bool
var destdir *string
var hostname string

func main() {

    cshostPtr = flag.String("cshost", "localhost", "cloudshark host")
    cstokenPtr = flag.String("cstoken", "xxxxxxx", "cloudshark api token")
    csschemePtr = flag.String("csscheme", "https", "cloudshark scheme http|https")
    redisnode := flag.String("redisnode", "127.0.0.1", "Hostname|IP of redis server.  Default localhost")
    redisport := flag.Int("redisport", 6379, "Port of redis server. Default 6379")
    redischannel := flag.String("redischannel", "capture", "Redis channel to subscribe to.  Default capture")

    upPtr = flag.Bool("upload", false, "Upload pcap")

    // flags for writing locally
    wLocal = flag.Bool("writelocal", false, "Write files locally.  Must set destdir")
    destdir = flag.String("destdir", "", "Destination directory locally for pcap files.  Requires -writelocal")

    // parse the flags
    flag.Parse()

    if(*wLocal)  {
        if _, err := os.Stat(*destdir); os.IsNotExist(err) {
            log.Fatal(*destdir + " does not exist");

        }  
    }

    hostname, _ = os.Hostname()
    // Channel for thread sync
    done := make(chan bool)

    go func()  {
        fmt.Println("Starting Redis Thread")
        subToRedis(*redisnode, *redisport, *redischannel)
        done <- true
    }()

    <- done
    fmt.Println("Exiting")
}
