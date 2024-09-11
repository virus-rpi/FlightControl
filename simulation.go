package main

import (
	"FlightControl/ThreeDView"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"image/color"
)

func simulationTab() fyne.CanvasObject {
	threeDEnv := ThreeDView.NewThreeDWidget()
	threeDEnv.SetBackgroundColor(color.RGBA{R: 135, G: 206, B: 235, A: 255})
	ThreeDView.NewPlane(1000, ThreeDView.Point3D{X: 0, Y: 0, Z: 0}, ThreeDView.Rotation3D{X: 0, Y: 0, Z: 0}, color.RGBA{G: 255, A: 255}, threeDEnv, 10)
	ThreeDView.NewRocket(300, ThreeDView.Point3D{X: 0, Y: 0, Z: 320}, ThreeDView.Rotation3D{X: 0, Y: 0, Z: 0}, color.RGBA{R: 255, G: 100, B: 10, A: 255}, threeDEnv, 2, 15)
	ThreeDView.NewOrientationObject(threeDEnv)

	camera := ThreeDView.NewCamera(ThreeDView.Point3D{Y: 500, Z: 200}, ThreeDView.Rotation3D{})
	manualController := ThreeDView.NewManualController()
	camera.SetController(manualController)
	rotationSlider := manualController.GetRotationSlider()
	positionButtons := manualController.GetPositionControl()
	infoLabel := manualController.GetInfoLabel()
	//orbitController := ThreeDView.NewOrbitController(ThreeDView.Point3D{X: 0, Y: 0, Z: 100})
	//camera.SetController(orbitController)
	//orbitController.PointAtTarget()
	threeDEnv.SetCamera(&camera)

	return container.NewBorder(nil, nil, nil, nil, threeDEnv, container.NewVBox(rotationSlider, positionButtons, infoLabel))
}
