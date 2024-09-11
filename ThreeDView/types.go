package ThreeDView

import (
	"image/color"
	"math"
)

type Point3D struct {
	X, Y, Z float64
}

func (point *Point3D) RotateX(pivot Point3D, degrees float64) {
	radians := degreesToRadians(degrees)
	translatedY := point.Y - pivot.Y
	translatedZ := point.Z - pivot.Z
	newY := translatedY*math.Cos(radians) - translatedZ*math.Sin(radians)
	newZ := translatedY*math.Sin(radians) + translatedZ*math.Cos(radians)
	point.Y = newY + pivot.Y
	point.Z = newZ + pivot.Z
}

func (point *Point3D) RotateY(pivot Point3D, degrees float64) {
	radians := degreesToRadians(degrees)
	translatedX := point.X - pivot.X
	translatedZ := point.Z - pivot.Z
	newX := translatedX*math.Cos(radians) + translatedZ*math.Sin(radians)
	newZ := -translatedX*math.Sin(radians) + translatedZ*math.Cos(radians)
	point.X = newX + pivot.X
	point.Z = newZ + pivot.Z
}

func (point *Point3D) RotateZ(pivot Point3D, degrees float64) {
	radians := degreesToRadians(degrees)
	translatedX := point.X - pivot.X
	translatedY := point.Y - pivot.Y
	newX := translatedX*math.Cos(radians) - translatedY*math.Sin(radians)
	newY := translatedX*math.Sin(radians) + translatedY*math.Cos(radians)
	point.X = newX + pivot.X
	point.Y = newY + pivot.Y
}

func (point *Point3D) Rotate(pivot Point3D, x, y, z float64) {
	point.RotateX(pivot, x)
	point.RotateY(pivot, y)
	point.RotateZ(pivot, z)
}

func (point *Point3D) Add(other Point3D) {
	point.X += other.X
	point.Y += other.Y
	point.Z += other.Z
}

func (point *Point3D) Distance(other Point3D) float64 {
	dx := point.X - other.X
	dy := point.Y - other.Y
	dz := point.Z - other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (point *Point3D) Magnitude() float64 {
	return math.Sqrt(point.X*point.X + point.Y*point.Y + point.Z*point.Z)
}

func (point *Point3D) Dot(other Point3D) float64 {
	return point.X*other.X + point.Y*other.Y + point.Z*other.Z
}

type Point2D struct {
	X, Y int64
}

func (point *Point2D) InBounds() bool {
	return point.X >= 0 && point.X < Width && point.Y >= 0 && point.Y < Height
}

type FaceData struct {
	face     [3]Point3D
	color    color.Color
	distance float64
}

type ProjectedFaceData struct {
	face     [3]Point2D
	color    color.Color
	distance float64
}

type ThreeDShape struct {
	faces    []FaceData
	Rotation Point3D
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

func (shape *ThreeDShape) RotateX(degrees float64) {
	shape.Rotation.X += degrees
}

func (shape *ThreeDShape) RotateY(degrees float64) {
	shape.Rotation.Y += degrees
}

func (shape *ThreeDShape) RotateZ(degrees float64) {
	shape.Rotation.Z += degrees
}

func (shape *ThreeDShape) Move(x, y, z float64) {
	shape.Position.X += x
	shape.Position.Y += y
	shape.Position.Z += z
}

type Camera struct {
	Position    Point3D
	FocalLength float64
	Fov         float64
	Pitch       float64
	Yaw         float64
	Roll        float64
	controller  CameraController
}

func NewCamera(position Point3D, rotation Point3D, focalLength float64) Camera {
	return Camera{Position: position, FocalLength: focalLength, Pitch: rotation.X, Yaw: rotation.Y, Roll: rotation.Z, Fov: 90}
}

func (camera *Camera) SetController(controller CameraController) {
	camera.controller = controller
	controller.setCamera(camera)
}

func (camera *Camera) Project(point Point3D) Point2D {
	point.RotateX(camera.Position, camera.Pitch)
	point.RotateY(camera.Position, camera.Yaw)
	point.RotateZ(camera.Position, camera.Roll)

	translatedX := point.X - camera.Position.X
	translatedY := point.Y - camera.Position.Y
	translatedZ := point.Z - camera.Position.Z

	if translatedZ == 0 {
		translatedZ = 0.000001
	}

	fovRadians := degreesToRadians(camera.Fov)
	scale := camera.FocalLength / math.Tan(fovRadians/2)

	x2D := (translatedX * scale / translatedZ) + float64(Width)/2
	y2D := (translatedY * scale / translatedZ) + float64(Height)/2

	x2D = math.Max(math.Min(x2D, float64(math.MaxInt64)), float64(math.MinInt64))
	y2D = math.Max(math.Min(y2D, float64(math.MaxInt64)), float64(math.MinInt64))

	return Point2D{int64(x2D), int64(y2D)}
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
