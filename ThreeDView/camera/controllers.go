package camera

import (
	. "FlightControl/ThreeDView/types"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

// OrbitController is a controller that allows the camera to orbit around a target
type OrbitController struct {
	BaseController
	target   Point3D    // The point the camera is orbiting around in world space
	rotation Rotation3D // The rotation of the camera around the target in world space in degrees from the perspective of the target
	distance Unit       // The distance of the camera from the target
}

func NewOrbitController(orbitCenter Point3D) *OrbitController {
	return &OrbitController{target: orbitCenter, distance: 500, rotation: Rotation3D{Y: 300}}
}

func (controller *OrbitController) setCamera(camera *Camera) {
	controller.BaseController.camera = camera
	controller.updatePosition()
}

func (controller *OrbitController) SetTarget(center Point3D) {
	controller.target = center
}

func (controller *OrbitController) Move(distance Unit) {
	controller.distance += distance
	controller.updatePosition()
}

func (controller *OrbitController) PointAtTarget() {
	direction := DirectionVector{Point3D: controller.target}
	direction.Subtract(controller.camera.Position)
	direction.Normalize()
	rotation := direction.ToRotation()
	controller.camera.Rotation.X = -rotation.X
	controller.camera.Rotation.Y = -rotation.Y

	controller.camera.Rotation.Z = controller.rotation.Z - 90
}

func (controller *OrbitController) Rotate(rotation Rotation3D) {
	controller.rotation.Add(rotation)
	controller.rotation.Normalize()
	controller.updatePosition()
}

func (controller *OrbitController) OnDrag(x, y float32) {
	controller.Rotate(Rotation3D{Y: Degrees(-y), Z: Degrees(x)})
}

func (controller *OrbitController) OnDragEnd() {}

func (controller *OrbitController) OnScroll(_, y float32) {
	controller.Move(Unit(y))
}

func (controller *OrbitController) updatePosition() {
	controller.camera.Position = controller.target
	controller.camera.Position.Add(Point3D{X: controller.distance})
	controller.camera.Position.Rotate(controller.target, controller.rotation)
	controller.PointAtTarget()
}

func (controller *OrbitController) GetRotationSlider() *fyne.Container {
	sliderYaw := widget.NewSlider(-360, 360)
	sliderYaw.Value = float64(controller.rotation.X)
	sliderYaw.OnChanged = func(value float64) {
		controller.rotation.X = Degrees(value)
		controller.updatePosition()
	}
	sliderPitch := widget.NewSlider(-360, 360)
	sliderPitch.Value = float64(controller.rotation.Y)
	sliderPitch.OnChanged = func(value float64) {
		controller.rotation.Y = Degrees(value)
		controller.updatePosition()
	}
	sliderRoll := widget.NewSlider(-360, 360)
	sliderRoll.Value = float64(controller.rotation.Z)
	sliderRoll.OnChanged = func(value float64) {
		controller.rotation.Z = Degrees(value)
		controller.updatePosition()
	}
	sliderContainer := container.NewVBox(sliderYaw, sliderPitch, sliderRoll)
	return sliderContainer
}

type ManualController struct {
	BaseController
}

func NewManualController() *ManualController {
	return &ManualController{}
}

func (controller *ManualController) GetRotationSlider() *fyne.Container {
	sliderYaw := widget.NewSlider(0, 360)
	sliderYaw.OnChanged = func(value float64) {
		controller.camera.Rotation.X = Degrees(value)
	}
	sliderPitch := widget.NewSlider(0, 360)
	sliderPitch.OnChanged = func(value float64) {
		controller.camera.Rotation.Y = Degrees(value)
	}
	sliderRoll := widget.NewSlider(0, 360)
	sliderRoll.OnChanged = func(value float64) {
		controller.camera.Rotation.Z = Degrees(value)
	}
	sliderContainer := container.NewVBox(sliderYaw, sliderPitch, sliderRoll)
	return sliderContainer
}

func (controller *ManualController) GetPositionControl() *fyne.Container {
	sliderX := widget.NewSlider(-100, 100)
	sliderX.OnChanged = func(value float64) {
		if value > 0 {
			controller.camera.Position.X += 10
		} else {
			controller.camera.Position.X -= 10
		}
	}
	sliderX.OnChangeEnded = func(value float64) {
		sliderX.Value = 0
	}

	sliderY := widget.NewSlider(-100, 100)
	sliderY.OnChanged = func(value float64) {
		if value > 0 {
			controller.camera.Position.Y += 10
		} else {
			controller.camera.Position.Y -= 10
		}
	}
	sliderY.OnChangeEnded = func(value float64) {
		sliderY.Value = 0
	}

	sliderZ := widget.NewSlider(-100, 100)
	sliderZ.OnChanged = func(value float64) {
		if value > 0 {
			controller.camera.Position.Z += 10
		} else {
			controller.camera.Position.Z -= 10
		}
	}
	sliderZ.OnChangeEnded = func(value float64) {
		sliderZ.Value = 0
	}

	buttonContainer := container.NewVBox(
		sliderX,
		sliderY,
		sliderZ,
	)
	return buttonContainer
}

func (controller *ManualController) GetInfoLabel() *widget.Label {
	label := widget.NewLabel("X: 0 Y: 0 Z: 0      Yaw: 0 Pitch: 0 Roll: 0")
	go func() {
		ticker := time.NewTicker(time.Second / 30)
		defer ticker.Stop()
		for range ticker.C {
			label.SetText(fmt.Sprintf("X: %d Y: %d Z: %d      Yaw: %d Pitch: %d Roll: %d",
				controller.camera.Position.X, controller.camera.Position.Y, controller.camera.Position.Z,
				controller.camera.Rotation.X, controller.camera.Rotation.Y, controller.camera.Rotation.Z))
			label.Refresh()
		}
	}()
	return label
}
