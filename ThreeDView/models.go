package ThreeDView

import "image/color"

func NewCube(size float64, position Point3D, rotation Point3D, color color.Color, w *ThreeDWidget) ThreeDShape {
	half := size / 2
	return ThreeDShape{
		Vertices: []Point3D{
			{-half, -half, -half}, {half, -half, -half},
			{half, half, -half}, {-half, half, -half},
			{-half, -half, half}, {half, -half, half},
			{half, half, half}, {-half, half, half},
		},
		Faces: [][3]int{
			{0, 1, 2}, {0, 2, 3},
			{4, 5, 6}, {4, 6, 7},
			{0, 1, 5}, {0, 5, 4},
			{3, 2, 6}, {3, 6, 7},
			{0, 3, 7}, {0, 7, 4},
			{1, 2, 6}, {1, 6, 5},
		},
		Position: position,
		Rotation: rotation,
		color:    color,
		w:        w,
	}
}

func NewPlane(size float64, position Point3D, rotation Point3D, color color.Color, w *ThreeDWidget) ThreeDShape {
	gridSize := int(size / 10)
	half := size / 2
	numVertices := (gridSize + 1) * (gridSize + 1)
	vertices := make([]Point3D, numVertices)
	faces := make([][3]int, gridSize*gridSize*2)

	index := 0
	for i := 0; i <= gridSize; i++ {
		for j := 0; j <= gridSize; j++ {
			x := -half + float64(i)*10
			y := -half + float64(j)*10
			vertices[index] = Point3D{X: x, Y: y, Z: 0}
			index++
		}
	}

	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			index = i*(gridSize+1) + j
			faces[i*gridSize*2+j*2] = [3]int{index, index + 1, index + gridSize + 1}
			faces[i*gridSize*2+j*2+1] = [3]int{index + 1, index + gridSize + 2, index + gridSize + 1}
		}
	}

	return ThreeDShape{
		Vertices: vertices,
		Faces:    faces,
		Position: position,
		Rotation: rotation,
		color:    color,
		w:        w,
	}
}
