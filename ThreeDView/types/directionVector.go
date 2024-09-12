package types

import "math"

// DirectionVector represents a vector in 3D space as a normalized vector
type DirectionVector struct {
	Point3D
}

// ToRotation converts a DirectionVector to a rotation in degree
func (point *DirectionVector) ToRotation() Rotation3D {
	return Rotation3D{
		X: Radians(math.Asin(float64(point.Y))).ToDegrees(),
		Y: Radians(math.Atan2(float64(point.X), float64(point.Z))).ToDegrees(),
	}
}
