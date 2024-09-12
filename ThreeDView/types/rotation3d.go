package types

import "math"

// Rotation3D represents a rotation in 3D space
type Rotation3D struct {
	X, Y, Z Degrees
}

// Minus negates the rotation in all axes and returns the negated rotation
func (rotation *Rotation3D) Minus() Rotation3D {
	return Rotation3D{-rotation.X, -rotation.Y, -rotation.Z}
}

// ToRotationMatrix creates a rotation matrix from the rotation
func (rotation *Rotation3D) ToRotationMatrix() RotationMatrix {
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

// Add adds another rotation to the rotation
func (rotation *Rotation3D) Add(other Rotation3D) {
	rotation.X += other.X
	rotation.Y += other.Y
	rotation.Z += other.Z
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
	rotation.X = Degrees(math.Mod(float64(rotation.X), 360))
	rotation.Y = Degrees(math.Mod(float64(rotation.Y), 360))
	rotation.Z = Degrees(math.Mod(float64(rotation.Z), 360))
}
