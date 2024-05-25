package utils

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tidwall/gjson"
)

const VERSION_MANIFEST_JSON = `https://piston-meta.mojang.com/mc/game/version_manifest.json`
const OBJECT_HASH_GET = "https://resources.download.minecraft.net/"
const OBJECT_HASH_SAVE_DIR = "./.minecraft/assets/objects/"

// version_manifest.json
func GetVersionJson() {
	log.Println("=> 开始下载 version_manifest.json")
	resp, err := http.Get(VERSION_MANIFEST_JSON)

	if err != nil {
		log.Println("下载失败:", err)
	}
	defer resp.Body.Close()

	manifestJson, errCreate := os.Create("./.gmcl/version_manifest.json")
	if errCreate != nil {
		log.Println("创建文件失败:", errCreate)
	}
	defer manifestJson.Close()

	_, errCopy := io.Copy(manifestJson, resp.Body)

	if errCopy != nil {
		log.Panicln("复制文件时出错:", errCopy)
	} else {
		log.Println("-> 下载完成")
	}
}

// Json 下载
func DownloadsGmae(Version string, forge bool, downWin fyne.Window) {
	log.Println("==> 开始下载, 选中版本:", Version)
	var downLog string
	entry_Down_Log := widget.NewEntry() // 下载日志
	entry_Down_Log.MultiLine = true     // 多行
	downLog = "==> Start downloading: " + time.Now().Format("15:04:05") + "\n" + "=> Version: " + Version + " | " + "Forge: " + strconv.FormatBool(forge)
	entry_Down_Log.SetText(downLog)

	progress_Down := widget.NewProgressBarInfinite() // 下载进度条

	content_Down_Start := container.NewVBox(entry_Down_Log, progress_Down)

	downWin.SetContent(content_Down_Start)

	jsonFile, errReadJson := os.ReadFile("./.gmcl/version_manifest.json")
	if errReadJson != nil {
		log.Println("读取失败", errReadJson)
	}

	gameUrl := gjson.Get(string(jsonFile), `versions.#(id=="`+Version+`")+.url`)

	gameJson, errGet := http.Get(gameUrl.String())
	if errGet != nil {
		log.Println("==> 下载版本 Json 失败", errGet)
	}
	defer gameJson.Body.Close()

	errMkDir := os.MkdirAll(".minecraft/versions/"+Version, 0777)
	if errMkDir != nil {
		log.Println("==> 创建目录失败", errMkDir)
	}

	versionJson, errCreate := os.Create(".minecraft/versions/" + Version + "/" + path.Base(gameUrl.String()))
	if errCreate != nil {
		log.Println("==> 创建 Json 失败", errCreate)
	}
	defer versionJson.Close()

	_, errCopy := io.Copy(versionJson, gameJson.Body)

	if errCopy != nil {
		log.Panicln("复制文件时出错:", errCopy)
	} else {
		log.Println("-> 下载" + path.Base(gameUrl.String()) + "完成")
		downLog = downLog + "\n" + "=> " + "Download: " + path.Base(gameUrl.String()) + "\n" + "-> Successfully"
		entry_Down_Log.SetText(downLog)
	}

	assetsDownload(".minecraft/versions/"+Version+"/"+path.Base(gameUrl.String()), entry_Down_Log, downLog)
	GetGameJar(Version, ".minecraft/versions/"+Version+"/"+path.Base(gameUrl.String()))
}

func assetsDownload(path string, entry_Down_Log *widget.Entry, downLog string) {
	// 创建资源目录
	errIndex := os.MkdirAll(".minecraft/assets/indexes", 0777)
	if errIndex != nil {
		log.Println("==> 创建目录失败", errIndex)
	}
	errObject := os.MkdirAll(".minecraft/assets/objects", 0777)
	if errObject != nil {
		log.Println("==> 创建目录失败", errObject)
	}

	jsonFile, errRead := os.ReadFile(path)
	if errRead != nil {
		log.Println("解析"+path+"失败", errRead)
	}

	assetIndex_ID := gjson.Get(string(jsonFile), `assetIndex.id`)
	assetIndex_Url := gjson.Get(string(jsonFile), `assetIndex.url`)

	log.Println("=> 开始下载 " + assetIndex_ID.String() + ".json")
	resp, err := http.Get(assetIndex_Url.String())

	if err != nil {
		log.Println("下载失败:", err)
	}
	defer resp.Body.Close()

	indexPath := ".minecraft/assets/indexes/" + assetIndex_ID.String() + ".json"
	indexJson, errCreate := os.Create(indexPath)
	if errCreate != nil {
		log.Println("创建文件失败:", errCreate)
	}
	defer indexJson.Close()

	_, errCopy := io.Copy(indexJson, resp.Body)

	if errCopy != nil {
		log.Println("复制文件时出错:", errCopy)
	} else {
		log.Println("-> 下载完成")
		downLog = downLog + "\n" + "=> " + "Download: " + assetIndex_ID.String() + ".json" + "\n" + "-> Successfully"
		entry_Down_Log.SetText(downLog)

	}

	indexFile, errIndexRead := os.ReadFile(indexPath)
	if errIndexRead != nil {
		log.Println("解析", indexFile, "失败", errIndexRead)
	}

	Object_Hash := gjson.Get(string(indexFile), `@dig:hash`) // 获取 Objects 内所有 hash 值
	for _, hash := range Object_Hash.Array() {
		url := OBJECT_HASH_GET + hash.String()[:2] + "/" + hash.String()
		dir := OBJECT_HASH_SAVE_DIR + hash.String()[:2] + "/" + hash.String()
		log.Println("=> Get:", url)

		resp, err := http.Get(url)
		if err != nil {
			log.Println("-> Get: 失败", err)
		}
		defer resp.Body.Close()

		errMkDir := os.MkdirAll(OBJECT_HASH_SAVE_DIR+hash.String()[:2]+"/", 0777)
		if errMkDir != nil {
			log.Println("-> Mkdir: 失败", errMkDir)
		}

		objectFiles, errCreate := os.Create(dir)
		if errCreate != nil {
			log.Println("-> Create: 失败", errCreate)
		}
		defer objectFiles.Close()

		_, errCopy := io.Copy(objectFiles, resp.Body)
		if errCopy != nil {
			log.Println("-> Copy: 失败", errCopy)
		} else {
			log.Println("-> 成功")
		}
	}

}

func GetGameJar(version, path string) {
	jsonFile, errRead := os.ReadFile(path)
	if errRead != nil {
		log.Println("解析"+path+"失败", errRead)
	}

	gameUrl := gjson.Get(string(jsonFile), "downloads.client.url")
	log.Println("=> Get:", gameUrl)

	resp, err := http.Get(gameUrl.String())
	if err != nil {
		log.Println("-> Get: 失败", err)
	}
	defer resp.Body.Close()

	gameFile, errCreate := os.Create("./.minecraft/versions/" + version + "/" + version + ".jar")
	if errCreate != nil {
		log.Println("-> Create: 失败", errCreate)
	}
	defer gameFile.Close()

	_, errCopy := io.Copy(gameFile, resp.Body)
	if errCopy != nil {
		log.Println("-> Copy: 失败", errCopy)
	} else {
		log.Println("-> 成功")
	}
}

func GetGameList() (gameListRelease *[]string, gameListSnapshot *[]string) {

	jsonFile, errReadJson := os.ReadFile("./.gmcl/version_manifest.json")
	if errReadJson != nil {
		log.Println("读取失败", errReadJson)
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
