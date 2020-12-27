package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"

	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/flopp/go-findfont"
)

func init() {

	// 加载字体
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		// 黑体:simhei.ttf
		// 宋体:Songti.ttc
		if strings.Contains(path, "Songti.ttc") {
			// 设置环境变量
			os.Setenv("FYNE_FONT", path)
			break
		}
	}

	log.Println("字体已加载")
}

func main() {
	a := app.New()
	w := a.NewWindow("Proxy")

	w.Resize(fyne.NewSize(400, 200))

	// w.SetContent(container.NewBorder(container.NewVBox(makeEntry(), makeButton()),
	// 	nil, nil, nil, makeScrollList()))

	w.SetContent(container.NewVBox(makeEntry(), makeButton()))

	w.ShowAndRun()

	// 取消环境变量
	os.Unsetenv("FYNE_FONT")
}

// 输入框
func makeEntry() fyne.CanvasObject {
	entryPort := widget.NewEntry()
	entryPort.SetPlaceHolder("端口")
	entryTargetHost := widget.NewEntry()
	entryTargetHost.SetPlaceHolder("目标主机")

	return container.NewVBox(
		entryPort,
		entryTargetHost)
}

// 按钮
func makeButton() fyne.CanvasObject {

	btn1 := widget.NewButton("Run Server", RunServer)

	btn2 := widget.NewButton("Reload", Reload)

	btn3 := widget.NewButton("Close Server", CloseServer)

	return container.NewHBox(btn1, btn2, btn3)
}

// 列表（带滚动条）
func makeScrollList() fyne.CanvasObject {

	var data []string
	for i := 0; i < 20; i++ {
		data = append(data, fmt.Sprintf("URL A %d", i))
	}
	for i := 0; i < 20; i++ {
		data = append(data, fmt.Sprintf("URL B %d", i))
	}

	icon := widget.NewIcon(nil)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("模板Label"),
				widget.NewCheck("使用代理", func(on bool) {})) // 每隔15个的Checkbox也会选中，奇葩BUG
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// 模板列
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
			// CheckBox列
			item.(*fyne.Container).Objects[2].(*widget.Check).OnChanged = func(on bool) {
				log.Println(data[id] + " " + strconv.FormatBool(on))
			}
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		icon.SetResource(theme.DocumentIcon())
	}
	list.OnUnselected = func(id widget.ListItemID) {
		icon.SetResource(nil)
	}

	scroll := container.NewScroll(list)

	return scroll
}
