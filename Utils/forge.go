package utils

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tidwall/gjson"
)

type DownInfo struct {
	MineVersion  string
	ForgeVersion string
}

// 获取适用于指定 minecraft 版本的 Forge 列表，后接 minecraft 版本，正序
const FORGE_LIST = `https://bmclapi2.bangbang93.com/forge/minecraft/`

// 根据 build 下载forge，后接 build
const FORGE_GET = `https://bmclapi2.bangbang93.com/forge/download/`

var RespRead []byte // 用于储存读取后的 resp.Body

// version string: forge版本
func InstallForge(version string) {
	if ok := JavaCheck(); ok {
		slog.Info("Start installing Forge")
		GetForge(version)
	} else {
		slog.Info("Installation failed")
	}
}

func JavaCheck() bool {
	java := exec.Command("java", " --version")
	var stdout, stderr bytes.Buffer
	java.Stdout = &stdout
	java.Stderr = &stderr
	err := java.Run()
	if err != nil {
		Glog("ERROR", "JavaCheck", "err", err)
		return false
	}

	javaVersion, errGet := stdout.String(), stderr.String()
	if errGet != "" {
		Glog("ERROR", "JavaCheck", "err", err)
		return false
	} else {
		fmt.Println(javaVersion)
		return true
	}
}

// version string: Forge 版本
func GetForge(version string) {
	buildForge := gjson.Get(string(RespRead), `#(version="`+version+`").build`)

	downUrl := FORGE_GET + buildForge.String()
	slog.Info("Get Forge from: " + downUrl + " Version: " + buildForge.String())

	resp, err := http.Get(downUrl)
	if err != nil {
		Glog("ERROR", "GetForge", "err", err)
	}
	defer resp.Body.Close()

	JarPath := "./forge-" + version + "-installer.jar"
	jar, errCreate := os.Create(JarPath)
	if errCreate != nil {
		Glog("ERROR", "GetForge", "errCreate", errCreate)
	}
	defer jar.Close()

	_, errCopy := io.Copy(jar, resp.Body)
	if errCopy != nil {
		Glog("ERROR", "GetForge", "errCopy", errCopy)
	} else {
		slog.Info("Download completed, Start to install Forge")
		InstallJar(JarPath)
	}

}

func InstallJar(path string) {
	install := exec.Command("java", "-jar", path, "nogui", "--installClient")
	var stdout, stderr bytes.Buffer
	install.Stdout = &stdout
	install.Stderr = &stderr
	err := install.Run()
	if err != nil {
		Glog("ERROR", "InstallJar", "err", err)
	}

	installOut, installErr := stdout.String(), stderr.String()
	if installErr != "" {
		slog.Error("Install Failed: " + installErr)
	} else {
		fmt.Println(installOut)
		slog.Info("Forge installed successfully")
	}

}

// version string: Minecraft 版本
func GetForgeList(version string) (ForgeList *[]string) {
	ListForge := &[]string{}
	slog.Info("Forge: Ture")
	Urllist := FORGE_LIST + version
	GlogINFO("Get forge list from: " + Urllist)

	resp, errGet := http.Get(Urllist)
	if errGet != nil {
		Glog("ERROR", "GetForgeLis", "errGet", errGet)
	}
	defer resp.Body.Close()

	respRead, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		Glog("ERROR", "GetForgeList", "errRead", errRead)
	} else {
		RespRead = respRead
	}

	ListForge = &[]string{}
	VersionForgeList := gjson.Get(string(respRead), `@dig:version`)
	for _, ForgeList := range VersionForgeList.Array() {
		*ListForge = append(*ListForge, ForgeList.Str)
	}

	return ListForge

}

// version string: Minecraft 版本
func ModSet(downWin fyne.Window, version string) {
	GameChoose := DownInfo{
		MineVersion: version,
	}
	ForgeList := GetForgeList(version)

	lable_ForgeVersion := widget.NewLabel("Forge")
	label_ForgeVersion := widget.NewLabel("Forge Version: " + GameChoose.ForgeVersion)

	Select_ForgeVersion := widget.NewSelect(*ForgeList, func(chooseForge string) {
		GameChoose.ForgeVersion = chooseForge
		label_ForgeVersion.SetText("Forge Version:" + GameChoose.ForgeVersion)
	})

	label_MinecraftVersion := widget.NewLabel("Minecraft Version: " + GameChoose.MineVersion)

	button_DownloadWithMods := widget.NewButton("Download", func() {
		//	DownloadsGmae(GameChoose.MineVersion, true, GameChoose.ForgeVersion, downWin)
		GetForge(GameChoose.ForgeVersion)
	})

	content_DownWithMods := container.NewVBox(lable_ForgeVersion, Select_ForgeVersion, label_MinecraftVersion, label_ForgeVersion, button_DownloadWithMods)
	downWin.SetContent(content_DownWithMods)

}
