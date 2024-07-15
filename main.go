package main

import (
	"errors"
	resources "go-mcl/Resources"
	utils "go-mcl/Utils"
	"log/slog"

	"net/url"
	"os"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
 	Exit launcher     
 `

const DownGame = `
 	Downloads Game    
 `
const HANDBOOK_FIRST = `
1. Download Version json
before Clicking Download 
Game.`

const HANDBOOK_SECOND = `
2. Select a version below 
before launching game.`

type UserToml struct {
	UserName   string
	LoginDate  string
	ThemeColor string
}

type GameChoose struct {
	Version      string
	Forge        bool
	ForgeVersion string
}

func init() {
	slog.Info("==> 少女祈祷中...")
	slog.Info("System: " + runtime.GOOS)
	slog.Info("Arch: " + runtime.GOARCH)
	path, err := os.Getwd()
	if err != nil {
		utils.Glog("ERROR", "init", "err", err)
	} else {
		slog.Info("Path: " + path)
	}

	errMkdir := os.MkdirAll("./.gmcl", 0777)
	if errMkdir != nil {
		utils.Glog("ERROR", "init", "errMkdir", errMkdir)
	}
}

func main() {

	gmcl := app.New()
	gmcl.SetIcon(resources.ResourceIconSvg) // 设置 icon.svg 图标

	// 读取设置的主题
	var UserTheme UserToml

	user_theme, err := os.ReadFile("./.gmcl/user.toml")
	if err != nil {
		utils.Glog("WARN", "main", "err", err)
	}

	errUnmarshal := toml.Unmarshal([]byte(user_theme), &UserTheme)
	if errUnmarshal != nil {
		utils.Glog("WARN", "main", "errUnmarshal", errUnmarshal)
	}

	switch UserTheme.ThemeColor {
	case "Forgive-Green":
		{
			gmcl.Settings().SetTheme(&Forgive_Green{})
		}
	case "Dark":
		{
			gmcl.Settings().SetTheme(&Dark{})
		}
	case "Bili-Pink":
		{
			gmcl.Settings().SetTheme(&Bili_Pink{})
		}
	default:
		{
			errSetTheme := errors.New("matched theme failed")
			utils.Glog("WARN", "main", "errSetTheme", errSetTheme)
			gmcl.Settings().SetTheme(theme.DefaultTheme())
		}
	}

	gmclWindow_Main := gmcl.NewWindow("GMCL - 1.0.0")
	gmclWindow_Main.SetMaster() // 设置为主窗口

	// Home页头像
	image_Author := canvas.NewImageFromResource(theme.AccountIcon())
	image_Author.FillMode = canvas.ImageFillContain
	image_Author.SetMinSize(fyne.NewSize(50, 50))

	// 用户登陆
	button_Login := container.NewVBox(widget.NewButton("Login in", func() {
		loginWin := gmcl.NewWindow("Login in")

		input_UserName := widget.NewEntry()
		input_UserName.SetPlaceHolder("UserName")

		label_LoginTime := widget.NewLabel("Login Time: ")

		label_LoginLog := widget.NewLabel("")

		button_Login := widget.NewButton("Login", func() {
			ifSuccess := CreateUserToml(input_UserName.Text)
			if ifSuccess {
				label_LoginTime.SetText("Login time: " + time.Now().Format("2006-01-02 15:04:05"))
				label_LoginLog.SetText("Login Succeeded" + "\n" + "Please restart GMCL")
			} else {
				label_LoginLog.SetText("Login failed")
			}
		})

		content_Login := container.NewVBox(input_UserName, label_LoginTime, button_Login, label_LoginLog)
		loginWin.SetContent(content_Login)
		loginWin.Resize(fyne.NewSize(250, 300))
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

	// 主页游戏版本扫描
	var LaunchVersion string
	VersionList := utils.VersionScan(".minecraft/versions/") // 获取扫描结果

	entry_GameList := widget.NewEntry()
	entry_GameList.MultiLine = true
	entry_GameList.SetText("HandBook" + "\n" + "-----------------" + HANDBOOK_FIRST + "\n" + HANDBOOK_SECOND)

	// 启动版本选择
	select_LaunchVersion := widget.NewSelect(*VersionList, func(launchVersionChoose string) {
		LaunchVersion = launchVersionChoose
		slog.Info("Launch version: " + LaunchVersion)
	})

	content_entry_GameList := container.NewVBox(container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 250)), entry_GameList), select_LaunchVersion)

	// 启动游戏
	button_LaunchGame := container.NewVBox(widget.NewButton(LAUNCH_GAME, func() {
		utils.LaunchCheck(LaunchVersion)
	}))

	// 下载 version_manifest.json
	button_Down_Version_Json := container.NewVBox(widget.NewButton(Downloads_Version_Json, utils.GetVersionJson))

	// 下载游戏
	button_Down_Game := container.NewVBox(widget.NewButton(DownGame, func() {
		dowmWin := gmcl.NewWindow("Download Game")

		// 游戏版本选择
		gameListRelease_Indicator, gameListSnapshot_Indicator := utils.GetGameList()
		var gameListRelease []string = *gameListRelease_Indicator // 指针类型转换
		var gameListSnapshot []string = *gameListSnapshot_Indicator

		GameVersionChoose := &GameChoose{}

		label_GameChoose := widget.NewLabel("Choose:")

		gameTypeName_Release := widget.NewLabel("Release")
		Select_ReleaseChoose := widget.NewSelect(gameListRelease, func(chooseRelease string) { // 发行版

			label_GameChoose.SetText("Choose: " + chooseRelease)
			GameVersionChoose.Version = chooseRelease
		})

		gameTypeName_Snapshot := widget.NewLabel("Snapshot")
		Select_SnapshotChoose := widget.NewSelect(gameListSnapshot, func(chooseSnapshot string) { // 快照

			label_GameChoose.SetText("Choose: " + chooseSnapshot)
			GameVersionChoose.Version = chooseSnapshot
		})

		button_DownWin := widget.NewButton("DownLoads", func() {
			utils.DownloadsGmae(GameVersionChoose.Version, false, "", dowmWin)
		})

		buuton_Mods := widget.NewButton("Mod Loader", func() {
			utils.ModSet(dowmWin, GameVersionChoose.Version)
		})

		content_Down := container.NewVBox(gameTypeName_Release, Select_ReleaseChoose, gameTypeName_Snapshot, Select_SnapshotChoose, label_GameChoose, button_DownWin, buuton_Mods)

		dowmWin.Resize(fyne.NewSize(400, 500))
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
		slog.Info("Start generating")
		lable_CreDefLauToml.SetText("Create Successfully")
	}))

	// 左边组件
	choices_Settings := container.NewVBox(widget.NewSelect([]string{"Dark", "Forgive-Green", "Bili-Pink"}, func(color string) {
		toml_Theme, err := os.OpenFile("./.gmcl/user.toml", os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			utils.Glog("ERROR", "main", "err", err)
		}

		// 从配置文件的末尾第 0 个字符开始写入
		// [x,0]: 相对文件原点第 x 个字符， [x,1]: 从上次写入位置游标第 x 个字符开始 [x,2]: 相对文件末尾 x 个字符开始
		_, errSeek := toml_Theme.Seek(0, 2)
		if errSeek != nil {
			utils.Glog("ERROR", "main", "errSeek", errSeek)
		}

		var theme_Set UserToml
		toml_ThemeFile, errRead := os.ReadFile("./.gmcl/user.toml")
		if errRead != nil {
			utils.Glog("ERROR", "main", "errRead", errRead)
		}

		errUnmarshal := toml.Unmarshal([]byte(toml_ThemeFile), &theme_Set)
		if errUnmarshal != nil {
			utils.Glog("ERROR", "main", "errUnmarshal", errUnmarshal)
		}

		if theme_Set.ThemeColor == "" {

			_, errWriteTheme := toml_Theme.WriteString("\n" + "ThemeColor = " + `"` + color + `"`)
			if errWriteTheme != nil {
				utils.Glog("ERROR", "main", "errWriteTheme", errWriteTheme)
			}
		} else {
			slog.Info("Find the old theme config, try to cover")
			_, errSeek = toml_Theme.Seek(0, 0)
			if errSeek != nil {
				utils.Glog("ERROR", "main", "errSeek", errSeek)
			}

			userName := theme_Set.UserName
			loginDate := theme_Set.LoginDate
			_, errWriteUserNmae := toml_Theme.WriteString("UserName = " + `"` + userName + `"` + "\n")
			if errWriteUserNmae != nil {
				utils.Glog("ERROR", "main", "errWriteUserNmae", errWriteUserNmae)
			}
			_, errWriteDate := toml_Theme.WriteString(`LoginDate = ` + `"` + loginDate + `"`)
			if errWriteDate != nil {
				utils.Glog("ERROR", "main", "errWriteDate", errWriteDate)
			}
			_, errWriteTheme := toml_Theme.WriteString("\n" + "ThemeColor = " + `"` + color + `"` + `         `)
			// `     ` - 空格作用: 避免 "Dark/Bili-Pink"主题长度不够导致未能完全覆盖 "Forgive-Green"主题
			if errWriteTheme != nil {
				utils.Glog("ERROR", "main", "errWriteTheme", errWriteTheme)
			}
		}

	}), lable_CreDefLauToml, button_CreDefLauToml) // 左边组建设置

	// 右边组件
	button_Settings := container.NewVBox(widget.NewButton("Set Theme", func() {
		var UserTheme UserToml

		user_theme, err := os.ReadFile("./.gmcl/user.toml")
		if err != nil {
			utils.Glog("ERROR", "main", "err", err)
		}

		errUnmarshal := toml.Unmarshal([]byte(user_theme), &UserTheme)
		if errUnmarshal != nil {
			utils.Glog("ERROR", "main", "errUnmarshal", errUnmarshal)
		}

		switch UserTheme.ThemeColor {
		case "Forgive-Green":
			{
				gmcl.Settings().SetTheme(&Forgive_Green{})
				slog.Info("Theme: " + UserTheme.ThemeColor + " Set success")
			}
		case "Dark":
			{
				gmcl.Settings().SetTheme(&Dark{})
				slog.Info("Theme: " + UserTheme.ThemeColor + " Set success")
			}
		case "Bili-Pink":
			{
				gmcl.Settings().SetTheme(&Bili_Pink{})
				slog.Info("Theme: " + UserTheme.ThemeColor + " Set success")
			}
		default:
			{
				errSetTheme_button := errors.New("matched theme failed")
				utils.Glog("WARN", "main", "errSetTheme_button", errSetTheme_button)
				gmcl.Settings().SetTheme(theme.DefaultTheme())
			}
		}
	}))

	// 设置-文件信息
	lable_UserConfigInfo := widget.NewLabel(`User Config: .gmcl/user.toml`)
	lable_LaunchConfigInfo := widget.NewLabel("Launch Config: .gmcl/launch.toml")
	lable_LaunchScript := widget.NewLabel("Launch Script: .gmcl/launch.sh")
	link_README := widget.NewHyperlink("README", parseURL("https://github.com/inuyasha-660/GMCL/blob/main/README.md"))
	lable_README := container.NewBorder(nil, nil, widget.NewLabel("More configuration:"), link_README)

	content_Left := container.NewBorder(nil, nil, choices_Settings, button_Settings)
	content_Right := container.NewVBox(lable_UserConfigInfo, lable_LaunchConfigInfo, lable_LaunchScript, lable_README)

	// 设置 - 启动器信息
	image_Icon := canvas.NewImageFromResource(resources.ResourceIconSvg)
	image_Icon.FillMode = canvas.ImageFillContain
	image_Icon.SetMinSize(fyne.NewSize(50, 50))

	link_Gmcl_Github := widget.NewHyperlink("Github", parseURL("https://github.com/inuyasha-660/GMCL"))
	label_Gmcl := container.NewBorder(nil, nil, widget.NewLabel("GMCL - 1.0.0"), link_Gmcl_Github)

	link_Author_Github := widget.NewHyperlink("Github", parseURL("https://github.com/inuyasha-660"))
	label_Author := container.NewBorder(nil, nil, widget.NewLabel("Inuyasha-660"), link_Author_Github)

	label_BuildDate := widget.NewLabel("Build Date: 2024-07-14 14:10") // TODO: 构建时修改

	label_AppInfo := container.NewVBox(label_Gmcl, label_Author, label_BuildDate)

	// 设置总布局
	content_AppInfo := container.NewVBox(image_Icon, label_AppInfo)
	content_Settings := container.NewBorder(nil, nil, content_Left, content_Right)

	// AppTabs
	homeTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), container.NewBorder(nil, nil,
			list_DeviceInfo,        // 用户/设备信息
			content_entry_GameList, // 运行日志
			button_All,             // 显示所有按钮
		)),
		// container.NewTabItem("Home", container.NewVBox(button_LaunchGame)),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), container.NewBorder(nil, nil, content_Settings, content_AppInfo)), // 注释为不带图标Tabs
		// container.NewTabItem("Settings", widget.NewLabel("Settings")),
	)

	homeTabs.SetTabLocation(container.TabLocationLeading)
	content := container.NewBorder(nil, nil, homeTabs, nil)

	gmclWindow_Main.SetContent(content)
	gmclWindow_Main.Resize(fyne.NewSize(800, 500))
	gmclWindow_Main.Show()
	gmcl.Run()
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
		utils.Glog("ERROR", "CreateUserToml", "err", err)
	}
	defer user_toml.Close()

	errChmod := os.Chmod("./.gmcl/user.toml", 0777)
	if errChmod != nil {
		utils.Glog("ERROR", "CreateUserToml", "errChmod", errChmod)
	}

	_, errWriteUserNmae := user_toml.WriteString("UserName = " + `"` + userName + `"` + "\n")
	if errWriteUserNmae != nil {
		utils.Glog("ERROR", "CreateUserToml", "errWriteUserNmae", errWriteUserNmae)
		return false
	} else {
		slog.Info("UserName written successfully")

		_, errWriteDate := user_toml.WriteString(`LoginDate = ` + `"` + time.Now().Format("2006-01-02 15:04:05") + `"`)
		if errWriteDate != nil {
			utils.Glog("ERROR", "CreateUserToml", "errWriteDate", errWriteDate)
			return false
		} else {
			slog.Info("LoginDate written successfully")
			slog.Info("Logined successfully, please restart GMCL")
			return true
		}
	}
}

// 读取用户配置
func ReadUserToml() (toml_UserNmae, toml_LoginDate string) {
	var UserToml UserToml

	toml_UserFile, errRead := os.ReadFile("./.gmcl/user.toml")
	if errRead != nil {
		utils.Glog("ERROR", "ReadUserToml", "errRead", errRead)
	}

	err := toml.Unmarshal([]byte(toml_UserFile), &UserToml)
	if err != nil {
		utils.Glog("ERROR", "ReadUserTom", "err", err)
	}

	return UserToml.UserName, UserToml.LoginDate

}

// 退出启动器
func Exit() {
	slog.Info("Exit with code: 0")
	os.Exit(0)
}
