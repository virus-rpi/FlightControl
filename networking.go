package main

import (
	"FlightControl/warp"
	"log"
)

func listenToNewData() {
	newDataChannel := ps.Sub("newData")
	for newData := range newDataChannel {
		newData := newData.(Data)
		updateVoltage(newData.voltage)
		updateStatus(string(newData.status))
		updateHeight(newData.altitude)
		updateMaxHeight(newData.maxAltitude)
	}
}

func getLog() {
	res, err := warp.Client.C.GetLog(*warp.Client.Ctx, &warp.Empty{})
	if err != nil {
		log.Println("Error getting log:", err)
		return
	}
	currentLog = parseCSVLog(res.Log)
}
