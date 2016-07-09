package main

import (
    "time"
)

type Capmsg struct  {
    Node            string          `json:"node,omitempty"`
    Nodere          string          `json:"nodere,omitempty"`
    Interface       string          `json:"interface,omitempty"`
    Interalias      string          `json:"interalias,omitempty"`
    Tags            string          `json:"tags,omitempty"`
    Bpf             string          `json:"bpf,omitempty"`
    Customer        string          `json:"customer,omitempty"`
    Snap            int             `json:"snap"`
    Packets         int             `json:"packets"`
    Alertid         int             `json:"alertid,omitempty"`
    Alertstr        int             `json:"alertstr,omitempty"`
    Timeout         time.Duration   `json:"timeout,omitempty"`
}

type Capmsgs []Capmsg

type tomlConfig struct {
    Gen     General             `toml:"general"`
    Cs      Cloudshark          `toml:"cloudshark"`
    R       Redis               `toml:"redis"`
    Ifmap   InterfaceAliases    `toml:"interface"`
    Log     Syslog              `toml:"syslog"`
}

type General struct  {
    Maxpackets      int         `toml:"maxpackets"`
    Writelocal      bool        `toml:"writelocal"`
    Localdir        string      `toml:"localdir"`
    Snap            int         `toml:"snaplength"`
}

type Cloudshark struct {
    Host        string          `toml:"host"`
    Port        int             `toml:"port"`
    Scheme      string          `toml:"scheme"`
    Token       string          `toml:"token"`
    Upload      bool            `toml:"upload"`
}

type InterfaceAlias struct {
    Name        string          `toml:"name"`
    Alias       []string        `toml:"alias"`
}

type InterfaceAliases []InterfaceAlias

type Redis struct  {
    Host        string      `toml:"host"`
    Port        int         `toml:"port"`
    Channel     string      `toml:"channel"`
}

type Syslog struct {
    Priority    string      `toml:"priority"`
    Tag         string      `toml:"tag"`
}
