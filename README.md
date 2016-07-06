# pcapdaemon

## Description
This is a daemon that will subscribe to a redis pub/sub channel for requests to capture.  It will capture and then optionally upload to Cloudshark.  It could really be adapted to upload anywhere but the key was that I wanted to be able to trigger captures based on any number of events (traps, log events etc) via a lightweight mechanism.  A design goal was to have it capture into a buffer in memory and post the buffer without adding any kind of filesystem/io dependency.  That all being said, it also has the ability to write the pcap files locally to a configurable directory.

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
        "timeout": 15
    }
    
 * node - node name to capture on (exact match) Use either node or nodere but not both and one is required.
 * nodere - node regex to capture on.  Use either node or nodere but not both and one is required.
 * interface - interface to capture on.
 * tags - additional metadata for the capture file.  Comma separated list.
 * bpf - A filter string to capture on
 * customer - Additional metadata field ( not required )
 * snap - Snaplength or amount to capture into the packet.  Integer.
 * packets - Number of packets to capture.  Integer.
 * alertid - An integer ID for the event or alert or whatever you're tracking (not required)
 * timeout - Number of seconds to let capture last should the number of packets not get hit.  Integer
