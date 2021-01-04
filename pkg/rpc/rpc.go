package rpc

import (
	"fmt"
	"github.com/myxy99/component/pkg/xnet"
	"time"
)

type Config struct {
	EtcdAddr         string        `mapStructure:"etcd_addr"`
	ServerIp         string        `mapStructure:"ip"`
	ServerPort       int           `mapStructure:"port"`
	RegisterTTL      time.Duration `mapStructure:"register_ttl"`
	RegisterInterval time.Duration `mapStructure:"register_interval"`
	Timeout          time.Duration `mapStructure:"timeout"`
}

func DefaultConfig() *Config {
	host, port, err := xnet.GetLocalMainIP()
	if err != nil {
		host = "localhost"
	}
	return &Config{
		EtcdAddr:         "127.0.0.1:2379",
		ServerIp:         host,
		ServerPort:       port,
		RegisterTTL:      30 * time.Second,
		RegisterInterval: 15 * time.Second,
		Timeout:          30 * time.Second,
	}
}

func (c Config) Addr() string {
	return fmt.Sprintf("%v:%v", c.ServerIp, c.ServerPort)
}
