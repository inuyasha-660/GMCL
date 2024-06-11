package utils

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/tidwall/gjson"
)

const VERSION_MANIFEST_JSON = `https://piston-meta.mojang.com/mc/game/version_manifest.json`
const OBJECT_HASH_GET = "https://resources.download.minecraft.net/"
const OBJECT_HASH_SAVE_DIR = "./.minecraft/assets/objects/"
const LIBRARISES = ".minecraft/libraries/"

// version_manifest.json
func GetVersionJson() {
	slog.Info("开始下载 version_manifest.json")
	resp, err := http.Get(VERSION_MANIFEST_JSON)

	if err != nil {
		slog.Error("下载失败:", err)
	}
	defer resp.Body.Close()

	manifestJson, errCreate := os.Create("./.gmcl/version_manifest.json")
	if errCreate != nil {
		slog.Error("创建文件失败:", errCreate)
	}
	defer manifestJson.Close()

	_, errCopy := io.Copy(manifestJson, resp.Body)

	if errCopy != nil {
		slog.Error("复制文件时出错:", errCopy)
	} else {
		slog.Info("下载完成")
	}
}

// Json 下载
func DownloadsGmae(Version string, forge bool, downWin fyne.Window) {
	slog.Info("开始下载, 选中版本: " + Version)
	var downsLog string
	entry_Down_sLog := widget.NewEntry() // 下载日志
	entry_Down_sLog.MultiLine = true     // 多行
	downsLog = "==> Start downloading: " + time.Now().Format("15:04:05") + "\n" + "=> Version: " + Version + " | " + "Forge: " + strconv.FormatBool(forge)
	entry_Down_sLog.SetText(downsLog)

	progress_Down := container.NewVBox(widget.NewProgressBarInfinite()) // 下载进度条

	content_Down_Start := container.NewVBox(container.New(layout.NewGridWrapLayout(fyne.NewSize(400, 400)), entry_Down_sLog), progress_Down)

	downWin.SetContent(content_Down_Start)

	jsonFile, errReadJson := os.ReadFile("./.gmcl/version_manifest.json")
	if errReadJson != nil {
		slog.Error("读取失败", errReadJson)
	}

	gameUrl := gjson.Get(string(jsonFile), `versions.#(id=="`+Version+`")+.url`)

	gameJson, errGet := http.Get(gameUrl.String())
	if errGet != nil {
		slog.Error("下载版本 Json 失败", errGet)
	}
	defer gameJson.Body.Close()

	errMkDir := os.MkdirAll(".minecraft/versions/"+Version, 0777)
	if errMkDir != nil {
		slog.Error("创建目录失败", errMkDir)
	}

	versionJson, errCreate := os.Create(".minecraft/versions/" + Version + "/" + path.Base(gameUrl.String()))
	if errCreate != nil {
		slog.Error("创建 Json 失败", errCreate)
	}
	defer versionJson.Close()

	_, errCopy := io.Copy(versionJson, gameJson.Body)

	if errCopy != nil {
		slog.Error("复制文件时出错:", errCopy)
	} else {
		slog.Info("下载" + path.Base(gameUrl.String()) + "完成")
		downsLog = downsLog + "\n" + "=> " + "Get: " + path.Base(gameUrl.String()) + "\n" + "-> Successfully"
		entry_Down_sLog.SetText(downsLog)
	}

	assetsDownload(".minecraft/versions/"+Version+"/"+path.Base(gameUrl.String()), entry_Down_sLog, downsLog)
	GetGameJar(Version, ".minecraft/versions/"+Version+"/"+path.Base(gameUrl.String()))
	GetLibraries(Version, ".minecraft/versions/"+Version+"/"+path.Base(gameUrl.String()))
}

func assetsDownload(path string, entry_Down_sLog *widget.Entry, downsLog string) {
	// 创建资源目录
	errIndex := os.MkdirAll(".minecraft/assets/indexes", 0777)
	if errIndex != nil {
		slog.Error("创建目录失败", errIndex)
	}
	errObject := os.MkdirAll(".minecraft/assets/objects", 0777)
	if errObject != nil {
		slog.Error("创建目录失败", errObject)
	}

	jsonFile, errRead := os.ReadFile(path)
	if errRead != nil {
		slog.Error("解析"+path+"失败", errRead)
	}

	assetIndex_ID := gjson.Get(string(jsonFile), `assetIndex.id`)
	assetIndex_Url := gjson.Get(string(jsonFile), `assetIndex.url`)

	slog.Info("开始下载 " + assetIndex_ID.String() + ".json")
	resp, err := http.Get(assetIndex_Url.String())

	if err != nil {
		slog.Error("下载失败:", err)
	}
	defer resp.Body.Close()

	indexPath := ".minecraft/assets/indexes/" + assetIndex_ID.String() + ".json"
	indexJson, errCreate := os.Create(indexPath)
	if errCreate != nil {
		slog.Error("创建文件失败:", errCreate)
	}
	defer indexJson.Close()

	_, errCopy := io.Copy(indexJson, resp.Body)

	if errCopy != nil {
		slog.Error("复制文件时出错:", errCopy)
	} else {
		slog.Info("下载完成")
		downsLog = downsLog + "\n" + "=> " + "Get: " + assetIndex_ID.String() + ".json" + "\n" + "-> Successfully"
		entry_Down_sLog.SetText(downsLog)

	}

	indexFile, errIndexRead := os.ReadFile(indexPath)
	if errIndexRead != nil {
		slog.Error("解析 "+indexPath, "失败", errIndexRead)
	}

	Object_Hash := gjson.Get(string(indexFile), `@dig:hash`) // 获取 Objects 内所有 hash 值
	for _, hash := range Object_Hash.Array() {
		url := OBJECT_HASH_GET + hash.String()[:2] + "/" + hash.String()
		dir := OBJECT_HASH_SAVE_DIR + hash.String()[:2] + "/" + hash.String()
		slog.Info("Get: " + url)
		downsLog = downsLog + "\n" + "Get: Objects"
		entry_Down_sLog.SetText(downsLog)

		resp, err := http.Get(url)
		if err != nil {
			slog.Error("Get: 失败", err)
		}
		defer resp.Body.Close()

		errMkDir := os.MkdirAll(OBJECT_HASH_SAVE_DIR+hash.String()[:2]+"/", 0777)
		if errMkDir != nil {
			slog.Error("Mkdir: 失败", errMkDir)
		}

		objectFiles, errCreate := os.Create(dir)
		if errCreate != nil {
			slog.Error("Create: 失败", errCreate)
		}
		defer objectFiles.Close()

		_, errCopy := io.Copy(objectFiles, resp.Body)
		if errCopy != nil {
			slog.Error("Copy: 失败", errCopy)
		} else {
			slog.Error("成功")
			downsLog = downsLog + "\n" + "-> Successfully"
			entry_Down_sLog.SetText(downsLog)
		}
	}

}

func GetGameJar(version, path string) {
	jsonFile, errRead := os.ReadFile(path)
	if errRead != nil {
		slog.Error("解析"+path+"失败", errRead)
	}

	gameUrl := gjson.Get(string(jsonFile), "downloads.client.url")
	slog.Error("Get: " + gameUrl.String())

	resp, err := http.Get(gameUrl.String())
	if err != nil {
		slog.Error("Get: 失败", err)
	}
	defer resp.Body.Close()

	gameFile, errCreate := os.Create("./.minecraft/versions/" + version + "/" + version + ".jar")
	if errCreate != nil {
		slog.Error("Create: 失败", errCreate)
	}
	defer gameFile.Close()

	_, errCopy := io.Copy(gameFile, resp.Body)
	if errCopy != nil {
		slog.Error("Copy: 失败", errCopy)
	} else {
		slog.Info("成功")
	}
}

func GetLibraries(version, path string) {
	jsonFile, errRead := os.ReadFile(path)
	if errRead != nil {
		slog.Error("解析"+path+"失败", errRead)
	}

	libUrl := gjson.Get(string(jsonFile), "libraries.@dig:url") // 获取资源 Url
	for _, url := range libUrl.Array() {
		libPath := gjson.Get(string(jsonFile), "libraries.@dig:path") // 获取资源 Path
		slog.Info("Get: " + url.String())
		slog.Info("Path: " + libPath.String())
	}
}

func GetGameList() (gameListRelease *[]string, gameListSnapshot *[]string) {

	jsonFile, errReadJson := os.ReadFile("./.gmcl/version_manifest.json")
	if errReadJson != nil {
		slog.Error("读取失败", errReadJson)
	}

	ListRelease := &[]string{}
	ListSnapshot := &[]string{}

	versionReleaseIDGet := gjson.Get(string(jsonFile), `versions.#(type=="release")#.id`)
	for _, versionsReleaseID := range versionReleaseIDGet.Array() {
		*ListRelease = append(*ListRelease, versionsReleaseID.Str)
	}

	versionSnapshotIDGet := gjson.Get(string(jsonFile), `latest.snapshot`) // 为性能只保留最新快照
	for _, versionSnapshotID := range versionSnapshotIDGet.Array() {
		*ListSnapshot = append(*ListSnapshot, versionSnapshotID.Str+" (Latest)")
	}

	return ListRelease, ListSnapshot
}
