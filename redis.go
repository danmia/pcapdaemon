package main

import (
    "os"
    "fmt"
    "log"
    "strconv"
    "encoding/json"
    "github.com/garyburd/redigo/redis"
)

var c redis.Conn
var gerr error

func subToRedis(server string, port int, subchannel string) {

    fmt.Println("Attempting connect to " + server) 
    c, gerr = redis.Dial("tcp", server + ":" + strconv.Itoa(port))
    if gerr != nil {
        fmt.Printf("Error connecting: %s\n", gerr)
        log.Printf("Error connecting to redis: %s\n", gerr)
        os.Exit(1) 
    }
    fmt.Println("Connected to " + server) 
    psc := redis.PubSubConn{c}
    psc.Subscribe(subchannel)

    for {
        var msg Capmsg
        switch v := psc.Receive().(type) {
        case redis.Message:
            fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
            if err := json.Unmarshal(v.Data, &msg); err != nil {
                fmt.Println(err)
            } else {
                if(len(msg.Interface) > 0)  {
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
                } else if(len(msg.Alias) > 0)  {
                    for _,v := range msg.Alias  {
                        if _, ok := almap[v]; ok  {
                            log.Println("Alias " + v + " exists in alias map")
                            fmt.Println("Alias " + v + " exists in alias map")
                            go captureToBuffer(msg, almap[v]);
                        } else {
                            log.Println("Alias " + v + " does not exist in alias map")
                            fmt.Println("Alias " + v + " does not exist in alias map")
                        }
                    }
                }
            }
        case redis.Subscription:
            fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
        case error:
            fmt.Printf("Error: %s\n", v)
        }

    }
    
}
