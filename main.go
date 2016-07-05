package main

import (
    "flag"
    "os"
    "fmt"
    "log" 
    "log/syslog" 
    "github.com/google/gopacket/pcap"
)

var cshostPtr *string
var cstokenPtr *string
var csschemePtr *string
var upPtr *bool
var wLocal *bool
var destdir *string
var hostname string
var ifmap = map[string]pcap.Interface{}

func updateInterfaceMap()  {
    
    x, ierr := pcap.FindAllDevs() 
    if(ierr != nil)  {
        fmt.Printf("Error loading interfaces: %s", ierr)
        log.Printf("Error loading interfaces: %s", ierr)
    }
    
    for _, v := range x  {
        fmt.Println("Found interface " + v.Name + " description: " + v.Description)
        log.Println("Found interface " + v.Name + " description: " + v.Description)
        ifmap[v.Name] = v
    }
}

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

    logwriter, e := syslog.New(syslog.LOG_NOTICE, "pcapdaemon")
    if e == nil {
        log.SetOutput(logwriter)
    }

    // create interface map
    updateInterfaceMap()

    if(*wLocal)  {
        if _, err := os.Stat(*destdir); os.IsNotExist(err) {
            log.Fatal(*destdir + " does not exist");

        }  
    }

    hostname, _ = os.Hostname()
    // Channel for thread sync
    done := make(chan bool)

    go func()  {
        log.Println("Starting Redis Thread")
        subToRedis(*redisnode, *redisport, *redischannel)
        done <- true
    }()

    <- done
    log.Println("Exiting")
}
