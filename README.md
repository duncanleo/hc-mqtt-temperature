# hc-mqtt-temperature
This is a CLI tool built with [hc](https://github.com/brutella/hc) that emulates a HomeKit temperature and (optional) humidity accessory.

### Usage
```shell
Usage of hc-mqtt-temperature:
  -brokerURI string
    	URI of the MQTT broker (default "127.0.0.1:1883")
  -clientID string
    	client ID for MQTT (default "hc-mqtt-temperature")
  -humJSONPath string
    	JSON path to the humidity value (default "humidity")
  -humidity
    	whether to enable humidity
  -manufacturer string
    	manufacturer of the sensor (default "Aosong Electronics")
  -model string
    	model number of the sensor (default "DHT22")
  -name string
    	name of the sensor to display in HomeKit (default "hc-mqtt-temperature")
  -pin string
    	PIN number to pair the HomeKit accessory (default "00102003")
  -port string
    	Port number for the HomeKit accessory
  -serial string
    	serial number of the sensor (default "0000")
  -storagePath string
    	path to store data (default "hc-mqtt-temperature-data")
  -tempJSONPath string
    	JSON path to the temperature value (default "temperature")
  -topicHum string
    	topic for humidity (default "humidity")
  -topicTemp string
    	topic for temperature (default "temp")
```

### JSON Path
The code uses the [gjson](https://github.com/tidwall/gjson) package to parse data freely from any JSON response. The key system is similar to `jq` but it omits the leading period (`.`). See this [playground](http://tidwall.com/gjson-play) for more info.
