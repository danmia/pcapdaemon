package main

import (
    "os"
    "fmt"
    "log"
    "strconv"
	"strings"
	"time"
    "encoding/json"
    "github.com/garyburd/redigo/redis"
)

var c redis.Conn
var gerr error

func subToRedis(server string, port int, subchannel string, auth string) {

    fmt.Println("Attempting connect to Redis: " + server) 
    c, gerr = redis.Dial("tcp", server + ":" + strconv.Itoa(port))
    if gerr != nil {
        fmt.Printf("Error connecting to redis: %s\n", gerr)
        log.Printf("Error connecting to redis: %s\n", gerr)
        os.Exit(1) 
    }
    fmt.Println("Connected to Redis: " + server) 

	
	// Handle Redis Auth if applicable
	if(auth != "")  {
		_, err := c.Do("AUTH", "blah")
		if err != nil {
			// handle error
			log.Println("Failed redis auth: ", err)
			fmt.Println("Failed redis auth: ", err)
		}
	}

    psc := redis.PubSubConn{c}
    psc.Subscribe(subchannel)

    for {
        var msg Capmsg
        switch v := psc.Receive().(type) {
        case redis.Message:

            if(config.Gen.LogRequests)  {
                fmt.Printf("Redis Request: %s: message: %s\n", v.Channel, v.Data)
                log.Printf("Redis Request: %s: message: %s\n", v.Channel, v.Data)
            }

            if err := json.Unmarshal(v.Data, &msg); err != nil {
                fmt.Println("Redis: ", err)
                log.Println("Redis: ", err)
            } else {

                // set AliasMatched to empty to ensure nobody passes it in hence breaking things
                msg.AliasMatched = ""

                if(msg.LogRequest && ! config.Gen.LogRequests)  {
                    fmt.Printf("Redis Request: %s: message: %s\n", v.Channel, v.Data)
                    log.Printf("Redis Request: %s: message: %s\n", v.Channel, v.Data)
                }
 
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
        case redis.Subscription:
            fmt.Printf("Redis:  %s: %s %d\n", v.Channel, v.Kind, v.Count)
            log.Printf("Redis:  %s: %s %d\n", v.Channel, v.Kind, v.Count)
        case error:
            fmt.Printf("Redis Error: %s\n", v)
            log.Printf("Redis Error: %s\n", v)
			if(strings.Contains(v.Error(), "network"))  {
				fmt.Println("Redis:  We have a network issue")
				for {
					time.Sleep(time.Second * 3)
					c, gerr = redis.Dial("tcp", server + ":" + strconv.Itoa(port))
					if gerr != nil {
						fmt.Printf("Redis:  Error reconnecting: %s\n", gerr)
						log.Printf("Redis:  Error reconnecting to redis: %s\n", gerr)
					} else {
						fmt.Println("Redis:  Reconnected to " + server)
						log.Println("Redis:  Reconnected to " + server)

						// Handle Redis Auth if applicable
						if(auth != "")  {
							_, err := c.Do("AUTH", "blah")
							if err != nil {
								// handle error
								log.Println("Redis:  Failed redis auth on reconnect: ", err)
								fmt.Println("Redis:  Failed redis auth on reconnect: ", err)
							}
						}

						psc = redis.PubSubConn{c}
						psc.Subscribe(subchannel)
						break
					}

				}
			}
        }

    }
    
}
