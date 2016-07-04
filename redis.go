package main

import (
    "fmt"
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
                if(*upPtr)  {
                    go func() {
                        captureToBuffer(msg);
                    }()
                }
            }
        case redis.Subscription:
            fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
        case error:
            fmt.Printf("Error: %s\n", v)
        }

    }
    
}
