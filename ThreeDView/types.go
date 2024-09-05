package ThreeDView

import (
	"image/color"
	"math"
)

type Point3D struct {
	X, Y, Z float64
}

type Point2D struct {
	X, Y int64
}

type ThreeDShape struct {
	Vertices []Point3D
	Faces    [][3]int
	Rotation Point3D
	Position Point3D
	color    color.Color
	w        *ThreeDWidget
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

func (shape *ThreeDShape) GetFaces() []FaceData {
	faces := make([]FaceData, len(shape.Faces))
	for i, face := range shape.Faces {
		p1 := shape.Vertices[face[0]]
		p2 := shape.Vertices[face[1]]
		p3 := shape.Vertices[face[2]]

		p1 = rotateX(p1, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.X)
		p1 = rotateY(p1, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.Y)
		p1 = rotateZ(p1, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.Z)

		p2 = rotateX(p2, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.X)
		p2 = rotateY(p2, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.Y)
		p2 = rotateZ(p2, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.Z)

		p3 = rotateX(p3, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.X)
		p3 = rotateY(p3, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.Y)
		p3 = rotateZ(p3, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.Z)

		p1.X += shape.Position.X
		p1.Y += shape.Position.Y
		p1.Z += shape.Position.Z

		p2.X += shape.Position.X
		p2.Y += shape.Position.Y
		p2.Z += shape.Position.Z

		p3.X += shape.Position.X
		p3.Y += shape.Position.Y
		p3.Z += shape.Position.Z

		faces[i] = FaceData{face: [3]Point3D{p1, p2, p3}, color: shape.color, distance: 0}
	}
	return faces
}

type Camera struct {
	Position    Point3D
	FocalLength float64
	Scale       float64
	Pitch       float64
	Yaw         float64
	Roll        float64
}

func NewCamera(position Point3D, rotation Point3D, focalLength, scale float64) Camera {
	return Camera{Position: position, FocalLength: focalLength, Scale: scale, Pitch: rotation.X, Yaw: rotation.Y, Roll: rotation.Z}
}

func (camera *Camera) Project(point Point3D) Point2D {
	point = rotateX(point, camera.Position, camera.Pitch)
	point = rotateY(point, camera.Position, camera.Yaw)
	point = rotateZ(point, camera.Position, camera.Roll)

	translatedX := point.X - camera.Position.X
	translatedY := point.Y - camera.Position.Y
	translatedZ := point.Z - camera.Position.Z

	if translatedZ == 0 {
		translatedZ = 0.000001
	}

	x2D := (translatedX*camera.FocalLength/translatedZ)*camera.Scale + float64(Width)/2
	y2D := (translatedY*camera.FocalLength/translatedZ)*camera.Scale + float64(Height)/2

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
	camera.Yaw = math.Atan2(direction.X, direction.Z) * (180 / math.Pi)
}

func (camera *Camera) MoveForward(distance float64) {
	direction := Point3D{
		X: -math.Sin(degreesToRadians(camera.Yaw)),
		Y: math.Sin(degreesToRadians(camera.Pitch)),
		Z: -math.Cos(degreesToRadians(camera.Yaw)),
	}

	camera.Position.X += direction.X * distance
	camera.Position.Y += direction.Y * distance
	camera.Position.Z += direction.Z * distance
}
