package main

import (
  "flag"
)

func main() {

    cshostPtr := flag.String("cshost", "localhost", "cloushark host")
    cstokenPtr := flag.String("cstoken", "xxxxxxx", "cloushark api token")
    csschemePtr := flag.String("csscheme", "https", "cloushark scheme http|https")
    upPtr := flag.Bool("upload", false, "Upload pcap")
    flag.Parse()

    if(*upPtr)  {
        pbuf := captureToBuffer();
        postBufferCloudshark(*csschemePtr, *cshostPtr, *cstokenPtr, pbuf) 
    }

    


}
