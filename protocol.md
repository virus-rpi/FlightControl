# WARP - Protocol
Acronym: ***Wa***ter ***R***ocket ***P***rotocol

This is a description of the protocol used by the Water-Rocket to communicate with the ground station.

The Water-Rocket uses gRPC to communicate with the ground station and the controller. 

## Service Definitions

The `WaterRocketService` service provides the following methods:

### Get-Endpoints

- `GetVoltage` - Get the current voltage that flows to the Water-Rocket controller in mV
- `GetStatus` - Get the current status of the Water-Rocket 
- `GetAltitude` - Get the current altitude of the Water-Rocket in m
- `GetAcceleration` - Get the current acceleration of the Water-Rocket in m/s^2
- `GetRotation` - Get the current rotation of the Water-Rocket in deg/s
- `GetSpacialData` - Get the current spacial data of the Water-Rocket. This includes the altitude, orientation, acceleration and velocity in x, y and z direction
- `GetMaxAltitude` - Get the maximum altitude of the Water-Rocket
- `GetMinAltitude` - Get the minimum altitude of the Water-Rocket
- `GetLog` - Get the entire last log of the Water-Rocket
- `GetLogs` - Get a list of all logs and their ids
- `GetLogById` - Get the log with the given id (id is a number)
- `GetLoggingStatus` - Get the current logging status of the Water-Rocket
  
### Post-Endpoints
- `Reset` - Reset the Water-Rocket
- `Arm` - Arm the Water-Rocket (only possible if the base station supports connection to the Water-Rocket)
- `Disarm` - Disarm the Water-Rocket (only possible if the base station supports connection to the Water-Rocket)
- `Launch` - Launch the Water-Rocket (only possible if the base station supports connection to the Water-Rocket)
- `Abort` - Abort the Water-Rocket and release pressure (only possible if the base station supports connection to the Water-Rocket)
- `DeployParachute` - Deploy the parachute 
- `DeployStage` - Deploy the next stage 
- `LogStart` - Start logging
- `LogStop` - Stop logging
- `RecalibrateGyroscope` - Recalibrate the gyroscope
- `RecalibrateAccelerometer` - Recalibrate the accelerometer
- `RecalibrateBarometer` - Recalibrate the barometer
- `RecalibrateMagnetometer` - Recalibrate the magnetometer
- `RecalibrateGPS` - Recalibrate the GPS
- `ResetMax` - Reset all maximum values
- `ResetMin` - Reset all minimum values
- `ResetGyroscope` - Reset the gyroscope
- `ResetAccelerometer` - Reset the accelerometer
- `ResetBarometer` - Reset the barometer
- `ResetMagnetometer` - Reset the magnetometer
- `ResetGPS` - Reset the GPS

The `ControlService` service provides the following methods:

### Post-Endpoints
- `UpdateLiveData` - Update the live data of the Water-Rocket including the timestamp, max altitude, status, voltage and spacial data

The `BaseStationService` service provides the following methods:

### Get-Endpoints
- `GetStatus` - Get the current status of the Base Station
- `GetPressure` - Get the current pressure of the Base Station
- `GetGoalPressure` - Get the goal pressure of the Base Station

### Post-Endpoints
- `SetGoalPressure` - Set the goal pressure of the Base Station
- `Abort` - Abort the Water-Rocket and release pressure 
- `Arm` - Arm the Water-Rocket
- `Disarm` - Disarm the Water-Rocket
- `Launch` - Launch the Water-Rocket
- `RecalibratePressureSensor` - Recalibrate the pressure sensor
- `ResetPressureSensor` - Reset the pressure sensor

## Data Types

The following data types are used:

- `Status` - enum - The status of the Water-Rocket

  Possible values:
  - `ROCKET_IDLE`: The Water-Rocket is idle
  - `ROCKET_ARMED`: The Water-Rocket is armed
  - `ROCKET_BOOSTED_ASCENT`: The Water-Rocket is ascending with the boosters
  - `ROCKET_POWERED_ASCENT`: The Water-Rocket is ascending with the main engine
  - `ROCKET_UNPOWERED_ASCENT`: The Water-Rocket is ascending without any engine
  - `ROCKET_DESCENT`: The Water-Rocket is descending
  - `ROCKET_PARACHUTE_DESCENT`: The Water-Rocket is descending with the parachute deployed
  - `ROCKET_LANDED`: The Water-Rocket has landed
  - `ROCKET_ERROR`: An error occurred

- `Velocity` - struct - The velocity in m/s in x, y and z direction
- `LoggingStatus` - enum - The current logging status of the Water-Rocket

  Possible values:
  - `LOGGING_IDLE`: The Water-Rocket is not logging
  - `LOGGING_ACTIVE`: The Water-Rocket is logging
  - `LOGGING_ERROR`: An error occurred
- `BaseStationStatus` - enum - The status of the Base Station

  Possible values:
  - `BASE_STATION_IDLE`: The Base Station is idle
  - `BASE_STATION_ARMED`: The Base Station reached the goal pressure, and the Water-Rocket is in armed state
  - `BASE_STATION_LAUNCHED`: The Water-Rocket has been launched
  - `BASE_STATION_ABORTED`: The Water-Rocket has been aborted (pressure released)
  - `BASE_STATION_UNDER_PRESSURE`: The Base Station is under pressure (pressure goal reached)
  - `BASE_STATION_ARMING`: The Base Station is arming (building up pressure)
  - `BASE_STATION_ERROR`: An error occurred


## MessageTypes 

The following message types are used:

- `Empty` - An empty message (used as prop for some requests)
- `AcknowlegedResponse` - An empty message that is sent as an acknowledgement when no data is needed 
- `VoltageResponse` - A message that contains the voltage in mV as `int64`
- `StatusResponse` - A message that contains the status of the Water-Rocket as `Status`
- `AltitudeResponse` - A message that contains the altitude in m as `float`
- `AccelerationResponse` - A message that contains the acceleration in m/s^2 as `float` in x, y and z direction 
- `RotationResponse` - A message that contains the rotation in deg/s as `float` in x, y and z direction
- `SpacialDataResponse` - A message that contains the spacial data containing `AltitudeResponse`, `RotationResponse`, `Velocity` and `AccelerationResponse`
- `LogResponse` - A message that contains the log as `string`
- `LogsResponse` - A message that contains the logs as `map<int32, string>` where the key is the id of the log and the timestamp is the value
- `LoggingStatusResponse` - A message that contains the logging status as `LoggingStatus`
- `LogByIdRequest` - A message that contains the id of a log as `int32`
- `UpdateLiveDataRequest` - A message that contains the live data of the Water-Rocket, including the timestamp, max altitude,  status, voltage and spacial data 
- `BaseStationStatusResponse` - A message that contains the status of the Base Station as `BaseStationStatus`
- `PressureResponse` - A message that contains the pressure of the Base Station as `float`
- `GoalPressureRequest` - A message that contains the new goal pressure of the Base Station as `float`