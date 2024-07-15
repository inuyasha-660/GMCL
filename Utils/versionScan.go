package utils

import (
	"io/fs"
	"log/slog"
	"path"
	"path/filepath"
)

// 游戏目录扫描器
func VersionScan(scanDir string) (VersionList *[]string) {
	ScanList := &[]string{}
	err := filepath.WalkDir(scanDir, func(result string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		if path.Ext(filepath.Base(result)) == ".jar" {
			*ScanList = append(*ScanList, filepath.Base(result))
			return nil
		}

		return nil
	})
	if err != nil {
		Glog("ERROR", "VersionScan", "err", err)
	}

	return ScanList
}
