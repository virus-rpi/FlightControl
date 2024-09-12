package main

import (
	"FlightControl/ThreeDView"
	"FlightControl/ThreeDView/camera"
	"FlightControl/ThreeDView/object"
	"FlightControl/ThreeDView/types"
	"fyne.io/fyne/v2"
	"image/color"
)

func simulationTab() fyne.CanvasObject {
	threeDEnv := ThreeDView.NewThreeDWidget()
	threeDEnv.SetBackgroundColor(color.RGBA{R: 135, G: 206, B: 235, A: 255})
	threeDEnv.SetFPSCap(60)
	threeDEnv.SetTPSCap(600)

	object.NewPlane(10000, types.Point3D{X: 0, Y: 0, Z: 0}, types.Rotation3D{X: 0, Y: 0, Z: 0}, color.RGBA{G: 255, A: 255}, threeDEnv, 50)
	rocket := object.NewRocket(300, types.Point3D{X: 0, Y: 0, Z: 320}, types.Rotation3D{X: 0, Y: 0, Z: 0}, color.RGBA{R: 255, G: 100, B: 10, A: 255}, threeDEnv, 2, 15)

	envCamera := camera.NewCamera(types.Point3D{Y: 500, Z: 200}, types.Rotation3D{})
	orbitController := camera.NewOrbitController(rocket)
	envCamera.SetController(orbitController)
	threeDEnv.SetCamera(&envCamera)

	threeDEnv.RegisterTickMethod(func() {
		rocket.Position.Z += 5
		orbitController.Update()
	})

	return threeDEnv
}
