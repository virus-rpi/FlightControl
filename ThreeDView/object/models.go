package object

import (
	. "FlightControl/ThreeDView/types"
	"image/color"
	"math"
)

func NewCube(size Unit, position Point3D, rotation Rotation3D, color color.Color, w ThreeDWidgetInterface) *Object {
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

		facesData[i] = FaceData{Face: [3]Point3D{p1, p2, p3}, Color: color}
	}

	cube := Object{
		Faces:    facesData,
		Position: position,
		Rotation: rotation,
		Widget:   w,
	}
	w.AddObject(&cube)
	return &cube
}

func NewPlane(size Unit, position Point3D, rotation Rotation3D, color color.Color, w ThreeDWidgetInterface, resolution int) *Object {
	half := size / 2
	step := size / Unit(resolution)
	var vertices []Point3D
	for i := 0; i <= resolution; i++ {
		for j := 0; j <= resolution; j++ {
			vertices = append(vertices, Point3D{
				X: -half + Unit(i)*step,
				Y: -half + Unit(j)*step,
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

		facesData[i] = FaceData{Face: [3]Point3D{p1, p2, p3}, Color: color}
	}

	plane := Object{
		Faces:    facesData,
		Position: position,
		Rotation: rotation,
		Widget:   w,
	}
	w.AddObject(&plane)
	return &plane
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
	Object
	Stages int
	Radius Unit
	Size   Unit
}

func NewRocket(size Unit, position Point3D, rotation Rotation3D, baseColor color.Color, w ThreeDWidgetInterface, stages int, radius Unit) *Rocket {
	faces := buildRocketFaces(size, radius, baseColor, stages)

	rocket := Rocket{
		Object: Object{
			Faces:    faces,
			Position: position,
			Rotation: rotation,
			Widget:   w,
		},
		Stages: stages,
		Radius: radius,
		Size:   size,
	}
	w.AddObject(&rocket.Object)
	return &rocket
}

func (rocket *Rocket) RemoveStage() {
	if rocket.Stages > 1 {
		rocket.Stages--
		rocket.Faces = buildRocketFaces(rocket.Size, rocket.Radius, rocket.Faces[0].Color, rocket.Stages)
	}
}

func buildRocketFaces(size Unit, radius Unit, baseColor color.Color, stages int) []FaceData {
	var vertices []Point3D
	var faces []FaceData

	tipHeight := size / 8
	for i := 0; i < 360; i += 10 {
		angle := Degrees(i).ToRadians()
		vertices = append(vertices, Point3D{
			X: radius * Unit(math.Cos(float64(angle))),
			Y: radius * Unit(math.Sin(float64(angle))),
			Z: -tipHeight / 2,
		})
	}
	tipVertexCount := len(vertices)
	vertices = append(vertices, Point3D{X: 0, Y: 0, Z: tipHeight / 2})
	for i := 0; i < tipVertexCount-1; i++ {
		p1 := vertices[i]
		p2 := vertices[i+1]
		p3 := vertices[tipVertexCount]
		faces = append(faces, FaceData{Face: [3]Point3D{p1, p2, p3}, Color: baseColor})
	}
	p1 := vertices[tipVertexCount-1]
	p2 := vertices[0]
	p3 := vertices[tipVertexCount]
	faces = append(faces, FaceData{Face: [3]Point3D{p1, p2, p3}, Color: baseColor})

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
			angle := Degrees(i).ToRadians()
			vertices = append(vertices, Point3D{
				X: radius * Unit(math.Cos(float64(angle))),
				Y: radius * Unit(math.Sin(float64(angle))),
				Z: -(Unit(stage)*bodyHeight + tipHeight/2),
			})
			vertices = append(vertices, Point3D{
				X: radius * Unit(math.Cos(float64(angle))),
				Y: radius * Unit(math.Sin(float64(angle))),
				Z: -(Unit(stage+1)*bodyHeight + tipHeight/2),
			})
		}
		bodyVertexCount := len(vertices) - startIndex

		stageColor := adjustColorBrightness(baseColorRGBA, factor, stage+1)

		for i := 0; i < bodyVertexCount-2; i += 2 {
			p1 := vertices[startIndex+i]
			p2 := vertices[startIndex+i+1]
			p3 := vertices[startIndex+i+2]
			faces = append(faces, FaceData{Face: [3]Point3D{p1, p2, p3}, Color: stageColor})
			p1 = vertices[startIndex+i+1]
			p2 = vertices[startIndex+i+3]
			p3 = vertices[startIndex+i+2]
			faces = append(faces, FaceData{Face: [3]Point3D{p1, p2, p3}, Color: stageColor})
		}
		p1 = vertices[startIndex+bodyVertexCount-2]
		p2 = vertices[startIndex+bodyVertexCount-1]
		p3 = vertices[startIndex]
		faces = append(faces, FaceData{Face: [3]Point3D{p1, p2, p3}, Color: stageColor})
		p1 = vertices[startIndex+bodyVertexCount-1]
		p2 = vertices[startIndex+1]
		p3 = vertices[startIndex]
		faces = append(faces, FaceData{Face: [3]Point3D{p1, p2, p3}, Color: stageColor})
	}
	return faces
}

func NewOrientationObject(w ThreeDWidgetInterface) *Object {
	size := Unit(2)
	thickness := size / 20

	faces := []FaceData{
		{
			Face: [3]Point3D{
				{X: 0, Y: -thickness, Z: -thickness},
				{X: size, Y: -thickness, Z: -thickness},
				{X: 0, Y: thickness, Z: -thickness},
			},
			Color: color.RGBA{R: 255, A: 255},
		},
		{
			Face: [3]Point3D{
				{X: size, Y: -thickness, Z: -thickness},
				{X: size, Y: thickness, Z: -thickness},
				{X: 0, Y: thickness, Z: -thickness},
			},
			Color: color.RGBA{R: 255, A: 255},
		},
		{
			Face: [3]Point3D{
				{X: -thickness, Y: 0, Z: -thickness},
				{X: -thickness, Y: size, Z: -thickness},
				{X: thickness, Y: 0, Z: -thickness},
			},
			Color: color.RGBA{R: 255, G: 255, A: 255},
		},
		{
			Face: [3]Point3D{
				{X: -thickness, Y: size, Z: -thickness},
				{X: thickness, Y: size, Z: -thickness},
				{X: thickness, Y: 0, Z: -thickness},
			},
			Color: color.RGBA{R: 255, G: 255, A: 255},
		},
		{
			Face: [3]Point3D{
				{X: -thickness, Y: -thickness, Z: 0},
				{X: -thickness, Y: -thickness, Z: size},
				{X: thickness, Y: -thickness, Z: 0},
			},
			Color: color.RGBA{B: 255, A: 255},
		},
		{
			Face: [3]Point3D{
				{X: -thickness, Y: -thickness, Z: size},
				{X: thickness, Y: -thickness, Z: size},
				{X: thickness, Y: -thickness, Z: 0},
			},
			Color: color.RGBA{B: 255, A: 255},
		},
	}

	orientationObject := Object{
		Faces:    faces,
		Position: Point3D{},
		Rotation: Rotation3D{},
		Widget:   w,
	}
	w.RegisterTickMethod(func() {
		orientationObject.Position = w.GetCamera().UnProject(Point2D{X: 60, Y: 120}, 20, w.GetWidth(), w.GetHeight())
	})
	w.AddObject(&orientationObject)
	return &orientationObject
}
