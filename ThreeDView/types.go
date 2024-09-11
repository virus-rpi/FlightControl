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
	Faces    []FaceData
	Rotation Point3D
	Position Point3D
	widget   *ThreeDWidget
}

func (shape *ThreeDShape) GetFaces() []FaceData {
	faces := make([]FaceData, len(shape.Faces))
	for i, face := range shape.Faces {
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
	Pitch       float64
	Yaw         float64
	Roll        float64
	OrbitCenter Point3D
	Fov         float64
}

func NewCamera(position Point3D, rotation Point3D, focalLength float64) Camera {
	return Camera{Position: position, FocalLength: focalLength, Pitch: rotation.X, Yaw: rotation.Y, Roll: rotation.Z, Fov: 90}
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

func (camera *Camera) PointAt(target Point3D) {
	direction := Point3D{
		X: target.X - camera.Position.X,
		Y: target.Y - camera.Position.Y,
		Z: target.Z - camera.Position.Z,
	}

	length := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y + direction.Z*direction.Z)
	direction.X /= length
	direction.Y /= length
	direction.Z /= length

	camera.Pitch = math.Asin(-direction.Y) * (180 / math.Pi)
	camera.Yaw = math.Atan2(direction.X, direction.Z) * (-180 / math.Pi)

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
	camera.Roll = math.Atan2(correctedUp.X, correctedUp.Y)*(-180/math.Pi) + 180
}

func (camera *Camera) MoveForward(distance float64) {
	direction := Point3D{
		X: math.Cos(degreesToRadians(camera.Pitch)) * math.Sin(degreesToRadians(camera.Yaw)),
		Y: math.Sin(degreesToRadians(camera.Pitch)),
		Z: math.Cos(degreesToRadians(camera.Pitch)) * math.Cos(degreesToRadians(camera.Yaw)),
	}

	camera.Position.X += direction.X * distance
	camera.Position.Y += direction.Y * distance
	camera.Position.Z += direction.Z * distance
}

func (camera *Camera) SetOrbitCenter(center Point3D) {
	camera.OrbitCenter = center
}

func (camera *Camera) IsPointInFrustum(point Point3D) bool {
	return true
}
