package main

import (
    "fmt"
    "log"
    "time"
    "bytes"
    "regexp"
    "strconv"
    "io/ioutil"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/pcapgo"
)


func captureToBuffer(req Capmsg, iface string)  {

    var (
        snapshotLen     int32  = 1500
        promiscuous     bool   = true
        err             error
        rerr            error
        timeout         time.Duration = 10 * time.Second
        handle          *pcap.Handle
        packetCount     int = 0
        fileName        string
        tagstr          string
        matchNode       bool = false
    )

    // Do sanity checking on max number of packets
    if(req.Packets == 0)  {
        fmt.Println("Invalid Capture size.  packets must be set to between 1 and " + strconv.Itoa(*maxpackets))
        log.Println("Invalid Capture size.  packets must be set to between 1 and " + strconv.Itoa(*maxpackets))
        return
    }
    if(req.Packets > config.Gen.Maxpackets)  {
        fmt.Println("Invalid Capture size.  packets cannot be > than maxpackets which is " + strconv.Itoa(config.Gen.Maxpackets))
        log.Println("Invalid Capture size.  packets cannot be > than maxpackets which is " + strconv.Itoa(config.Gen.Maxpackets))
        return
    }

    // Check the node against the message to see if we match either node or nodere 
    if(req.Node != "" && req.Nodere != "")  {
        fmt.Println("Invalid msg:  both node and nodere are set.  Use one or the other")
        log.Println("Invalid msg:  both node and nodere are set.  Use one or the other")
        return
    }

    if(req.Node == "" && req.Nodere == "")  {
        fmt.Println("Invalid msg:  both node and nodere are missing.  Use one or the other")
        log.Println("Invalid msg:  both node and nodere are missing.  Use one or the other")
        return
    }

    if(req.Node != "")  {
        if(req.Node == hostname)  {
            fmt.Println("Matched node: " + req.Node)
            log.Println("Matched node: " + req.Node)
            matchNode = true
        } else if(req.Node == "any")  {
            fmt.Println("Matched node: any")
            log.Println("Matched node: any")
            matchNode = true
        } 
    } else if(req.Nodere != "")  {
        matchNode, rerr = regexp.MatchString(req.Nodere, hostname)
        if(rerr != nil)  {
            fmt.Printf("Error applying regex:  %s\n", rerr)
            log.Printf("Error applying regex:  %s\n", rerr)
        } 

        if(matchNode)  {
            fmt.Println("Node regex match:  " + req.Nodere + " against " + hostname)
            log.Println("Node regex match:  " + req.Nodere + " against " + hostname)
        }
    }
        
    if(matchNode == false)  {
        fmt.Println("We didn't match via node or nodere " + hostname)
        log.Println("We didn't match via node or nodere " + hostname)
        return
    }
    // END OF NODE MATCHING

    if(req.Timeout != 0)  {
        timeout = req.Timeout * time.Second
    }

    // If snaplength is not overridden in the message then use the system default
    if(req.Snap == 0)  {
        req.Snap = config.Gen.Snap
    }

    log.Println("Capturing " + strconv.Itoa(req.Packets) + " packets on interface " + iface + " with a snaplength of " + strconv.Itoa(req.Snap))
    fmt.Println("Capturing " + strconv.Itoa(req.Packets) + " packets on interface " + iface + " with a snaplength of " + strconv.Itoa(req.Snap))

    fileName = hostname + "-" + iface + "-" + strconv.FormatInt(time.Now().Unix(), 10) + ".pcap"

    var f bytes.Buffer
    w := pcapgo.NewWriter(&f)
    w.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet)
    
    // Open the device for capturing
    handle, err = pcap.OpenLive(iface, int32(req.Snap), promiscuous, timeout)
    if err != nil {
        fmt.Printf("Error opening device %s: %v", iface, err)
        log.Printf("Error opening device %s: %v", iface, err)
    }
    if(req.Bpf != "")  {
        err := handle.SetBPFFilter(req.Bpf); 
        if(err != nil)  {
            fmt.Printf("Error compiling BPF Filter:[%s]  %s\n", req.Bpf, err)
            log.Printf("Error compiling BPF Filter:[%s]  %s\n", req.Bpf, err)
            return
        } else {
            fmt.Printf("Successfully compiled BPF Filter: [%s]\n", req.Bpf)
            log.Printf("Successfully compiled BPF Filter: [%s]\n", req.Bpf)
        } 
    }

    defer handle.Close()

    // Start processing packets
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
        // Process packet here
        // fmt.Println(packet)
        w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
        packetCount++
        
        // Only capture a fixed amount of packets
        if packetCount >= req.Packets {
            break
        }
    }


    // Handle Tags 
    // First time I'm touching tagstr which is why I don't check for empty
    if(req.Customer != "")  {
        tagstr = "customer:" + req.Customer
    }

    if(req.Alertid != 0 && tagstr == "")  {
        tagstr = "alertid:" + strconv.Itoa(req.Alertid)
    } else if(req.Alertid != 0 && tagstr != "") {
        tagstr = tagstr + ",alertid:" + strconv.Itoa(req.Alertid)
    }

    if(req.Tags != "" &&  tagstr == "")  {
        tagstr = req.Tags 
    } else if(req.Tags != "" &&  tagstr != "") {
        tagstr = tagstr + "," + req.Tags 
    }

    if(tagstr == "")  {
        tagstr = "node:" + hostname + ",interface:" + iface + ",snaplength:" + strconv.Itoa(req.Snap)
    } else {
        tagstr = tagstr + ",node:" + hostname + ",interface:" + iface + ",snaplength:" + strconv.Itoa(req.Snap)
    }

    if(config.Gen.Writelocal)  {
        ferr := ioutil.WriteFile(config.Gen.Localdir + "/" + fileName, f.Bytes(), 0644)
        if(ferr != nil)  {
            fmt.Printf("Error writing file: %s", ferr)
            log.Printf("Error writing file: %s", ferr)
        }
    }

    if(config.Cs.Upload)  {
        postBufferCloudshark(config.Cs.Scheme, config.Cs.Host, config.Cs.Port, config.Cs.Token, f, fileName, tagstr)
    }

	if(config.Aws.Upload)  {
		var msgfolder string
		var msgbucket string
		var msgacl string
		var msgregion string
		var msgep string

		if(req.Bucket != "")  {
			msgbucket = req.Bucket
		} else {
			msgbucket = *config.Aws.Bucket
		}
	
		if(req.Folder != "")  {
            msgfolder = req.Folder
        } else {
            msgfolder = *config.Aws.Folder
        }

		if(req.Acl != "")  {
            msgacl = req.Acl
        } else {
            msgacl = *config.Aws.Acl
        }

		if(req.Region != "")  {
            msgregion = req.Region
        } else {
            msgregion = *config.Aws.Region
        }
	
		if(req.Endpoint != "")  {
            msgep = req.Endpoint
        } else {
            msgep = *config.Aws.Endpoint
        }

		msgconfig := awsconfig
		msgconfig.Region = &msgregion	
        msgconfig.Endpoint = &msgep

		postS3(*msgconfig, msgbucket, f, fileName, tagstr, msgfolder, msgacl)
	}

    fmt.Println("Returning from capture")
    return
}
