package utils

import (
	"fmt"
	"time"
)

func Glog(level string, function string, varName string, err error) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), level, function, "->", varName, ":", err)
}

func GlogINFO(info string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "INFO", info)
}
