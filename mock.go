package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/net/websocket"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var mockLogger = log.New(log.Writer(), "[Mockserver] ", log.LstdFlags)

func mockTab() fyne.CanvasObject {
	noticeLabel := widget.NewLabel("This is a mock tab. You can use it to mock a rocket.")

	var status Status = StatusIdle
	mockServer := &MockServer{
		ip: "localhost:8080",
	}
	mockServer.generateMockData(status)

	serverStatusLabel := widget.NewLabel("Server status: stopped")

	startServerButton := widget.NewButton("Start server", func() {
		mockServer.start()
		serverStatusLabel.SetText("Server status: running")
	})

	serverStartStopContainer := container.NewVBox(startServerButton, serverStatusLabel)

	statusSelector := widget.NewSelect([]string{
		"Idle",
		"Armed",
		"Boosted ascent",
		"Powered ascent",
		"Unpowered ascent",
		"Descent",
		"Parachute descent",
		"Landed",
		"Error",
	}, func(selected string) {
		mockServer.data = Data{}
		switch selected {
		case "Idle":
			status = StatusIdle
		case "Armed":
			status = StatusArmed
		case "Boosted ascent":
			status = StatusBoostedAscent
		case "Powered ascent":
			status = StatusPoweredAscent
		case "Unpowered ascent":
			status = StatusUnpoweredAscent
		case "Descent":
			status = StatusDescent
		case "Parachute descent":
			status = StatusParachuteDescent
		case "Landed":
			status = StatusLanded
		case "Error":
			status = StatusError
		}
		mockServer.generateMockData(status)
	})

	autoGenerateMockDataButton := widget.NewButton("Auto-generate mock data", func() {
		go func() {
			for {
				mockServer.generateMockData(status)
				time.Sleep(100 * time.Millisecond)
			}
		}()
	})

	tabContainer := container.NewVBox(noticeLabel, widget.NewSeparator(), serverStartStopContainer, widget.NewSeparator(), statusSelector, autoGenerateMockDataButton)

	return tabContainer
}

type MockServer struct {
	running   bool
	ip        string
	mu        sync.Mutex
	data      Data
	clients   map[*websocket.Conn]bool
	broadcast chan Data
}

func (s *MockServer) start() {
	if s.running {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	mockLogger.Println("Starting mock server...")
	s.running = true
	s.clients = make(map[*websocket.Conn]bool)
	s.broadcast = make(chan Data)

	http.HandleFunc("/websocket", s.handleWebsocket)
	go func() {
		mockLogger.Fatal(http.ListenAndServe(s.ip, nil))
	}()
	mockLogger.Println("Mock server started.")

	go s.sendMockData()
}

func (s *MockServer) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(func(conn *websocket.Conn) {
		mockLogger.Println("New WebSocket connection from", conn.RemoteAddr())
		defer func() {
			s.mu.Lock()
			delete(s.clients, conn)
			s.mu.Unlock()
			err := conn.Close()
			if err != nil {
				mockLogger.Println("Error closing WebSocket connection:", err)
			}
		}()

		s.mu.Lock()
		s.clients[conn] = true
		s.mu.Unlock()

		select {}
	}).ServeHTTP(w, r)
}

func (s *MockServer) sendMockData() {
	mockLogger.Println("Sending mock data...")
	for {
		s.mu.Lock()
		if !s.running {
			s.mu.Unlock()
			return
		}
		s.mu.Unlock()

		dataString := fmt.Sprintf("%d,%f,%f,%d,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f",
			time.Now().Unix(),
			s.data.altitude,
			s.data.maxAltitude,
			s.data.status.toIndex(),
			s.data.voltage,
			s.data.xRotation,
			s.data.yRotation,
			s.data.zRotation,
			s.data.xRotationSpeed,
			s.data.yRotationSpeed,
			s.data.zRotationSpeed,
			s.data.xAcceleration,
			s.data.yAcceleration,
			s.data.zAcceleration,
			s.data.xVelocity,
			s.data.yVelocity,
			s.data.zVelocity,
		)

		s.mu.Lock()
		for client := range s.clients {
			_, err := client.Write([]byte(dataString))
			if err != nil {
				mockLogger.Println(err)
				err := client.Close()
				if err != nil {
					return
				}
				delete(s.clients, client)
			}
		}
		s.mu.Unlock()

		time.Sleep(10 * time.Millisecond)
	}
}

func (s *MockServer) stop() {
	if !s.running {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
	close(s.broadcast)
	for client := range s.clients {
		err := client.Close()
		if err != nil {
			mockLogger.Println("Error closing client connection:", err)
		}
	}
	mockLogger.Println("Mock server stopped.")
}

func (s *MockServer) generateMockData(status Status) {
	switch status {
	case StatusIdle, StatusArmed:
		s.data = Data{
			timestamp:      strconv.Itoa(0),
			altitude:       1.7,
			maxAltitude:    1.7,
			status:         status,
			voltage:        rand.Float64()*0.5 + 4.0,
			xRotation:      0,
			yRotation:      0,
			zRotation:      0,
			xRotationSpeed: 0,
			yRotationSpeed: 0,
			zRotationSpeed: 0,
			xAcceleration:  0,
			yAcceleration:  0,
			zAcceleration:  0,
			xVelocity:      0,
			yVelocity:      0,
			zVelocity:      0,
		}
	case StatusBoostedAscent:
		var zAcceleration float64
		if s.data.zAcceleration == 0 {
			zAcceleration = 50
		} else if s.data.zAcceleration > 0 {
			zAcceleration = s.data.zAcceleration - 0.1
		} else {
			zAcceleration = 0
		}
		s.data = Data{
			timestamp:      strconv.Itoa(0),
			altitude:       s.data.altitude + s.data.zVelocity*0.01,
			maxAltitude:    s.data.altitude + s.data.zVelocity*0.01,
			status:         status,
			voltage:        rand.Float64()*0.5 + 4.0,
			xRotation:      0,
			yRotation:      0,
			zRotation:      0,
			xRotationSpeed: 0,
			yRotationSpeed: 0,
			zRotationSpeed: 0,
			xAcceleration:  0,
			yAcceleration:  0,
			zAcceleration:  zAcceleration,
			xVelocity:      0,
			yVelocity:      0,
			zVelocity:      s.data.zVelocity + zAcceleration*0.01 - 0.1, // -0.1 is the drag
		}
	case StatusPoweredAscent:
		var zAcceleration float64
		if s.data.zAcceleration == 0 {
			zAcceleration = 100
		} else if s.data.zAcceleration > 0 {
			zAcceleration = s.data.zAcceleration - 0.2
		} else {
			zAcceleration = 0
		}
		var altitude float64
		if s.data.altitude == 0 {
			altitude = 100
		} else {
			altitude = s.data.altitude
		}
		s.data = Data{
			timestamp:      strconv.Itoa(0),
			altitude:       altitude + s.data.zVelocity*0.01,
			maxAltitude:    s.data.altitude + s.data.zVelocity*0.01,
			status:         status,
			voltage:        rand.Float64()*0.5 + 4.0,
			xRotation:      0,
			yRotation:      0,
			zRotation:      0,
			xRotationSpeed: 0,
			yRotationSpeed: 0,
			zRotationSpeed: 0,
			xAcceleration:  0,
			yAcceleration:  0,
			zAcceleration:  zAcceleration,
			xVelocity:      0,
			yVelocity:      0,
			zVelocity:      s.data.zVelocity + zAcceleration*0.01 - 0.1, // -0.1 is the drag
		}
	case StatusUnpoweredAscent:
		var zAcceleration float64
		if s.data.zAcceleration == 0 {
			zAcceleration = 0
		} else if s.data.zAcceleration > 0 {
			zAcceleration = s.data.zAcceleration - 0.1
		} else {
			zAcceleration = 0
		}
		var altitude float64
		if s.data.altitude == 0 {
			altitude = 300
		} else {
			altitude = s.data.altitude
		}
		s.data = Data{
			timestamp:      strconv.Itoa(0),
			altitude:       altitude + s.data.zVelocity*0.01,
			maxAltitude:    altitude + s.data.zVelocity*0.01,
			status:         status,
			voltage:        rand.Float64()*0.5 + 4.0,
			xRotation:      0,
			yRotation:      0,
			zRotation:      0,
			xRotationSpeed: 0,
			yRotationSpeed: 0,
			zRotationSpeed: 0,
			xAcceleration:  0,
			yAcceleration:  0,
			zAcceleration:  zAcceleration,
			xVelocity:      0,
			yVelocity:      0,
			zVelocity:      s.data.zVelocity + zAcceleration*0.01 - 0.1, // -0.1 is the drag
		}
	case StatusDescent:
		var altitude float64
		if s.data.altitude == 0 {
			altitude = 350
		} else {
			altitude = s.data.altitude
		}

		s.data = Data{
			timestamp:      strconv.Itoa(0),
			altitude:       altitude - s.data.zVelocity*0.01,
			maxAltitude:    350,
			status:         status,
			voltage:        rand.Float64()*0.5 + 4.0,
			xRotation:      0,
			yRotation:      0,
			zRotation:      0,
			xRotationSpeed: 0,
			yRotationSpeed: 0,
			zRotationSpeed: 0,
			xAcceleration:  0,
			yAcceleration:  0,
			zAcceleration:  -9.81,
			xVelocity:      0,
			yVelocity:      0,
			zVelocity:      s.data.zVelocity - 9.81*0.01 + 0.1, // 0.1 is the drag
		}
	case StatusParachuteDescent:
		var altitude float64
		if s.data.altitude == 0 {
			altitude = 200
		} else {
			altitude = s.data.altitude
		}

		s.data = Data{
			timestamp:      strconv.Itoa(0),
			altitude:       altitude - s.data.zVelocity*0.01,
			maxAltitude:    350,
			status:         status,
			voltage:        rand.Float64()*0.5 + 4.0,
			xRotation:      0,
			yRotation:      0,
			zRotation:      0,
			xRotationSpeed: 0,
			yRotationSpeed: 0,
			zRotationSpeed: 0,
			xAcceleration:  0,
			yAcceleration:  0,
			zAcceleration:  -9.81,
			xVelocity:      0,
			yVelocity:      0,
			zVelocity:      s.data.zVelocity - 9.81*0.01 + 0.4, // 0.4 is the drag of the parachute
		}
	case StatusLanded:
		s.data = Data{
			timestamp:      strconv.Itoa(0),
			altitude:       0,
			maxAltitude:    350,
			status:         status,
			voltage:        rand.Float64()*0.5 + 3.0,
			xRotation:      90,
			yRotation:      0,
			zRotation:      0,
			xRotationSpeed: 0,
			yRotationSpeed: 0,
			zRotationSpeed: 0,
			xAcceleration:  0,
			yAcceleration:  0,
			zAcceleration:  0,
			xVelocity:      0,
			yVelocity:      0,
			zVelocity:      0,
		}
	default:
		s.data = Data{}
	}
}
