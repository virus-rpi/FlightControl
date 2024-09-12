package main

import (
	"FlightControl/ThreeDView/object"
	"FlightControl/ThreeDView/types"
	"image/color"
	"time"
)

type Rocket struct {
	objects   []*object.Object
	Data      Data
	rotation  types.Rotation3D
	position  types.Point3D
	seperated bool
}

func NewTwoStageRocket(position types.Point3D, rotation types.Rotation3D, w object.ThreeDWidgetInterface) *Rocket {
	rocket := Rocket{
		objects:   make([]*object.Object, 3),
		Data:      Data{},
		rotation:  rotation,
		position:  position,
		seperated: false,
	}
	rocket.position.Z += 180

	tip := object.NewCone(
		rocket.position,
		rocket.rotation,
		color.RGBA{R: 200, G: 200, B: 200, A: 255},
		w,
		types.Unit(30),
		types.Unit(10),
	)

	stage1 := object.NewCylinder(
		types.Point3D{X: rocket.position.X, Y: rocket.position.Y, Z: rocket.position.Z - 45},
		rocket.rotation,
		color.RGBA{R: 150, G: 150, B: 150, A: 255},
		w,
		types.Unit(60),
		types.Unit(10),
	)

	stage2 := object.NewCylinder(
		types.Point3D{X: rocket.position.X, Y: rocket.position.Y, Z: rocket.position.Z - 105},
		rocket.rotation,
		color.RGBA{R: 100, G: 100, B: 100, A: 255},
		w,
		types.Unit(60),
		types.Unit(10),
	)

	rocket.objects[0] = tip
	rocket.objects[1] = stage1
	rocket.objects[2] = stage2

	return &rocket
}

func (rocket *Rocket) GetPosition() types.Point3D {
	return rocket.position
}

func (rocket *Rocket) Move(position types.Point3D) {
	rocket.position.Add(position)
	for _, obj := range rocket.objects {
		obj.Position.Add(position)
	}
}

func (rocket *Rocket) SeparateStage() {
	if rocket.seperated {
		return
	}
	seperatedStage := rocket.objects[2]
	rocket.objects[2] = object.NewEmpty(seperatedStage.Widget, seperatedStage.Position)
	rocket.seperated = true
	go func() {
		for {
			if seperatedStage.Position.Z <= 0 {
				break
			}
			seperatedStage.Position.Z -= 2
			seperatedStage.Rotation.Add(types.Rotation3D{X: 1, Y: 1, Z: 1})
			time.Sleep(time.Millisecond * 10)
		}
		seperatedStage.Rotation.X = 90
		seperatedStage.Position.Z = 15
	}()
}
