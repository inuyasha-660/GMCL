package main

import (
	resources "go-mcl/Resources"
	utils "go-mcl/Utils"
	"log"
	"net/url"
	"os"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pelletier/go-toml/v2"
)

const LAUNCH_GAME = `
     Launch Game     
 `

const Downloads_Version_Json = `
     Dowmloads & Update Version Json     
 `

const EXIT = `
 	EXit launcher     
 `

const DownGame = `
 	Downloads Game    
 `

const RUNNING_LOG_LINE = `---------------------------`

type UserToml struct {
	UserName   string
	LoginDate  string
	ThemeColor string
}

func init() {
	log.Println("==> 少女祈祷中...")

	errMkdir := os.MkdirAll("./.gmcl", 0777)
	if errMkdir != nil {
		log.Println("==> 创建目录失败:", errMkdir)
	}
}

func main() {

	hmcl := app.New()
	hmcl.SetIcon(resources.ResourceIconSvg) // 设置 icon.svg 图标

	// 读取设置的主题
	var UserTheme UserToml

	user_theme, err := os.ReadFile("./.gmcl/user.toml")
	if err != nil {
		log.Println("读取用户主题失败")
	}

	errUnmarshal := toml.Unmarshal([]byte(user_theme), &UserTheme)
	if errUnmarshal != nil {
		log.Println("解析用户主题失败")
	}

	switch UserTheme.ThemeColor {
	case "Forgive-Green":
		{
			hmcl.Settings().SetTheme(&Forgive_Green{})
		}
	case "Dark":
		{
			hmcl.Settings().SetTheme(&Dark{})
		}
	case "Bili-Pink":
		{
			hmcl.Settings().SetTheme(&Bili_Pink{})
		}
	default:
		{
			log.Println("读取失败, 重置为默认主题")
			hmcl.Settings().SetTheme(theme.DefaultTheme())
		}
	}

	hmclWindow_Main := hmcl.NewWindow("HMCL - 1.0.0")
	hmclWindow_Main.SetMaster() // 设置为主窗口

	// Home页头像
	image_Author := canvas.NewImageFromResource(theme.AccountIcon())
	image_Author.FillMode = canvas.ImageFillContain
	image_Author.SetMinSize(fyne.NewSize(50, 50))

	// 用户登陆
	button_Login := container.NewVBox(widget.NewButton("Login in", func() {
		loginWin := hmcl.NewWindow("Login in")

		input_UserName := widget.NewEntry()
		input_UserName.SetPlaceHolder("UserName")

		input_UUID := widget.NewEntry()
		input_UUID.SetPlaceHolder("UUID")

		label_LoginLog := widget.NewLabel("")

		button_Login := widget.NewButton("Login", func() {
			ifSuccess := CreateUserToml(input_UserName.Text)
			if ifSuccess {
				label_LoginLog.SetText("Login Succeeded" + "\n" + "Please restart GMCL")
			} else {
				label_LoginLog.SetText("Login failed")
			}
		})

		content_Login := container.NewVBox(input_UserName, input_UUID, button_Login, label_LoginLog)
		loginWin.SetContent(content_Login)
		loginWin.Resize(fyne.NewSize(200, 250))
		loginWin.Show()
	}))

	// 用户信息
	toml_UserNmae, toml_LoginTime := ReadUserToml()
	lable_User := container.NewVBox(
		widget.NewLabel("User: "+toml_UserNmae),
		widget.NewLabel("Login Date: "+toml_LoginTime),
		widget.NewLabel(`Device Info `+"\n"+`System: `+runtime.GOOS+"\n"+`Arch: `+runtime.GOARCH),
		button_Login)

	list_DeviceInfo := container.NewVBox(image_Author, lable_User)

	// Home页运行日志
	lable_Log_Path := widget.NewLabel("Running Path: " + "\n" + GetPATH()) // 获取运行目录
	lable_Log := container.NewVBox(widget.NewLabel("Running Log"+"\n"+RUNNING_LOG_LINE), lable_Log_Path)

	// 启动游戏
	button_LaunchGame := container.NewVBox(widget.NewButton(LAUNCH_GAME, utils.LaunchCheck))

	// 下载 version_manifest.json
	button_Down_Version_Json := container.NewVBox(widget.NewButton(Downloads_Version_Json, utils.GetVersionJson))

	// 下载游戏
	button_Down_Game := container.NewVBox(widget.NewButton(DownGame, func() {
		dowmWin := hmcl.NewWindow("Download Game")

		// 游戏版本选择
		gameListRelease_Indicator, gameListSnapshot_Indicator := utils.GetGameList()
		var gameListRelease []string = *gameListRelease_Indicator // 指针类型转换
		var gameListSnapshot []string = *gameListSnapshot_Indicator

		label_GameChoose := widget.NewLabel("Choose ==>")

		gameTypeName_Release := widget.NewLabel("Release")
		Select_ReleaseChoose := widget.NewSelect(gameListRelease, func(chooseRelease string) { // 发行版
			log.Println(chooseRelease)
			label_GameChoose.SetText("Choose => " + chooseRelease)
		})

		gameTypeName_Snapshot := widget.NewLabel("Snapshot")
		Select_SnapshotChoose := widget.NewSelect(gameListSnapshot, func(chooseSnapshot string) { // 快照
			log.Println(chooseSnapshot)
			label_GameChoose.SetText("Choose => " + chooseSnapshot)
		})

		check_Forge := widget.NewCheck("Forge", func(forge bool) {
			log.Println("Forge:", forge)
		})

		button_DownWin := widget.NewButton("DownLoads", func() {
			lable_Down_Log := widget.NewLabel(time.Now().Format("15:04:05") + " => Start Download") // 下载日志

			progress_Down := widget.NewProgressBarInfinite() // 下载进度条

			content_Down_Start := container.NewVBox(lable_Down_Log, progress_Down)

			dowmWin.SetContent(content_Down_Start)
		})

		content_Down := container.NewVBox(gameTypeName_Release, Select_ReleaseChoose, gameTypeName_Snapshot, Select_SnapshotChoose, check_Forge, label_GameChoose, button_DownWin)

		dowmWin.Resize(fyne.NewSize(300, 500))
		dowmWin.SetContent(content_Down)
		dowmWin.Show()
	}))

	// 退出
	button_Exit := container.NewVBox(widget.NewButton(EXIT, Exit))

	button_All := container.NewVBox(button_Down_Version_Json, button_Down_Game, button_LaunchGame, button_Exit)

	// 设置界面

	// 生成默认配置
	lable_CreDefLauToml := widget.NewLabel("")

	button_CreDefLauToml := container.NewVBox(widget.NewButton("Create Launch toml", func() {
		log.Println("开始生成")
		lable_CreDefLauToml.SetText("Create Successfully")
	}))

	// 左边组件
	choices_Settings := container.NewVBox(widget.NewSelect([]string{"Dark", "Forgive-Green", "Bili-Pink"}, func(color string) {
		toml_Theme, err := os.OpenFile("./.gmcl/user.toml", os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			log.Println("主题 - 打开用户配置失败", err)
		}
		toml_Theme.Seek(0, 2) // 从配置文件的末尾第 0 个字符开始写入
		// [x,0]: 相对文件原点第 x 个字符， [x,1]: 从上次写入位置游标第 x 个字符开始 [x,2]: 相对文件末尾 x 个字符开始

		var theme_Set UserToml
		toml_ThemeFile, errRead := os.ReadFile("./.gmcl/user.toml")
		if errRead != nil {
			log.Println("ReadFile: 读取用户配置失败", errRead)
		}

		errUnmarshal := toml.Unmarshal([]byte(toml_ThemeFile), &theme_Set)
		if errUnmarshal != nil {
			log.Println("读取用户配置失败", err)
		}

		if theme_Set.ThemeColor == "" {

			_, errWriteTheme := toml_Theme.WriteString("\n" + "ThemeColor = " + `"` + color + `"`)
			if errWriteTheme != nil {
				log.Println("主题配置写入失败", err)
			}
		} else {
			log.Println("发现存在的主题配置, 尝试覆盖中")
			toml_Theme.Seek(0, 0)
			userName := theme_Set.UserName
			loginDate := theme_Set.LoginDate
			_, errWriteUserNmae := toml_Theme.WriteString("UserName = " + `"` + userName + `"` + "\n")
			if errWriteUserNmae != nil {
				log.Println("覆盖用户名失败", errWriteUserNmae)
			}
			_, errWriteDate := toml_Theme.WriteString(`LoginDate = ` + `"` + loginDate + `"`)
			if errWriteDate != nil {
				log.Println("覆盖日期失败", errWriteDate)
			}
			_, errWriteTheme := toml_Theme.WriteString("\n" + "ThemeColor = " + `"` + color + `"` + `         `)
			// `     ` - 空格作用: 避免 "Dark/Bili-Pink"主题长度不够导致未能完全覆盖 "Forgive-Green"主题
			if errWriteTheme != nil {
				log.Println("主题配置写入失败", err)
			}
		}

	}), lable_CreDefLauToml, button_CreDefLauToml) // 左边组建设置

	// 右边组件
	button_Settings := container.NewVBox(widget.NewButton("Set Theme", func() {
		var UserTheme UserToml

		user_theme, err := os.ReadFile("./.gmcl/user.toml")
		if err != nil {
			log.Println("读取用户主题失败")
		}

		errUnmarshal := toml.Unmarshal([]byte(user_theme), &UserTheme)
		if errUnmarshal != nil {
			log.Println("解析用户主题失败")
		}

		switch UserTheme.ThemeColor {
		case "Forgive-Green":
			{
				hmcl.Settings().SetTheme(&Forgive_Green{})
				log.Println("主题: " + UserTheme.ThemeColor + " 设置成功")
			}
		case "Dark":
			{
				hmcl.Settings().SetTheme(&Dark{})
				log.Println("主题: " + UserTheme.ThemeColor + " 设置成功")
			}
		case "Bili-Pink":
			{
				hmcl.Settings().SetTheme(&Bili_Pink{})
				log.Println("主题: " + UserTheme.ThemeColor + " 设置成功")
			}
		default:
			{
				log.Println("读取失败, 重置为默认主题")
				hmcl.Settings().SetTheme(theme.DefaultTheme())
			}
		}
	}))

	// 设置 - 启动器信息
	image_Icon := canvas.NewImageFromResource(resources.ResourceIconSvg)
	image_Icon.FillMode = canvas.ImageFillContain
	image_Icon.SetMinSize(fyne.NewSize(50, 50))

	link_Gmcl_Github := widget.NewHyperlink("Github", parseURL("https://github.com/inuyasha-660/GMCL"))
	label_Gmcl := container.NewBorder(nil, nil, widget.NewLabel("GMCL - 1.0.0"), link_Gmcl_Github)

	link_Author_Github := widget.NewHyperlink("Github", parseURL("https://github.com/inuyasha-660"))
	label_Author := container.NewBorder(nil, nil, widget.NewLabel("Inuyasha-660"), link_Author_Github)

	label_BuildDate := widget.NewLabel("Build Date: 2024-04-26 20:36") // TODO: 构建时修改

	label_AppInfo := container.NewVBox(label_Gmcl, label_Author, label_BuildDate)

	// 设置总布局
	content_AppInfo := container.NewVBox(image_Icon, label_AppInfo)
	content_Settings := container.NewBorder(nil, nil, choices_Settings, button_Settings)

	// AppTabs
	homeTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), container.NewBorder(nil, nil,
			list_DeviceInfo, // 用户/设备信息
			lable_Log,       // 运行日志
			button_All,      // 显示所有按钮
		)),
		// container.NewTabItem("Home", container.NewVBox(button_LaunchGame)),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), container.NewBorder(nil, nil, content_Settings, content_AppInfo)), // 注释为不带图标Tabs
		// container.NewTabItem("Settings", widget.NewLabel("Settings")),
	)

	homeTabs.SetTabLocation(container.TabLocationLeading)
	content := container.NewBorder(nil, nil, homeTabs, nil)

	hmclWindow_Main.SetContent(content)
	hmclWindow_Main.Resize(fyne.NewSize(800, 500))
	hmclWindow_Main.Show()
	hmcl.Run()
}

// 显示 Url
func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func GetPATH() string {
	Path, err := os.Getwd()
	if err != nil {
		return "Fail to get Path: "
	} else {
		return Path
	}
}

// 创建用户配置
func CreateUserToml(userName string) bool {
	user_toml, err := os.Create("./.gmcl/user.toml")
	if err != nil {
		log.Println("创建用户配置失败")
	}
	defer user_toml.Close()

	os.Chmod("./.gmcl/user.toml", 0777)

	_, errWriteUserNmae := user_toml.WriteString("UserName = " + `"` + userName + `"` + "\n")
	if errWriteUserNmae != nil {
		log.Println("用户名写入失败:", errWriteUserNmae)
		return false
	} else {
		log.Println("用户名写入成功")

		_, errWriteDate := user_toml.WriteString(`LoginDate = ` + `"` + time.Now().Format("2006-01-02 15:04:05") + `"`)
		if errWriteDate != nil {
			log.Println("登陆日期写入失败", errWriteDate)
			return false
		} else {
			log.Println("登陆日期写入成功")
			log.Println("登陆成功, 请重启启动器")
			return true
		}
	}
}

// 读取用户配置
func ReadUserToml() (toml_UserNmae string, toml_LoginDate string) {
	var UserToml UserToml

	toml_UserFile, errRead := os.ReadFile("./.gmcl/user.toml")
	if errRead != nil {
		log.Println("ReadFile: 读取用户配置失败", errRead)
	}

	err := toml.Unmarshal([]byte(toml_UserFile), &UserToml)
	if err != nil {
		log.Println("读取用户配置失败", err)
	}

	return UserToml.UserName, UserToml.LoginDate

}

// 退出启动器
func Exit() {
	log.Println("Exit with code: 0")
	os.Exit(0)
}
