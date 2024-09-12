package types

import (
	"math"
)

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

// Normalize normalizes the point
func (point *Point3D) Normalize() {
	magnitude := point.Magnitude()
	if magnitude == 0 {
		return
	}
	point.X /= magnitude
	point.Y /= magnitude
	point.Z /= magnitude
}

// Cross returns the cross product of the point with another point
func (point *Point3D) Cross(other Point3D) Point3D {
	return Point3D{
		X: point.Y*other.Z - point.Z*other.Y,
		Y: point.Z*other.X - point.X*other.Z,
		Z: point.X*other.Y - point.Y*other.X,
	}
}
