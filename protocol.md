# WARP - Protocol
Acronym: ***Wa***ter ***R***ocket ***P***rotocol

This is a description of the protocol used by the Water-Rocket to communicate with the ground station.

## Endpoints of the Water-Rocket

### Get-Endpoints

- `/get/voltage` - Get the current voltage that flows to the Water-Rocket controller in mV
- `/get/status` - Get the current status of the Water-Rocket 

   Possible values:
   - `idle`: The Water-Rocket is idle
   - `armed`: The Water-Rocket is armed
   - `boosted-ascent`: The Water-Rocket is ascending with the boosters
   - `powered-ascent`: The Water-Rocket is ascending with the main engine
   - `unpowered-ascent`: The Water-Rocket is ascending without any engine 
   - `descent`: The Water-Rocket is descending
   - `parachute-descent`: The parachute has been deployed and the Water-Rocket is descending 
   - `landed`: The Water-Rocket has landed
   - `error`: An error occurred (e.g. the parachute did not deploy)

- `/get/altitude` - Get the current altitude of the Water-Rocket in m
- `/get/acceleration/x` - Get the current x-acceleration of the Water-Rocket in m/s^2
- `/get/acceleration/y` - Get the current y-acceleration of the Water-Rocket in m/s^2
- `/get/acceleration/z` - Get the current z-acceleration of the Water-Rocket in m/s^2
- `/get/rotation/x` - Get the current x-rotation of the Water-Rocket in deg/s
- `/get/rotation/y` - Get the current y-rotation of the Water-Rocket in deg/s
- `/get/rotation/z` - Get the current z-rotation of the Water-Rocket in deg/s
- `/get/spacial-data` - Get the current spacial data of the Water-Rocket. This includes the altitude, orientation, acceleration and velocity in x, y and z direction as json.
   Format:
   ```json
   {
       "altitude": "0.0",
       "x-rotation": "0.0",
       "y-rotation": "0.0",
       "z-rotation": "0.0",
       "x-rotation-speed": "0.0",
       "y-rotation-speed": "0.0",
       "z-rotation-speed": "0.0",
       "x-acceleration": "0.0",
       "y-acceleration": "0.0",
       "z-acceleration": "0.0",
       "x-velocity": "0.0",
       "y-velocity": "0.0",
       "z-velocity": "0.0"
   }
   ```
- `/get/max/altitude` - Get the maximum altitude of the Water-Rocket
- `/get/min/altitude` - Get the minimum altitude of the Water-Rocket
- `/get/log` - Get the entire last log of the Water-Rocket
   
   Format: 
   A long string with all data points in the same format as the websocket message seperated by newlines
- `/get/log/{id}` - Get the log with the given id (id is a number)
   
    Format: 
    Same as `/get/log`
- `/get/logs` - Get a list of all logs and their ids
  ```csv
  id,timestamp\n
  id,timestamp\n
  id,timestamp\n
   ```
- `/get/logging-status` - Get the current logging status of the Water-Rocket

   Possible values:
   - `idle`: The Water-Rocket is not logging
   - `logging`: The Water-Rocket is logging
   - `error`: An error occurred while logging
- `/get/websocket` - Get the websocket address of the Water-Rocket

### Post-Endpoints
- `/post/reset` - Reset the Water-Rocket
- `/post/arm` - Arm the Water-Rocket (only possible if the base station supports connection to the Water-Rocket)
- `/post/disarm` - Disarm the Water-Rocket (only possible if the base station supports connection to the Water-Rocket)
- `/post/launch` - Launch the Water-Rocket (only possible if the base station supports connection to the Water-Rocket)
- `/post/abort` - Abort the Water-Rocket and release pressure (only possible if the base station supports connection to the Water-Rocket)
- `/post/deploy/parachute` - Deploy the parachute
- `/post/deploy/stage` - Deploy the next stage
- `/post/log/start` - Start logging data
- `/post/log/stop` - Stop logging data
- `/post/recalibrate/gyroscope` - Recalibrate the gyroscope
- `/post/recalibrate/accelerometer` - Recalibrate the accelerometer
- `/post/recalibrate/barometer` - Recalibrate the barometer
- `/post/recalibrate/gps` - Recalibrate the gps
- `/post/reset/max` - Reset all maximum values
- `/post/reset/min` - Reset all minimum values
- `/post/reset/gyroscope` - Reset the gyroscope
- `/post/reset/accelerometer` - Reset the accelerometer
- `/post/reset/barometer` - Reset the barometer
- `/post/reset/gps` - Reset the gps


### Websocket-Message Format
The Water-Rocket sends live sensor readings to the flight control client in following format:
```csv
timestamp,altitude,max-altitude,status (index),voltage,x-rotation,y-rotation,z-rotation,x-rotation-speed,y-rotation-speed,z-rotation-speed,x-acceleration,y-acceleration,z-acceleration,x-velocity,y-velocity,z-velocity
```


## Endpoints of the Base Station
### Get-Endpoints
- `/get/status` - Get the current status of the Base Station

   Possible values:
   - `idle`: The Water-Rocket is idle
   - `armed`: The Water-Rocket is armed
   - `launched`: The Water-Rocket is launched
   - `aborted`: The Water-Rocket is aborted (pressure released)
   - `under-pressure`: The Water-Rocket is under pressure (goal pressure reached)
   - `arming`: The Water-Rocket is arming (building up pressure)
- `/get/pressure` - Get the current pressure of the Base Station
- `/get/goal-pressure` - Get the goal pressure of the Base Station

### Post-Endpoints
- `/post/set/goal-pressure` - Set the goal pressure of the Base Station
- `/post/abort` - Abort the Water-Rocket and release pressure
- `/post/arm` - Arm the Water-Rocket
- `/post/disarm` - Disarm the Water-Rocket
- `/post/launch` - Launch the Water-Rocket
- `/post/recalibrate/pressure-sensor` - Recalibrate the pressure sensor
- `/post/reset/pressure-sensor` - Reset the pressure sensor