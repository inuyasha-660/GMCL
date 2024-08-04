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

	"github.com/pelletier/go-toml/v2"
	"github.com/tidwall/gjson"
)

const SCRIPT_BASH = `#!/usr/bin/env bash`

const LAUNCH_TOML_DEFALUT_JVM = `java -Xmx1024m -Xmn128m -XX:+UseG1GC -XX:-UseAdaptiveSizePolicy -XX:-OmitStackTraceInFastThrow`

const LAUNCH_TOML_LAUNCH_NAME = ` -Dminecraft.launcher.brand=GMCL -Dminecraft.launcher.version=0.9.0 `

const LAUNCH_TOML_LOG_FILE = ` -Dlog4j.configurationFile=.minecraft/versions/`

const MINECRAFT_NO_MOD = ` net.minecraft.client.main.Main `

const MINECRAFT_WITH_MOD = `net.minecraft.launchwrapper.Launch `

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

type LauncConfig struct {
	Xmx                       string
	Xmn                       string
	UseG1GC                   bool
	UseAdaptiveSizePolicy     bool
	OmitStackTraceInFastThrow bool

	Width  string
	Height string
	UUID   string
}

// 依赖检测
func LaunchCheck(VersionChoose, userName string) {
	if okVersion := CheckVersion(VersionChoose); okVersion {
		if okJava := CheckJava(); okJava {
			dir := strings.TrimSuffix(VersionChoose, path.Ext(VersionChoose))
			MkdirPath := ".minecraft/versions/" + dir + "/natives-" + runtime.GOOS + "-" + runtime.GOARCH
			if ok := CheckIfExist(MkdirPath); ok {
				readLaunchToml(dir, MkdirPath, userName)
			} else {
				UnzipJar(VersionChoose, userName)
			}
		} else {
			slog.Error("Abort to launch with errors in checking Java")
		}
	} else {
		slog.Error("Abort to launch with errors in checking version")
	}

}

// TODO: 检测 Natives库文件是否存在，存在则跳过解压否则进行解压

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
func UnzipJar(version, userName string) {
	system := runtime.GOOS
	arch := runtime.GOARCH

	dir := strings.TrimSuffix(version, path.Ext(version))                         // 去除拓展名的游戏Jar文件，类似: 1.21(1.21.jar)
	MkdirPath := ".minecraft/versions/" + dir + "/natives-" + system + "-" + arch // Natives目录
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

				readLaunchToml(dir, MkdirPath, userName)

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

				readLaunchToml(dir, MkdirPath, userName)

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

				readLaunchToml(dir, MkdirPath, userName)

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

// 读取用户配置，若不存在则使用默认配置
func readLaunchToml(dir, MkdirPath, userName string) {
	slog.Info("Create/Read Configuration")
	var launchToml string

	if ok := CheckIfExist("./.gmcl/launch.toml"); ok {
		slog.Info("Read configuration from ./.gmcl/launch.toml")
		var launchConfig LauncConfig
		tomlRead, errRead := os.ReadFile("./.gmcl/launch.toml")
		if errRead != nil {
			Glog("ERROR", "readLaunchToml", "errRead", errRead)
		}

		errUnmarshal := toml.Unmarshal([]byte(tomlRead), &launchConfig)
		if errUnmarshal != nil {
			Glog("ERROR", "readLaunchToml", "errUnmarshal", errUnmarshal)
		}

		cpArtifact := readArtifact(dir)
		var UseG1GC string
		var UseAdaptiveSizePolicy string
		var OmitStackTraceInFastThrow string

		if launchConfig.UseG1GC {
			UseG1GC = "-XX:+UseG1GC"
		}

		if launchConfig.UseAdaptiveSizePolicy {
			UseAdaptiveSizePolicy = "-XX:-UseAdaptiveSizePolicy"
		}

		if launchConfig.OmitStackTraceInFastThrow {
			OmitStackTraceInFastThrow = "-XX:-OmitStackTraceInFastThrow"
		}

		launchToml = "java -Xmx" + launchConfig.Xmx + " -Xmn" + launchConfig.Xmn + " " + UseG1GC + " " + UseAdaptiveSizePolicy + " " + OmitStackTraceInFastThrow +
			LAUNCH_TOML_LAUNCH_NAME + LAUNCH_TOML_LOG_FILE + dir + "/client-1.12.xml" + " -Djava.library.path=" + MkdirPath + " -cp " + cpArtifact + ".minecraft/versions/" + dir + "/" + dir + ".jar " + MINECRAFT_NO_MOD + " --username " + userName + " --version " + dir + " --gameDir .minecraft/versions/" +
			dir + " --assetsDir .minecraft/assets/ " + "--accessToken 123456qwert" + " --assetIndex 17 " + "--uuid " + launchConfig.UUID + " --width " + launchConfig.Width + " --height " + launchConfig.Height

		createScript(launchToml)
	} else {
		slog.Info("Using the default configuration")
		userNameMd5 := Md5Create(userName)
		// TODO: Mod 加载器启动支持
		uuid := userNameMd5[:8] + "-" + userNameMd5[8:12] + "-" + userNameMd5[12:16] + "-" + userNameMd5[16:20] + "-" + userNameMd5[20:]
		cpArtifact := readArtifact(dir)
		launchToml = LAUNCH_TOML_DEFALUT_JVM + LAUNCH_TOML_LAUNCH_NAME + LAUNCH_TOML_LOG_FILE + dir + "/client-1.12.xml" +
			" -Djava.library.path=" + MkdirPath + " -cp " + cpArtifact + ".minecraft/versions/" + dir + "/" + dir + ".jar " + MINECRAFT_NO_MOD + " --username " + userName + " --version " + dir + " --gameDir .minecraft/versions/" +
			dir + " --assetsDir .minecraft/assets/ " + "--accessToken 123456qwert" + " --assetIndex 17 " + "--uuid " + uuid + " --width 1800 " + "--height 1000"
		createScript(launchToml)
	}

}

// Version: 不含 .jar 拓展名
func readArtifact(version string) string {
	var pathGet string
	versionJson, err := os.ReadFile(".minecraft/versions/" + version + "/" + version + ".json")
	if err != nil {
		Glog("ERROR", "readArtifact", "err", err)
	}

	libPath := gjson.Get(string(versionJson), "libraries.@dig:path")
	for _, path := range libPath.Array() {
		pathGet = pathGet + ".minecraft/libraries/" + path.String() + ":"
	}
	return pathGet
}

// 创建启动脚本
func createScript(toml string) {
	slog.Info("Create and write launch script")
	launchSh, err := os.Create("./.gmcl/launch.sh")
	if err != nil {
		Glog("ERROR", "createScrip", "err", err)
	} else {
		slog.Info("Create Success")
	}
	defer launchSh.Close()

	os.Chmod("./.gmcl/launch.sh", 0777)

	_, errWrite := launchSh.WriteString(toml)

	if errWrite != nil {
		Glog("ERROR", "createScrip", "errWrite", errWrite)
	} else {
		slog.Info("write Success")
		launchGame()
	}
}

// 执行脚本
func launchGame() {
	slog.Info("Start launching game")

	launch := exec.Command("bash", ".gmcl/launch.sh")

	err := launch.Run()
	if err != nil {
		Glog("ERROR", "launchGame", "err", err)
	}

}
