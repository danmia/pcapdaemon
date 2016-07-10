# pcapdaemon

## Description
This is a daemon that will subscribe to a redis pub/sub channel for requests to capture.  It will capture and then optionally upload to Cloudshark or save to the local filesystem.  It could really be adapted to upload anywhere but the key was that I wanted to be able to trigger captures based on any number of events (traps, log events etc) via a lightweight mechanism.  A design goal was to have it capture into a buffer in memory and post the buffer without adding any kind of filesystem/io dependency.  

## Options
    -cshost string          cloudshark host (default "localhost")
    -csscheme string        cloudshark scheme http|https (default "https")
    -cstoken string         cloudshark api token (default "xxxxxxx")
    -csport int             cloudshark port
    -redischannel string    Redis channel to subscribe to.  Default capture (default "capture")
    -redisnode string       Hostname|IP of redis server.  Default localhost (default "127.0.0.1")
    -redisport int          Port of redis server. Default 6379 (default 6379)
    -upload                 Upload pcap
    -writelocal             Write pcap files locally.  Requires setting destdir
    -destdir                Directory to store locally written pcap files in
    -maxpackets             Maximum number of packets per capture.  Default 50000.
    -config string          Path to configuration file.  No default.
    
## Message format
    {
        "node": "node name",
        "nodere": "node regex",
        "interface": ["bond1","bond2"],
        "alias": ["local","public"],
        "tags": "blah,tagme,stuff",
        "bpf": "dst ip 10.0.0.1",
        "customer": "importantcustomer",
        "snap": 1500,
        "packets": 50,
        "alertid": 655443,
        "timeout": 15
    }
    
 * node - node name to capture on (exact match) Use either node or nodere but not both and one is required.
 * nodere - node regex to capture on.  Use either node or nodere but not both and one is required.
 * interface - an array of interfaces to dump on. (Must supply either interface or alias but not both)
 * alias - an array of interface aliases to dump on.  (Must supply either interface or alias but not both)
 * tags - additional metadata for the capture file.  Comma separated list.
 * bpf - A filter string to capture on
 * customer - Additional metadata field ( not required )
 * snap - Snaplength or amount to capture into the packet.  Integer.
 * packets - Number of packets to capture.  Integer.
 * alertid - An integer ID for the event or alert or whatever you're tracking (not required)
 * timeout - Number of seconds to let capture last should the number of packets not get hit.  Integer.

## Configuration File format (toml)
``` 
## Config file
[general]
maxpackets  = 50000
writelocal  = false
localdir    = "/tmp"
snaplength  = 500

[cloudshark]
host        = "cloudshark.org"
scheme      = "https"
port        = 443
token       = "fffffffffffffffffff"
upload      = true

[redis]
host        = "node.running.redis.net"
port        = 6379
channel     = "capture"

# consult /usr/include/sys/syslog.h for Priority which is a combination 
# of facility and severity
[syslog]
priority    = 85
tag         = pcapdaemon

[[interface]]
name        = "eth0"
alias       = ["main", "public"]

[[interface]]
name        = "lo"
alias       = ["local"]
```
## Installing / Running 
```
#go install
#sudo daemonize -o stdout.log -e stderr.log /path/pcapdaemon -config /path/to/your/config
