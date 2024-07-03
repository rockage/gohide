package main

// https://developer.fyne.io/api/v2.0/layout
import (
	"gohide/Global"
	"gohide/Keyboard"
	"io"
	"os"

	"encoding/json"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
	"github.com/golang/freetype/truetype"
)

var (
	myApp             fyne.App // Fyne 总框架
	boxMain           *fyne.Container
	myWindow          fyne.Window     // 窗口
	content           *fyne.Container // 容器
	newWindowkeyNames *canvas.Text
	labelKInfo2       = widget.NewLabel("尚未设置")
	newWindowCreated  bool
	allKeyData        Global.Data
	keyStatus         int = 0 // 按键修改标识

)

// 中文字体设置
func setFont() {
	fontFilePath := "simhei.ttf"
	fontPath, _ := findfont.Find(fontFilePath)
	fontData, _ := os.ReadFile(fontPath)
	truetype.Parse(fontData)
	os.Setenv("FYNE_FONT", fontPath)
	os.Setenv("FYNE_FONT_MONOSPACE", fontPath)
}

func createUi() {
	allKeyData = readJSONFromFile() // Read JSON from file
	K1 := getKeyString(allKeyData.K1)
	K2 := getKeyString(allKeyData.K2)
	K3 := getKeyString(allKeyData.K3)

	// UI元素：
	var (
		boxKeySet1    *fyne.Container
		boxKeySet2    *fyne.Container
		boxKeySet3    *fyne.Container
		boxInfomation *fyne.Container
		button1       = widget.NewButton("设置", nil)
		button2       = widget.NewButton("设置", nil)
		button3       = widget.NewButton("设置", nil)
		buttonOk      = widget.NewButton("确定", nil)
		space         = layout.NewSpacer()
		labelK1       = widget.NewLabel("   隐藏&显示窗口:")
		labelK2       = widget.NewLabel("设置需要隐藏窗口:")
		labelK3       = widget.NewLabel(" 隐藏&显示本程序:")
		labelKInfo1   = widget.NewLabel("当前隐蔽窗口名称:")
		keyName1      = widget.NewLabel(K1)
		keyName2      = widget.NewLabel(K2)
		keyName3      = widget.NewLabel(K3)
	)
	// UI布局：

	boxKeySet1 = container.New(layout.NewGridLayout(5), space, labelK1, keyName1, button1, space)
	boxKeySet2 = container.New(layout.NewGridLayout(5), space, labelK2, keyName2, button2, space)
	boxKeySet3 = container.New(layout.NewGridLayout(5), space, labelK3, keyName3, button3, space)
	boxInfomation = container.New(layout.NewGridLayout(5), space, labelKInfo1, labelKInfo2, space, space)
	boxMain = container.New(layout.NewVBoxLayout(), boxKeySet1, boxKeySet2, boxKeySet3, boxInfomation)
	myWindow.Resize(fyne.NewSize(640, 200))
	myWindow.CenterOnScreen()
	myWindow.SetContent(boxMain)

	OnButtonTapped := func(button *widget.Button) { // 按钮 Tapped(即Click) 事件回调
		mainWindowKeyNames := widget.NewLabel("") // 主界面 KeyName 文字
		if !newWindowCreated {                    // 防止重复创建新窗口
			Keyboard.SetNoTopMost() // 暂时取消程序顶层状态
			switch button {
			case button1:
				keyStatus = 1
				K1 := getKeyString(allKeyData.K1)
				newWindowkeyNames.Text = K1
				mainWindowKeyNames = keyName1
			case button2:
				keyStatus = 2
				K2 := getKeyString(allKeyData.K2)
				newWindowkeyNames.Text = K2
				mainWindowKeyNames = keyName2
			case button3:
				keyStatus = 3
				K3 := getKeyString(allKeyData.K3)
				newWindowkeyNames.Text = K3
				mainWindowKeyNames = keyName3
			}

			newWindow := myApp.NewWindow("在键盘上按下新的快捷键") // 创建新的弹出窗口

			container_text := container.NewHBox(space, newWindowkeyNames, space)
			content = container.NewVBox(space, container_text, space, buttonOk)
			newWindow.SetContent(container.NewVBox(space, content, space))
			newWindow.Resize(fyne.NewSize(300, 100))
			newWindow.CenterOnScreen()
			newWindow.Show()
			newWindowCreated = true
			buttonOk.OnTapped = func() { // 确定按钮 Tap(Click） 事件回调
				// 同时刷新主界面文字 ** 注：Label 和 container 需要同时刷新，UI才能更新 **
				mainWindowKeyNames.Text = newWindowkeyNames.Text // 主界面的 KeyName Label 刷新
				boxMain.Refresh()                                // 主界面的 container 刷新
				newWindow.Close()
				saveJSONToFile(allKeyData)       // 存盘
				Keyboard.AllKeyData = allKeyData // 新设置立即生效

			}
			newWindow.SetOnClosed(func() { // 窗口关闭事件回调
				Keyboard.SetTopMost() // 恢复程序顶层状态
				keyStatus = 0
				newWindowCreated = false
			})
		}
	}

	button1.OnTapped = func() {
		OnButtonTapped(button1)
	}
	button2.OnTapped = func() {
		OnButtonTapped(button2)
	}
	button3.OnTapped = func() {
		OnButtonTapped(button3)
	}

	myWindow.ShowAndRun()
}

func getKeyString(s Global.KeyData) string {
	var keyString string

	for _, v := range s.Cons { // 控制按键：CTRL ALT SHIFT
		keyString += v.Name + " + "
	}
	keyString += s.Key.Name // 普通按键
	return keyString

}

func WindowName(s string) {
	labelKInfo2.Text = s
	boxMain.Refresh()
}

func KeyProce(s Global.KeyData) { // 按键回调处理
	switch keyStatus {
	case 1:
		allKeyData.K1 = s
		newWindowkeyNames.Text = getKeyString(s) // 刷新弹出窗口的文字
		content.Refresh()                        // 刷新弹出窗口的容器
	case 2:
		allKeyData.K2 = s
		newWindowkeyNames.Text = getKeyString(s)
		content.Refresh()
	case 3:
		allKeyData.K3 = s
		newWindowkeyNames.Text = getKeyString(s)
		content.Refresh()
	default:

	}

}

func main() {

	Keyboard.AllKeyData = readJSONFromFile()
	Keyboard.Callback1 = KeyProce
	Keyboard.Callback2 = WindowName
	defer Keyboard.ReleaseHook() // 释放钩子
	go Keyboard.Start()          // 将钩子设为协程
	setFont()                    // 设置字体，否则中文乱码
	myApp = app.New()
	myWindow = myApp.NewWindow("GoHide 摸鱼伴侣@rockage")
	newWindowkeyNames = canvas.NewText("", nil)
	createUi()
}

func readJSONFromFile() Global.Data {
	var data Global.Data
	jsonFile, _ := os.Open("setting.json")
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &data)
	return data
}
func saveJSONToFile(data Global.Data) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile("setting.json", jsonData, 0644)

}
