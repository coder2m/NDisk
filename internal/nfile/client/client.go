package xclient

import (
	xoss "github.com/myxy99/component/xinvoker/oss"
	"github.com/myxy99/component/xinvoker/oss/standard"
)

var (
	ossFile standard.Oss
	ossAli  standard.Oss
)

func OssFile() standard.Oss {
	if ossFile == nil {
		ossFile = xoss.Invoker("file")
	}
	return ossFile
}

func OssAli() standard.Oss {
	if ossAli == nil {
		ossAli = xoss.Invoker("aliyun")
	}
	return ossAli
}
