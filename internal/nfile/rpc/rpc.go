package rpc

import (
	"fmt"
	"github.com/myxy99/component/pkg/xnet"
	NFilePb "github.com/myxy99/ndisk/pkg/pb/nfile"
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

type Server struct{}

func (s Server) FileUpload(server NFilePb.NFileService_FileUploadServer) error {
	panic("implement me")
}

func (s Server) FileDownload(info *NFilePb.FileInfo, server NFilePb.NFileService_FileDownloadServer) error {
	panic("implement me")
}
