package xclient

import (
	xoss "github.com/myxy99/component/xinvoker/oss"
	"github.com/myxy99/component/xinvoker/oss/standard"
)

var (
	disk     standard.Oss
	oss      standard.Oss
	cos      standard.Oss
	tfs      standard.Oss
	gfs      standard.Oss
	sevenNiu standard.Oss
)

func Disk() standard.Oss {
	if disk == nil {
		disk = xoss.Invoker("file")
	}
	return disk
}

func Oss() standard.Oss {
	if oss == nil {
		oss = xoss.Invoker("aliyun")
	}
	return oss
}

func Cos() standard.Oss {
	if cos == nil {
		cos = xoss.Invoker("cos")
	}
	return cos
}

func Tfs() standard.Oss {
	if tfs == nil {
		tfs = xoss.Invoker("tfs")
	}
	return tfs
}

func Gfs() standard.Oss {
	if gfs == nil {
		gfs = xoss.Invoker("gfs")
	}
	return gfs
}

func SevenNiu() standard.Oss {
	if sevenNiu == nil {
		sevenNiu = xoss.Invoker("sevenNiu")
	}
	return sevenNiu
}
