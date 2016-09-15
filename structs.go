package main

import (
    "time"
)

type Capmsg struct  {
    Node            string          `json:"node,omitempty"`
    Nodere          string          `json:"nodere,omitempty"`
    Interface       []string        `json:"interface,omitempty"`
    Alias           []string        `json:"alias,omitempty"`
    Tags            string          `json:"tags,omitempty"`
    Bpf             string          `json:"bpf,omitempty"`
    Customer        string          `json:"customer,omitempty"`
    Snap            int             `json:"snap"`
    Packets         int             `json:"packets"`
    Alertid         int             `json:"alertid,omitempty"`
    Alertstr        int             `json:"alertstr,omitempty"`
    Timeout         time.Duration   `json:"timeout,omitempty"`
    Folder			string          `json:"folder,omitempty"`
    Bucket			string          `json:"bucket,omitempty"`
	Acl				string			`json:"acl,omitempty"`	
	Region			string			`json:"region,omitempty"`
	Endpoint		string			`json:"endpoint,omitempty"`
	Encryption		bool			`json:"encryption,omitempty"`
}

type Capmsgs []Capmsg

type tomlConfig struct {
    Gen     General             `toml:"general"`
    Aws		S3					`toml:"s3"`
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
    Priority    int         `toml:"priority"`
    Tag         string      `toml:"tag"`
}

type S3 struct {
	AccessId	*string		`toml:"accessid"`
	AccessKey	*string		`toml:"accesskey"`
	Endpoint	*string		`toml:"endpoint"`
	Region		*string		`toml:"region"`
	Bucket		*string		`toml:"bucket"`
	Folder		*string		`toml:"pcaps"`
	Upload		bool		`toml:"upload"`
	Acl			*string		`toml:"acl"`
	Encryption	*bool		`toml:"encryption"`
}
