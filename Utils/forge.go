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

const USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36"

var RespRead []byte // 用于储存读取后的 resp.Body

// version string: Forge 版本
func GetForge(version string) {
	buildForge := gjson.Get(string(RespRead), `#(version="`+version+`").build`)

	downUrl := FORGE_GET + buildForge.String()
	slog.Info("Get Forge from: " + downUrl + " Version: " + buildForge.String())

	client := &http.Client{}
	requ, errNewReq := http.NewRequest("GET", downUrl, nil)
	if errNewReq != nil {
		Glog("ERROR", "GetForge", "errNewReq", errNewReq)
	}
	requ.Header.Set("User-Agent", USER_AGENT)

	resp, errDo := client.Do(requ)
	if errDo != nil {
		Glog("ERROR", "GetForge", "errDo", errDo)
	}
	defer resp.Body.Close()

	JarPath := ".minecraft/forge-" + version + "-installer.jar"
	jar, errCreate := os.Create(JarPath)
	if errCreate != nil {
		Glog("ERROR", "GetForge", "errCreate", errCreate)
	}

	_, errCopy := io.Copy(jar, resp.Body)
	if errCopy != nil {
		Glog("ERROR", "GetForge", "errCopy", errCopy)
	} else {
		slog.Info("Download completed, Start to Write launcher_profiles.json")
		ProfilesWriter(JarPath)
	}

}

func ProfilesWriter(JarPath string) {
	file, errCreate := os.Create(".minecraft/launcher_profiles.json")
	if errCreate != nil {
		Glog("ERROR", "ProfilesWriter", "errCreate", errCreate)
	}

	_, err := file.WriteString(`{
    "profiles": {
        "(Default)": {
            "name": "(Default)"
        }
    },
    "selectedProfileName": "(Default)"
}
	`)
	if err != nil {
		Glog("ERROR", "ProfilesWriter", "err", err)
	} else {

		slog.Info("Write Completed, Start to Install Forge")
		InstallJar(JarPath)
	}
}

func InstallJar(path string) {
	installPath, errGetwd := os.Getwd()
	if errGetwd != nil {
		Glog("ERROR", "InstallJar", "errGetwd", errGetwd)
	}

	install := exec.Command("java", "-jar", path, "nogui", "--installClient", installPath+"/.minecraft")
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
		DownloadsGmae(GameChoose.MineVersion, true, GameChoose.ForgeVersion, downWin)
	})

	content_DownWithMods := container.NewVBox(lable_ForgeVersion, Select_ForgeVersion, label_MinecraftVersion, label_ForgeVersion, button_DownloadWithMods)
	downWin.SetContent(content_DownWithMods)

}
