package types

import "math"

// Rotation3D represents a rotation in 3D space in degrees
type Rotation3D struct {
	Roll, Pitch, Yaw Degrees
}

// Minus negates the rotation in all axes and returns the negated rotation
func (rotation *Rotation3D) Minus() Rotation3D {
	return Rotation3D{-rotation.Roll, -rotation.Pitch, -rotation.Yaw}
}

// ToRotationMatrix creates a rotation matrix from the rotation
func (rotation *Rotation3D) ToRotationMatrix() RotationMatrix {
	rx := float64(rotation.Roll.ToRadians())
	ry := float64(rotation.Pitch.ToRadians())
	rz := float64(rotation.Yaw.ToRadians())

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

// Add adds another rotation to the rotation
func (rotation *Rotation3D) Add(other Rotation3D) {
	rotation.Roll += other.Roll
	rotation.Pitch += other.Pitch
	rotation.Yaw += other.Yaw
}

// ToDirectionVector converts the rotation to a normalized direction vector
func (rotation *Rotation3D) ToDirectionVector() DirectionVector {
	rotationMatrix := rotation.ToRotationMatrix()
	directionVector := DirectionVector{Point3D{
		X: Unit(rotationMatrix[0][2]),
		Y: Unit(rotationMatrix[1][2]),
		Z: Unit(rotationMatrix[2][2]),
	}}
	directionVector.Normalize()
	return directionVector
}

// Normalize normalizes the rotation to be within 0-360 degrees
func (rotation *Rotation3D) Normalize() {
	rotation.Roll = Degrees(math.Mod(float64(rotation.Roll), 360))
	rotation.Pitch = Degrees(math.Mod(float64(rotation.Pitch), 360))
	rotation.Yaw = Degrees(math.Mod(float64(rotation.Yaw), 360))
}
