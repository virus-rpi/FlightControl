package main

import (
	"FlightControl/ThreeDView"
	"FlightControl/ThreeDView/camera"
	"FlightControl/ThreeDView/object"
	"FlightControl/ThreeDView/types"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

func simulationTab() fyne.CanvasObject {
	threeDEnv := ThreeDView.NewThreeDWidget()
	threeDEnv.SetBackgroundColor(color.RGBA{R: 135, G: 206, B: 235, A: 255})
	threeDEnv.SetFPSCap(60)
	threeDEnv.SetTPSCap(600)

	object.NewPlane(10000, types.Point3D{X: 0, Y: 0, Z: 0}, types.Rotation3D{X: 0, Y: 0, Z: 0}, color.RGBA{G: 255, A: 255}, threeDEnv, 50)
	rocket := NewTwoStageRocket(types.Point3D{X: 0, Y: 0, Z: 0}, types.Rotation3D{X: 0, Y: 0, Z: 0}, threeDEnv)

	envCamera := camera.NewCamera(types.Point3D{Y: 500, Z: 200}, types.Rotation3D{})
	orbitController := camera.NewOrbitController(rocket)
	envCamera.SetController(orbitController)
	threeDEnv.SetCamera(&envCamera)

	threeDEnv.RegisterTickMethod(func() {
		rocket.Move(types.Point3D{X: 0, Y: 0, Z: 1})
		orbitController.Update()
	})

	separateButton := widget.NewButton("Separate", func() {
		rocket.SeparateStage()
	})
	separateButton.Resize(fyne.NewSize(100, 50))
	buttonContainer := container.NewVBox(separateButton)

	return container.NewBorder(nil, nil, nil, nil, threeDEnv, buttonContainer)
}
