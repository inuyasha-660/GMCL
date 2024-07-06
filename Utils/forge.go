package utils

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
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

func GetForge(version string) {
	slog.Info(version)
}

// version string: Minecraft 版本
func GetForgeList(version string) (ForgeList *[]string) {
	// TODO: 设置中添加是否自动下载 List 内最新 Froge 版本或者手动选择版本
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
		DownloadsGmae(GameChoose.MineVersion, true, GameChoose.ForgeVersion, downWin)
	})

	content_DownWithMods := container.NewVBox(lable_ForgeVersion, Select_ForgeVersion, label_MinecraftVersion, label_ForgeVersion, button_DownloadWithMods)
	downWin.SetContent(content_DownWithMods)

}
