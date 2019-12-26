package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	"github.com/tidwall/gjson"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	temp     float64 = 0
	humidity float64 = 0
)

func connect(clientID string, uri *url.URL) (mqtt.Client, error) {
	var opts = mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientID)

	var client = mqtt.NewClient(opts)
	var token = client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	return client, token.Error()
}

func main() {
	var name = flag.String("name", "hc-mqtt-temperature", "name of the sensor to display in HomeKit")
	var manufacturer = flag.String("manufacturer", "Aosong Electronics", "manufacturer of the sensor")
	var model = flag.String("model", "DHT22", "model number of the sensor")
	var serial = flag.String("serial", "0000", "serial number of the sensor")
	var pin = flag.String("pin", "00102003", "PIN number to pair the HomeKit accessory")
	var port = flag.String("port", "", "Port number for the HomeKit accessory")
	var storagePath = flag.String("storagePath", "hc-mqtt-temperature-data", "path to store data")

	var brokerURI = flag.String("brokerURI", "127.0.0.1:1883", "URI of the MQTT broker")
	var clientID = flag.String("clientID", "hc-mqtt-temperature", "client ID for MQTT")

	var topicTemp = flag.String("topicTemp", "temp", "topic for temperature")
	var topicHum = flag.String("topicHum", "humidity", "topic for humidity")

	var isHumidityEnabled = flag.Bool("humidity", false, "whether to enable humidity")
	var tempJSONPath = flag.String("tempJSONPath", "temperature", "JSON path to the temperature value")
	var humJSONPath = flag.String("humJSONPath", "humidity", "JSON path to the humidity value")

	flag.Parse()

	mqttURI, err := url.Parse(*brokerURI)
	if err != nil {
		log.Fatal(err)
	}

	info := accessory.Info{
		Name:         *name,
		Manufacturer: *manufacturer,
		Model:        *model,
		SerialNumber: *serial,
	}

	ac := accessory.New(info, accessory.TypeSensor)

	tempSensor := service.NewTemperatureSensor()

	tempStatusActive := characteristic.NewStatusActive()
	tempStatusActive.SetValue(false)
	tempSensor.AddCharacteristic(tempStatusActive.Characteristic)

	tempStatusFault := characteristic.NewStatusFault()
	tempStatusFault.SetValue(characteristic.StatusFaultGeneralFault)
	tempSensor.AddCharacteristic(tempStatusFault.Characteristic)

	tempSensor.CurrentTemperature.OnValueGet(func() interface{} {
		log.Println("tempSensor.CurrentTemperature.OnValueGet")
		return temp
	})

	ac.AddService(tempSensor.Service)

	client, err := connect(*clientID, mqttURI)
	if err != nil {
		log.Fatal(err)
	}

	var updateHumidity mqtt.MessageHandler

	if *isHumidityEnabled {
		humiditySensor := service.NewHumiditySensor()

		humidityStatusFault := characteristic.NewStatusFault()
		humiditySensor.AddCharacteristic(humidityStatusFault.Characteristic)

		humidityStatusActive := characteristic.NewStatusActive()
		humiditySensor.AddCharacteristic(humidityStatusActive.Characteristic)

		humiditySensor.CurrentRelativeHumidity.OnValueGet(func() interface{} {
			log.Println("humiditySensor.CurrentRelativeHumidity.OnValueGet")
			return humidity
		})

		ac.AddService(humiditySensor.Service)

		updateHumidity = func(client mqtt.Client, msg mqtt.Message) {
			log.Printf("[%s]: %s\n", *topicHum, string(msg.Payload()))
			humidity = gjson.Get(string(msg.Payload()), *humJSONPath).Float()
			humiditySensor.CurrentRelativeHumidity.UpdateValue(humidity)

			humidityStatusActive.UpdateValue(true)
			humidityStatusFault.UpdateValue(characteristic.StatusFaultNoFault)
		}
	}

	client.Subscribe(*topicTemp, 0, func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("[%s]: %s\n", *topicTemp, string(msg.Payload()))
		temp = gjson.Get(string(msg.Payload()), *tempJSONPath).Float()
		log.Println(temp)
		tempSensor.CurrentTemperature.UpdateValue(temp)
		tempStatusActive.UpdateValue(true)
		tempStatusFault.UpdateValue(characteristic.StatusFaultNoFault)

		if *isHumidityEnabled {
			updateHumidity(client, msg)
		}
	})

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
		<-t.Stop()
	})

	t.Start()
}
