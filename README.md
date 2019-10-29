# hc-temp-sensor
This is a CLI tool built with [hc](https://github.com/brutella/hc) that emulates a HomeKit temperature and (optional) humidity accessory.

### Usage
```shell
Usage of hc-http-temperature
  -humJSONPath string
    	JSON path to the humidity value (default ".humidity")
  -humidity
    	whether to enable humidity
  -manufacturer string
    	manufacturer of the sensor (default "Aosong Electronics")
  -model string
    	model number of the sensor (default "DHT22")
  -name string
    	name of the sensor to display in HomeKit (default "hc-http-temperature")
  -pin string
    	PIN number to pair the HomeKit accessory (default "00102003")
  -port string
    	Port number for the HomeKit accessory (default "50004")
  -serial string
    	serial number of the sensor (default "0000")
  -storagePath string
    	path to store data (default "hc-http-temperature-data")
  -tempJSONPath string
    	JSON path to the temperature value (default ".temperature")
  -url string
    	URL to fetch temperature
```

### JSON Path
The code uses the [gjson](https://github.com/tidwall/gjson) package to parse data freely from any JSON response. The key system is similar to `jq` but it omits the leading period (`.`). See this [playground](http://tidwall.com/gjson-play) for more info.