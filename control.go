package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func updateVoltage(voltage float64, currentLabel *widget.Label) {
	currentLabel.SetText("Current: " + fmt.Sprintf("%f", voltage) + " V")
}

func updateStatus(status string, statusLabel *widget.Label) {
	statusLabel.SetText("Status: " + status)
}

func updateHeight(height float64, heightLabel *widget.Label) {
	heightLabel.SetText("Height: " + fmt.Sprintf("%f", height) + " m")
}

func updateMaxHeight(maxHeight float64, maxHeightLabel *widget.Label) {
	maxHeightLabel.SetText("Max height: " + fmt.Sprintf("%f", maxHeight) + " m")
}

func controlTab(App fyne.App) fyne.CanvasObject {
	ipLabel := widget.NewLabel("WaRa IP: " + App.Preferences().StringWithFallback("WaRaIP", "Not set"))
	voltageLabel := widget.NewLabel("Voltage: 0 V")
	statusLabel := widget.NewLabel("Status: Not connected")
	heightLabel := widget.NewLabel("Height: 0 m")
	maxHeightLabel := widget.NewLabel("Max height: 0 m")

	App.Preferences().AddChangeListener(func() {
		ipLabel.SetText("WaRa IP: " + App.Preferences().StringWithFallback("WaRaIP", "Not set"))
	})

	resetButton := widget.NewButton("Reset", func() {})
	openHatchButton := widget.NewButton("Open hatch", func() {})
	startLoggingButton := widget.NewButton("Start logging", func() {})
	recalculateGyroButton := widget.NewButton("Recalculate gyro offset", func() {})
	resetGyroButton := widget.NewButton("Reset gyro", func() {})

	infoContainer := container.NewVBox(
		ipLabel, voltageLabel, statusLabel, heightLabel, maxHeightLabel,
	)

	content := container.NewVBox(
		infoContainer,
		container.NewGridWithColumns(2,
			resetButton, openHatchButton, startLoggingButton,
			recalculateGyroButton, resetGyroButton,
		),
	)

	return content
}
