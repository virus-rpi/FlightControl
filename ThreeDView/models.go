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
		widget:   w,
	}
}

func NewPlane(size float64, position Point3D, rotation Point3D, color color.Color, w *ThreeDWidget, resolution int) ThreeDShape {
	half := size / 2
	step := size / float64(resolution)
	var vertices []Point3D
	for i := 0; i <= resolution; i++ {
		for j := 0; j <= resolution; j++ {
			vertices = append(vertices, Point3D{
				X: -half + float64(i)*step,
				Y: -half + float64(j)*step,
				Z: 0,
			})
		}
	}

	var faces [][3]int
	for i := 0; i < resolution; i++ {
		for j := 0; j < resolution; j++ {
			topLeft := i*(resolution+1) + j
			topRight := topLeft + 1
			bottomLeft := topLeft + (resolution + 1)
			bottomRight := bottomLeft + 1

			faces = append(faces, [3]int{topLeft, topRight, bottomRight})
			faces = append(faces, [3]int{topLeft, bottomRight, bottomLeft})
		}
	}

	return ThreeDShape{
		Vertices: vertices,
		Faces:    faces,
		Position: position,
		Rotation: rotation,
		color:    color,
		widget:   w,
	}
}
