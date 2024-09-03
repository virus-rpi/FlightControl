package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var voltageLabel *widget.Label
var statusLabel *widget.Label
var heightLabel *widget.Label
var maxHeightLabel *widget.Label

func updateVoltage(voltage float64) {
	voltageLabel.SetText("Voltage: " + fmt.Sprintf("%f", voltage) + " V")
}

func updateStatus(status string) {
	statusLabel.SetText("Status: " + status)
}

func updateHeight(height float64) {
	heightLabel.SetText("Height: " + fmt.Sprintf("%f", height) + " m")
}

func updateMaxHeight(maxHeight float64) {
	maxHeightLabel.SetText("Max height: " + fmt.Sprintf("%f", maxHeight) + " m")
}

func controlTab(App fyne.App, MainWindow fyne.Window) fyne.CanvasObject {
	ipLabel := widget.NewLabel("WaRa IP: " + App.Preferences().StringWithFallback("WaRaIP", "Not set"))
	voltageLabel = widget.NewLabel("Voltage: N/A V")
	statusLabel = widget.NewLabel("Status: Not connected")
	heightLabel = widget.NewLabel("Height: N/A m")
	maxHeightLabel = widget.NewLabel("Max height: N/A m")

	App.Preferences().AddChangeListener(func() {
		ipLabel.SetText("WaRa IP: " + App.Preferences().StringWithFallback("WaRaIP", "Not set"))
	})

	resetButton := widget.NewButton("Reset", func() {})
	openHatchButton := widget.NewButton("Open hatch", func() {})
	startLoggingButton := widget.NewButton("Start logging", func() {})
	recalculateGyroButton := widget.NewButton("Recalculate gyro offset", func() {})
	resetGyroButton := widget.NewButton("Reset gyro", func() {})
	getLogButton := widget.NewButton("Get log", func() {})

	ipEditButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		ipEntry := widget.NewEntry()
		ipEntry.SetPlaceHolder("Enter IP")
		dialog.ShowForm("Set WaRa IP", "OK", "Cancel", []*widget.FormItem{
			widget.NewFormItem("IP", ipEntry),
		}, func(ok bool) {
			if !ok {
				return
			}
			App.Preferences().SetString("WaRaIP", ipEntry.Text)
			updateWebsocket(App)
		}, MainWindow)
	})

	infoContainer := container.NewVBox(
		container.NewHBox(ipLabel, ipEditButton), voltageLabel, statusLabel, heightLabel, maxHeightLabel,
	)

	var buttonContainer *fyne.Container
	if fyne.CurrentDevice().IsMobile() {
		buttonContainer = container.NewVBox(
			resetButton, openHatchButton, startLoggingButton,
			recalculateGyroButton, resetGyroButton, getLogButton,
		)
	} else {
		buttonContainer = container.NewGridWithColumns(2,
			resetButton, openHatchButton, startLoggingButton,
			recalculateGyroButton, resetGyroButton, getLogButton,
		)
	}

	content := container.NewVBox(
		infoContainer,
		buttonContainer,
	)

	return content
}
