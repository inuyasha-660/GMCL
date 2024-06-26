package utils

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"
)

// Java 版本检测
func LaunchCheck() {
	slog.Info("依赖检测")
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
		slog.Info("Java版本")
		fmt.Println(javaVersionOut)
		createScript()
	}

}

// 创建启动脚本
func createScript() {
	slog.Info("创建并写入启动脚本")
	launchSh, err := os.Create("./.gmcl/launch.sh")
	if err != nil {
		slog.Error("创建失败:", err)
	} else {
		slog.Info("创建成功")
	}
	defer launchSh.Close()

	os.Chmod("./.gmcl/launch.sh", 0777)

	_, errWrite := launchSh.WriteString(`# Create Date: ` + time.Now().Format("2006-01-02 15:04:05") + "\n" + ` echo "Hello world"`)
	if errWrite != nil {
		slog.Error("写入失败:", errWrite)
	} else {
		slog.Info("写入成功")
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
		slog.Info("启动成功")
		fmt.Println(launchOut)
	}
}
