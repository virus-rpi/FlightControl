package ThreeDView

import (
	"image/color"
	"math"
)

type Degrees float64

func (degrees Degrees) ToRadians() Radians {
	return Radians(degrees * math.Pi / 180)
}

type Radians float64

func (radians Radians) ToDegrees() Degrees {
	return Degrees(radians * 180 / math.Pi)
}

type Rotation3D struct {
	X, Y, Z Degrees
}

type Unit float64

type Pixel int64

type Point3D struct {
	X, Y, Z Unit
}

func (point *Point3D) RotateX(pivot Point3D, degrees Degrees) {
	radians := degrees.ToRadians()
	translatedY := point.Y - pivot.Y
	translatedZ := point.Z - pivot.Z
	newY := Unit(float64(translatedY)*math.Cos(float64(radians)) - float64(translatedZ)*math.Sin(float64(radians)))
	newZ := Unit(float64(translatedY)*math.Sin(float64(radians)) + float64(translatedZ)*math.Cos(float64(radians)))
	point.Y = newY + pivot.Y
	point.Z = newZ + pivot.Z
}

func (point *Point3D) RotateY(pivot Point3D, degrees Degrees) {
	radians := degrees.ToRadians()
	translatedX := point.X - pivot.X
	translatedZ := point.Z - pivot.Z
	newX := Unit(float64(translatedX)*math.Cos(float64(radians)) + float64(translatedZ)*math.Sin(float64(radians)))
	newZ := Unit(float64(-translatedX)*math.Sin(float64(radians)) + float64(translatedZ)*math.Cos(float64(radians)))
	point.X = newX + pivot.X
	point.Z = newZ + pivot.Z
}

func (point *Point3D) RotateZ(pivot Point3D, degrees Degrees) {
	radians := degrees.ToRadians()
	translatedX := point.X - pivot.X
	translatedY := point.Y - pivot.Y
	newX := Unit(float64(translatedX)*math.Cos(float64(radians)) - float64(translatedY)*math.Sin(float64(radians)))
	newY := Unit(float64(translatedX)*math.Sin(float64(radians)) + float64(translatedY)*math.Cos(float64(radians)))
	point.X = newX + pivot.X
	point.Y = newY + pivot.Y
}

func (point *Point3D) Rotate(pivot Point3D, x, y, z Degrees) {
	point.RotateX(pivot, x)
	point.RotateY(pivot, y)
	point.RotateZ(pivot, z)
}

func (point *Point3D) Add(other Point3D) {
	point.X += other.X
	point.Y += other.Y
	point.Z += other.Z
}

func (point *Point3D) Subtract(other Point3D) {
	point.X -= other.X
	point.Y -= other.Y
	point.Z -= other.Z
}

func (point *Point3D) Distance(other Point3D) Unit {
	dx := point.X - other.X
	dy := point.Y - other.Y
	dz := point.Z - other.Z
	return Unit(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

func (point *Point3D) Magnitude() Unit {
	return Unit(math.Sqrt(float64(point.X*point.X + point.Y*point.Y + point.Z*point.Z)))
}

func (point *Point3D) Dot(other Point3D) Unit {
	return point.X*other.X + point.Y*other.Y + point.Z*other.Z
}

type Point2D struct {
	X, Y Pixel
}

func (point *Point2D) InBounds() bool {
	return point.X >= 0 && point.X < Width && point.Y >= 0 && point.Y < Height
}

type FaceData struct {
	face     [3]Point3D
	color    color.Color
	distance Unit
}

type ProjectedFaceData struct {
	face     [3]Point2D
	color    color.Color
	distance Unit
}

type ThreeDShape struct {
	faces    []FaceData
	Rotation Rotation3D
	Position Point3D
	widget   *ThreeDWidget
}

func (shape *ThreeDShape) GetFaces() []FaceData {
	faces := make([]FaceData, len(shape.faces))
	for i, face := range shape.faces {
		p1 := face.face[0]
		p2 := face.face[1]
		p3 := face.face[2]

		p1.Rotate(Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.X, shape.Rotation.Y, shape.Rotation.Z)
		p2.Rotate(Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.X, shape.Rotation.Y, shape.Rotation.Z)
		p3.Rotate(Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.X, shape.Rotation.Y, shape.Rotation.Z)

		p1.Add(shape.Position)
		p2.Add(shape.Position)
		p3.Add(shape.Position)

		faces[i] = FaceData{face: [3]Point3D{p1, p2, p3}, color: face.color, distance: 0}
	}
	return faces
}

func (shape *ThreeDShape) RotateX(degrees Degrees) {
	shape.Rotation.X += degrees
}

func (shape *ThreeDShape) RotateY(degrees Degrees) {
	shape.Rotation.Y += degrees
}

func (shape *ThreeDShape) RotateZ(degrees Degrees) {
	shape.Rotation.Z += degrees
}

func (shape *ThreeDShape) Move(x, y, z Unit) {
	shape.Position.X += x
	shape.Position.Y += y
	shape.Position.Z += z
}

type Camera struct {
	Position   Point3D
	Fov        Degrees
	Pitch      Degrees
	Yaw        Degrees
	Roll       Degrees
	controller CameraController
}

func NewCamera(position Point3D, rotation Rotation3D) Camera {
	return Camera{Position: position, Pitch: rotation.X, Yaw: rotation.Y, Roll: rotation.Z, Fov: 90}
}

func (camera *Camera) SetController(controller CameraController) {
	camera.controller = controller
	controller.setCamera(camera)
}

func (camera *Camera) Project(point Point3D) Point2D {
	translatedPoint := Point3D{
		X: point.X - camera.Position.X,
		Y: point.Y - camera.Position.Y,
		Z: point.Z - camera.Position.Z,
	}

	translatedPoint.Rotate(camera.Position, -camera.Yaw, -camera.Pitch, -camera.Roll)

	epsilon := 1e-6
	if math.Abs(float64(translatedPoint.Z)) < epsilon {
		translatedPoint.Z = Unit(epsilon)
	}

	fovRadians := camera.Fov.ToRadians()
	scale := Unit(float64(Width) / (2 * math.Tan(float64(fovRadians/2))))

	x2D := (translatedPoint.X * scale / translatedPoint.Z) + Unit(Width)/2
	y2D := (translatedPoint.Y * scale / translatedPoint.Z) + Unit(Height)/2

	return Point2D{Pixel(x2D), Pixel(y2D)}
}

func (camera *Camera) UnProject(point2d Point2D, distance Unit) Point3D {
	x := Unit((float64(point2d.X) - float64(Width)/2) / (float64(Width) / 2))
	y := Unit((float64(point2d.Y) - float64(Height)/2) / (float64(Height) / 2))

	fovRadians := camera.Fov.ToRadians()
	scale := Unit(math.Tan(float64(fovRadians)/2) * float64(distance))

	worldX := x * scale
	worldY := y * scale
	worldZ := -distance

	pointInCameraSpace := Point3D{X: worldX, Y: worldY, Z: worldZ}

	pointInCameraSpace.Rotate(camera.Position, camera.Yaw, camera.Pitch, camera.Roll)

	pointInWorldSpace := pointInCameraSpace
	pointInWorldSpace.Add(camera.Position)

	return pointInWorldSpace
}

// IsPointInFrustum TODO: Implement this function
func (camera *Camera) IsPointInFrustum(point Point3D) bool {
	return true
}

type CameraController interface {
	setCamera(*Camera)
}

type BaseController struct {
	camera *Camera
}

func (controller *BaseController) setCamera(camera *Camera) {
	controller.camera = camera
}

type DragController interface {
	onDrag(float32, float32)
	onDragEnd()
}

type ScrollController interface {
	onScroll(float32, float32)
}
