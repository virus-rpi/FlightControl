package ThreeDView

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"math"
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
	controller.camera.Rotation = controller.rotation
}

func (controller *OrbitController) Rotate(rotation Rotation3D) {
	controller.rotation.Add(rotation)
	controller.rotation.Normalize()
	controller.updatePosition()
}

func (controller *OrbitController) onDrag(x, y float32) {
	// TODO: drag x: always rotate around the GLOBAL z axis
	// TODO: drag y: rotate around the GLOBAL x and y axis based on the GLOBAL y axis and resulting percentage of x and y on how much they contribute to the rotation in the direction up/down in the viewport
	log.Println("drag", x, y)
	upward := Point3D{X: 0, Y: 1, Z: 0}
	upward.Normalize()

	totalUpward := math.Abs(float64(upward.X)) + math.Abs(float64(upward.Y))
	percentUpwardX := math.Abs(float64(upward.X)) / totalUpward
	percentUpwardY := math.Abs(float64(upward.Y)) / totalUpward

	log.Println("u", percentUpwardX, percentUpwardY)

	rotation := Rotation3D{
		X: Degrees(float64(y) * percentUpwardX),
		Y: Degrees(float64(y) * percentUpwardY),
		Z: Degrees(x),
	}

	log.Println(rotation, controller.rotation)

	controller.Rotate(rotation)
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

	controller.PointAtTarget()
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
