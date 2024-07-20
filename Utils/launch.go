package utils

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"
)

const SCRIPT_BASH = `#!/usr/bin/env bash`

var NATIVES_JAR_LINUX_X64 = map[string]string{
	`libglfw.so`:         `.minecraft/libraries/org/lwjgl/lwjgl-glfw/*/lwjgl-glfw-*-natives-linux.jar`,
	`libjemalloc.so`:     `.minecraft/libraries/org/lwjgl/lwjgl-jemalloc/*/lwjgl-jemalloc-*-natives-linux.jar`,
	`liblwjgl.so`:        `.minecraft/libraries/org/lwjgl/lwjgl/*/lwjgl-*-natives-linux.jar`,
	`liblwjgl_opengl.so`: `.minecraft/libraries/org/lwjgl/lwjgl-opengl/*/lwjgl-opengl-*-natives-linux.jar`,
	`liblwjgl_stb.so`:    `.minecraft/libraries/org/lwjgl/lwjgl-stb/*/lwjgl-stb-*-natives-linux.jar`,
	`libopenal.so`:       `.minecraft/libraries/org/lwjgl/lwjgl-openal/*/lwjgl-openal-*-natives-linux.jar`,
}

var NATIVES_JAR_MACOS = map[string]string{
	`libglfw.so`:         `.minecraft/libraries/org/lwjgl/lwjgl-glfw/*/lwjgl-glfw-*-natives-macos.jar`,
	`libjemalloc.so`:     `.minecraft/libraries/org/lwjgl/lwjgl-jemalloc/*/lwjgl-jemalloc-*-natives-macos.jar`,
	`liblwjgl.so`:        `.minecraft/libraries/org/lwjgl/lwjgl/*/lwjgl-*-natives-macos.jar`,
	`liblwjgl_opengl.so`: `.minecraft/libraries/org/lwjgl/lwjgl-opengl/*/lwjgl-opengl-*-natives-macos.jar`,
	`liblwjgl_stb.so`:    `.minecraft/libraries/org/lwjgl/lwjgl-stb/*/lwjgl-stb-*-natives-macos.jar`,
	`libopenal.so`:       `.minecraft/libraries/org/lwjgl/lwjgl-openal/*/lwjgl-openal-*-natives-macos.jar`,
}

var NATIVES_JAR_MACOS_ARM64 = map[string]string{
	`libglfw.so`:         `.minecraft/libraries/org/lwjgl/lwjgl-glfw/*/lwjgl-glfw-*-natives-macos-arm64.jar`,
	`libjemalloc.so`:     `.minecraft/libraries/org/lwjgl/lwjgl-jemalloc/*/lwjgl-jemalloc-*-natives-macos-arm64.jar`,
	`liblwjgl.so`:        `.minecraft/libraries/org/lwjgl/lwjgl/*/lwjgl-*-natives-macos-arm64.jar`,
	`liblwjgl_opengl.so`: `.minecraft/libraries/org/lwjgl/lwjgl-opengl/*/lwjgl-opengl-*-natives-macos-arm64.jar`,
	`liblwjgl_stb.so`:    `.minecraft/libraries/org/lwjgl/lwjgl-stb/*/lwjgl-stb-*-natives-macos-arm64.jar`,
	`libopenal.so`:       `.minecraft/libraries/org/lwjgl/lwjgl-openal/*/lwjgl-openal-*-natives-macos-arm64.jar`,
}

// 依赖检测
func LaunchCheck(VersionChoose string) {
	if okVersion := CheckVersion(VersionChoose); okVersion {
		if okJava := CheckJava(); okJava {
			UnzipJar(VersionChoose)
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

// 解压所需 natives 文件
func UnzipJar(version string) {
	system := runtime.GOOS
	arch := runtime.GOARCH

	dir := strings.TrimSuffix(version, path.Ext(version))
	MkdirPath := ".minecraft/versions/" + dir + "/natives-" + system + "-" + arch
	errMkdir := os.MkdirAll(MkdirPath, 0777)
	if errMkdir != nil {
		Glog("ERROR", "UnzipCmd", "errMkdir", errMkdir)
	}

	switch system {
	case "linux":
		{
			if arch == "amd64" {
				slog.Info("Create unzip natives script")
				script, errCreate := os.Create(".gmcl/launch-natives.sh")
				if errCreate != nil {
					Glog("ERROR", "createNativesScript", "errCreate", errCreate)
				}

				_, errWrite := script.WriteString(SCRIPT_BASH + "\n" + "")
				if errWrite != nil {
					Glog("ERROR", "createNativesScript", "errWrite", errWrite)
				}

				for lib, path := range NATIVES_JAR_LINUX_X64 {
					createNativesScript(lib, path, "linux", script)
				}
				_, err := script.WriteString("\n" + "mv linux/x64/org/lwjgl/*.so " + MkdirPath + "\n" +
					"mv linux/x64/org/lwjgl/*/*.so " + MkdirPath + "\n" + "rm -rf linux")
				if err != nil {
					Glog("ERROR", "UnzipJar", "err", err)
				}

				slog.Info("Run the unzip script")
				UnzipCmd()

				createScript()

			} else {
				slog.Error("Unsupported Architecture" + arch)
			}
		}
	case "darwin":
		{
			if arch == "arm64" {
				slog.Info("Create unzip natives script")
				script, errCreate := os.Create(".gmcl/launch-natives.sh")
				if errCreate != nil {
					Glog("ERROR", "createNativesScript", "errCreate", errCreate)
				}

				_, errWrite := script.WriteString(SCRIPT_BASH + "\n" + "")
				if errWrite != nil {
					Glog("ERROR", "createNativesScript", "errWrite", errWrite)
				}

				for lib, path := range NATIVES_JAR_MACOS_ARM64 {
					createNativesScript(lib, path, "macos", script)
				}
				_, err := script.WriteString("\n" + "mv macos/x64/org/lwjgl/*.so " + MkdirPath + "\n" +
					"mv macos/x64/org/lwjgl/*/*.so " + MkdirPath + "\n" + "rm -rf macos")
				if err != nil {
					Glog("ERROR", "UnzipJar", "err", err)
				}

				slog.Info("Run the unzip script")
				UnzipCmd()

				createScript()

			} else {
				slog.Info("Create unzip natives script")
				script, errCreate := os.Create(".gmcl/launch-natives.sh")
				if errCreate != nil {
					Glog("ERROR", "createNativesScript", "errCreate", errCreate)
				}

				_, errWrite := script.WriteString(SCRIPT_BASH + "\n" + "")
				if errWrite != nil {
					Glog("ERROR", "createNativesScript", "errWrite", errWrite)
				}

				for lib, path := range NATIVES_JAR_MACOS {
					createNativesScript(lib, path, "macos", script)
				}
				_, err := script.WriteString("\n" + "mv macos/arm64/org/lwjgl/*.so " + MkdirPath + "\n" +
					"mv macos/arm64/org/lwjgl/*/*.so " + MkdirPath + "\n" + "rm-rf macos")
				if err != nil {
					Glog("ERROR", "UnzipJar", "err", err)
				}

				slog.Info("Run the unzip script")
				UnzipCmd()

				createScript()

			}
		}
	default:
		{
			slog.Error("Unsupported Systems" + system)
		}
	}

}

func createNativesScript(lib, libPath, system string, script *os.File) {
	if lib == "liblwjgl.so" {
		_, err := script.WriteString("\n" + `/usr/bin/unzip -n ` + libPath + ` "` + system + `/*/*/*/*.so"`)
		if err != nil {
			Glog("ERROR", "createNativesScript-liblwjgl.so", "err", err)
		}
	} else {
		_, err := script.WriteString("\n" + `/usr/bin/unzip -n ` + libPath + ` "` + system + `/*/*/*/*/*.so"`)
		if err != nil {
			Glog("ERROR", "createNativesScript", "err", err)
		}
	}
}

// 执行解压命令
func UnzipCmd() {
	os.Chmod("./.gmcl/launch-natives.sh", 0777)
	unzip := exec.Command("bash", ".gmcl/launch-natives.sh")
	var stdout, stderr bytes.Buffer
	unzip.Stdout = &stdout
	unzip.Stderr = &stderr
	err := unzip.Run()
	if err != nil {
		slog.Error(err.Error())
	}

	launchOut, launchErr := stdout.String(), stderr.String()
	if launchErr != "" {
		slog.Error(launchErr)
	} else {
		fmt.Println(launchOut)
	}

}

// 创建启动脚本
func createScript() {
	slog.Info("Create and write launch script")
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
	launch := exec.Command("bash", ".gmcl/launch.sh")
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
