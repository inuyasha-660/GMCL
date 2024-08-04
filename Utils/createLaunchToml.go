package utils

import (
	"crypto/md5"
	"fmt"
	"os"
)

// JVM 参数
const TOML_JVM = `Xmx = "1024m"
Xmn = "128m"
UseG1GC = true
UseAdaptiveSizePolicy = true
OmitStackTraceInFastThrow = true
`

// Minecraft 参数
const TOML_MINECRAFT = `
Width = "1800"
Height = "1000"
`

func CreateLaunchToml(userName string) bool {
	toml, err := os.Create("./.gmcl/launch.toml")
	if err != nil {
		Glog("ERROR", "CreateLaunchToml", "err", err)
		return false
	} else {
		os.Chmod("./.gmcl/launch.toml", 0777)
		userNameMd5 := Md5Create(userName)
		uuid := userNameMd5[:8] + "-" + userNameMd5[8:12] + "-" + userNameMd5[12:16] + "-" + userNameMd5[16:20] + "-" + userNameMd5[20:]
		_, errWrite := toml.WriteString(TOML_JVM + TOML_MINECRAFT + "UUID =" + `"` + uuid + `"`)
		if errWrite != nil {
			Glog("ERROR", "CreateLaunchToml", "errWrite", errWrite)
			return false
		} else {
			return true
		}
	}

}

func Md5Create(userName string) string {
	nameByte := []byte(userName)
	md5 := md5.Sum(nameByte)
	result := fmt.Sprintf("%x", md5)
	return result
}
