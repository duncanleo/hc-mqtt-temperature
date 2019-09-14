package main

import (
	"flag"
	"log"

	"github.com/brutella/hc/service"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	hcLog "github.com/brutella/hc/log"
	"github.com/duncanleo/hc-http-fan/config"
)

func main() {
	hcLog.Debug.Enable()

	cfg, err := config.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	var name = flag.String("name", "DHT22", "name of the sensor to display in HomeKit")

	info := accessory.Info{
		Name: *name,
		// Manufacturer: cfg.Manufacturer,
		// Model:        cfg.Model,
		// SerialNumber: cfg.Serial,
	}

	ac := accessory.New(info, accessory.TypeThermostat)

	tempSensor := service.NewTemperatureSensor()
	tempSensor.CurrentTemperature.OnValueGet(func() interface{} {
		return 99
	})
	tempSensor.CurrentTemperature.OnValueRemoteUpdate(func(v float64) {
		log.Println(v)
	})

	ac.AddService(tempSensor.Service)

	hcConfig := hc.Config{
		Pin:         cfg.Pin,
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
