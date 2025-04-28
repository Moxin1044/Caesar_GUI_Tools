package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// 对文本进行凯撒加解密或暴力破解
func caesar(text string, shift int, mode string, brute bool) []string {
	if brute {
		results := make([]string, 0, 25)
		for i := 1; i <= 25; i++ {
			results = append(results,
				fmt.Sprintf("Shift %2d: %s", i, processText(text, i, "decrypt")))
		}
		return results
	}
	return []string{processText(text, shift, mode)}
}

func processText(text string, shift int, mode string) string {
	var result []rune
	for _, ch := range text {
		if idx := strings.IndexRune(letters, ch); idx >= 0 {
			base := 0
			if idx >= 26 {
				base = 26
			}
			offset := idx - base
			if mode == "decrypt" {
				offset = (offset - shift + 26) % 26
			} else {
				offset = (offset + shift) % 26
			}
			result = append(result, rune(letters[base+offset]))
		} else {
			result = append(result, ch)
		}
	}
	return string(result)
}

func parseShift(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 || v > 25 {
		return 0, fmt.Errorf("移位值需为 1–25 之间的整数")
	}
	return v, nil
}

func main() {
	myApp := app.New()
	win := myApp.NewWindow("凯撒密码工具 v1.0")
	win.Resize(fyne.NewSize(800, 600))
	win.SetFixedSize(true)

	// 多行编辑框：输入、输出
	inputEntry := widget.NewMultiLineEntry()
	inputEntry.Wrapping = fyne.TextWrapOff
	inputScroll := container.NewScroll(inputEntry)
	inputScroll.SetMinSize(fyne.NewSize(0, 200))

	outputEntry := widget.NewMultiLineEntry()
	outputEntry.Wrapping = fyne.TextWrapOff
	outputScroll := container.NewScroll(outputEntry)
	outputScroll.SetMinSize(fyne.NewSize(0, 200))

	// 控件：模式、移位、暴力破解
	shiftEntry := widget.NewEntry()
	modeSelect := widget.NewSelect([]string{"加密", "解密"}, nil)
	bruteCheck := widget.NewCheck("暴力破解", nil)

	// 执行按钮
	runBtn := widget.NewButton("执行加解密", func() {
		go func() {
			mode := "encrypt"
			if modeSelect.Selected == "解密" {
				mode = "decrypt"
			}
			var res []string
			if bruteCheck.Checked {
				res = caesar(inputEntry.Text, 0, "", true)
			} else {
				shiftVal, err := parseShift(shiftEntry.Text)
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				res = caesar(inputEntry.Text, shiftVal, mode, false)
			}
			outputEntry.SetText(strings.Join(res, "\n"))
		}()
	})

	// 布局：工具栏
	toolbar := container.NewHBox(
		widget.NewLabel("操作模式:"), modeSelect,
		widget.NewLabel("移位值:"), shiftEntry,
		bruteCheck, runBtn,
	)

	// 主界面：文本处理
	content := container.NewVBox(
		toolbar,
		widget.NewLabel("输入文本："), inputScroll,
		widget.NewLabel("输出结果："), outputScroll,
	)

	win.SetContent(content)
	win.ShowAndRun()
}
