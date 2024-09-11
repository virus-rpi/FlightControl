package main

import (
	"FlightControl/ThreeDView"
	"fyne.io/fyne/v2"
	"image/color"
)

func simulationTab() fyne.CanvasObject {
	threeDEnv := ThreeDView.NewThreeDWidget()
	threeDEnv.SetBackgroundColor(color.RGBA{R: 135, G: 206, B: 235, A: 255})
	plane := ThreeDView.NewPlane(1000, ThreeDView.Point3D{X: 0, Y: 0, Z: 0}, ThreeDView.Point3D{X: 0, Y: 0, Z: 0}, color.RGBA{G: 255, A: 255}, threeDEnv, 10)
	rocket := ThreeDView.NewRocket(300, ThreeDView.Point3D{X: 0, Y: 0, Z: 320}, ThreeDView.Point3D{X: 0, Y: 0, Z: 0}, color.RGBA{R: 255, G: 100, B: 10, A: 255}, threeDEnv, 2, 15)
	threeDEnv.AddObject(&plane)
	threeDEnv.AddObject(&rocket.ThreeDShape)

	camera := ThreeDView.NewCamera(ThreeDView.Point3D{Y: 500, Z: 200}, ThreeDView.Point3D{}, 300)
	orbitController := ThreeDView.NewOrbitController(ThreeDView.Point3D{X: 0, Y: 0, Z: 100})
	camera.SetController(orbitController)
	orbitController.PointAt(ThreeDView.Point3D{X: 0, Y: 0, Z: 100})
	threeDEnv.SetCamera(&camera)

	return threeDEnv
}
