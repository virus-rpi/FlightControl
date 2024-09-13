package main

import (
	"FlightControl/ThreeDView"
	"FlightControl/ThreeDView/camera"
	"FlightControl/ThreeDView/types"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"math"
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
	voltageLabel = widget.NewLabel("Voltage: N/A")
	statusLabel = widget.NewLabel("Status: Not connected")
	heightLabel = widget.NewLabel("Height: N/A")
	maxHeightLabel = widget.NewLabel("Max height: N/A")

	App.Preferences().AddChangeListener(func() {
		ipLabel.SetText("WaRa IP: " + App.Preferences().StringWithFallback("WaRaIP", "Not set"))
	})

	threeDVisualisation, rocket := threeDVisualisation()

	resetButton := widget.NewButton("Reset", func() { go reset(App) })
	deployParachuteButton := widget.NewButton("Deploy parachute", func() { go deployParachute(App) })
	deployStageButton := widget.NewButton("Deploy stage", func() { go deployStage(App) })
	startLoggingButton := widget.NewButton("Start logging", func() { go startLogging(App) })
	stopLoggingButton := widget.NewButton("Stop logging", func() { go stopLogging(App) })
	recalibrateGyroButton := widget.NewButton("Recalibrate gyro", func() { go recalibrateGyro(App) })
	recalibrateAccelerometerButton := widget.NewButton("Recalibrate accelerometer", func() { go recalibrateAccelerometer(App) })
	recalibrateBarometerButton := widget.NewButton("Recalibrate barometer", func() { go recalibrateBarometer(App) })
	resetMaxButton := widget.NewButton("Reset max", func() { go resetMax(App) })
	resetMinButton := widget.NewButton("Reset min", func() { go resetMin(App) })
	resetGyroButton := widget.NewButton("Reset gyro", func() { go resetGyro(App) })
	resetAccelerometerButton := widget.NewButton("Reset accelerometer", func() { go resetAccelerometer(App) })
	resetBarometerButton := widget.NewButton("Reset barometer", func() { go resetBarometer(App) })
	getLogButton := widget.NewButton("Get log", func() { go getLog(App) })

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

	infoLabelContainer := container.NewVBox(
		container.NewHBox(ipLabel, ipEditButton), voltageLabel, statusLabel, heightLabel, maxHeightLabel,
	)

	infoContainer := container.NewGridWithColumns(3,
		infoLabelContainer,
		container.NewCenter(),
		threeDVisualisation,
	)

	buttons := []fyne.CanvasObject{
		resetButton,
		deployParachuteButton,
		deployStageButton,
		startLoggingButton,
		stopLoggingButton,
		recalibrateGyroButton,
		recalibrateAccelerometerButton,
		recalibrateBarometerButton,
		resetMaxButton,
		resetMinButton,
		resetGyroButton,
		resetAccelerometerButton,
		resetBarometerButton,
		getLogButton,
	}

	var buttonContainer *fyne.Container
	if fyne.CurrentDevice().IsMobile() && fyne.CurrentDevice().Orientation() == 0 {
		buttonContainer = container.NewVBox(
			buttons...,
		)
	} else {
		buttonContainer = container.NewGridWithColumns(3,
			buttons...,
		)
	}

	content := container.NewVBox(
		infoContainer,
		buttonContainer,
	)

	go func() {
		selectedTabChannel := ps.Sub("selectedTab")
		for selectedTab := range selectedTabChannel {
			selectedTab := selectedTab.(*container.TabItem)
			if selectedTab.Text == "Control" {
				threeDVisualisation.Show()
			} else {
				threeDVisualisation.Hide()
			}
		}
	}()
	go func() {
		newestDataChannel := ps.Sub("newData")
		for newestData := range newestDataChannel {
			newestData := newestData.(Data)
			updateVoltage(newestData.voltage)
			updateStatus(string(newestData.status))
			updateHeight(newestData.altitude)
			updateMaxHeight(newestData.maxAltitude)

			rocket.DataChannel <- newestData
		}
	}()

	return content
}

func threeDVisualisation() (fyne.CanvasObject, *Rocket) {
	threeDEnv := ThreeDView.NewThreeDWidget()
	if fyne.CurrentDevice().IsMobile() {
		threeDEnv.SetFPSCap(30)
	} else {
		threeDEnv.SetFPSCap(60)
	}

	rocket := NewTwoStageRocket(types.Point3D{X: 0, Y: 0, Z: 0}, types.Rotation3D{X: 0, Y: 0, Z: 0}, threeDEnv)
	envCamera := camera.NewCamera(types.Point3D{}, types.Rotation3D{})
	orbitController := camera.NewOrbitController(rocket)
	orbitController.SetControlsEnabled(false)
	orbitController.SetRotation(types.Rotation3D{X: 0, Y: 0, Z: 0})
	orbitController.SetOffset(types.Point3D{X: 0, Y: 0, Z: -55})
	envCamera.SetController(orbitController)
	threeDEnv.SetCamera(&envCamera)

	updateDistance := func() {
		orbitController.SetDistance(types.Unit(math.Max(float64(threeDEnv.Size().Width), float64(threeDEnv.Size().Height)) / 2))
	}

	threeDEnv.RegisterTickMethod(func() {
		updateDistance()
	})

	return threeDEnv, rocket
}
