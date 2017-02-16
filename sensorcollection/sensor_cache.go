package sensorcollection

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type SensorCacheStation struct {
	StationId        int
	NumberOfDevices     int
	OldestAndNewestDates  string
	ObservationCount int
	ParameterIds []int
	DeviceIds [] int
}

type SensorCacheDevice struct {
	StationId        int
	EnhancedParameterId     int
	ParameterId  int
	Discriminant string
	Obs []Sov
}

type Sov struct {
	Units string
	ObservationCount int
	Depth float64
	Dates []int64
	Values []float64
}

type GetStationsOwner struct {
	Id int
	Label string
	Url string
	Operator_Sector string
	Country string
}

type GetStationsSensor struct {
	Id int
	Label string
	ParameterIds []int
}

type GetStationsStation struct {
	Id int
	Label string
	Urn string
	Url string
	StartDate int64
	EndDate int64
	Latitude float64
	Longitude float64
	SourceId int
	PublisherId int
	PlatformId int
	Sensors []GetStationsStationSensor
	Parameters []GetStationsStationParameter
}

type EnhancedParameter struct {
	Id            int
	ParameterId   int
	CellMethods   string
	Interval      string
	VerticalDatum string
}

func (s GetStationsStation) String() string {
	return "Station( Id: " + strconv.Itoa(s.Id) + ", Label: " + s.Label + ", lat: " +
		strconv.FormatFloat(s.Latitude, 'E', -3, 64) + ", lon: " +
		strconv.FormatFloat(s.Longitude, 'E', -3, 64) + ")"
}

type GetStationsStationParameterDevice struct {
	Id int
	EnhancedParameterId int
	StartDate int64
	EndDate int64
	Discriminant string
	DepthMin float64
	DepthMax float64
}

type GetStationsStationParameter struct {
	Id int
	StartDate int64
	EndDate int64
	Urn string
	Devices []GetStationsStationParameterDevice
}

type GetStationsStationSensor struct {
	Id int
	StartDate int64
	EndDate int64
	ParameterIds []int
}

type GetStationsResponse struct {
	CurrentTime int64
	Stations []GetStationsStation
	Parameters map[string] string
	Platforms map[string] string
	EnhancedParameters map[string] EnhancedParameter
	Sensors map[string] GetStationsSensor
	Sources map[string] GetStationsOwner
}

func GetSensorServiceGetStationsResponse() GetStationsResponse{
	resp, err := http.Get("http://pdx.axiomalaska.com/stationsensorservice/getDataValues?method=GetStationsResultSetRowsJSON&verbose=true&version=3&realtimeonly=true&appregion=cencoos&jsoncallback=false")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var s GetStationsResponse

	err = json.Unmarshal(body, &s)

	if err != nil {
		log.Fatal(err)
	}

	return s
}

func GetSensorCacheDevice(deviceId int) SensorCacheDevice{
	resp, err := http.Get("http://pdx.axiomalaska.com/sensor-cache-service/devices/" + strconv.Itoa(deviceId) + "/obs?start_time=0")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var s SensorCacheDevice

	err = json.Unmarshal(body, &s)

	if err != nil {
		log.Fatal(err)
	}

	return s
}

func GetSensorCacheStation(stationId int) SensorCacheStation {
	resp, err := http.Get("http://pdx.axiomalaska.com/sensor-cache-service/stations/" + strconv.Itoa(stationId))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var s SensorCacheStation

	err = json.Unmarshal(body, &s)

	if err != nil {
		log.Fatal(err)
	}

	return s
}
