package main

import (
    "fmt"
    "os"
    "time"
    "bytes"

    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/pcapgo"
)


func captureToBuffer(req Capmsg)  {

    var (
        deviceName  string = "eth0"
        snapshotLen int32  = 1500
        promiscuous bool   = false
        err         error
        // timeout     time.Duration = -1 * time.Second
        timeout     time.Duration = 10 * time.Second
        handle      *pcap.Handle
        packetCount int = 0
    )

    fmt.Println("Capturing on interface: " + req.Interface)

    var f bytes.Buffer
    w := pcapgo.NewWriter(&f)
    w.WriteFileHeader(uint32(snapshotLen), layers.LinkTypeEthernet)

    // Open the device for capturing
    handle, err = pcap.OpenLive(deviceName, snapshotLen, promiscuous, timeout)
    if err != nil {
        fmt.Printf("Error opening device %s: %v", deviceName, err)
        os.Exit(1)
    }
    defer handle.Close()

    // Start processing packets
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
        // Process packet here
        fmt.Println(packet)
        w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
        packetCount++
        
        // Only capture 100 and then stop
        if packetCount >= 100 {
            break
        }
    }

    fmt.Println("Returning from capture")
    postBufferCloudshark(*csschemePtr, *cshostPtr, *cstokenPtr, f)
    return
}
