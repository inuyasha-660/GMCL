## GO-minecraft-launcher
一个使用``Go``编写的，适用于``*nix``的 Gui 离线我的世界启动器

## Features
- [x] 用户登陆
- [x] 主题设置
- [x] Mod loader支持(Forge)
- [x] 广泛的游戏版本可选
- [x] 自定义启动参数

## 截图
![Home](./Resources/GMCL-Home-v0.9.0.png)

**Text in GMCL is translated by Google.**

## 安装
``unzip``和``bash``将被用于解压/执行脚本，确保系统已安装并添加到``PATH``

### 使用预编译包
到[Releases](https://github.com/inuyasha-660/GMCL/releases)处选择版本并下载对应架构压缩包

``````bash
tar -zxvf go-mcl-[System]-[Arch]-[Version].tar.gz # 请根据实际替换[]及其内容
./go-mcl
``````

### 从源码安装

### 依赖安装
根据[Getting Started](https://docs.fyne.io/started/)安装相应依赖

### 编译
``````bash
git clone https://github.com/inuyasha-660/GMCL.git && cd GMCL
go mod tidy
go run . # 直接运行
go build . && ./go-mcl # 编译运行
``````

## 自定义配置
点击设置页``Create Launch toml``生成配置

| Key      |  值类型   | 含义 |
|----------|----------|-----|
| Xmx      | string | 最大堆大小 |
| Xmn      | string | 新生代堆大小 |
| UseG1GC  |  bool  | 开启G1 |
| UseAdaptiveSizePolicy | bool | 自动选择年轻代区大小和相应的Survivor区比例 |
| OmitStackTraceInFastThrow | bool |  省略异常栈信息从而快速抛出 |
| Width | string | 窗口宽度 |
| Height | string | 窗口高度 |
| UUID | string | 用户UUID |


## 启动
出现启动窗口后即可关闭启动器，游戏日志位于``./logs``目录
