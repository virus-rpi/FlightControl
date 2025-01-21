package main

import (
	"FlightControl/ThreeDView/object"
	"FlightControl/ThreeDView/types"
	"image/color"
	"time"
)

const (
	tipHeight   = types.Unit(30)
	stageHeight = types.Unit(60)
	radius      = types.Unit(10)
)

type Rocket struct {
	objects     []*object.Object
	rotation    types.Rotation3D
	position    types.Point3D
	seperated   bool
	DataChannel chan Data
}

func NewTwoStageRocket(position types.Point3D, rotation types.Rotation3D, w object.ThreeDWidgetInterface) *Rocket {
	rocket := Rocket{
		objects:     make([]*object.Object, 3),
		rotation:    rotation,
		position:    position,
		seperated:   false,
		DataChannel: make(chan Data),
	}
	rocket.position.Z += 180

	tip := object.NewCone(
		rocket.position,
		rocket.rotation,
		color.RGBA{R: 200, G: 200, B: 200, A: 255},
		w,
		tipHeight,
		radius,
	)

	stage1 := object.NewCylinder(
		types.Point3D{X: rocket.position.X, Y: rocket.position.Y, Z: rocket.position.Z - tipHeight*1.5},
		rocket.rotation,
		color.RGBA{R: 150, G: 150, B: 150, A: 255},
		w,
		types.Unit(60),
		types.Unit(10),
	)

	stage2 := object.NewCylinder(
		types.Point3D{X: rocket.position.X, Y: rocket.position.Y, Z: rocket.position.Z - tipHeight*1.5 - stageHeight},
		rocket.rotation,
		color.RGBA{R: 100, G: 100, B: 100, A: 255},
		w,
		types.Unit(60),
		types.Unit(10),
	)

	rocket.objects[0] = tip
	rocket.objects[1] = stage1
	rocket.objects[2] = stage2

	go rocket.listenForData()

	return &rocket
}

func (rocket *Rocket) GetSensorPosition() types.Point3D {
	return rocket.position
}

func (rocket *Rocket) GetPosition() types.Point3D {
	sensorPosition := rocket.position
	if rocket.seperated {
		sensorPosition.Z -= (stageHeight) / 2
	} else {
		sensorPosition.Z -= (stageHeight * 2) / 2
	}
	sensorPosition.Rotate(rocket.position, rocket.rotation)
	return sensorPosition
}

func (rocket *Rocket) GetRotation() types.Rotation3D {
	return rocket.rotation
}

func (rocket *Rocket) Move(position types.Point3D) {
	rocket.position.Add(position)
	for _, obj := range rocket.objects {
		obj.Position.Add(position)
	}
}

func (rocket *Rocket) SetPosition(position types.Point3D) {
	rocket.position = position
	rocket.objects[0].Position = rocket.position
	rocket.objects[1].Position = types.Point3D{X: rocket.position.X, Y: rocket.position.Y, Z: rocket.position.Z - tipHeight*1.5}
	rocket.objects[2].Position = types.Point3D{X: rocket.position.X, Y: rocket.position.Y, Z: rocket.position.Z - tipHeight*1.5 - stageHeight}
}

func (rocket *Rocket) SetRotation(rotation types.Rotation3D) {
	rocket.rotation = rotation
	for _, obj := range rocket.objects {
		obj.Rotation = rotation
		obj.Position.Rotate(rocket.position, rotation)
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
			seperatedStage.Rotation.Add(types.Rotation3D{Roll: 1, Pitch: 1, Yaw: 1})
			time.Sleep(time.Millisecond * 10)
		}
		seperatedStage.Rotation.Roll = 90
		seperatedStage.Position.Z = 15
	}()
}

func (rocket *Rocket) listenForData() {
	for data := range rocket.DataChannel {
		rocket.SetPosition(types.Point3D{X: rocket.position.X, Y: rocket.position.Y, Z: types.Unit(data.altitude) * 100})
		rocket.SetRotation(types.Rotation3D{Roll: types.Degrees(data.xRotation), Pitch: types.Degrees(data.yRotation), Yaw: types.Degrees(data.zRotation)})
		if data.status.toIndex() > Status(StatusBoostedAscent).toIndex() && data.status != StatusError {
			rocket.SeparateStage()
		}
	}
}
