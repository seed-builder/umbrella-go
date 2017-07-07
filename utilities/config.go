package utilities

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	HttpPort string `yaml:"http_port"`
	TcpPort string `yaml:"tcp_port"`
	TcpTestTimeout int32 `yaml:"tcp_test_timeout"`
	TcpTestMax int32 `yaml:"tcp_test_max"`
	Debug bool `yaml:"debug"`
	DbDialect string `yaml:"db_dialect"`
	DbServer string `yaml:"db_server"`
	DbUser string `yaml:"db_user"`
	DbPassword string `yaml:"db_password"`
	DbDatabase string `yaml:"db_database"`
	Salt string `yaml:"salt"`
}

var SysConfig Config

func init() {
	path := "./config.yml"
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err1 := yaml.Unmarshal(configFile, &SysConfig)
	if err1 != nil {
		log.Fatalf("error: %v", err1)
	}
}