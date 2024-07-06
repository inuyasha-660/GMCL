## GO-minecraft-launcher
一个使用``Go``编写的，适用于``*nix``的 Gui 离线我的世界启动器

## 截图


## 安装

### 使用预编译包
到[Releases](https://github.com/inuyasha-660/GMCL/releases)处选择版本并下载对应架构压缩包

``````bash
tar -zxvf go-mcl-[System]-[Arch]-[Version] # 请根据实际替换[]及其内容
./go-mcl
``````

### 从源码安装

``````bash
git clone https://github.com/inuyasha-660/GMCL.git && cd GMCL
go mod download
go run . # 直接运行
go build . && ./go-mcl # 编译运行
``````