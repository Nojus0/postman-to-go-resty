package main

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"os"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Postman to Go Resty")
	w.SetMaster()
	// w.Resize(fyne.Size{Width: 300, Height: 300})
	// w.SetFixedSize(true)

	requestNameBox := widget.NewEntry()
	requestNameBox.PlaceHolder = "Request Name"

	httpBox := widget.NewEntry()
	httpBox.PlaceHolder = "HTTP Request etc.\n POST /something HTTP1.1"
	httpBox.MultiLine = true

	outputJson := widget.NewEntry()
	outputJson.PlaceHolder = `{"success":true}`
	outputJson.MultiLine = true

	inProgBox := container.NewVBox(
		widget.NewLabel("Generating..."),
	)

	var sucessBox *fyne.Container

	mainBox := container.NewPadded(container.NewVBox(
		requestNameBox,
		httpBox,
		outputJson,
		widget.NewSeparator(),
		widget.NewButton("Generate", func() {
			w.SetContent(inProgBox)
			w.RequestFocus()
			Run(requestNameBox.Text, outputJson.Text, httpBox.Text)
			requestNameBox.SetText("")
			outputJson.SetText("")
			httpBox.SetText("")
			w.SetContent(sucessBox)
			w.RequestFocus()
		}),
	))

	sucessBox = container.NewVBox(
		widget.NewLabel("Sucessfully Generated Go Resty Request"),
		widget.NewButton("OK", func() {
			w.SetContent(mainBox)
			w.RequestFocus()
		}),
	)

	w.Resize(fyne.NewSize(500, 252))

	w.SetContent(mainBox)

	w.ShowAndRun()

}

func Run(pkgName, outputJson, rawHttp string) {

	cwd, _ := os.Getwd()

	dirPath := path.Join(cwd, pkgName)

	os.Mkdir(dirPath, os.ModePerm)

	request, _ := http.ReadRequest(
		bufio.NewReader(
			bytes.NewBufferString(rawHttp),
		),
	)

	funcName := "Post" + pkgName

	restyCode := HttpRequestToResty(request, pkgName, funcName)

	os.WriteFile(path.Join(dirPath, funcName+".go"), []byte(restyCode), os.ModePerm)
	bodyBuf, _ := io.ReadAll(request.Body)

	if string(bodyBuf) != "" {
		inputTypes, _ := JsonToGolang(string(bodyBuf), "Input", pkgName)
		os.WriteFile(path.Join(dirPath, "Input.go"), []byte(inputTypes), os.ModePerm)
	}

	if outputJson != "" {
		outputTypes, _ := JsonToGolang(outputJson, "Output", pkgName)
		os.WriteFile(path.Join(dirPath, "Output.go"), []byte(outputTypes), os.ModePerm)
	}

}
