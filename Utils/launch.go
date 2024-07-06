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
func LaunchCheck() {
	slog.Info("Check Java")
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
	} else {
		fmt.Println(javaVersionOut)
		createScript()
	}

}

// 创建启动脚本
func createScript() {
	slog.Info("Create and write script")
	launchSh, err := os.Create("./.gmcl/launch.sh")
	if err != nil {
		Glog("ERROR", "createScrip", "err", err)
	} else {
		slog.Info("Success")
	}
	defer launchSh.Close()

	os.Chmod("./.gmcl/launch.sh", 0777)

	_, errWrite := launchSh.WriteString(`# Create Date: ` + time.Now().Format("2006-01-02 15:04:05") + "\n" + ` echo "Hello world"`)
	if errWrite != nil {
		Glog("ERROR", "createScrip", "errWrite", errWrite)
	} else {
		slog.Info("Success")
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
