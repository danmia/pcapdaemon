package main

import (
//    "os"
    "fmt"
    "log"
	"strings"
    "encoding/json"
    "github.com/Shopify/sarama"
)


func subToKafka(server []string, subchannel string) {

    kconfig := sarama.NewConfig()
    kconfig.Consumer.Return.Errors = true

    fmt.Println("Attempting connect to Kafka: " + strings.Join(server[:], ",")) 
    consumer, err := sarama.NewConsumer(server, kconfig)
    if err != nil {
        panic(err)
    }

    defer func() {
        if err := consumer.Close(); err != nil {
            log.Fatalln(err)
        }
    }()

    partitionConsumer, err := consumer.ConsumePartition(subchannel, 0, sarama.OffsetNewest)
    if err != nil {
        panic(err)
    }

    defer func() {
        if err := partitionConsumer.Close(); err != nil {
            log.Fatalln(err)
        }
    }()

    for {
        var msg Capmsg
        select  {
            case strmsg := <-partitionConsumer.Messages():

                if(config.Gen.LogRequests)  {
                    fmt.Printf("Kafka Request: %s: message: %s\n", subchannel, strmsg.Value)
                    log.Printf("Kafka Request: %s: message: %s\n", subchannel, strmsg.Value)
                }

                if err := json.Unmarshal(strmsg.Value, &msg); err != nil {
                    fmt.Println("Kafka: ", err)
                    log.Println("Kafka: ", err)
                } else {

                    // set AliasMatched to empty to ensure nobody passes it in hence breaking things
                    msg.AliasMatched = ""

                    if(msg.LogRequest && ! config.Gen.LogRequests)  {
                        fmt.Printf("Kafka Request: %s: message: %s\n", subchannel, strmsg.Value)
                        log.Printf("Kafka Request: %s: message: %s\n", subchannel, strmsg.Value)
                    }

                    if(len(msg.Interface) == 0 && len(msg.Alias) == 0) {
                        log.Println("Invalid msg:  both interface and alias are missing.  Use one or the other")
                        fmt.Println("Invalid msg:  both interface and alias are missing.  Use one or the other")
                    } else if(len(msg.Interface) > 0 && len(msg.Alias) > 0) {
                        log.Println("Invalid msg:  both interface and alias are set.  Use one or the other")
                        fmt.Println("Invalid msg:  both interface and alias are set.  Use one or the other")
                    } else if(len(msg.Interface) > 0)  {
                        for _, v := range msg.Interface  {
                            if _, ok := ifmap[v]; ok  {
                                log.Println("Interface " + v + " exists in interface map")
                                fmt.Println("Interface " + v + " exists in interface map")
                                go captureToBuffer(msg, v);
                            } else {
                                log.Println("Interface " + v + " does not exist in interface map")
                                fmt.Println("Interface " + v + " does not exist in interface map")
                            }            
                        }
                    } else {
                        for _,v := range msg.Alias  {
                            if _, ok := almap[v]; ok  {
                                for _, dname := range almap[v]  {
                                    log.Println("Alias " + v + " exists in alias map for device " + dname)
                                    fmt.Println("Alias " + v + " exists in alias map for device " + dname)
                                    msg.AliasMatched = v
                                    if _, ok := ifmap[dname]; ok  {
                                        go captureToBuffer(msg, dname);
                                    } else {
                                        log.Println("Alias " + v + " maps to interface " + dname + " which doesn't exist")
                                        fmt.Println("Alias " + v + " maps to interface " + dname + " which doesn't exist")
                                    } 
                                }
                            } else {
                                log.Println("Alias " + v + " does not exist in alias map")
                                fmt.Println("Alias " + v + " does not exist in alias map")
                            }
                        }
                    }
			    }
            case conserr := <-partitionConsumer.Errors():
                log.Println("Kafka Error Channel:  " + conserr.Topic + " Error: " + conserr.Err.Error())
                fmt.Println("Kafka Error Channel:  " + conserr.Topic + " Error: " + conserr.Err.Error())
        }
    }
    
}
