package sensorcollection

import (
	"testing"
	"fmt"
)

func TestConnection(t *testing.T){

	getStationResponse := GetSensorServiceGetStationsResponse()

	for _,s := range getStationResponse.Stations {
		for _, p := range s.Parameters {
			fmt.Printf("%v\n", p.Devices)
		}
	}
}
