package ThreeDView

import (
	"image/color"
	"math"
)

func NewCube(size float64, position Point3D, rotation Point3D, color color.Color, w *ThreeDWidget) ThreeDShape {
	half := size / 2
	vertices := []Point3D{
		{X: -half, Y: -half, Z: -half},
		{X: half, Y: -half, Z: -half},
		{X: half, Y: half, Z: -half},
		{X: -half, Y: half, Z: -half},
		{X: -half, Y: -half, Z: half},
		{X: half, Y: -half, Z: half},
		{X: half, Y: half, Z: half},
		{X: -half, Y: half, Z: half},
	}
	faces := [][3]int{
		{0, 1, 2}, {0, 2, 3},
		{4, 5, 6}, {4, 6, 7},
		{0, 1, 5}, {0, 5, 4},
		{2, 3, 7}, {2, 7, 6},
		{0, 3, 7}, {0, 7, 4},
		{1, 2, 6}, {1, 6, 5},
	}

	var facesData = make([]FaceData, len(faces))
	for i := 0; i < len(faces); i++ {
		face := faces[i]
		p1 := vertices[face[0]]
		p2 := vertices[face[1]]
		p3 := vertices[face[2]]

		facesData[i] = FaceData{face: [3]Point3D{p1, p2, p3}, color: color}
	}

	return ThreeDShape{
		Faces:    facesData,
		Position: position,
		Rotation: rotation,
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

	var facesData = make([]FaceData, len(faces))
	for i := 0; i < len(faces); i++ {
		face := faces[i]
		p1 := vertices[face[0]]
		p2 := vertices[face[1]]
		p3 := vertices[face[2]]

		facesData[i] = FaceData{face: [3]Point3D{p1, p2, p3}, color: color}
	}

	return ThreeDShape{
		Faces:    facesData,
		Position: position,
		Rotation: rotation,
		widget:   w,
	}
}

func adjustColorBrightness(c color.RGBA, factor float64, stage int) color.RGBA {
	adjust := func(value uint8, factor float64) uint8 {
		newValue := float64(value)
		for i := 0; i < stage; i++ {
			newValue *= factor
		}
		if newValue > 255 {
			return 255
		} else if newValue < 0 {
			return 0
		}
		return uint8(newValue)
	}

	return color.RGBA{
		R: adjust(c.R, factor),
		G: adjust(c.G, factor),
		B: adjust(c.B, factor),
		A: c.A,
	}
}

type Rocket struct {
	ThreeDShape
	Stages int
	Radius float64
	Size   float64
}

func NewRocket(size float64, position Point3D, rotation Point3D, baseColor color.Color, w *ThreeDWidget, stages int, radius float64) Rocket {
	faces := buildRocketFaces(size, radius, baseColor, stages)

	return Rocket{
		ThreeDShape: ThreeDShape{
			Faces:    faces,
			Position: position,
			Rotation: rotation,
			widget:   w,
		},
		Stages: stages,
		Radius: radius,
		Size:   size,
	}
}

func (rocket *Rocket) RemoveStage() {
	if rocket.Stages > 1 {
		rocket.Stages--
		rocket.Faces = buildRocketFaces(rocket.Size, rocket.Radius, rocket.Faces[0].color, rocket.Stages)
	}
}

func buildRocketFaces(size float64, radius float64, baseColor color.Color, stages int) []FaceData {
	var vertices []Point3D
	var faces []FaceData

	tipHeight := size / 8
	for i := 0; i < 360; i += 10 {
		angle := degreesToRadians(float64(i))
		vertices = append(vertices, Point3D{
			X: radius * math.Cos(angle),
			Y: radius * math.Sin(angle),
			Z: -tipHeight / 2,
		})
	}
	tipVertexCount := len(vertices)
	vertices = append(vertices, Point3D{X: 0, Y: 0, Z: tipHeight / 2})
	for i := 0; i < tipVertexCount-1; i++ {
		p1 := vertices[i]
		p2 := vertices[i+1]
		p3 := vertices[tipVertexCount]
		faces = append(faces, FaceData{face: [3]Point3D{p1, p2, p3}, color: baseColor})
	}
	p1 := vertices[tipVertexCount-1]
	p2 := vertices[0]
	p3 := vertices[tipVertexCount]
	faces = append(faces, FaceData{face: [3]Point3D{p1, p2, p3}, color: baseColor})

	bodyHeight := size / 2
	baseColorRGBA := baseColor.(color.RGBA)
	isLightColor := (float64(baseColorRGBA.R)*299+float64(baseColorRGBA.G)*587+float64(baseColorRGBA.B)*114)/1000 > 128
	factor := 0.7
	if isLightColor {
		factor = 1.3
	}

	for stage := 0; stage < stages; stage++ {
		startIndex := len(vertices)
		for i := 0; i < 360; i += 10 {
			angle := degreesToRadians(float64(i))
			vertices = append(vertices, Point3D{
				X: radius * math.Cos(angle),
				Y: radius * math.Sin(angle),
				Z: -(float64(stage)*bodyHeight + tipHeight/2),
			})
			vertices = append(vertices, Point3D{
				X: radius * math.Cos(angle),
				Y: radius * math.Sin(angle),
				Z: -(float64(stage+1)*bodyHeight + tipHeight/2),
			})
		}
		bodyVertexCount := len(vertices) - startIndex

		stageColor := adjustColorBrightness(baseColorRGBA, factor, stage+1)

		for i := 0; i < bodyVertexCount-2; i += 2 {
			p1 := vertices[startIndex+i]
			p2 := vertices[startIndex+i+1]
			p3 := vertices[startIndex+i+2]
			faces = append(faces, FaceData{face: [3]Point3D{p1, p2, p3}, color: stageColor})
			p1 = vertices[startIndex+i+1]
			p2 = vertices[startIndex+i+3]
			p3 = vertices[startIndex+i+2]
			faces = append(faces, FaceData{face: [3]Point3D{p1, p2, p3}, color: stageColor})
		}
		p1 = vertices[startIndex+bodyVertexCount-2]
		p2 = vertices[startIndex+bodyVertexCount-1]
		p3 = vertices[startIndex]
		faces = append(faces, FaceData{face: [3]Point3D{p1, p2, p3}, color: stageColor})
		p1 = vertices[startIndex+bodyVertexCount-1]
		p2 = vertices[startIndex+1]
		p3 = vertices[startIndex]
		faces = append(faces, FaceData{face: [3]Point3D{p1, p2, p3}, color: stageColor})
	}
	return faces
}
