package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"
)

const (
	Width  = 800
	Height = 600
)

type ThreeDWidget struct {
	widget.BaseWidget
	image  *canvas.Image
	angleX float64
	angleY float64
	angleZ float64
	camera Camera
}

func NewThreeDWidget() *ThreeDWidget {
	w := &ThreeDWidget{}
	w.ExtendBaseWidget(w)
	w.camera = NewCamera(Point3D{X: 0, Y: 0, Z: 250}, Point3D{X: 0}, 1, 100)
	w.image = canvas.NewImageFromImage(w.render())
	go w.animate()
	return w
}

func (w *ThreeDWidget) animate() {
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for range ticker.C {
		w.angleZ += 1
		w.image.Image = w.render()
		canvas.Refresh(w.image)
	}
}

func (w *ThreeDWidget) CreateRenderer() fyne.WidgetRenderer {
	return &threeDRenderer{image: w.image}
}

func (w *ThreeDWidget) render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{A: 0}}, image.Point{}, draw.Src)

	plane := NewPlane(1000, Point3D{X: 0, Y: 0, Z: 0}, Point3D{X: 0, Y: 0, Z: 0}, color.RGBA{G: 255, A: 255}, w)
	cube := NewCube(100, Point3D{X: 0, Y: 0, Z: 100}, Point3D{X: w.angleX, Y: w.angleY, Z: w.angleZ}, color.RGBA{R: 255, A: 255}, w)

	plane.Draw(img)
	cube.Draw(img)

	return img
}

func (w *ThreeDWidget) Dragged(event *fyne.DragEvent) {
	w.RotateCameraAroundPoint(Point3D{X: 0, Y: 0, Z: 0}, float64(event.Dragged.DY/10), float64(event.Dragged.DX/10), 0)
	w.PointCameraAt(Point3D{X: 0, Y: 0, Z: 0})
}

func (w *ThreeDWidget) DragEnd() {}

func (w *ThreeDWidget) Scrolled(event *fyne.ScrollEvent) {
	w.camera.MoveForward(float64(event.Scrolled.DY) / 3)
	w.PointCameraAt(Point3D{X: 0, Y: 0, Z: 0})
	w.image.Image = w.render()
	canvas.Refresh(w.image)
}

type threeDRenderer struct {
	image *canvas.Image
}

func (r *threeDRenderer) Layout(size fyne.Size) {
	r.image.Resize(size)
}

func (r *threeDRenderer) MinSize() fyne.Size {
	return r.image.MinSize()
}

func (r *threeDRenderer) Refresh() {
	canvas.Refresh(r.image)
}

func (r *threeDRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *threeDRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.image}
}

func (r *threeDRenderer) Destroy() {}

func (w *ThreeDWidget) MoveCamera(dx, dy, dz float64) {
	w.camera.Position.X += dx
	w.camera.Position.Y += dy
	w.camera.Position.Z += dz
	w.image.Image = w.render()
	canvas.Refresh(w.image)
}

func (w *ThreeDWidget) RotateCamera(dPitch, dYaw, dRoll float64) {
	w.camera.Pitch += dPitch
	w.camera.Yaw += dYaw
	w.camera.Roll += dRoll
	w.image.Image = w.render()
	canvas.Refresh(w.image)
}

func (w *ThreeDWidget) RotateCameraAroundPoint(point Point3D, dPitch, dYaw, dRoll float64) {
	w.camera.Position = rotateX(w.camera.Position, point, dPitch)
	w.camera.Position = rotateY(w.camera.Position, point, dYaw)
	w.camera.Position = rotateZ(w.camera.Position, point, dRoll)
	w.image.Image = w.render()
	canvas.Refresh(w.image)
}

func (w *ThreeDWidget) PointCameraAt(target Point3D) {
	w.camera.PointAt(target)
	w.image.Image = w.render()
	canvas.Refresh(w.image)
}

func drawLine(img *image.RGBA, point1, point2 Point2D, col color.Color, retry bool) {
	x0 := int(point1.X)
	y0 := int(point1.Y)
	x1 := int(point2.X)
	y1 := int(point2.Y)
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy

	maxIterations := dx + dy + 1

	for i := 0; i < maxIterations; i++ {
		img.Set(x0, y0, col)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
		if (sx > 0 && x0 > x1) || (sx < 0 && x0 < x1) || (sy > 0 && y0 > y1) || (sy < 0 && y0 < y1) {
			break
		}
		if x0 < 0 || x0 >= img.Bounds().Dx() || y0 < 0 || y0 >= img.Bounds().Dy() {
			if !retry {
				drawLine(img, point2, point1, col, true)
			}
			break
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Camera struct {
	Position    Point3D
	FocalLength float64
	Scale       float64
	Pitch       float64
	Yaw         float64
	Roll        float64
}

func NewCamera(position Point3D, rotation Point3D, focalLength, scale float64) Camera {
	return Camera{Position: position, FocalLength: focalLength, Scale: scale, Pitch: rotation.X, Yaw: rotation.Y, Roll: rotation.Z}
}

func (camera *Camera) Project(point Point3D) Point2D {
	point = rotateX(point, camera.Position, camera.Pitch)
	point = rotateY(point, camera.Position, camera.Yaw)
	point = rotateZ(point, camera.Position, camera.Roll)

	translatedX := point.X - camera.Position.X
	translatedY := point.Y - camera.Position.Y
	translatedZ := point.Z - camera.Position.Z

	if translatedZ == 0 {
		translatedZ = 0.1
	}

	x2D := (translatedX*camera.FocalLength/translatedZ)*camera.Scale + Width/2
	y2D := (translatedY*camera.FocalLength/translatedZ)*camera.Scale + Height/2

	return Point2D{x2D, y2D}
}

func (camera *Camera) PointAt(target Point3D) {
	direction := Point3D{
		X: target.X - camera.Position.X,
		Y: target.Y - camera.Position.Y,
		Z: target.Z - camera.Position.Z,
	}

	length := math.Sqrt(direction.X*direction.X + direction.Y*direction.Y + direction.Z*direction.Z)
	direction.X /= length
	direction.Y /= length
	direction.Z /= length

	camera.Pitch = math.Asin(-direction.Y) * (180 / math.Pi)
	camera.Yaw = math.Atan2(direction.X, direction.Z) * (180 / math.Pi)
}

func (camera *Camera) MoveForward(distance float64) {
	direction := Point3D{
		X: -math.Sin(degreesToRadians(camera.Yaw)),
		Y: math.Sin(degreesToRadians(camera.Pitch)),
		Z: -math.Cos(degreesToRadians(camera.Yaw)),
	}

	camera.Position.X += direction.X * distance
	camera.Position.Y += direction.Y * distance
	camera.Position.Z += direction.Z * distance
}

type Point3D struct {
	X, Y, Z float64
}

type Point2D struct {
	X, Y float64
}

type ThreeDShape struct {
	Vertices [8]Point3D
	Edges    [12][2]int
	Rotation Point3D
	Position Point3D
	color    color.Color
	w        *ThreeDWidget
}

func NewCube(size float64, position Point3D, rotation Point3D, color color.Color, w *ThreeDWidget) ThreeDShape {
	half := size / 2
	return ThreeDShape{
		Vertices: [8]Point3D{
			{-half, -half, -half}, {half, -half, -half},
			{half, half, -half}, {-half, half, -half},
			{-half, -half, half}, {half, -half, half},
			{half, half, half}, {-half, half, half},
		},
		Edges: [12][2]int{
			{0, 1}, {1, 2}, {2, 3}, {3, 0},
			{4, 5}, {5, 6}, {6, 7}, {7, 4},
			{0, 4}, {1, 5}, {2, 6}, {3, 7},
		},
		Position: position,
		Rotation: rotation,
		color:    color,
		w:        w,
	}
}

func NewPlane(size float64, position Point3D, rotation Point3D, color color.Color, w *ThreeDWidget) ThreeDShape {
	half := size / 2
	return ThreeDShape{
		Vertices: [8]Point3D{
			{-half, -half, 0}, {half, -half, 0},
			{half, half, 0}, {-half, half, 0},
			{-half, -half, 0}, {half, -half, 0},
			{half, half, 0}, {-half, half, 0},
		},
		Edges: [12][2]int{
			{0, 1}, {1, 2}, {2, 3}, {3, 0},
			{4, 5}, {5, 6}, {6, 7}, {7, 4},
			{0, 4}, {1, 5}, {2, 6}, {3, 7},
		},
		Position: position,
		Rotation: rotation,
		color:    color,
		w:        w,
	}
}

func (shape *ThreeDShape) Draw(img *image.RGBA) {
	for _, edge := range shape.Edges {
		start := shape.Vertices[edge[0]]
		end := shape.Vertices[edge[1]]

		start = rotateX(start, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.X)
		start = rotateY(start, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.Y)
		start = rotateZ(start, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.Z)

		end = rotateX(end, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.X)
		end = rotateY(end, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.Y)
		end = rotateZ(end, Point3D{X: 0, Y: 0, Z: 0}, shape.Rotation.Z)

		start.X += shape.Position.X
		start.Y += shape.Position.Y
		start.Z += shape.Position.Z

		end.X += shape.Position.X
		end.Y += shape.Position.Y
		end.Z += shape.Position.Z

		point1 := shape.w.camera.Project(start)
		point2 := shape.w.camera.Project(end)

		drawLine(img, point1, point2, shape.color, false)
	}
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func rotateX(point, pivot Point3D, degrees float64) Point3D {
	radians := degreesToRadians(degrees)
	translatedY := point.Y - pivot.Y
	translatedZ := point.Z - pivot.Z
	newY := translatedY*math.Cos(radians) - translatedZ*math.Sin(radians)
	newZ := translatedY*math.Sin(radians) + translatedZ*math.Cos(radians)
	newY += pivot.Y
	newZ += pivot.Z

	return Point3D{X: point.X, Y: newY, Z: newZ}
}

func rotateY(point, pivot Point3D, degrees float64) Point3D {
	radians := degreesToRadians(degrees)
	translatedX := point.X - pivot.X
	translatedZ := point.Z - pivot.Z
	newX := translatedX*math.Cos(radians) + translatedZ*math.Sin(radians)
	newZ := -translatedX*math.Sin(radians) + translatedZ*math.Cos(radians)
	newX += pivot.X
	newZ += pivot.Z

	return Point3D{X: newX, Y: point.Y, Z: newZ}
}

func rotateZ(point, pivot Point3D, degrees float64) Point3D {
	radians := degreesToRadians(degrees)
	translatedX := point.X - pivot.X
	translatedY := point.Y - pivot.Y
	newX := translatedX*math.Cos(radians) - translatedY*math.Sin(radians)
	newY := translatedX*math.Sin(radians) + translatedY*math.Cos(radians)
	newX += pivot.X
	newY += pivot.Y

	return Point3D{X: newX, Y: newY, Z: point.Z}
}
