# pcapdaemon

## Description
This is a daemon that will subscribe to a redis pub/sub channel or amazon SQS queue for requests to capture.  It will capture and then optionally upload to Cloudshark, Amazon S3 or save to the local filesystem.  It could really be adapted to upload anywhere but the key was that I wanted to be able to trigger captures based on any number of events (traps, log events etc) via a lightweight mechanism.  A design goal was to have it capture into a buffer in memory and post the buffer without adding any kind of filesystem/io dependency.  

## Capture Controls
There are 4 controls that determine when/why the capture will exit and upload.  duration, timeout, packets and bytes.  The first of those parameters to hit wins out and the capture will complete and upload.  It's worth noting that the capture size will not be exact as that would require slicing packets into parts.  Every packet it looks at the number of bytes and if the previous packet pushed us over the limit (with pcap format overhead of course), then we complete and upload.

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
        "timeout": 15,
        "duration": 60,
        "bytes": 1000,
        "packetdebug":false,
		"folder": "myfolder",
		"bucket": "mybucket",
		"acl": "public-read",
		"region": "us-east-1",
		"endpoint": "s3.amazonaws.com",
		"encryption": false
    }
    
 * node - node name to capture on (exact match) Use either node or nodere but not both and one is required. Keyword: "any" matches any host
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
 * duration - Max amount of time to capture for
 * bytes - Max bytes to capture.  Note this will not be exact as that would require slicing a packet in half.
 * packetdebug - Print / log captured packet metadata.  True or false.  Defaults to false.
 * folder - S3 folder inside your bucket // S3 ONLY
 * bucket - S3 bucket // S3 ONLY
 * folder - S3 ACL // S3 ONLY
 * folder - S3 Endpoint // S3 ONLY
 * encryption - S3 Server side encryption AES256 // S3 ONLY
 * acl - S3 ACL // S3 ONLY
 * endpoint - S3 Endpoint // S3 ONLY

## Configuration File format (toml)
 * Defining interfaces is optional.  You only need to do it if you'd like to use an alias.  The basic use case is to group catpure interfaces across several nodes that may have different physical names for a variety of reasons.
 * Redis auth is optional as are ALL of the S3 options 
 * You must set listen to "true" for either Redis or SQS.  You can use both.
 * You must enable one of the following Cloudshark, S3 or writelocal
 * SQS support is implemented using long poll
 * SQS chunksize is the number of messages to process at once.  This could end up being the number of simultaneous captures so use with care
 * Given the time sensitive nature of capture messages, I recommend setting Default visibility timeout to 10 seconds and setting the message retention period to no more than 1 minute (these are in queue configuration in AWS SQS gui)
 * defaulttimeout is the default timeout to wait for packets during a capture.  IE.  If this amount of time passes between receiving packets, the capture will exit and upload unless number of packets is zero.
 * maxtimeout is the upper bound of the timeout explained above.  Since this can be overridden per message, an upper allowable bound seemed like a reasonable control to put in place.
 * maxduration is the maximum allowable duration you can have in a capture message.  This is to prevent someone doing somethhing awful.
 * maxbytes is the maximum bytes that can bet set in the capture message.  Again, an attempt to let the sysadmin protect the system from "bad" messages.
 * packetdebug is a boolean that enables logging and printing to STDOUT packet metadata that was captured.  Defaults to false. Use with care.
``` 
## Config file
[general]
maxpackets      = 50000
writelocal      = false
localdir        = "/tmp"
snaplength      = 500
defaulttimeout  = 10
maxtimeout      = 3600
maxduration     = 3600
maxbytes        = 100000000
packetdebug     = false

[cloudshark]
host        = "www.cloudshark.org"
scheme      = "https"
port        = 443
token       = "fffffffffffffffffff"
upload      = true

[redis]
listen		= true
host        = "node.running.redis.net"
port        = 6379
channel     = "capture"
auth		= "password"

# consult /usr/include/sys/syslog.h for Priority which is a combination 
# of facility and severity
[syslog]
priority    = 25
tag         = "pcapdaemon"

[[interface]]
name        = "eth0"
alias       = ["main", "public"]

[[interface]]
name        = "lo"
alias       = ["local"]

[s3]
accessid	= "xxxxxxxxxxxxxxxxxxxxxxxxx"
accesskey	= "xxxxxxxxxxxxxxxxxxxxxxx"
endpoint	= "s3.amazonaws.com"
bucket		= "pcapdaemon"
folder		= "pcaps"
upload		= true
region		= "us-east-1"
acl			= "private"
encryption	= false

[sqs]
listen		= true
region		= "us-east-1"
accessid	= "xxxxxxxxxxxxxxxxxxxxxxxxx"
accesskey	= "xxxxxxxxxxxxxxxxxxxxxxx"
url			= "https://example.amazone.com/asdfasdf/myqueue"
waitseconds = 20
chunksize	= 10
```
## Installing / Running 
```
#go get
#go install
#sudo daemonize -o stdout.log -e stderr.log /path/pcapdaemon -config /path/to/your/config
