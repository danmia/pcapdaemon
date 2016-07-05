# pcapdaemon

## Description
This is a daemon that will subscribe to a redis pub/sub channel for requests to capture.  It will capture and then optionally upload to Cloudshark.  It could really be adapted to upload anywhere but the key was that I wanted to be able to trigger captures based on any number of events (traps, log events etc) via a lightweight mechanism.  A design goal was to have it capture into a buffer in memory and post the buffer without adding any kind of filesystem/io dependency.

## Options
    -cshost string          cloudshark host (default "localhost")
    -csscheme string        cloudshark scheme http|https (default "https")
    -cstoken string         cloudshark api token (default "xxxxxxx")
    -redischannel string    Redis channel to subscribe to.  Default capture (default "capture")
    -redisnode string       Hostname|IP of redis server.  Default localhost (default "127.0.0.1")
    -redisport int          Port of redis server. Default 6379 (default 6379)
    -upload                 Upload pcap
    -writelocal             Write pcap files locally.  Requires setting destdir
    -destdir                Directory to store locally written pcap files in
    
## Message format
    {
        "node": "node name",
        "nodere": "node regex",
        "interface": "bond1",
        "tags": "blah,tagme,stuff",
        "bpf": "dst ip 10.0.0.1",
        "customer": "importantcustomer",
        "snap": 1500,
        "packets": 50,
        "alertid": 655443,
        "duration": 15
    }
