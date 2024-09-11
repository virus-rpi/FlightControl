package ThreeDView

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

type OrbitController struct {
	BaseController
	target   Point3D
	rotation Rotation3D
	distance Unit
}

func NewOrbitController(orbitCenter Point3D) *OrbitController {
	return &OrbitController{target: orbitCenter, distance: 500}
}

func (controller *OrbitController) SetTarget(center Point3D) {
	controller.target = center
}

func (controller *OrbitController) Move(distance Unit) {
	controller.distance += distance
	controller.updatePosition()
}

func (controller *OrbitController) PointAtTarget() {
	// TODO: Point the camera at the target
}

func (controller *OrbitController) Rotate(rotation Rotation3D) {
	controller.rotation.Add(rotation)
	controller.updatePosition()
}

func (controller *OrbitController) onDrag(x, y float32) {
	// TODO: Translate 2D drag to 2D rotation to 3D rotation
	controller.Rotate(Rotation3D{X: Degrees(y / 20), Y: Degrees(x / 20)})
}

func (controller *OrbitController) onDragEnd() {}

func (controller *OrbitController) onScroll(_, y float32) {
	controller.Move(Unit(y / 3))
}

func (controller *OrbitController) updatePosition() {
	direction := controller.rotation.ToDirectionVector()
	direction.Normalize()

	controller.camera.Position.X = controller.target.X + controller.distance*direction.X
	controller.camera.Position.Y = controller.target.Y + controller.distance*direction.Y
	controller.camera.Position.Z = controller.target.Z + controller.distance*direction.Z
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
			label.SetText(fmt.Sprintf("X: %.2f Y: %.2f Z: %.2f      Yaw: %.2f Pitch: %.2f Roll: %.2f",
				controller.camera.Position.X, controller.camera.Position.Y, controller.camera.Position.Z,
				controller.camera.Rotation.X, controller.camera.Rotation.Y, controller.camera.Rotation.Z))
			label.Refresh()
		}
	}()
	return label
}
