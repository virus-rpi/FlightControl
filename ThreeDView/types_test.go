package ThreeDView

import (
	"image/color"
	"testing"
)

func TestGetFaces(t *testing.T) {
	shape := ThreeDShape{
		Vertices: []Point3D{
			{X: 0, Y: 0, Z: 0},
			{X: 1, Y: 0, Z: 0},
			{X: 0, Y: 1, Z: 0},
		},
		Faces:    [][3]int{{0, 1, 2}},
		Rotation: Point3D{X: 0, Y: 0, Z: 0},
		Position: Point3D{X: 0, Y: 0, Z: 0},
		color:    color.RGBA{R: 255, A: 255},
	}

	faces := shape.GetFaces()

	if len(faces) != 1 {
		t.Errorf("Expected 1 face, got %d", len(faces))
	}

	expectedFace := FaceData{
		face:     [3]Point3D{{X: 0, Y: 0, Z: 0}, {X: 1, Y: 0, Z: 0}, {X: 0, Y: 1, Z: 0}},
		color:    color.RGBA{R: 255, A: 255},
		distance: 0,
	}

	if faces[0] != expectedFace {
		t.Errorf("Expected face %v, got %v", expectedFace, faces[0])
	}
}

func TestProject(t *testing.T) {
	camera := NewCamera(Point3D{X: 0, Y: 0, Z: 0}, Point3D{X: 0, Y: 0, Z: 0}, 1, 1)
	point := Point3D{X: 1, Y: 1, Z: 1}

	projected := camera.Project(point)

	expected := Point2D{X: Width/2 + 1, Y: Height/2 + 1}

	if projected != expected {
		t.Errorf("Expected %v, got %v", expected, projected)
	}
}

func TestPointAt(t *testing.T) {
	camera := NewCamera(Point3D{X: 0, Y: 0, Z: 0}, Point3D{X: 0, Y: 0, Z: 0}, 1, 1)
	target := Point3D{X: 0, Y: 0, Z: -1}

	camera.PointAt(target)

	if camera.Pitch != 0 || camera.Yaw != 0 {
		t.Errorf("Expected Pitch and Yaw to be 0, got Pitch: %f, Yaw: %f", camera.Pitch, camera.Yaw)
	}
}

func TestMoveForward(t *testing.T) {
	camera := NewCamera(Point3D{X: 0, Y: 0, Z: 0}, Point3D{X: 0, Y: 0, Z: 0}, 1, 1)
	camera.MoveForward(1)

	expected := Point3D{X: 0, Y: 0, Z: -1}

	if camera.Position != expected {
		t.Errorf("Expected position %v, got %v", expected, camera.Position)
	}
}

func TestIsPointInFrustum(t *testing.T) {
	camera := NewCamera(Point3D{X: 0, Y: 0, Z: 0}, Point3D{X: 0, Y: 0, Z: 0}, 1, 1)
	point := Point3D{X: 0, Y: 0, Z: -1}

	inFrustum := camera.IsPointInFrustum(point)

	if !inFrustum {
		t.Errorf("Expected point to be in frustum")
	}

	point = Point3D{X: 1000, Y: 1000, Z: -1}

	inFrustum = camera.IsPointInFrustum(point)

	if inFrustum {
		t.Errorf("Expected point to be out of frustum")
	}
}
