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
        timeout         time.Duration = 3 * time.Second
        handle          *pcap.Handle
        packetCount     int = 0
        fileName        string
        tagstr          string
        matchNode       bool = false
        captimeout      time.Duration
        capduration     time.Duration
        capbytes        int
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


    // Timeout management is to break out of a capture if long periods of time pass
    // without capturing any packets
    if(req.Timeout != 0)  {
        if(req.Timeout > config.Gen.Maxtimeout)  {
            log.Printf("Error:  Max timeout %d is greater than max allowable timeout %d\n", req.Timeout, config.Gen.Maxtimeout)
            fmt.Printf("Error:  Max timeout %d is greater than max allowable timeout %d\n", req.Timeout, config.Gen.Maxtimeout)
            return
        }
        captimeout = req.Timeout * time.Second
    } else {
        captimeout = config.Gen.Deftimeout * time.Second
    }

    // Duration managment is to put a cap on how long to capture for no matter what is going on
    if(req.Duration != 0)  {
        if(req.Duration > config.Gen.Maxduration)  {
            log.Printf("Error:  Max duration %d is greater than max allowable duration %d\n", req.Duration, config.Gen.Maxduration)
            fmt.Printf("Error:  Max duration %d is greater than max allowable duration %d\n", req.Duration, config.Gen.Maxduration)
            return
        }
        capduration = req.Duration * time.Second
    } else {
        capduration = config.Gen.Maxduration * time.Second
    }
    
    // Byte management is to break out after a certain number of bytes
    if(req.Bytes != 0)  {
        if(req.Bytes > config.Gen.Maxbytes)  {
            log.Printf("Error:  message bytes %d is greater than max allowable bytes %d\n", req.Bytes, config.Gen.Maxbytes)
            fmt.Printf("Error:  message bytes %d is greater than max allowable bytes %d\n", req.Bytes, config.Gen.Maxbytes)
            return
        }
        capbytes = req.Bytes
    } else {
        capbytes = config.Gen.Maxbytes
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
    packetSource.DecodeOptions = gopacket.DecodeOptions{Lazy: false, NoCopy: false, SkipDecodeRecovery: true}

    packetchan := packetSource.Packets()
    captimer := time.NewTimer(capduration)

    C:
    for  {
        select  {
            case packet := <-packetchan:
                // Process packet here

                // Global packet debug
                if(config.Gen.PacketDebug || req.PacketDebug)  {
                    fmt.Println(packet)
                    log.Println(packet)
                }

                w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
                packetCount++
            
                // Only capture a fixed amount of packets
                if packetCount >= req.Packets {
                    fmt.Printf("Packet count %d hit for capture %s, size: %d bytes\n", req.Packets, fileName, f.Len())
                    log.Printf("Packet count %d hit for capture %s, size: %d bytes\n", req.Packets, fileName, f.Len())
                    break C
                }

                // Only capture a fixed amount of packets
                if capbytes <= f.Len() {
                    fmt.Printf("Size limit %d bytes hit for capture %s, size: %d bytes\n", capbytes, fileName, f.Len())
                    log.Printf("Size limit %d bytes hit for capture %s, size: %d bytes\n", capbytes, fileName, f.Len())
                    break C
                }

            case <-time.After(captimeout):
                fmt.Printf("Packet timeout %s hit for capture %s, captured %d packets, size: %d\n", captimeout.String(), fileName, packetCount, f.Len())
                log.Printf("Packet timeout %s hit for capture %s, captured %d packets, size: %d\n", captimeout.String(), fileName, packetCount, f.Len())

                // If there are no packets before the timeout then return without uploading
                if(packetCount == 0)  {
                    log.Printf("Packet timeout %s hit for capture %s and packet count is 0 so returning without uploading\n", captimeout.String(), fileName)
                    fmt.Printf("Packet timeout %s hit for capture %s and packet count is 0 so returning without uploading\n", captimeout.String(), fileName)
                    return
                } 
                break C

            case <- captimer.C:
                fmt.Printf("Capture duration %s hit for capture %s, captured %d packets, size: %d\n", capduration.String(), fileName, packetCount, f.Len())
                log.Printf("Capture duration %s hit for capture %s, captured %d packets, size: %d\n", capduration.String(), fileName, packetCount, f.Len())
                
                // If there are no packets before the total duration hits then return without uploading
                if(packetCount == 0)  {
                    log.Printf("Capture duration %s hit for capture %s and packet count is 0 so returning without uploading\n", capduration.String(), fileName)
                    fmt.Printf("Capture duration %s hit for capture %s and packet count is 0 so returning without uploading\n", capduration.String(), fileName)
                    return
                }
                break C
    
        }
    }

    // Handle Tags 
    // First tim I'm touching tagstr which is why I don't check for empty
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
        } else {
            log.Printf("Written locally file: %s", config.Gen.Localdir + "/" + fileName)
            fmt.Printf("Written locally file: %s", config.Gen.Localdir + "/" + fileName)
        }

    }

    if(config.Cs.Upload)  {
        log.Printf("Uploading to Cloudshark file: %s\n", fileName)
        fmt.Printf("Uploading to Cloudshark file: %s\n", fileName)
        postBufferCloudshark(config.Cs.Scheme, config.Cs.Host, config.Cs.Port, config.Cs.Token, f, fileName, tagstr)
    }

	if(config.Aws.Upload)  {
        log.Printf("Uploading to S3 file: %s\n", fileName)
        fmt.Printf("Uploading to S3 file: %s\n", fileName)

		var msgfolder string
		var msgbucket string
		var msgacl string
		var msgregion string
		var msgep string
		var msgenc bool

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

		fmt.Println("Enc: ", req.Encryption)
		if(req.Encryption)  {
            msgenc = req.Encryption
        } else {
            msgenc = *config.Aws.Encryption
        }

		msgconfig := awsconfig
		msgconfig.Region = &msgregion	
        msgconfig.Endpoint = &msgep

        log.Printf("Uploading to S3 file: %s", msgbucket + ":" + msgfolder + "/" + fileName)
        fmt.Printf("Uploading to S3 file: %s", msgbucket + ":" + msgfolder + "/" + fileName)

		postS3(*msgconfig, msgbucket, f, fileName, tagstr, msgfolder, msgacl, msgenc)
	}

    fmt.Println("Returning from capture")
    return
}
