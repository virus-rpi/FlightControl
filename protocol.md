# Protocol

This is a description of the protocol used by the Water-Rocket to communicate with the ground station.

## Endpoints of the Water-Rocket

### Get-Endpoints

- `/get/voltage` - Get the current voltage that flows to the Water-Rocket controller in mV
- `/get/status` - Get the current status of the Water-Rocket 

   Possible values:
   - `idle`: The Water-Rocket is idle
   - `armed`: The Water-Rocket is armed
   - `ascent`: The Water-Rocket is ascending
   - `descent`: The Water-Rocket is descending
   - `landed`: The Water-Rocket has landed
   - `error`: An error occurred (e.g. the parachute did not deploy)

- `/get/altitude` - Get the current altitude of the Water-Rocket in m
- `/get/x-acceleration` - Get the current x-acceleration of the Water-Rocket in m/s^2
- `/get/y-acceleration` - Get the current y-acceleration of the Water-Rocket in m/s^2
- `/get/z-acceleration` - Get the current z-acceleration of the Water-Rocket in m/s^2
- `/get/x-rotation` - Get the current x-rotation of the Water-Rocket in deg/s
- `/get/y-rotation` - Get the current y-rotation of the Water-Rocket in deg/s
- `/get/z-rotation` - Get the current z-rotation of the Water-Rocket in deg/s
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
- `/get/log/id` - Get the log with the given id
- `/get/logs` - Get a list of all logs and their ids
- `/get/logging-status` - Get the current logging status of the Water-Rocket

   Possible values:
   - `idle`: The Water-Rocket is not logging
   - `logging`: The Water-Rocket is logging
   - `error`: An error occurred while logging

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