/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/5 19:20
 **/
package jwt

import (
	"github.com/myxy99/component/xcfg"
	"time"
)

const DefaultSecret = "kih**&hgyshq##js"

var cfg = xcfg.UnmarshalWithExpect("jwt", defaultConfig()).(*config)

type config struct {
	Secret string        `mapStructure:"secret"`
	Time   time.Duration `mapStructure:"time"`
	Issuer string        `mapStructure:"issuer"`
}

func defaultConfig() *config {
	return &config{
		Secret: DefaultSecret,
		Time:   72 * time.Hour,
		Issuer: "NDiskUser",
	}
}
