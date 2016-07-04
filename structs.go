package main

import (
    "time"
)

type Capmsg struct  {
    Node            string          `json:"node,omitempty"`
    Interface       string          `json:"interface,omitempty"`
    Tags            string          `json:"tags,omitempty"`
    Bpf             string          `json:"bpf,omitempty"`
    Customer        string          `json:"customer,omitempty"`
    Snap            int             `json:"snap"`
    Packets         int             `json:"packets"`
    Alertid         int             `json:"alertid,omitempty"`
    Duration        time.Duration   `json:"duration,omitempty"`
}

type Capmsgs []Capmsg
