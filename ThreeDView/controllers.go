package ThreeDView

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"math"
	"time"
)

type OrbitController struct {
	BaseController
	target Point3D
}

func NewOrbitController(orbitCenter Point3D) *OrbitController {
	return &OrbitController{target: orbitCenter}
}

func (controller *OrbitController) SetTarget(center Point3D) {
	controller.target = center
}

func (controller *OrbitController) MoveForward(distance Unit) {
	direction := Point3D{
		X: controller.target.X - controller.camera.Position.X,
		Y: controller.target.Y - controller.camera.Position.Y,
		Z: controller.target.Z - controller.camera.Position.Z,
	}

	length := Unit(math.Sqrt(float64(direction.X*direction.X + direction.Y*direction.Y + direction.Z*direction.Z)))
	direction.X /= length
	direction.Y /= length
	direction.Z /= length

	controller.camera.Position.X += direction.X * distance
	controller.camera.Position.Y += direction.Y * distance
	controller.camera.Position.Z += direction.Z * distance
}

func (controller *OrbitController) PointAtTarget() {
	direction := controller.target
	direction.Subtract(controller.camera.Position)

	controller.camera.Rotation.X = Radians(math.Atan2(float64(direction.X), float64(direction.Z))).ToDegrees()

	distanceXZ := math.Sqrt(float64(direction.X*direction.X + direction.Z*direction.Z))

	controller.camera.Rotation.Y = Radians(math.Atan2(float64(direction.Y), distanceXZ)).ToDegrees()
	controller.camera.Rotation.Z = 0
}

func (controller *OrbitController) Rotate(yaw, pitch, roll Degrees) {
	yawRadians := yaw.ToRadians()
	pitchRadians := pitch.ToRadians()
	rollRadians := roll.ToRadians()

	qYaw := Quaternion{
		W: math.Cos(float64(yawRadians / 2)),
		X: 0,
		Y: math.Sin(float64(yawRadians / 2)),
		Z: 0,
	}

	qPitch := Quaternion{
		W: math.Cos(float64(pitchRadians / 2)),
		X: math.Sin(float64(pitchRadians / 2)),
		Y: 0,
		Z: 0,
	}

	qRoll := Quaternion{
		W: math.Cos(float64(rollRadians / 2)),
		X: 0,
		Y: 0,
		Z: math.Sin(float64(rollRadians / 2)),
	}

	combined := qYaw.Multiply(qPitch).Multiply(qRoll)
	rotationMatrix := combined.ToRotationMatrix()
	controller.camera.Position = applyRotationMatrix(controller.camera.Position, controller.target, rotationMatrix)
	controller.PointAtTarget()
}

func (controller *OrbitController) onDrag(x, y float32) {
	controller.Rotate(Degrees(x/10), Degrees(y/10), 0)
}

func (controller *OrbitController) onDragEnd() {}

func (controller *OrbitController) onScroll(_, y float32) {
	controller.MoveForward(Unit(y / 3))
}

type Quaternion struct {
	W, X, Y, Z float64
}

func (q Quaternion) Multiply(q2 Quaternion) Quaternion {
	return Quaternion{
		W: q.W*q2.W - q.X*q2.X - q.Y*q2.Y - q.Z*q2.Z,
		X: q.W*q2.X + q.X*q2.W + q.Y*q2.Z - q.Z*q2.Y,
		Y: q.W*q2.Y - q.X*q2.Z + q.Y*q2.W + q.Z*q2.X,
		Z: q.W*q2.Z + q.X*q2.Y - q.Y*q2.X + q.Z*q2.W,
	}
}

func (q Quaternion) ToRotationMatrix() [3][3]float64 {
	return [3][3]float64{
		{1 - 2*q.Y*q.Y - 2*q.Z*q.Z, 2*q.X*q.Y - 2*q.Z*q.W, 2*q.X*q.Z + 2*q.Y*q.W},
		{2*q.X*q.Y + 2*q.Z*q.W, 1 - 2*q.X*q.X - 2*q.Z*q.Z, 2*q.Y*q.Z - 2*q.X*q.W},
		{2*q.X*q.Z - 2*q.Y*q.W, 2*q.Y*q.Z + 2*q.X*q.W, 1 - 2*q.X*q.X - 2*q.Y*q.Y},
	}
}

func applyRotationMatrix(point, center Point3D, rotationMatrix [3][3]float64) Point3D {
	return Point3D{
		X: Unit(rotationMatrix[0][0]*float64(point.X-center.X) + rotationMatrix[0][1]*float64(point.Y-center.Y) + rotationMatrix[0][2]*float64(point.Z-center.Z) + float64(center.X)),
		Y: Unit(rotationMatrix[1][0]*float64(point.X-center.X) + rotationMatrix[1][1]*float64(point.Y-center.Y) + rotationMatrix[1][2]*float64(point.Z-center.Z) + float64(center.Y)),
		Z: Unit(rotationMatrix[2][0]*float64(point.X-center.X) + rotationMatrix[2][1]*float64(point.Y-center.Y) + rotationMatrix[2][2]*float64(point.Z-center.Z) + float64(center.Z)),
	}
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
