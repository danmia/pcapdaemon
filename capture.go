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


func captureToBuffer(req Capmsg)  {

    var (
        snapshotLen     int32  = 1500
        promiscuous     bool   = true
        err             error
        rerr            error
        timeout         time.Duration = 10 * time.Second
        handle          *pcap.Handle
        packetCount     int = 0
        fileName        string
        matchNode       bool = false
    )

    if _, ok := ifmap[req.Interface]; ok  {
        log.Println("Interface " + req.Interface + " exists in interface map")
        fmt.Println("Interface " + req.Interface + " exists in interface map")
    } else {
        log.Println("Interface " + req.Interface + " does not exist in interface map")
        fmt.Println("Interface " + req.Interface + " does not exist in interface map")
        return;
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

    log.Println("Capturing " + strconv.Itoa(req.Packets) + " packets on interface " + req.Interface + " with a snaplength of " + strconv.Itoa(req.Snap))
    fmt.Println("Capturing " + strconv.Itoa(req.Packets) + " packets on interface " + req.Interface + " with a snaplength of " + strconv.Itoa(req.Snap))

    fileName = hostname + "-" + req.Interface + "-" + strconv.FormatInt(time.Now().Unix(), 10) + ".pcap"

    var f bytes.Buffer
    w := pcapgo.NewWriter(&f)
    w.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet)
    
    // Open the device for capturing
    handle, err = pcap.OpenLive(req.Interface, int32(req.Snap), promiscuous, timeout)
    if err != nil {
        fmt.Printf("Error opening device %s: %v", req.Interface, err)
        log.Printf("Error opening device %s: %v", req.Interface, err)
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

    fmt.Println("Returning from capture")

    if(*wLocal)  {
        ferr := ioutil.WriteFile(*destdir + "/" + fileName, f.Bytes(), 0644)
        if(ferr != nil)  {
            fmt.Printf("Error writing file: %s", ferr)
            log.Printf("Error writing file: %s", ferr)
        }
    }

    if(*upPtr)  {
        postBufferCloudshark(*csschemePtr, *cshostPtr, *cstokenPtr, f, fileName)
    }

    return
}
