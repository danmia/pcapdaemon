package main

import (
	"time"
)

// Capmsg defines all the options that define a pcap capture request
type Capmsg struct {
	Node         string        `json:"node,omitempty"`
	Nodere       string        `json:"nodere,omitempty"`
	Interface    []string      `json:"interface,omitempty"`
	Alias        []string      `json:"alias,omitempty"`
	AliasMatched string        `json:"aliasmatched,omitempty"`
	Tags         string        `json:"tags,omitempty"`
	Bpf          string        `json:"bpf,omitempty"`
	Customer     string        `json:"customer,omitempty"`
	Snap         int           `json:"snap"`
	Packets      int           `json:"packets"`
	Alertid      int           `json:"alertid,omitempty"`
	Alertstr     int           `json:"alertstr,omitempty"`
	Timeout      time.Duration `json:"timeout,omitempty"`
	Duration     time.Duration `json:"duration,omitempty"`
	Bytes        int           `json:"bytes,omitempty"`
	PacketDebug  bool          `json:"packetdebug,omitempty"`
	LogRequest   bool          `json:"logrequest,omitempty"`
	Folder       string        `json:"folder,omitempty"`
	Bucket       string        `json:"bucket,omitempty"`
	ACL          string        `json:"acl,omitempty"`
	Region       string        `json:"region,omitempty"`
	Endpoint     string        `json:"endpoint,omitempty"`
	Encryption   bool          `json:"encryption,omitempty"`
}

// Capmsgs is simply an array of Capmsg
type Capmsgs []Capmsg

type tomlConfig struct {
	Gen    General          `toml:"general"`
	Aws    S3               `toml:"s3"`
	AwsSqs Sqs              `toml:"sqs"`
	Cs     Cloudshark       `toml:"cloudshark"`
	R      Redis            `toml:"redis"`
	K      Kafka            `toml:"kafka"`
	Ifmap  InterfaceAliases `toml:"interface"`
	Log    Syslog           `toml:"syslog"`
}

// General defines the top level "general" section of the the pcapdaemon config file
type General struct {
	Maxpackets  int           `toml:"maxpackets"`
	Maxbytes    int           `toml:"maxbytes"`
	Maxtimeout  time.Duration `toml:"maxtimeout"`
	Maxduration time.Duration `toml:"maxduration"`
	Deftimeout  time.Duration `toml:"defaulttimeout"`
	Writelocal  bool          `toml:"writelocal"`
	Localdir    string        `toml:"localdir"`
	Snap        int           `toml:"snaplength"`
	PacketDebug bool          `toml:"packetdebug"`
	LogRequests bool          `toml:"logrequests"`
}

// Cloudshark definese the cloudshark section of the config file
type Cloudshark struct {
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	Scheme  string `toml:"scheme"`
	Timeout int    `toml:"timeout"`
	Token   string `toml:"token"`
	Upload  bool   `toml:"upload"`
}

// InterfaceAlias defines an alias array that maps to a physical interface
type InterfaceAlias struct {
	Name  string   `toml:"name"`
	Alias []string `toml:"alias"`
}

// InterfaceAliases defines an array of InterfaceAlias
type InterfaceAliases []InterfaceAlias

// Redis defines the redis section of the config file
type Redis struct {
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	Channel string `toml:"channel"`
	Auth    string `toml:"auth"`
	Listen  bool   `toml:"listen"`
}

// Kafka defines the kafka section of the config file
type Kafka struct {
	Server []string `toml:"server"`
	Topic  string   `toml:"topic"`
	Listen bool     `toml:"listen"`
}

// Syslog defines the syslog section of the config file options
type Syslog struct {
	Priority int    `toml:"priority"`
	Tag      string `toml:"tag"`
}

// S3 defines the options necessary to use an S3 bucket as a destination for the pcap file
type S3 struct {
	AccessID   *string `toml:"accessid"`
	AccessKey  *string `toml:"accesskey"`
	Endpoint   *string `toml:"endpoint"`
	Region     *string `toml:"region"`
	Bucket     *string `toml:"bucket"`
	Folder     *string `toml:"pcaps"`
	Upload     bool    `toml:"upload"`
	ACL        *string `toml:"acl"`
	Encryption *bool   `toml:"encryption"`
}

// Sqs defines a struct that describes the options necessary to pull Capmsgs off of an Amazon SQS queue
type Sqs struct {
	AccessID    *string `toml:"accessid"`
	AccessKey   *string `toml:"accesskey"`
	Region      *string `toml:"region"`
	URL         *string `toml:"url"`
	Waitseconds *int64  `toml:"waitseconds"`
	Chunksize   *int64  `toml:"chunksize"`
	Listen      bool    `toml:"listen"`
}

// CsSuccess defines the object that is returned when a pcap is succuessfully posted to Cloudshark
type CsSuccess struct {
	Filename string `json:"filename,omitempty"`
	ID       string `json:"id,omitempty"`
}

// CsFail defines the object that is returned when an object fails to post to Cloudshark
type CsFail struct {
	Status     int      `json:"status,omitempty"`
	Exceptions []string `json:"exceptions,omitempty"`
}
