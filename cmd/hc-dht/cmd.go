package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/MichaelS11/go-dht"
	"github.com/brutella/hc/service"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
)

func main() {
	var name = flag.String("name", "DHT22", "name of the sensor to display in HomeKit")
	var manufacturer = flag.String("manufacturer", "Aosong Electronics", "manufacturer of the sensor")
	var model = flag.String("model", "DHT22", "model number of the sensor")
	var serial = flag.String("serial", "0000", "serial number of the sensor")
	var pin = flag.String("pin", "00102003", "PIN number to pair the HomeKit accessory")
	var gpioPin = flag.Int("gpioPin", 4, "GPIO pin of the DHT22")
	var storagePath = flag.String("storagePath", "hc-dht-data", "path to store data")
	var useFahrenheit = flag.Bool("f", false, "add this flag to use Fahrenheit")

	info := accessory.Info{
		Name:         *name,
		Manufacturer: *manufacturer,
		Model:        *model,
		SerialNumber: *serial,
	}

	var tempUnit = dht.Celsius
	if *useFahrenheit {
		tempUnit = dht.Fahrenheit
	}

	// DHT Setup
	err := dht.HostInit()
	if err != nil {
		log.Panic("HostInit error: ", err)
		return
	}
	dhtSensor, err := dht.NewDHT(fmt.Sprintf("GPIO%d", *gpioPin), tempUnit, "")
	if err != nil {
		log.Panic("DHT init error: ", err)
		return
	}

	var dhtStop = make(chan struct{})
	var dhtStopped = make(chan struct{})
	var humidity float64
	var temperature float64

	go dhtSensor.ReadBackground(&humidity, &temperature, 20*time.Second, dhtStop, dhtStopped)

	ac := accessory.New(info, accessory.TypeSensor)

	tempSensor := service.NewTemperatureSensor()
	tempSensor.CurrentTemperature.OnValueGet(func() interface{} {
		return temperature
	})

	ac.AddService(tempSensor.Service)

	humiditySensor := service.NewHumiditySensor()
	humiditySensor.CurrentRelativeHumidity.OnValueGet(func() interface{} {
		return humidity
	})

	ac.AddService(humiditySensor.Service)

	hcConfig := hc.Config{
		Pin:         *pin,
		StoragePath: *storagePath,
	}

	t, err := hc.NewIPTransport(hcConfig, ac)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()

	close(dhtStop)

	<-dhtStopped
}
