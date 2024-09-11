package ThreeDView

import "math"

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

func (controller *OrbitController) PointAt(target Point3D) {
	direction := Point3D{
		X: target.X - controller.camera.Position.X,
		Y: target.Y - controller.camera.Position.Y,
		Z: target.Z - controller.camera.Position.Z,
	}

	length := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y + direction.Z*direction.Z)
	direction.X /= length
	direction.Y /= length
	direction.Z /= length

	controller.camera.Pitch = math.Asin(-direction.Y) * (180 / math.Pi)
	controller.camera.Yaw = math.Atan2(direction.X, direction.Z) * (-180 / math.Pi)

	up := Point3D{X: 0, Y: 0, Z: 1}
	right := Point3D{
		X: direction.Y*up.Z - direction.Z*up.Y,
		Y: direction.Z*up.X - direction.X*up.Z,
		Z: direction.X*up.Y - direction.Y*up.X,
	}
	rightLength := math.Sqrt(right.X*right.X + right.Y*right.Y + right.Z*right.Z)
	right.X /= rightLength
	right.Y /= rightLength
	right.Z /= rightLength
	correctedUp := Point3D{
		X: right.Y*direction.Z - right.Z*direction.Y,
		Y: right.Z*direction.X - right.X*direction.Z,
		Z: right.X*direction.Y - right.Y*direction.X,
	}
	controller.camera.Roll = math.Atan2(correctedUp.X, correctedUp.Y)*(-180/math.Pi) + 180
}

func (controller *OrbitController) MoveForward(distance float64) {
	direction := Point3D{
		X: controller.target.X - controller.camera.Position.X,
		Y: controller.target.Y - controller.camera.Position.Y,
		Z: controller.target.Z - controller.camera.Position.Z,
	}

	length := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y + direction.Z*direction.Z)
	direction.X /= length
	direction.Y /= length
	direction.Z /= length

	controller.camera.Position.X += direction.X * distance
	controller.camera.Position.Y += direction.Y * distance
	controller.camera.Position.Z += direction.Z * distance
}

func (controller *OrbitController) Rotate(yaw, pitch, roll float32) {
	controller.camera.Position.Rotate(controller.target, float64(pitch), float64(yaw), float64(roll))
}

func (controller *OrbitController) onDrag(x, y float32) {
	controller.Rotate(y/20, y, x/10)
	controller.PointAt(controller.target)
}

func (controller *OrbitController) onDragEnd() {}

func (controller *OrbitController) onScroll(_, y float32) {
	controller.MoveForward(float64(y / 3))
}
