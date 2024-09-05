package ThreeDView

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"image"
	"image/color"
	"image/draw"
	"math"
	"sort"
	"sync"
	"time"
)

var (
	Width  = int64(800)
	Height = int64(600)
)

type ThreeDWidget struct {
	widget.BaseWidget
	image              *canvas.Image
	camera             Camera
	objects            []ThreeDShape
	animations         []func()
	bgColor            color.Color
	renderFaceOutlines bool
}

func NewThreeDWidget() *ThreeDWidget {
	w := &ThreeDWidget{}
	w.ExtendBaseWidget(w)
	w.bgColor = color.Transparent
	w.renderFaceOutlines = false
	w.camera = NewCamera(Point3D{}, Point3D{}, 0, 0)
	w.objects = []ThreeDShape{}
	w.image = canvas.NewImageFromImage(w.render())
	go w.animate()
	return w
}

func (w *ThreeDWidget) animate() {
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for range ticker.C {
		for _, animation := range w.animations {
			animation()
		}
		w.image.Image = w.render()
		canvas.Refresh(w.image)
	}
}

func (w *ThreeDWidget) RegisterAnimation(animation func()) {
	w.animations = append(w.animations, animation)
}

func (w *ThreeDWidget) AddObject(object ThreeDShape) {
	w.objects = append(w.objects, object)
}

func (w *ThreeDWidget) SetCamera(camera Camera) {
	w.camera = camera
}

func (w *ThreeDWidget) SetBackgroundColor(color color.Color) {
	w.bgColor = color
}

func (w *ThreeDWidget) SetRenderFaceOutlines(newVal bool) {
	w.renderFaceOutlines = newVal
}

func (w *ThreeDWidget) CreateRenderer() fyne.WidgetRenderer {
	return &threeDRenderer{image: w.image}
}

func (w *ThreeDWidget) render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(Width), int(Height)))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: w.bgColor}, image.Point{}, draw.Src)

	depthBuffer := make([][]float64, Height)
	for i := range depthBuffer {
		depthBuffer[i] = make([]float64, Width)
		for j := range depthBuffer[i] {
			depthBuffer[i][j] = math.MaxFloat64
		}
	}

	var faces []FaceData
	var wg3d sync.WaitGroup
	wg3d.Add(len(w.objects))
	for _, object := range w.objects {
		go func(object ThreeDShape) {
			defer wg3d.Done()
			objectFaces := object.GetFaces()
			for _, face := range objectFaces {
				faces = append(faces, FaceData{face: face.face, color: face.color, distance: w.faceDistance(face.face)})
			}
		}(object)
	}
	wg3d.Wait()

	var projectedFaces []ProjectedFaceData
	var wg2d sync.WaitGroup
	wg2d.Add(len(faces))
	var mu sync.Mutex
	for _, face := range faces {
		go func(face FaceData) {
			defer wg2d.Done()
			p1 := w.camera.Project(face.face[0])
			p2 := w.camera.Project(face.face[1])
			p3 := w.camera.Project(face.face[2])

			p1InBounds := p1.X >= 0 && p1.X < Width && p1.Y >= 0 && p1.Y < Height
			p2InBounds := p2.X >= 0 && p2.X < Width && p2.Y >= 0 && p2.Y < Height
			p3InBounds := p3.X >= 0 && p3.X < Width && p3.Y >= 0 && p3.Y < Height

			if !p1InBounds && !p2InBounds && !p3InBounds {
				return
			}

			mu.Lock()
			projectedFaces = append(projectedFaces, ProjectedFaceData{face: [3]Point2D{p1, p2, p3}, color: face.color, distance: face.distance})
			mu.Unlock()
		}(face)
	}
	wg2d.Wait()

	sort.Slice(projectedFaces, func(i, j int) bool {
		return projectedFaces[i].distance > projectedFaces[j].distance
	})

	for _, face := range projectedFaces {
		drawFace(img, face, depthBuffer, w.renderFaceOutlines)
	}

	return img
}

func (w *ThreeDWidget) faceDistance(face [3]Point3D) float64 {
	cameraPos := w.camera.Position

	normalPos := Point3D{
		X: (face[0].X + face[1].X + face[2].X) / 3,
		Y: (face[0].Y + face[1].Y + face[2].Y) / 3,
		Z: (face[0].Z + face[1].Z + face[2].Z) / 3,
	}

	return math.Sqrt(math.Pow(normalPos.X-cameraPos.X, 2) + math.Pow(normalPos.Y-cameraPos.Y, 2) + math.Pow(normalPos.Z-cameraPos.Z, 2))
}

func (w *ThreeDWidget) Dragged(event *fyne.DragEvent) {
	w.RotateCameraAroundPoint(Point3D{X: 0, Y: 0, Z: 100}, float64(event.Dragged.DY/10), 0, float64(event.Dragged.DX/10))
	w.camera.PointAt(Point3D{X: 0, Y: 0, Z: 100})
}

func (w *ThreeDWidget) DragEnd() {}

func (w *ThreeDWidget) Scrolled(event *fyne.ScrollEvent) {
	w.camera.MoveForward(float64(event.Scrolled.DY) / 3)
	w.camera.PointAt(Point3D{X: 0, Y: 0, Z: 100})
}

type threeDRenderer struct {
	image *canvas.Image
}

func (r *threeDRenderer) Layout(size fyne.Size) {
	r.image.Resize(size)
	Width = int64(size.Width)
	Height = int64(size.Height)
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
}

func (w *ThreeDWidget) RotateCamera(dPitch, dYaw, dRoll float64) {
	w.camera.Pitch += dPitch
	w.camera.Yaw += dYaw
	w.camera.Roll += dRoll
}

func (w *ThreeDWidget) RotateCameraAroundPoint(point Point3D, x, y, z float64) {
	w.camera.Position = rotateX(w.camera.Position, point, x)
	w.camera.Position = rotateY(w.camera.Position, point, y)
	w.camera.Position = rotateZ(w.camera.Position, point, z)
}

func drawFace(img *image.RGBA, face ProjectedFaceData, depthBuffer [][]float64, renderFaceOutlines bool) {
	drawFilledTriangle(img, face.face[0], face.face[1], face.face[2], face.color, depthBuffer, face.distance)

	if !renderFaceOutlines {
		return
	}
	point1 := face.face[0]
	point2 := face.face[1]
	point3 := face.face[2]
	drawLine(img, point1, point2, color.Black, depthBuffer, face.distance)
	drawLine(img, point2, point3, color.Black, depthBuffer, face.distance)
	drawLine(img, point3, point1, color.Black, depthBuffer, face.distance)
}

func drawLine(img *image.RGBA, point1, point2 Point2D, lineColor color.Color, depthBuffer [][]float64, depth float64) {
	x0 := point1.X
	y0 := point1.Y
	x1 := point2.X
	y1 := point2.Y
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := int64(-1)
	if x0 < x1 {
		sx = int64(1)
	}
	sy := int64(-1)
	if y0 < y1 {
		sy = int64(1)
	}
	err := dx - dy

	for {
		if x0 >= 0 && x0 < Width && y0 >= 0 && y0 < Height {
			if depth < depthBuffer[y0][x0] {
				img.Set(int(x0), int(y0), lineColor)
				depthBuffer[y0][x0] = depth
			}
		}
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

func drawFilledTriangle(img *image.RGBA, p1, p2, p3 Point2D, fillColor color.Color, depthBuffer [][]float64, depth float64) {
	if p2.Y < p1.Y {
		p1, p2 = p2, p1
	}
	if p3.Y < p1.Y {
		p1, p3 = p3, p1
	}
	if p3.Y < p2.Y {
		p2, p3 = p3, p2
	}

	drawHorizontalLine := func(y, x1, x2 int64, color color.Color) {
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		for x := x1; x <= x2; x++ {
			if y >= 0 && y < Height && x >= 0 && x < Width {
				if depth < depthBuffer[y][x] {
					img.Set(int(x), int(y), color)
					depthBuffer[y][x] = depth
				}
			}
		}
	}

	interpolateX := func(y, y1, y2, x1, x2 int64) int64 {
		if y1 == y2 {
			return x1
		}
		return x1 + (x2-x1)*(y-y1)/(y2-y1)
	}

	for y := p1.Y; y <= p2.Y; y++ {
		x1 := interpolateX(y, p1.Y, p2.Y, p1.X, p2.X)
		x2 := interpolateX(y, p1.Y, p3.Y, p1.X, p3.X)
		drawHorizontalLine(y, x1, x2, fillColor)
	}

	for y := p2.Y; y <= p3.Y; y++ {
		x1 := interpolateX(y, p2.Y, p3.Y, p2.X, p3.X)
		x2 := interpolateX(y, p1.Y, p3.Y, p1.X, p3.X)
		drawHorizontalLine(y, x1, x2, fillColor)
	}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
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
