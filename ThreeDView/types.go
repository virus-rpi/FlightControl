package ThreeDView

import (
	"image/color"
	"math"
	"sync"
)

// Degrees represents an angle in degrees
type Degrees float64

// ToRadians converts degrees to radians
func (degrees Degrees) ToRadians() Radians {
	return Radians(degrees * math.Pi / 180)
}

// Radians represents an angle in radians
type Radians float64

// ToDegrees converts radians to degrees
func (radians Radians) ToDegrees() Degrees {
	return Degrees(radians * 180 / math.Pi)
}

// Rotation3D represents a rotation in 3D space
type Rotation3D struct {
	X, Y, Z Degrees
}

// Minus negates the rotation in all axes and returns the negated rotation
func (rotation *Rotation3D) Minus() Rotation3D {
	return Rotation3D{-rotation.X, -rotation.Y, -rotation.Z}
}

// CreateRotationMatrix creates a rotation matrix from the rotation
func (rotation *Rotation3D) CreateRotationMatrix() RotationMatrix {
	rx := float64(rotation.X.ToRadians())
	ry := float64(rotation.Y.ToRadians())
	rz := float64(rotation.Z.ToRadians())

	cosX, sinX := math.Cos(rx), math.Sin(rx)
	cosY, sinY := math.Cos(ry), math.Sin(ry)
	cosZ, sinZ := math.Cos(rz), math.Sin(rz)

	return RotationMatrix{
		{
			cosY * cosZ,
			cosY * sinZ,
			-sinY,
		},
		{
			sinX*sinY*cosZ - cosX*sinZ,
			sinX*sinY*sinZ + cosX*cosZ,
			sinX * cosY,
		},
		{
			cosX*sinY*cosZ + sinX*sinZ,
			cosX*sinY*sinZ - sinX*cosZ,
			cosX * cosY,
		},
	}
}

// RotationMatrix represents a 3x3 rotation matrix
type RotationMatrix [3][3]float64

// ApplyInverseRotationMatrix applies the inverse of the rotation matrix to a point
func (rotationMatrix *RotationMatrix) ApplyInverseRotationMatrix(point Point3D) Point3D {
	return Point3D{
		X: Unit(rotationMatrix[0][0]*float64(point.X) + rotationMatrix[0][1]*float64(point.Y) + rotationMatrix[0][2]*float64(point.Z)),
		Y: Unit(rotationMatrix[1][0]*float64(point.X) + rotationMatrix[1][1]*float64(point.Y) + rotationMatrix[1][2]*float64(point.Z)),
		Z: Unit(rotationMatrix[2][0]*float64(point.X) + rotationMatrix[2][1]*float64(point.Y) + rotationMatrix[2][2]*float64(point.Z)),
	}
}

// Unit is the unit for distance in 3D space
type Unit float64

// Pixel is the unit for distance in 2D space
type Pixel int64

// Point3D represents a point in 3D space
type Point3D struct {
	X, Y, Z Unit
}

// RotateX rotates the point around a pivot point by the given rotation in the X axis
func (point *Point3D) RotateX(pivot Point3D, degrees Degrees) {
	radians := degrees.ToRadians()
	translatedY := point.Y - pivot.Y
	translatedZ := point.Z - pivot.Z
	newY := Unit(float64(translatedY)*math.Cos(float64(radians)) - float64(translatedZ)*math.Sin(float64(radians)))
	newZ := Unit(float64(translatedY)*math.Sin(float64(radians)) + float64(translatedZ)*math.Cos(float64(radians)))
	point.Y = newY + pivot.Y
	point.Z = newZ + pivot.Z
}

// RotateY rotates the point around a pivot point by the given rotation in the Y axis
func (point *Point3D) RotateY(pivot Point3D, degrees Degrees) {
	radians := degrees.ToRadians()
	translatedX := point.X - pivot.X
	translatedZ := point.Z - pivot.Z
	newX := Unit(float64(translatedX)*math.Cos(float64(radians)) + float64(translatedZ)*math.Sin(float64(radians)))
	newZ := Unit(float64(-translatedX)*math.Sin(float64(radians)) + float64(translatedZ)*math.Cos(float64(radians)))
	point.X = newX + pivot.X
	point.Z = newZ + pivot.Z
}

// RotateZ rotates the point around a pivot point by the given rotation in the Z axis
func (point *Point3D) RotateZ(pivot Point3D, degrees Degrees) {
	radians := degrees.ToRadians()
	translatedX := point.X - pivot.X
	translatedY := point.Y - pivot.Y
	newX := Unit(float64(translatedX)*math.Cos(float64(radians)) - float64(translatedY)*math.Sin(float64(radians)))
	newY := Unit(float64(translatedX)*math.Sin(float64(radians)) + float64(translatedY)*math.Cos(float64(radians)))
	point.X = newX + pivot.X
	point.Y = newY + pivot.Y
}

// Rotate rotates the point around a pivot point by the given rotation
func (point *Point3D) Rotate(pivot Point3D, rotation Rotation3D) {
	point.RotateX(pivot, rotation.X)
	point.RotateY(pivot, rotation.Y)
	point.RotateZ(pivot, rotation.Z)
}

// Add adds another point to the point
func (point *Point3D) Add(other Point3D) {
	point.X += other.X
	point.Y += other.Y
	point.Z += other.Z
}

// Subtract subtracts another point from the point
func (point *Point3D) Subtract(other Point3D) {
	point.X -= other.X
	point.Y -= other.Y
	point.Z -= other.Z
}

// DistanceTo returns the distance between the point and another point
func (point *Point3D) DistanceTo(other Point3D) Unit {
	dx := point.X - other.X
	dy := point.Y - other.Y
	dz := point.Z - other.Z
	return Unit(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

// Magnitude returns the magnitude of the point (distance from origin)
func (point *Point3D) Magnitude() Unit {
	return Unit(math.Sqrt(float64(point.X*point.X + point.Y*point.Y + point.Z*point.Z)))
}

// Dot returns the dot product of the point with another point
func (point *Point3D) Dot(other Point3D) Unit {
	return point.X*other.X + point.Y*other.Y + point.Z*other.Z
}

// Point2D represents a point in 2D space
type Point2D struct {
	X, Y Pixel
}

// InBounds returns true if the point is in bounds of the screen
func (point *Point2D) InBounds() bool {
	return point.X >= 0 && point.X < Width && point.Y >= 0 && point.Y < Height
}

// Face represents a face in 3D space as 3 3D points
type Face [3]Point3D

// Rotate rotates the face around a pivot point by the given rotation
func (face *Face) Rotate(pivot Point3D, rotation Rotation3D) {
	face[0].Rotate(pivot, rotation)
	face[1].Rotate(pivot, rotation)
	face[2].Rotate(pivot, rotation)
}

// Add adds another point to the face
func (face *Face) Add(other Point3D) {
	face[0].Add(other)
	face[1].Add(other)
	face[2].Add(other)
}

// DistanceTo returns the distance between the face and a point
func (face *Face) DistanceTo(point Point3D) Unit {
	return (face[0].DistanceTo(point) + face[1].DistanceTo(point) + face[2].DistanceTo(point)) / 3
}

// FaceData represents a face in 3D space
type FaceData struct {
	face     Face        // The face in 3D space as a Face
	color    color.Color // The color of the face
	distance Unit        // The distance of the face from the camera 3d world space
}

// ProjectedFaceData represents a face projected to 2D space
type ProjectedFaceData struct {
	face     [3]Point2D  // The face in 2D space as 3 2d points
	color    color.Color // The color of the face
	distance Unit        // The distance of the un-projected face from the camera in 3d world space
}

// ThreeDShape represents a 3D shape in world space
type ThreeDShape struct {
	faces    []FaceData    // Faces of the shape in local space
	Rotation Rotation3D    // Rotation of the shape in world space
	Position Point3D       // Position of the shape in world space
	widget   *ThreeDWidget // The widget the shape is in
}

// GetFaces returns the faces of the shape in world space as FaceData
func (shape *ThreeDShape) GetFaces() []FaceData {
	faces := make([]FaceData, len(shape.faces))
	var wg sync.WaitGroup
	wg.Add(len(shape.faces))
	for i, face := range shape.faces {
		go func(i int, face FaceData) {
			defer wg.Done()
			face.face.Rotate(Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation)
			face.face.Add(shape.Position)
			face.distance = face.face.DistanceTo(shape.widget.camera.Position)
			faces[i] = face
		}(i, face)
	}
	wg.Wait()
	return faces
}

// Camera represents a camera in 3D space
type Camera struct {
	Position   Point3D          // Camera position in world space in units
	Fov        Degrees          // Field of view in degrees
	Rotation   Rotation3D       // Camera rotation in camera space
	controller CameraController // Camera controller
}

// NewCamera creates a new camera at the given position in world space and rotation in camera space
func NewCamera(position Point3D, rotation Rotation3D) Camera {
	return Camera{Position: position, Rotation: rotation, Fov: 90}
}

// SetController sets the controller for the camera. It has to implement the CameraController interface
func (camera *Camera) SetController(controller CameraController) {
	camera.controller = controller
	controller.setCamera(camera)
}

// Project projects a 3D point to a 2D point on the screen
func (camera *Camera) Project(point Point3D) Point2D {
	translatedPoint := point
	translatedPoint.Subtract(camera.Position)

	translatedPoint.Rotate(Point3D{}, camera.Rotation)

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

// UnProject un-projects a 2D point on the screen to a 3D point in world space
func (camera *Camera) UnProject(point2d Point2D, distance Unit) Point3D {
	fovRadians := camera.Fov.ToRadians()
	scale := Unit(math.Tan(float64(fovRadians)/2) * float64(distance))

	pointInCameraSpace := Point3D{
		X: Unit((float64(point2d.X)-float64(Width)/2)/(float64(Width)/2)) * scale,
		Y: Unit((float64(point2d.Y)-float64(Height)/2)/(float64(Height)/2)) * scale,
		Z: -distance,
	}

	rotationMatrix := camera.Rotation.CreateRotationMatrix()

	pointInWorldSpace := rotationMatrix.ApplyInverseRotationMatrix(pointInCameraSpace)
	pointInWorldSpace.Add(camera.Position)

	return pointInWorldSpace
}

// CameraController is an interface for camera controllers to implement
type CameraController interface {
	setCamera(*Camera)
}

// BaseController is a base controller for camera controllers
type BaseController struct {
	camera *Camera
}

// setCamera sets the camera for the controller
func (controller *BaseController) setCamera(camera *Camera) {
	controller.camera = camera
}

// DragController is an interface for CameraController that supports dragging
type DragController interface {
	onDrag(float32, float32)
	onDragEnd()
}

// ScrollController is an interface for CameraController that supports scrolling
type ScrollController interface {
	onScroll(float32, float32)
}
