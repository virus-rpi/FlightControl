package camera

import (
	. "FlightControl/ThreeDView/types"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

type ObjectInterface interface {
	GetPosition() Point3D
}

// OrbitController is a controller that allows the camera to orbit around a target Object
type OrbitController struct {
	BaseController
	target          ObjectInterface // The Object the camera is orbiting around in world space
	rotation        Rotation3D      // The rotation of the camera around the target in world space in degrees from the perspective of the target
	distance        Unit            // The distance of the camera from the target
	controlsEnabled bool            // Whether the controls are enabled (dragging, scrolling)
}

// NewOrbitController creates a new OrbitController with the target Object
func NewOrbitController(target ObjectInterface) *OrbitController {
	return &OrbitController{target: target, distance: 500, rotation: Rotation3D{Y: 300}, controlsEnabled: true}
}

func (controller *OrbitController) setCamera(camera *Camera) {
	controller.BaseController.camera = camera
	controller.Update()
}

// SetControlsEnabled sets whether the controls are enabled
func (controller *OrbitController) SetControlsEnabled(enabled bool) {
	controller.controlsEnabled = enabled
}

// SetTarget sets the target Object for the camera to orbit around
func (controller *OrbitController) SetTarget(target ObjectInterface) {
	controller.target = target
	controller.Update()
}

// SetDistance sets the distance of the camera from the target Object
func (controller *OrbitController) SetDistance(distance Unit) {
	controller.distance = distance
	controller.Update()
}

// Move moves the camera closer or further from the target Object by the given distance
func (controller *OrbitController) Move(distance Unit) {
	controller.distance += distance
	controller.Update()
}

// Rotate rotates the camera around the target Object by the given rotation
func (controller *OrbitController) Rotate(rotation Rotation3D) {
	controller.rotation.Add(rotation)
	controller.rotation.Normalize()
	controller.Update()
}

// SetRotation sets the rotation of the camera around the target Object
func (controller *OrbitController) SetRotation(rotation Rotation3D) {
	controller.rotation = rotation
	controller.rotation.Normalize()
	controller.Update()
}

// OnDrag is called when the user drags the camera. DO NOT CALL THIS FUNCTION MANUALLY
func (controller *OrbitController) OnDrag(x, y float32) {
	if !controller.controlsEnabled {
		return
	}
	controller.Rotate(Rotation3D{Y: Degrees(-y), Z: Degrees(x)})
}

// OnDragEnd is called when the user stops dragging the camera. DO NOT CALL THIS FUNCTION MANUALLY
func (controller *OrbitController) OnDragEnd() {}

// OnScroll is called when the user scrolls the camera. DO NOT CALL THIS FUNCTION MANUALLY
func (controller *OrbitController) OnScroll(_, y float32) {
	if !controller.controlsEnabled {
		return
	}
	controller.Move(Unit(y))
}

// Update updates the position and rotation of the camera. Call after changing the targets position
func (controller *OrbitController) Update() {
	controller.updatePosition()
	controller.pointAtTarget()
}

func (controller *OrbitController) updatePosition() {
	if controller.camera == nil {
		return
	}
	newPosition := controller.target.GetPosition()
	newPosition.Add(Point3D{X: controller.distance})
	newPosition.Rotate(controller.target.GetPosition(), controller.rotation)
	controller.camera.Position = newPosition
}

func (controller *OrbitController) pointAtTarget() {
	if controller.camera == nil {
		return
	}
	targetPosition := controller.target.GetPosition()
	direction := DirectionVector{Point3D: targetPosition}
	direction.Subtract(controller.camera.Position)
	direction.Normalize()
	rotation := direction.ToRotation()
	controller.camera.Rotation.X = -rotation.X
	controller.camera.Rotation.Y = -rotation.Y

	controller.camera.Rotation.Z = controller.rotation.Z - 90
}

// ManualController is a controller that allows the camera to be manually controlled. Useful for debugging
type ManualController struct {
	BaseController
}

// NewManualController creates a new ManualController
func NewManualController() *ManualController {
	return &ManualController{}
}

// GetRotationSlider returns a container with sliders for controlling the rotation of the camera
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

// GetPositionControl returns a container with sliders for controlling the position of the camera
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

// GetInfoLabel returns a label that displays the position and rotation of the camera
func (controller *ManualController) GetInfoLabel() *widget.Label {
	label := widget.NewLabel("X: 0 Y: 0 Z: 0      Yaw: 0 Pitch: 0 Roll: 0")
	go func() {
		ticker := time.NewTicker(time.Second / 30)
		defer ticker.Stop()
		for range ticker.C {
			label.SetText(fmt.Sprintf("X: %.2f Y: %.2f Z: %.2f      Yaw: %d Pitch: %d Roll: %d",
				controller.camera.Position.X, controller.camera.Position.Y, controller.camera.Position.Z,
				controller.camera.Rotation.X, controller.camera.Rotation.Y, controller.camera.Rotation.Z))
			label.Refresh()
		}
	}()
	return label
}
