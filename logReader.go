package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

type LogHeader struct {
	NextPageAddress uint32
	PagesPerTick    uint16
	TickDuration    uint16
	StartTime       int64
}

type BaroConfig struct {
	PressureOversampleCount    uint8
	TemperatureOversampleCount uint8
	AltitudeScaleFactor        float32
	SeaLevelPressure           float32
}

type IMUConfig struct {
	FifoODR                     uint16
	GyroODR                     uint16
	AccelerometerODR            uint16
	GyroDecimation              uint8
	AccelerometerDecimation     uint8
	TemperatureSensorDecimation uint8
	GyroFullRangeScale          uint16
	AccelerometerFullRangeScale uint16
}

type TickData struct {
	TimeSinceBoot       int64
	RocketState         uint32
	AltitudeRelSeaLevel float32
	AltitudeRelGround   float32
	Pressure            float32
	Temperature         float32
	IMUData             IMUData
}

type IMUData struct {
	TemperatureAtTickStart int16
	AngularVelocity        [3]int16
	Acceleration           [3]int16
	Alignment              uint16
	FifoPatternIndex       uint16
	FifoValueSampleCount   uint16
	FifoValues             []int16
}

func read() {
	file, err := os.Open("logfile.bin")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	var header LogHeader
	err = binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		fmt.Println("Error reading header:", err)
		return
	}

	var baroConfig BaroConfig
	err = binary.Read(file, binary.LittleEndian, &baroConfig)
	if err != nil {
		fmt.Println("Error reading baro config:", err)
		return
	}

	var imuConfig IMUConfig
	err = binary.Read(file, binary.LittleEndian, &imuConfig)
	if err != nil {
		fmt.Println("Error reading IMU config:", err)
		return
	}

	fmt.Printf("Log Start Time: %s\n", time.Unix(0, header.StartTime*int64(time.Millisecond)).Format(time.RFC3339))
	fmt.Printf("Tick Duration: %d ms\n", header.TickDuration)
	fmt.Printf("Baro Config: %+v\n", baroConfig)
	fmt.Printf("IMU Config: %+v\n", imuConfig)

	for {
		var tickData TickData
		err = binary.Read(file, binary.LittleEndian, &tickData.TimeSinceBoot)
		if err != nil {
			break
		}
		err = binary.Read(file, binary.LittleEndian, &tickData.RocketState)
		if err != nil {
			break
		}
		err = binary.Read(file, binary.LittleEndian, &tickData.AltitudeRelSeaLevel)
		if err != nil {
			break
		}
		err = binary.Read(file, binary.LittleEndian, &tickData.AltitudeRelGround)
		if err != nil {
			break
		}
		err = binary.Read(file, binary.LittleEndian, &tickData.Pressure)
		if err != nil {
			break
		}
		err = binary.Read(file, binary.LittleEndian, &tickData.Temperature)
		if err != nil {
			break
		}

		var imuData IMUData
		err = binary.Read(file, binary.LittleEndian, &imuData.TemperatureAtTickStart)
		if err != nil {
			break
		}
		err = binary.Read(file, binary.LittleEndian, &imuData.AngularVelocity)
		if err != nil {
			break
		}
		err = binary.Read(file, binary.LittleEndian, &imuData.Acceleration)
		if err != nil {
			break
		}
		err = binary.Read(file, binary.LittleEndian, &imuData.Alignment)
		if err != nil {
			break
		}
		err = binary.Read(file, binary.LittleEndian, &imuData.FifoPatternIndex)
		if err != nil {
			break
		}
		err = binary.Read(file, binary.LittleEndian, &imuData.FifoValueSampleCount)
		if err != nil {
			break
		}

		imuData.FifoValues = make([]int16, imuData.FifoValueSampleCount)
		err = binary.Read(file, binary.LittleEndian, &imuData.FifoValues)
		if err != nil {
			break
		}

		tickData.IMUData = imuData

		fmt.Printf("Tick Data: %+v\n", tickData)
	}
}
