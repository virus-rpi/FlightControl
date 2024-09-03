package main

import (
	"encoding/csv"
	"fyne.io/fyne/v2"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
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

var ws *websocket.Conn
var done chan struct{}

var liveData struct {
	timestamp      string
	voltage        float64
	status         Status
	altitude       float64
	maxAltitude    float64
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

func initWebsocket(App fyne.App) {
	wsUrl := url.URL{Scheme: "ws", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/websocket"}
	log.Println("Connecting to " + wsUrl.String())

	var err error
	ws, err = websocket.Dial(wsUrl.String(), "", "http://localhost/")
	if err != nil {
		log.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Println(err)
		}
	}(ws)

	done = make(chan struct{})

	go func() {
		defer close(done)

		for {
			var msg string
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("Received: %s\n", msg)

			reader := csv.NewReader(strings.NewReader(msg))
			record, err := reader.Read()
			if err != nil {
				log.Println(err)
				continue
			}

			if len(record) == 17 {
				liveData.timestamp = record[0]
				liveData.altitude, _ = strconv.ParseFloat(record[1], 64)
				liveData.maxAltitude, _ = strconv.ParseFloat(record[2], 64)
				liveData.status = toStatus(record[3])
				liveData.voltage, _ = strconv.ParseFloat(record[4], 64)
				liveData.xRotation, _ = strconv.ParseFloat(record[5], 64)
				liveData.yRotation, _ = strconv.ParseFloat(record[6], 64)
				liveData.zRotation, _ = strconv.ParseFloat(record[7], 64)
				liveData.xRotationSpeed, _ = strconv.ParseFloat(record[8], 64)
				liveData.yRotationSpeed, _ = strconv.ParseFloat(record[9], 64)
				liveData.zRotationSpeed, _ = strconv.ParseFloat(record[10], 64)
				liveData.xAcceleration, _ = strconv.ParseFloat(record[11], 64)
				liveData.yAcceleration, _ = strconv.ParseFloat(record[12], 64)
				liveData.zAcceleration, _ = strconv.ParseFloat(record[13], 64)
				liveData.xVelocity, _ = strconv.ParseFloat(record[14], 64)
				liveData.yVelocity, _ = strconv.ParseFloat(record[15], 64)
				liveData.zVelocity, _ = strconv.ParseFloat(record[16], 64)
			} else {
				log.Println("Invalid number of fields")
			}

			updateVoltage(liveData.voltage)
			updateStatus(string(liveData.status))
			updateHeight(liveData.altitude)
			updateMaxHeight(liveData.maxAltitude)
		}
	}()
}

func updateWebsocket(App fyne.App) {
	if ws != nil {
		err := ws.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}
	initWebsocket(App)
}

func getVoltage() float64 {
	return liveData.voltage
}

func getStatus() Status {
	return liveData.status
}

func getHeight() float64 {
	return liveData.altitude
}

func getMaxHeight() float64 {
	return liveData.maxAltitude
}

func post(postUrl url.URL) {
	response, err := http.Post(postUrl.String(), "application/json", nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(response.Body)
}

func reset(App fyne.App) {
	resetUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/reset"}
	post(resetUrl)
}

func deployParachute(App fyne.App) {
	deployParachuteUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/deploy/parachute"}
	post(deployParachuteUrl)
}

func deployStage(App fyne.App) {
	deployStageUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/deploy/stage"}
	post(deployStageUrl)
}

func startLogging(App fyne.App) {
	startLoggingUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/log/start"}
	post(startLoggingUrl)
}

func stopLogging(App fyne.App) {
	stopLoggingUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/log/stop"}
	post(stopLoggingUrl)
}

func recalibrateGyro(App fyne.App) {
	recalibrateGyroUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/recalibrate/gyroscope"}
	post(recalibrateGyroUrl)
}

func recalibrateAccelerometer(App fyne.App) {
	recalibrateAccelerometerUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/recalibrate/accelerometer"}
	post(recalibrateAccelerometerUrl)
}

func recalibrateBarometer(App fyne.App) {
	recalibrateBarometerUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/recalibrate/barometer"}
	post(recalibrateBarometerUrl)
}

func resetMax(App fyne.App) {
	resetMaxUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/reset/max"}
	post(resetMaxUrl)
}

func resetMin(App fyne.App) {
	resetMinUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/reset/min"}
	post(resetMinUrl)
}

func resetGyro(App fyne.App) {
	resetGyroUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/reset/gyroscope"}
	post(resetGyroUrl)
}

func resetAccelerometer(App fyne.App) {
	resetAccelerometerUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/reset/accelerometer"}
	post(resetAccelerometerUrl)
}

func resetBarometer(App fyne.App) {
	resetBarometerUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/post/reset/barometer"}
	post(resetBarometerUrl)
}
