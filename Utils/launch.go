package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

// Java 版本检测
func LaunchCheck() {
	log.Println("==> 依赖检测")
	launch := exec.Command("java", "--version")
	var stdout, stderr bytes.Buffer
	launch.Stdout = &stdout
	launch.Stderr = &stderr
	err := launch.Run()
	if err != nil {
		log.Println(err)
	}

	javaVersionOut, javaVersionErr := stdout.String(), stderr.String()
	if javaVersionErr != "" {
		log.Println(javaVersionErr)
	} else {
		log.Println("=> Java版本")
		fmt.Println(javaVersionOut)
		createScript()
	}

}

// 创建启动脚本
func createScript() {
	log.Print("=> 创建并写入启动脚本")
	launchSh, err := os.Create("./.gmcl/launch.sh")
	if err != nil {
		log.Println("-> 创建失败:", err)
	} else {
		log.Print("-> 创建成功")
	}
	defer launchSh.Close()

	os.Chmod("./.gmcl/launch.sh", 0777)

	_, errWrite := launchSh.WriteString(`# Create Date: ` + time.Now().Format("2006-01-02 15:04:05") + "\n" + ` echo "Hello world"`)
	if errWrite != nil {
		log.Println("-> 写入失败:", errWrite)
	} else {
		log.Println("-> 写入成功")
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
		log.Println(err)
	}

	launchOut, launchErr := stdout.String(), stderr.String()
	if launchErr != "" {
		log.Println(launchErr)
	} else {
		log.Println("-> 启动成功")
		fmt.Println(launchOut)
		go func() {
			log.Println("-> 自动任务: 关闭启动器") // 启动完成后关闭启动器
			os.Exit(0)
		}()
	}
}
