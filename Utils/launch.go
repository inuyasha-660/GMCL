package utils

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"
)

// 依赖检测
func LaunchCheck(VersionChoose string) {
	if okVersion := CheckVersion(VersionChoose); okVersion {
		if okJava := CheckJava(); okJava {
			createScript()
		} else {
			slog.Error("Abort to launch with errors in checking Java")
		}
	} else {
		slog.Error("Abort to launch with errors in checking version")
	}

}

// 启动版本选择检测
func CheckVersion(VersionChoose string) bool {
	if VersionChoose == "" {
		slog.Error("Fail to read launch version: " + VersionChoose)
		return false
	} else {
		slog.Info("Start to launch, version: " + VersionChoose)
		return true
	}

}

// Java检测
func CheckJava() bool {
	launch := exec.Command("java", "--version")
	var stdout, stderr bytes.Buffer
	launch.Stdout = &stdout
	launch.Stderr = &stderr
	err := launch.Run()
	if err != nil {
		slog.Error(err.Error())
	}

	javaVersionOut, javaVersionErr := stdout.String(), stderr.String()
	if javaVersionErr != "" {
		slog.Info(javaVersionErr)
		return false
	} else {
		fmt.Println(javaVersionOut)
		return true
	}
}

// 创建启动脚本
func createScript() {
	slog.Info("Create and write script")
	launchSh, err := os.Create("./.gmcl/launch.sh")
	if err != nil {
		Glog("ERROR", "createScrip", "err", err)
	} else {
		slog.Info("Create Success")
	}
	defer launchSh.Close()

	os.Chmod("./.gmcl/launch.sh", 0777)

	_, errWrite := launchSh.WriteString(`# Create Date: ` + time.Now().Format("2006-01-02 15:04:05") + "\n" + ` echo "Hello world"`)
	if errWrite != nil {
		Glog("ERROR", "createScrip", "errWrite", errWrite)
	} else {
		slog.Info("write Success")
		launchGame()
	}
}

// 执行脚本
func launchGame() {
	launch := exec.Command("bash", "./.gmcl/launch.sh")
	var stdout, stderr bytes.Buffer
	launch.Stdout = &stdout
	launch.Stderr = &stderr
	err := launch.Run()
	if err != nil {
		slog.Error(err.Error())
	}

	launchOut, launchErr := stdout.String(), stderr.String()
	if launchErr != "" {
		slog.Error(launchErr)
	} else {
		slog.Info("Launch success")
		fmt.Println(launchOut)
	}
}
