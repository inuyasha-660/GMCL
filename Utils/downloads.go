package utils

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

const VERSION_MANIFEST_JSON = `https://piston-meta.mojang.com/mc/game/version_manifest.json`

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

// 游戏本体
func DownloadsGmae(version string) {
	log.Println("开始下载, 选中版本:", version)
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
