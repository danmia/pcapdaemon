package main

import (
    "flag"
    "os"
    "fmt"
    "log" 
    "log/syslog" 
    "github.com/google/gopacket/pcap"
    "github.com/BurntSushi/toml"
)

var cshostPtr *string
var cstokenPtr *string
var csschemePtr *string
var upPtr *bool
var wLocal *bool
var destdir *string
var hostname string
var maxpackets *int
var config tomlConfig

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

func validateOptions(c tomlConfig)  {

    // Redis configs
    if(c.R.Host == "")  {
        log.Fatal("You must supply a redis host")
    }
    if(c.R.Channel == "")  {
        log.Fatal("You must supply a redis channel to subscribe to")
    }
    if(c.R.Port == 0)  {
        c.R.Port = 6379
    }

    // Cloudshark configs
    if(c.Cs.Upload)  {
        if(c.Cs.Host == "")  {
            log.Fatal("If using Cloudshark, you must supply a host")
        }
        if(c.Cs.Token == "")  {
            log.Fatal("If using Cloudshark, you must supply an API token")
        }
        if(c.Cs.Scheme == "")  {
            c.Cs.Scheme = "https"
        }
        if(c.Cs.Port == 0)  {
            c.Cs.Port = 443
        }
    } else {
        fmt.Println("No cloudshark settings")
        log.Println("No cloudshark settings")
    }

    // General configs
    if(c.Gen.Maxpackets == 0)  {
        c.Gen.Maxpackets = 10000
    }
    if(c.Gen.Snap == 0)  {
        c.Gen.Snap = 512
    }

    if(c.Gen.Writelocal)  {
        if _, err := os.Stat(c.Gen.Localdir); os.IsNotExist(err) {
            log.Fatal(c.Gen.Localdir + " does not exist");
        }
    } 

    for _,v := range c.Ifmap  {
        if(v.Name == "")  {
            log.Fatal("Interface definition missing name property")
        }
        if(len(v.Alias) == 0)  {
            fmt.Println("Warning:  Interface [" + v.Name + "] has no aliases");
            log.Println("Warning:  Interface [" + v.Name + "] has no aliases");
        }
    }

    
}

func main() {

    cshostPtr = flag.String("cshost", "", "cloudshark host")
    cstokenPtr = flag.String("cstoken", "", "cloudshark api token")
    csschemePtr = flag.String("csscheme", "", "cloudshark scheme http|https")
    csportPtr := flag.Int("csport", 0, "cloudshark port")
    configfile := flag.String("config", "", "/path/to/configfile")
    redisnode := flag.String("redisnode", "", "Hostname|IP of redis server.  Default localhost")
    redisport := flag.Int("redisport", 0, "Port of redis server. Default 6379")
    redischannel := flag.String("redischannel", "", "Redis channel to subscribe to.  Default capture")
    maxpackets = flag.Int("maxpackets", 0, "Maximum number of packets per capture.  Default 50000")

    upPtr = flag.Bool("upload", false, "Upload pcap")

    // flags for writing locally
    wLocal = flag.Bool("writelocal", false, "Write files locally.  Must set destdir")
    destdir = flag.String("destdir", "", "Destination directory locally for pcap files.  Requires -writelocal")

    // parse the flags
    flag.Parse()

    if(*configfile != "")  {
        if _, err := os.Stat(*configfile); os.IsNotExist(err) {
            log.Fatal(*configfile + " does not exist");
        } else  {
            if _, err := toml.DecodeFile(*configfile, &config); err != nil {
                log.Fatal(err)
            }
        }
    }

    logwriter, e := syslog.New(syslog.LOG_NOTICE, "pcapdaemon")
    if e == nil {
        log.SetOutput(logwriter)
    }

    // create interface map
    updateInterfaceMap()

    // Overrides.  Command line parameters, if set, override the config file
    if(*redisnode != "")  {
        config.R.Host = *redisnode
    }
    if(*redisport != 0)  {
        config.R.Port = *redisport
    }
    if(*redischannel != "")  {
        config.R.Channel = *redischannel
    }

    if(*cshostPtr != "")  {
        config.Cs.Host = *cshostPtr    
    }
    if(*cstokenPtr != "")  {
        config.Cs.Token = *cstokenPtr
    }
    if(*csschemePtr != "")  {
        config.Cs.Scheme = *csschemePtr
    } 
    if(*csportPtr != 0)  {
        config.Cs.Port = *csportPtr
    }
    if(*upPtr)  {
        config.Cs.Upload = true
    }

    if(*maxpackets != 0)  {
        config.Gen.Maxpackets = *maxpackets
    }

    if(*wLocal)  {
        config.Gen.Writelocal = true
        if(*destdir != "")  {
            config.Gen.Localdir = *destdir
        }
    }

    // Run the validator AFTER defaults and config have been processed
    validateOptions(config)

    hostname, _ = os.Hostname()
    // Channel for thread sync
    done := make(chan bool)

    go func()  {
        log.Println("Starting Redis Thread")
        subToRedis(config.R.Host, config.R.Port, config.R.Channel)
        done <- true
    }()

    <- done
    log.Println("Exiting")
}
