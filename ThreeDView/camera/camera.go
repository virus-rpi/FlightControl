package camera

import (
	. "FlightControl/ThreeDView/types"
	"math"
)

// Camera represents a camera in 3D space
type Camera struct {
	Position   Point3D    // Camera position in world space in units
	Fov        Degrees    // Field of view in degrees
	Rotation   Rotation3D // Camera rotation in camera space
	Controller Controller // Camera Controller
}

// NewCamera creates a new camera at the given position in world space and rotation in camera space
func NewCamera(position Point3D, rotation Rotation3D) Camera {
	return Camera{Position: position, Rotation: rotation, Fov: 90}
}

// SetController sets the controller for the camera. It has to implement the Controller interface
func (camera *Camera) SetController(controller Controller) {
	camera.Controller = controller
	controller.setCamera(camera)
}

// Project projects a 3D point to a 2D point on the screen
func (camera *Camera) Project(point Point3D, width, height Pixel) Point2D {
	translatedPoint := point
	translatedPoint.Subtract(camera.Position)

	translatedPoint.Rotate(Point3D{}, camera.Rotation)

	epsilon := Unit(0.0001)
	if math.Abs(float64(translatedPoint.Z)) < float64(epsilon) {
		translatedPoint.Z = epsilon
	}

	fovRadians := camera.Fov.ToRadians()
	scale := Unit(float64(width) / (2 * math.Tan(float64(fovRadians/2))))

	x2D := (translatedPoint.X * scale / translatedPoint.Z) + Unit(width)/2
	y2D := (translatedPoint.Y * scale / translatedPoint.Z) + Unit(height)/2

	return Point2D{X: Pixel(x2D), Y: Pixel(y2D)}
}

// UnProject un-projects a 2D point on the screen to a 3D point in world space
func (camera *Camera) UnProject(point2d Point2D, distance Unit, width, height Pixel) Point3D {
	fovRadians := camera.Fov.ToRadians()
	halfWidth := float64(width) / 2
	halfHeight := float64(height) / 2
	scale := math.Tan(float64(fovRadians)/2) * float64(distance)

	pointInCameraSpace := Point3D{
		X: Unit((float64(point2d.X) - halfWidth) / halfWidth * scale),
		Y: Unit((float64(point2d.Y) - halfHeight) / halfHeight * scale),
		Z: -distance,
	}

	rotationMatrix := camera.Rotation.ToRotationMatrix()
	pointInWorldSpace := rotationMatrix.ApplyInverseRotationMatrix(pointInCameraSpace)
	pointInWorldSpace.Add(camera.Position)

	return pointInWorldSpace
}

// IsInFrustum checks if a point is in the camera's frustum
func (camera *Camera) IsInFrustum(point Point3D) bool {
	translatedPoint := point
	translatedPoint.Subtract(camera.Position)
	translatedPoint.Rotate(Point3D{}, camera.Rotation)

	fovRadians := camera.Fov.ToRadians()
	aspectRatio := 1.0
	tanFovOver2 := math.Tan(float64(fovRadians) / 2)

	if translatedPoint.Z < Unit(0.1) {
		return false
	}

	rightPlaneX := translatedPoint.Z * Unit(tanFovOver2*aspectRatio)
	if translatedPoint.X < -rightPlaneX || translatedPoint.X > rightPlaneX {
		return false
	}

	topPlaneY := translatedPoint.Z * Unit(tanFovOver2)
	if translatedPoint.Y < -topPlaneY || translatedPoint.Y > topPlaneY {
		return false
	}

	return true
}
