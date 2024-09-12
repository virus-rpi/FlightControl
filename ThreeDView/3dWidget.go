package ThreeDView

import (
	. "FlightControl/ThreeDView/camera"
	. "FlightControl/ThreeDView/object"
	. "FlightControl/ThreeDView/types"
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
	Width  = Pixel(800)
	Height = Pixel(600)
)

// ThreeDWidget is a widget that displays 3D objects
type ThreeDWidget struct {
	widget.BaseWidget
	image              *canvas.Image // The image that is rendered on
	camera             *Camera       // The camera of the 3D widget
	objects            []*Object     // The objects in the 3D widget
	tickMethods        []func()      // The methods that are called every frame
	bgColor            color.Color   // The background color of the 3D widget
	renderFaceOutlines bool          // Whether the faces should be rendered with outlines
	renderFaceColors   bool          // Whether the faces should be rendered with colors
	fpsCap             float64       // The maximum frames per second the widget should render at
	tpsCap             float64       // The maximum ticks per second the widget should tick at
}

// NewThreeDWidget creates a new 3D widget
func NewThreeDWidget() *ThreeDWidget {
	w := &ThreeDWidget{}
	w.ExtendBaseWidget(w)
	w.bgColor = color.Transparent
	w.renderFaceOutlines = false
	w.renderFaceColors = true
	standardCamera := NewCamera(Point3D{}, Rotation3D{})
	w.camera = &standardCamera
	w.objects = []*Object{}
	w.image = canvas.NewImageFromImage(w.render())
	w.fpsCap = math.Inf(1)
	w.tpsCap = math.Inf(1)
	go w.renderLoop()
	go w.tickLoop()
	return w
}

func (w *ThreeDWidget) tickLoop() {
	for {
		startTime := time.Now()
		tickDuration := time.Second / time.Duration(w.tpsCap)

		for _, tickMethod := range w.tickMethods {
			go tickMethod()
		}

		elapsedTime := time.Since(startTime)
		if elapsedTime < tickDuration {
			time.Sleep(tickDuration - elapsedTime)
		}
	}
}

func (w *ThreeDWidget) renderLoop() {
	for {
		startTime := time.Now()
		frameDuration := time.Second / time.Duration(w.fpsCap)

		w.image.Image = w.render()
		go canvas.Refresh(w.image)

		elapsedTime := time.Since(startTime)
		if elapsedTime < frameDuration {
			time.Sleep(frameDuration - elapsedTime)
		}
	}
}

// RegisterTickMethod registers an animation function to be called every frame
func (w *ThreeDWidget) RegisterTickMethod(tick func()) {
	w.tickMethods = append(w.tickMethods, tick)
}

// AddObject adds a 3D object as Object to the widget. This should be called in the method that creates the object
func (w *ThreeDWidget) AddObject(object *Object) {
	w.objects = append(w.objects, object)
}

func (w *ThreeDWidget) GetCamera() *Camera {
	return w.camera
}

func (w *ThreeDWidget) GetWidth() Pixel {
	return Width
}

func (w *ThreeDWidget) GetHeight() Pixel {
	return Height
}

// SetCamera sets the camera of the 3D widget
func (w *ThreeDWidget) SetCamera(camera *Camera) {
	w.camera = camera
}

// SetBackgroundColor sets the background color of the 3D widget
func (w *ThreeDWidget) SetBackgroundColor(color color.Color) {
	w.bgColor = color
}

// SetFPSCap sets the maximum frames per second the widget should render at
func (w *ThreeDWidget) SetFPSCap(fps float64) {
	w.fpsCap = fps
}

// SetTPSCap sets the maximum ticks per second the widget should update at. Animations are triggered at this rate
func (w *ThreeDWidget) SetTPSCap(tps float64) {
	w.tpsCap = tps
}

// SetRenderFaceOutlines sets whether the faces should be rendered with outlines.
// If false, only colors will be rendered. If colors are also false, nothing will be rendered.
// If true, the faces will be rendered with black outlines or the color of the face if face colors are disabled.
// Default is false
func (w *ThreeDWidget) SetRenderFaceOutlines(newVal bool) {
	w.renderFaceOutlines = newVal
}

// SetRenderFaceColors sets whether the faces should be rendered with colors.
// If false, only outlines will be rendered. If outline is also false, nothing will be rendered.
// Default is true
func (w *ThreeDWidget) SetRenderFaceColors(newVal bool) {
	w.renderFaceColors = newVal
}

func (w *ThreeDWidget) CreateRenderer() fyne.WidgetRenderer {
	return &threeDRenderer{image: w.image}
}

func (w *ThreeDWidget) render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(Width), int(Height)))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: w.bgColor}, image.Point{}, draw.Src)

	depthBuffer := make([][]Unit, Height)
	for i := range depthBuffer {
		depthBuffer[i] = make([]Unit, Width)
		for j := range depthBuffer[i] {
			depthBuffer[i][j] = Unit(math.Inf(1))
		}
	}

	var faces []FaceData
	var wg3d sync.WaitGroup
	var mu3d sync.Mutex
	wg3d.Add(len(w.objects))
	for _, object := range w.objects {
		go func(object *Object) {
			defer wg3d.Done()
			objectFaces := object.GetFaces()
			mu3d.Lock()
			faces = append(faces, objectFaces...)
			mu3d.Unlock()
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
			p1 := w.camera.Project(face.Face[0], Width, Height)
			p2 := w.camera.Project(face.Face[1], Width, Height)
			p3 := w.camera.Project(face.Face[2], Width, Height)

			if !p1.InBounds(0, 0, Width, Height) && !p2.InBounds(0, 0, Width, Height) && !p3.InBounds(0, 0, Width, Height) {
				return
			}

			mu.Lock()
			projectedFaces = append(projectedFaces, ProjectedFaceData{Face: [3]Point2D{p1, p2, p3}, Color: face.Color, Distance: face.Distance})
			mu.Unlock()
		}(face)
	}
	wg2d.Wait()

	sort.Slice(projectedFaces, func(i, j int) bool {
		return projectedFaces[i].Distance > projectedFaces[j].Distance
	})

	for _, face := range projectedFaces {
		drawFace(img, face, depthBuffer, w.renderFaceOutlines, w.renderFaceColors)
	}

	return img
}

func (w *ThreeDWidget) Dragged(event *fyne.DragEvent) {
	if controller, ok := w.camera.Controller.(DragController); ok {
		controller.OnDrag(event.Dragged.DX, event.Dragged.DY)
	}
}

func (w *ThreeDWidget) DragEnd() {
	if controller, ok := w.camera.Controller.(DragController); ok {
		controller.OnDragEnd()
	}
}

func (w *ThreeDWidget) Scrolled(event *fyne.ScrollEvent) {
	if controller, ok := w.camera.Controller.(ScrollController); ok {
		controller.OnScroll(event.Scrolled.DX, event.Scrolled.DY)
	}
}

type threeDRenderer struct {
	image *canvas.Image
}

// Layout resizes the widget to the given size
func (r *threeDRenderer) Layout(size fyne.Size) {
	r.image.Resize(size)
	Width = Pixel(size.Width)
	Height = Pixel(size.Height)
}

// MinSize returns the minimum size of the widget
func (r *threeDRenderer) MinSize() fyne.Size {
	return r.image.MinSize()
}

// Refresh refreshes the widget
func (r *threeDRenderer) Refresh() {
	canvas.Refresh(r.image)
}

// Objects returns the objects of the widget. This will be only the image that is rendered on
func (r *threeDRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.image}
}

func (r *threeDRenderer) Destroy() {}

func drawFace(img *image.RGBA, face ProjectedFaceData, depthBuffer [][]Unit, renderFaceOutlines bool, renderFaceColors bool) {
	if renderFaceColors {
		drawFilledTriangle(img, face.Face[0], face.Face[1], face.Face[2], face.Color, depthBuffer, face.Distance)
	}

	if !renderFaceOutlines {
		return
	}
	var outlineColor color.Color
	if !renderFaceColors {
		outlineColor = face.Color
	} else {
		outlineColor = color.Black
	}
	point1 := face.Face[0]
	point2 := face.Face[1]
	point3 := face.Face[2]
	drawLine(img, point1, point2, outlineColor, depthBuffer, face.Distance)
	drawLine(img, point2, point3, outlineColor, depthBuffer, face.Distance)
	drawLine(img, point3, point1, outlineColor, depthBuffer, face.Distance)
}

func drawLine(img *image.RGBA, point1, point2 Point2D, lineColor color.Color, depthBuffer [][]Unit, depth Unit) {
	x0 := point1.X
	y0 := point1.Y
	x1 := point2.X
	y1 := point2.Y
	dx := math.Abs(float64(x1 - x0))
	dy := math.Abs(float64(y1 - y0))
	sx := Pixel(-1)
	if x0 < x1 {
		sx = Pixel(1)
	}
	sy := Pixel(-1)
	if y0 < y1 {
		sy = Pixel(1)
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

func drawFilledTriangle(img *image.RGBA, p1, p2, p3 Point2D, fillColor color.Color, depthBuffer [][]Unit, depth Unit) {
	if p2.Y < p1.Y {
		p1, p2 = p2, p1
	}
	if p3.Y < p1.Y {
		p1, p3 = p3, p1
	}
	if p3.Y < p2.Y {
		p2, p3 = p3, p2
	}

	drawHorizontalLine := func(y, x1, x2 Pixel, color color.Color) {
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

	interpolateX := func(y, y1, y2, x1, x2 Pixel) Pixel {
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
