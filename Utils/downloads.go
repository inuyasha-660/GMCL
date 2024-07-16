package utils

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
	slog.Info("Start downloading version_manifest.json")
	resp, err := http.Get(VERSION_MANIFEST_JSON)

	if err != nil {
		Glog("ERROR", "GetVersionJson", "err", err)
	}
	defer resp.Body.Close()

	manifestJson, errCreate := os.Create("./.gmcl/version_manifest.json")
	if errCreate != nil {
		Glog("ERROR", "GetVersionJson", "errCreate", errCreate)
	}
	defer manifestJson.Close()

	_, errCopy := io.Copy(manifestJson, resp.Body)

	if errCopy != nil {
		Glog("ERROR", "GetVersionJson", "errCopy", errCopy)
	} else {
		slog.Info("Download completed")
	}
}

// Json 下载
func DownloadsGmae(Version string, forge bool, forgeVersion string, downWin fyne.Window) {
	slog.Info("Start Downloading")
	slog.Info("Target: " + "Minecraft Version: " + Version + " Forge: " + strconv.FormatBool(forge) + " Forge Version: " + forgeVersion)
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
		Glog("ERROR", "DownloadsGmae", "errReadJson", errReadJson)
	}

	gameUrl := gjson.Get(string(jsonFile), `versions.#(id=="`+Version+`")+.url`)

	gameJson, errGet := http.Get(gameUrl.String())
	if errGet != nil {
		Glog("ERROR", "DownloadsGmae", "errGet", errGet)
	}
	defer gameJson.Body.Close()

	errMkDir := os.MkdirAll(".minecraft/versions/"+Version, 0777)
	if errMkDir != nil {
		Glog("ERROR", "DownloadsGmae", "errMkDir", errMkDir)
	}

	versionJson, errCreate := os.Create(".minecraft/versions/" + Version + "/" + path.Base(gameUrl.String()))
	if errCreate != nil {
		Glog("ERROR", "DownloadsGmae", "errCreate", errCreate)
	}
	defer versionJson.Close()

	_, errCopy := io.Copy(versionJson, gameJson.Body)

	if errCopy != nil {
		Glog("ERROR", "DownloadsGmae", "errCopy", errCopy)
	} else {
		slog.Info("Download " + path.Base(gameUrl.String()) + " completed")
		downsLog = downsLog + "\n" + "=> " + "Get: " + path.Base(gameUrl.String()) + "\n" + "-> Successly"
		entry_Down_sLog.SetText(downsLog)
	}

	assetsDownload(".minecraft/versions/"+Version+"/"+path.Base(gameUrl.String()), entry_Down_sLog, downsLog)
	GetGameJar(Version, ".minecraft/versions/"+Version+"/"+path.Base(gameUrl.String()), entry_Down_sLog, downsLog)
	GetLibraries(Version, ".minecraft/versions/"+Version+"/"+path.Base(gameUrl.String()), entry_Down_sLog, downsLog)

	if forge {
		GetForge(forgeVersion)
	}
}

func assetsDownload(path string, entry_Down_sLog *widget.Entry, downsLog string) {
	// 创建资源目录
	errIndex := os.MkdirAll(".minecraft/assets/indexes", 0777)
	if errIndex != nil {
		Glog("ERROR", "assetsDownload", "errIndex", errIndex)
	}
	errObject := os.MkdirAll(".minecraft/assets/objects", 0777)
	if errObject != nil {
		Glog("ERROR", "assetsDownload", "errObject", errObject)
	}

	jsonFile, errRead := os.ReadFile(path)
	if errRead != nil {
		Glog("ERROR", "assetsDownload", "errRead", errRead)
	}

	assetIndex_ID := gjson.Get(string(jsonFile), `assetIndex.id`)
	assetIndex_Url := gjson.Get(string(jsonFile), `assetIndex.url`)

	slog.Info("Start Downloading " + assetIndex_ID.String() + ".json")
	resp, err := http.Get(assetIndex_Url.String())

	if err != nil {
		Glog("ERROR", "assetsDownload", "err", err)
	}
	defer resp.Body.Close()

	indexPath := ".minecraft/assets/indexes/" + assetIndex_ID.String() + ".json"
	indexJson, errCreate := os.Create(indexPath)
	if errCreate != nil {
		Glog("ERROR", "assetsDownload", "errCreate", errCreate)
	}
	defer indexJson.Close()

	_, errCopy := io.Copy(indexJson, resp.Body)

	if errCopy != nil {
		Glog("ERROR", "assetsDownload", "errCopy", errCopy)
	} else {
		slog.Info("Download completed")
		downsLog = downsLog + "\n" + "=> " + "Get: " + assetIndex_ID.String() + ".json" + "\n" + "-> Success"
		entry_Down_sLog.SetText(downsLog)

	}

	indexFile, errIndexRead := os.ReadFile(indexPath)
	if errIndexRead != nil {
		slog.Error("Read"+indexPath, "failed:", errIndexRead)
	}
	downsLog = downsLog + "\n" + "Start Downloading Object"
	entry_Down_sLog.SetText(downsLog)

	Object_Hash := gjson.Get(string(indexFile), `@dig:hash`) // 获取 Objects 内所有 hash 值
	for _, hash := range Object_Hash.Array() {
		url := OBJECT_HASH_GET + hash.String()[:2] + "/" + hash.String()
		dir := OBJECT_HASH_SAVE_DIR + hash.String()[:2] + "/" + hash.String()
		slog.Info("Get: " + url)

		resp, err := http.Get(url)
		if err != nil {
			Glog("ERROR", "assetsDownload", "err", err)
		}
		defer resp.Body.Close()

		errMkDir := os.MkdirAll(OBJECT_HASH_SAVE_DIR+hash.String()[:2]+"/", 0777)
		if errMkDir != nil {
			Glog("ERROR", "assetsDownload", "errMkDir", errMkDir)
		}

		objectFiles, errCreate := os.Create(dir)
		if errCreate != nil {
			Glog("ERROR", "assetsDownload", "errCreate", errCreate)
		}
		defer objectFiles.Close()

		_, errCopy := io.Copy(objectFiles, resp.Body)
		if errCopy != nil {
			Glog("ERROR", "assetsDownload", "errCopy", errCopy)
		} else {
			slog.Info("Success")
		}
	}
	downsLog = downsLog + "\n" + "Download success"
	entry_Down_sLog.SetText(downsLog)

}

func GetGameJar(version, path string, entry_Down_sLog *widget.Entry, downsLog string) {
	downsLog = downsLog + "\n" + "Start Downloading Game Jar"
	entry_Down_sLog.SetText(downsLog)
	jsonFile, errRead := os.ReadFile(path)
	if errRead != nil {
		Glog("ERROR", "GetGameJar", "errRead", errRead)
	}

	gameUrl := gjson.Get(string(jsonFile), "downloads.client.url")
	slog.Error("Get: " + gameUrl.String())

	resp, err := http.Get(gameUrl.String())
	if err != nil {
		Glog("ERROR", "GetGameJar", "err", err)
	}
	defer resp.Body.Close()

	gameFile, errCreate := os.Create("./.minecraft/versions/" + version + "/" + version + ".jar")
	if errCreate != nil {
		Glog("ERROR", "GetGameJar", "errCreate", errCreate)
	}
	defer gameFile.Close()

	_, errCopy := io.Copy(gameFile, resp.Body)
	if errCopy != nil {
		Glog("ERROR", "GetGameJar", "errCopy", errCopy)
	} else {
		slog.Info("Success")
		downsLog = downsLog + "\n" + "Success"
		entry_Down_sLog.SetText(downsLog)
	}
}

func GetLibraries(version, path string, entry_Down_sLog *widget.Entry, downsLog string) {
	downsLog = downsLog + "\n" + "Start Downloading Libraries"
	entry_Down_sLog.SetText(downsLog)
	jsonFile, errRead := os.ReadFile(path)
	if errRead != nil {
		Glog("ERROR", "GetLibraries", "errRead", errRead)
	}

	errMkdir := os.Mkdir(".minecraft/libraries/", 0777)
	if errMkdir != nil {
		Glog("ERROR", "GetLibraries", "errMkdir", errMkdir)
	}

	libUrl := gjson.Get(string(jsonFile), "libraries.@dig:url") // 获取资源 Url

	for _, url := range libUrl.Array() {
		slog.Info("Get: " + url.String())

		libPath := gjson.Get(string(jsonFile), `libraries.#(downloads.artifact.url="`+url.String()+`").downloads.artifact.path`) // 获取Url对应的Path
		slog.Info("Path: " + libPath.String())

		resp, errGet := http.Get(url.String())
		if errGet != nil {
			Glog("ERROR", "GetLibraries", "errGet", errGet)
		}
		defer resp.Body.Close()

		mkdirPath, _ := filepath.Split(".minecraft/libraries/" + libPath.String())
		errMkdir := os.MkdirAll(mkdirPath, 0777)
		if errMkdir != nil {
			Glog("ERROR", "GetLibraries-for", "errMkdir", errMkdir)
		}

		libDown, errCreate := os.Create(".minecraft/libraries/" + libPath.String())
		if errCreate != nil {
			Glog("ERROR", "GetLibraries", "errCreate", errCreate)
		}

		_, errCopy := io.Copy(libDown, resp.Body)
		if errCopy != nil {
			Glog("ERRPR", "GetLibraries", "errCopy", errCopy)
		}

	}
	downsLog = downsLog + "\n" + "Download completed" + "\n" + "You can close the window"
	entry_Down_sLog.SetText(downsLog)

}

func GetGameList() (gameListRelease *[]string, gameListSnapshot *[]string) {

	jsonFile, errReadJson := os.ReadFile("./.gmcl/version_manifest.json")
	if errReadJson != nil {
		Glog("ERROR", "GetGameList", "errReadJson", errReadJson)
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
