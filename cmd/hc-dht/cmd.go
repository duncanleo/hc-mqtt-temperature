package main

import (
	"flag"
	"log"

	"github.com/brutella/hc/service"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	hcLog "github.com/brutella/hc/log"
)

func main() {
	hcLog.Debug.Enable()

	var name = flag.String("name", "DHT22", "name of the sensor to display in HomeKit")
	var manufacturer = flag.String("manufacturer", "Aosong Electronics", "manufacturer of the sensor")
	var model = flag.String("model", "DHT22", "model number of the sensor")
	var serial = flag.String("serial", "0000", "serial number of the sensor")
	var pin = flag.String("pin", "00102003", "PIN number to pair the HomeKit accessory")

	info := accessory.Info{
		Name:         *name,
		Manufacturer: *manufacturer,
		Model:        *model,
		SerialNumber: *serial,
	}

	ac := accessory.New(info, accessory.TypeSensor)

	tempSensor := service.NewTemperatureSensor()
	tempSensor.CurrentTemperature.OnValueGet(func() interface{} {
		return 99
	})

	ac.AddService(tempSensor.Service)

	humiditySensor := service.NewHumiditySensor()
	humiditySensor.CurrentRelativeHumidity.OnValueGet(func() interface{} {
		return 2
	})

	ac.AddService(humiditySensor.Service)

	hcConfig := hc.Config{
		Pin:         *pin,
		StoragePath: "storage",
	}

	t, err := hc.NewIPTransport(hcConfig, ac)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}
