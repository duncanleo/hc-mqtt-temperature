package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	"github.com/tidwall/gjson"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
)

func makeHTTPRequest(url string) ([]byte, error) {
	log.Printf("GET '%s'", url)
	resp, err := http.Get(url)
	if err != nil {
		return make([]byte, 0), nil
	}
	log.Println(resp.Status)
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func main() {
	var name = flag.String("name", "hc-http-temperature", "name of the sensor to display in HomeKit")
	var manufacturer = flag.String("manufacturer", "Aosong Electronics", "manufacturer of the sensor")
	var model = flag.String("model", "DHT22", "model number of the sensor")
	var serial = flag.String("serial", "0000", "serial number of the sensor")
	var pin = flag.String("pin", "00102003", "PIN number to pair the HomeKit accessory")
	var port = flag.String("port", "50004", "Port number for the HomeKit accessory")
	var storagePath = flag.String("storagePath", "hc-http-temperature-data", "path to store data")

	var url = flag.String("url", "", "URL to fetch temperature")
	var isHumidityEnabled = flag.Bool("humidity", false, "whether to enable humidity")
	var tempJSONPath = flag.String("tempJSONPath", ".temperature", "JSON path to the temperature value")
	var humJSONPath = flag.String("humJSONPath", ".humidity", "JSON path to the humidity value")

	var fetchInterval = flag.Int("fetchInterval", 10, "time interval in seconds to fetch, default 2s")

	flag.Parse()

	var tempIntervalTicker = time.NewTicker(time.Second * time.Duration(*fetchInterval))
	var tempIntervalTimerChan = make(chan bool)
	var humIntervalTicker = time.NewTicker(time.Second * time.Duration(*fetchInterval))
	var humIntervalTimerChan = make(chan bool)

	info := accessory.Info{
		Name:         *name,
		Manufacturer: *manufacturer,
		Model:        *model,
		SerialNumber: *serial,
	}

	ac := accessory.New(info, accessory.TypeSensor)

	tempSensor := service.NewTemperatureSensor()

	tempStatusActive := characteristic.NewStatusActive()
	tempSensor.AddCharacteristic(tempStatusActive.Characteristic)

	tempStatusFault := characteristic.NewStatusFault()
	tempSensor.AddCharacteristic(tempStatusFault.Characteristic)

	var fetchTemperature = func() interface{} {
		data, err := makeHTTPRequest(*url)
		tempStatusActive.UpdateValue(true)
		if err != nil {
			log.Println(err)
			tempStatusActive.UpdateValue(false)
			tempStatusFault.UpdateValue(characteristic.StatusFaultGeneralFault)
			return nil
		}
		tempStatusFault.UpdateValue(characteristic.StatusFaultNoFault)
		return gjson.Get(string(data), *tempJSONPath).Float()
	}

	tempSensor.CurrentTemperature.OnValueGet(func() interface{} {
		log.Println("tempSensor.CurrentTemperature.OnValueGet")
		return fetchTemperature()
	})

	go func() {
		for {
			select {
			case <-tempIntervalTimerChan:
				return
			case <-tempIntervalTicker.C:
				tempSensor.CurrentTemperature.UpdateValue(fetchTemperature())
				break
			}
		}
	}()

	ac.AddService(tempSensor.Service)

	if *isHumidityEnabled {
		humiditySensor := service.NewHumiditySensor()

		humidityStatusFault := characteristic.NewStatusFault()
		humiditySensor.AddCharacteristic(humidityStatusFault.Characteristic)

		humidityStatusActive := characteristic.NewStatusActive()
		humiditySensor.AddCharacteristic(humidityStatusActive.Characteristic)

		var fetchHumidity = func() interface{} {
			data, err := makeHTTPRequest(*url)
			humidityStatusActive.UpdateValue(true)
			if err != nil {
				log.Println(err)
				humidityStatusActive.UpdateValue(false)
				humidityStatusFault.UpdateValue(characteristic.StatusFaultGeneralFault)
				return nil
			}
			humidityStatusFault.UpdateValue(characteristic.StatusFaultNoFault)
			return gjson.Get(string(data), *humJSONPath).Float()
		}

		humiditySensor.CurrentRelativeHumidity.OnValueGet(func() interface{} {
			log.Println("humiditySensor.CurrentRelativeHumidity.OnValueGet")
			return fetchHumidity()
		})

		go func() {
			for {
				select {
				case <-humIntervalTimerChan:
					return
				case <-humIntervalTicker.C:
					humiditySensor.CurrentRelativeHumidity.UpdateValue(fetchHumidity())
					break
				}
			}
		}()

		ac.AddService(humiditySensor.Service)
	}

	hcConfig := hc.Config{
		Pin:         *pin,
		StoragePath: *storagePath,
		Port:        *port,
	}

	t, err := hc.NewIPTransport(hcConfig, ac)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		tempIntervalTimerChan <- true
		humIntervalTimerChan <- true
		<-t.Stop()
	})

	t.Start()
}
