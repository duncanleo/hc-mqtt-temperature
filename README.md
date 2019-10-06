# hc-dht
This is a CLI tool built with [hc](https://github.com/brutella/hc) that emulates a HomeKit temperature and humidity accessory for the DHT22 sensor.

### Usage
Run this program. While the default options should work, you will probably need to configure the following:

`hc-dht -gpioPin=X -pin=00102003`

Set the GPIO Pin to the corresponding number for your Raspberry Pi/equivalent.
