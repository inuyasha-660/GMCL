package utils

import (
	"io"
	"log"
	"net/http"
	"os"
)

const VERSION_MANIFEST_JSON = `https://piston-meta.mojang.com/mc/game/version_manifest.json`

// 单线程下载 version_manifest.json
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

// 多线程下载游戏
func DownloadsGmae(version string) {
	log.Println("开始下载, 选中版本:", version)
}
