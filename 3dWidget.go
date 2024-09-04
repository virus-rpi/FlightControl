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

type ThreeDWidget struct {
	widget.BaseWidget
	image  *canvas.Image
	angleX float64
	angleY float64
	angleZ float64
}

func NewThreeDWidget() *ThreeDWidget {
	w := &ThreeDWidget{}
	w.ExtendBaseWidget(w)
	w.image = canvas.NewImageFromImage(w.render())
	go w.animate()
	return w
}

func (w *ThreeDWidget) animate() {
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for range ticker.C {
		w.angleX += 1
		w.angleY += 1
		w.angleZ += 1
		w.image.Image = w.render()
		canvas.Refresh(w.image)
	}
}

func (w *ThreeDWidget) CreateRenderer() fyne.WidgetRenderer {
	return &threeDRenderer{image: w.image}
}

func (w *ThreeDWidget) render() image.Image {
	width, height := 800, 600
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{A: 255}}, image.Point{}, draw.Src)

	cube := NewCube(200)
	fov := 256.0
	viewerDistance := 4.0

	for _, edge := range cube.Edges {
		start := cube.Vertices[edge[0]]
		end := cube.Vertices[edge[1]]

		start = rotateX(start, w.angleX)
		start = rotateY(start, w.angleY)
		start = rotateZ(start, w.angleZ)

		end = rotateX(end, w.angleX)
		end = rotateY(end, w.angleY)
		end = rotateZ(end, w.angleZ)

		x0, y0 := project(start, float64(width), float64(height), fov, viewerDistance)
		x1, y1 := project(end, float64(width), float64(height), fov, viewerDistance)

		drawLine(img, x0, y0, x1, y1, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	}

	return img
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

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, col color.Color) {
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

	for {
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
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func project(p Point3D, width, height, fov, viewerDistance float64) (int, int) {
	factor := fov / (viewerDistance + p.Z)
	x := p.X*factor + width/2
	y := -p.Y*factor + height/2
	return int(x), int(y)
}

type Point3D struct {
	X, Y, Z float64
}

type Cube struct {
	Vertices [8]Point3D
	Edges    [12][2]int
}

func NewCube(size float64) Cube {
	half := size / 2
	return Cube{
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
	}
}
func rotateX(p Point3D, angle float64) Point3D {
	rad := angle * math.Pi / 180
	cosa := math.Cos(rad)
	sina := math.Sin(rad)
	y := p.Y*cosa - p.Z*sina
	z := p.Y*sina + p.Z*cosa
	return Point3D{p.X, y, z}
}

func rotateY(p Point3D, angle float64) Point3D {
	rad := angle * math.Pi / 180
	cosa := math.Cos(rad)
	sina := math.Sin(rad)
	x := p.X*cosa + p.Z*sina
	z := -p.X*sina + p.Z*cosa
	return Point3D{x, p.Y, z}
}

func rotateZ(p Point3D, angle float64) Point3D {
	rad := angle * math.Pi / 180
	cosa := math.Cos(rad)
	sina := math.Sin(rad)
	x := p.X*cosa - p.Y*sina
	y := p.X*sina + p.Y*cosa
	return Point3D{x, y, p.Z}
}
