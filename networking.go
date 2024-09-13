package main

import (
	"encoding/csv"
	"encoding/json"
	"fyne.io/fyne/v2"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var ws *websocket.Conn
var done chan struct{}
var netLogger = log.New(log.Writer(), "[Networking] ", log.LstdFlags)

func initWebsocket(App fyne.App) {
	wsUrl := url.URL{Scheme: "ws", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/websocket"}
	netLogger.Println("Connecting to " + wsUrl.String())

	var err error
	for {
		ws, err = websocket.Dial(wsUrl.String(), "", "http://localhost/")
		if err == nil {
			break
		}
		netLogger.Println("WebSocket connection failed, retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	done = make(chan struct{})

	go func() {
		defer func() {
			close(done)
			if err := ws.Close(); err != nil {
				netLogger.Println("Error closing WebSocket:", err)
			}
		}()

		for {
			var msg string
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				if err == io.EOF {
					netLogger.Println("WebSocket read error: EOF")
				} else {
					netLogger.Println("WebSocket read error:", err)
				}
				return
			}
			//netLogger.Printf("Received: %s\n", msg)

			var newestData Data
			parseCSVData(msg, &newestData)
			ps.Pub(newestData, "newData")

			updateVoltage(newestData.voltage)
			updateStatus(string(newestData.status))
			updateHeight(newestData.altitude)
			updateMaxHeight(newestData.maxAltitude)
		}
	}()
}

func updateWebsocket(App fyne.App) {
	if ws != nil {
		err := ws.Close()
		if err != nil {
			netLogger.Println(err)
			return
		}
	}
	initWebsocket(App)
}

func post(postUrl url.URL) {
	response, err := http.Post(postUrl.String(), "application/json", nil)
	if err != nil {
		netLogger.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			netLogger.Println(err)
		}
	}(response.Body)
}

func get(getUrl url.URL) string {
	response, err := http.Get(getUrl.String())
	if err != nil {
		netLogger.Println(err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			netLogger.Println(err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		netLogger.Println(err)
		return ""
	}

	return string(body)
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

func getLog(App fyne.App) {
	getLogUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/log"}
	currentLog = parseCSVLog(get(getLogUrl))
}

func getLogById(App fyne.App, id int) {
	getLogUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/log/" + string(rune(id))}
	currentLog = parseCSVLog(get(getLogUrl))
}

func getLogList(App fyne.App) LogList {
	getLogListUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/logs"}
	logListString := get(getLogListUrl)
	lines := strings.Split(logListString, "\n")
	logList := make(LogList, len(lines))
	for i, line := range lines {
		reader := csv.NewReader(strings.NewReader(line))
		record, err := reader.Read()
		if err != nil {
			netLogger.Println(err)
			return nil
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			netLogger.Println(err)
			return nil
		}
		logList[i].id = id
		logList[i].timestamp = record[1]
	}
	return logList
}

func getLoggingStatus(App fyne.App) LoggingStatus {
	getLoggingStatusUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/log/status"}
	return toLoggingStatus(get(getLoggingStatusUrl))
}

func getVoltage(App fyne.App) float64 {
	getVoltageUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/voltage"}
	voltageString := get(getVoltageUrl)
	voltage, err := strconv.ParseFloat(voltageString, 64)
	if err != nil {
		netLogger.Println(err)
		return 0
	}
	return voltage
}

func getStatus(App fyne.App) Status {
	getStatusUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/status"}
	return toStatus(get(getStatusUrl))
}

func getHeight(App fyne.App) float64 {
	getHeightUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/altitude"}
	heightString := get(getHeightUrl)
	height, err := strconv.ParseFloat(heightString, 64)
	if err != nil {
		netLogger.Println(err)
		return 0
	}
	return height
}

func getXAcceleration(App fyne.App) float64 {
	getXAccelerationUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/acceleration/x"}
	xAccelerationString := get(getXAccelerationUrl)
	xAcceleration, err := strconv.ParseFloat(xAccelerationString, 64)
	if err != nil {
		netLogger.Println(err)
		return 0
	}
	return xAcceleration
}

func getYAcceleration(App fyne.App) float64 {
	getYAccelerationUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/acceleration/y"}
	yAccelerationString := get(getYAccelerationUrl)
	yAcceleration, err := strconv.ParseFloat(yAccelerationString, 64)
	if err != nil {
		netLogger.Println(err)
		return 0
	}
	return yAcceleration
}

func getZAcceleration(App fyne.App) float64 {
	getZAccelerationUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/acceleration/z"}
	zAccelerationString := get(getZAccelerationUrl)
	zAcceleration, err := strconv.ParseFloat(zAccelerationString, 64)
	if err != nil {
		netLogger.Println(err)
		return 0
	}
	return zAcceleration
}

func getXRotation(App fyne.App) float64 {
	getXRotationUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/rotation/x"}
	xRotationString := get(getXRotationUrl)
	xRotation, err := strconv.ParseFloat(xRotationString, 64)
	if err != nil {
		netLogger.Println(err)
		return 0
	}
	return xRotation
}

func getYRotation(App fyne.App) float64 {
	getYRotationUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/rotation/y"}
	yRotationString := get(getYRotationUrl)
	yRotation, err := strconv.ParseFloat(yRotationString, 64)
	if err != nil {
		netLogger.Println(err)
		return 0
	}
	return yRotation
}

func getZRotation(App fyne.App) float64 {
	getZRotationUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/rotation/z"}
	zRotationString := get(getZRotationUrl)
	zRotation, err := strconv.ParseFloat(zRotationString, 64)
	if err != nil {
		netLogger.Println(err)
		return 0
	}
	return zRotation
}

func getSpacialData(App fyne.App) SpacialData {
	getSpacialDataUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/spacial-data"}
	spacialDataJsonString := get(getSpacialDataUrl)
	var spacialData SpacialData
	err := json.Unmarshal([]byte(spacialDataJsonString), &spacialData)
	if err != nil {
		netLogger.Println(err)
		return SpacialData{}
	}
	return spacialData
}

func getMaxAltitude(App fyne.App) float64 {
	getMaxAltitudeUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/max/altitude"}
	maxAltitudeString := get(getMaxAltitudeUrl)
	maxAltitude, err := strconv.ParseFloat(maxAltitudeString, 64)
	if err != nil {
		netLogger.Println(err)
		return 0
	}
	return maxAltitude
}

func getMinAltitude(App fyne.App) float64 {
	getMinAltitudeUrl := url.URL{Scheme: "http", Host: App.Preferences().StringWithFallback("WaRaIP", "Not set"), Path: "/get/min/altitude"}
	minAltitudeString := get(getMinAltitudeUrl)
	minAltitude, err := strconv.ParseFloat(minAltitudeString, 64)
	if err != nil {
		netLogger.Println(err)
		return 0
	}
	return minAltitude
}
