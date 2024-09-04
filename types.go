package main

import (
	"encoding/csv"
	"strconv"
	"strings"
)

const (
	StatusIdle             = "idle"
	StatusArmed            = "armed"
	StatusBoostedAscent    = "boosted-ascent"
	StatusPoweredAscent    = "powered-ascent"
	StatusUnpoweredAscent  = "unpowered-ascent"
	StatusDescent          = "descent"
	StatusParachuteDescent = "parachute-descent"
	StatusLanded           = "landed"
	StatusError            = "error"
)

type Status string

func toStatus(indexStr string) Status {
	index, _ := strconv.Atoi(indexStr)
	switch index {
	case 0:
		return StatusIdle
	case 1:
		return StatusArmed
	case 2:
		return StatusBoostedAscent
	case 3:
		return StatusPoweredAscent
	case 4:
		return StatusUnpoweredAscent
	case 5:
		return StatusDescent
	case 6:
		return StatusParachuteDescent
	case 7:
		return StatusLanded
	case 8:
		return StatusError
	default:
		return StatusError
	}
}

const (
	LoggingStatusIdle    = "idle"
	LoggingStatusLogging = "logging"
	LoggingStatusError   = "error"
)

type LoggingStatus string

func toLoggingStatus(statusString string) LoggingStatus {
	switch statusString {
	case "idle":
		return LoggingStatusIdle
	case "logging":
		return LoggingStatusLogging
	case "error":
		return LoggingStatusError
	default:
		return LoggingStatusError
	}
}

type Data struct {
	timestamp      string
	altitude       float64
	maxAltitude    float64
	status         Status
	voltage        float64
	xRotation      float64
	yRotation      float64
	zRotation      float64
	xRotationSpeed float64
	yRotationSpeed float64
	zRotationSpeed float64
	xAcceleration  float64
	yAcceleration  float64
	zAcceleration  float64
	xVelocity      float64
	yVelocity      float64
	zVelocity      float64
}

func parseCSVData(csvString string, dataVar *Data) {
	reader := csv.NewReader(strings.NewReader(csvString))
	record, err := reader.Read()
	if err != nil {
		println(err)
		return
	}

	if len(record) == 17 {
		dataVar.timestamp = record[0]
		dataVar.altitude, _ = strconv.ParseFloat(record[1], 64)
		dataVar.maxAltitude, _ = strconv.ParseFloat(record[2], 64)
		dataVar.status = toStatus(record[3])
		dataVar.voltage, _ = strconv.ParseFloat(record[4], 64)
		dataVar.xRotation, _ = strconv.ParseFloat(record[5], 64)
		dataVar.yRotation, _ = strconv.ParseFloat(record[6], 64)
		dataVar.zRotation, _ = strconv.ParseFloat(record[7], 64)
		dataVar.xRotationSpeed, _ = strconv.ParseFloat(record[8], 64)
		dataVar.yRotationSpeed, _ = strconv.ParseFloat(record[9], 64)
		dataVar.zRotationSpeed, _ = strconv.ParseFloat(record[10], 64)
		dataVar.xAcceleration, _ = strconv.ParseFloat(record[11], 64)
		dataVar.yAcceleration, _ = strconv.ParseFloat(record[12], 64)
		dataVar.zAcceleration, _ = strconv.ParseFloat(record[13], 64)
		dataVar.xVelocity, _ = strconv.ParseFloat(record[14], 64)
		dataVar.yVelocity, _ = strconv.ParseFloat(record[15], 64)
		dataVar.zVelocity, _ = strconv.ParseFloat(record[16], 64)
	} else {
		println("Invalid number of fields")
	}
}

type Log []Data

func parseCSVLog(csvString string) Log {
	lines := strings.Split(csvString, "\n")
	log := make(Log, len(lines))
	for i, line := range lines {
		parseCSVData(line, &log[i])
	}
	return log
}

type LogListEntry struct {
	id        int
	timestamp string
}

type LogList []LogListEntry

type SpacialData struct {
	altitude       float64
	xRotation      float64
	yRotation      float64
	zRotation      float64
	xRotationSpeed float64
	yRotationSpeed float64
	zRotationSpeed float64
	xAcceleration  float64
	yAcceleration  float64
	zAcceleration  float64
	xVelocity      float64
	yVelocity      float64
	zVelocity      float64
}
