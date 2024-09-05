package main

import (
	"FlightControl/ThreeDView"
	"fyne.io/fyne/v2"
	"image/color"
)

func simulationTab() fyne.CanvasObject {
	threeDEnv := ThreeDView.NewThreeDWidget()
	threeDEnv.SetBackgroundColor(color.RGBA{R: 135, G: 206, B: 235, A: 255})
	plane := ThreeDView.NewPlane(1000, ThreeDView.Point3D{X: 0, Y: 0, Z: 0}, ThreeDView.Point3D{X: 0, Y: 0, Z: 0}, color.RGBA{G: 255, A: 255}, threeDEnv)
	cube := ThreeDView.NewCube(100, ThreeDView.Point3D{X: 0, Y: 0, Z: 100}, ThreeDView.Point3D{X: 0, Y: 0, Z: 0}, color.RGBA{B: 255, A: 255}, threeDEnv)
	camera := ThreeDView.NewCamera(ThreeDView.Point3D{Y: 500, Z: 200}, ThreeDView.Point3D{}, 30, 10)
	threeDEnv.AddObject(&plane)
	threeDEnv.AddObject(&cube)
	threeDEnv.SetCamera(&camera)
	camera.PointAt(ThreeDView.Point3D{X: 0, Y: 0, Z: 100})

	threeDEnv.RegisterAnimation(func() {
		cube.Rotation.Z += 1
		threeDEnv.Refresh()
	})

	return threeDEnv
}
