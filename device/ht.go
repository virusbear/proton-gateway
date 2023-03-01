package device

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"proton-gateway/homeassistant"
	"proton-gateway/message"
	"proton-gateway/packet"
	"time"
)

type payload struct {
	Temperature      float32 `json:"temperature,omitempty"`
	Humidity         float32 `json:"humidity,omitempty"`
	Voltage          float32 `json:"battery_voltage,omitempty"`
	Current          float32 `json:"battery_current,omitempty"`
	AbsoluteHumidity float32 `json:"absolute_humidity"`
	DewPoint         float32 `json:"dew_point"`
	Level            float32 `json:"battery_level"`
}

type ProtonHT struct {
	messageTimestamps map[string]time.Time
}

func NewProtonHT() Device {
	return &ProtonHT{
		messageTimestamps: make(map[string]time.Time),
	}
}

func (dev ProtonHT) deviceConfig(mac string) *homeassistant.DeviceConfig {
	conf := homeassistant.NewDeviceConfig()
	conf.AddIdentifier(dev.Id(mac))
	conf.SetManufacturer("espressif")
	conf.SetModel("lolin32-lite")
	conf.SetName(dev.Id(mac))
	conf.SetSoftwareVersion("v0.0.1")

	return conf
}

func (dev ProtonHT) entityConfig(mac string, entity string) *homeassistant.EntityConfig {
	conf := homeassistant.NewEntityConfig()
	conf.SetAvailabilityTopic(dev.availabilityTopic(mac))
	conf.SetObjectId(fmt.Sprintf("%s_%s", dev.Id(mac), entity))
	conf.SetUniqueId(fmt.Sprintf("%s_%s", dev.Id(mac), entity))
	conf.SetPayloadAvailable("online")
	conf.SetPayloadNotAvailable("offline")
	conf.Device = dev.deviceConfig(mac)

	return conf
}

func (dev ProtonHT) temperatureConfig(mac string) message.Message {
	conf := homeassistant.NewSensorConfig(dev.entityConfig(mac, "temperature"))

	conf.SetDeviceClass("temperature")
	conf.SetValueTemplate("{{ value_json.temperature | round(1) }}")
	conf.SetStateClass("measurement")
	conf.SetUnitOfMeasurement("°C")
	conf.SetName("Temperature")
	conf.SetStateTopic(dev.stateTopic(mac))

	return dev.configToMessage(
		homeassistant.AutoDiscoveryTopic(homeassistant.EntityTypeSensor, dev.Id(mac), "temperature"),
		&conf,
	)
}

func (dev ProtonHT) humidityConfig(mac string) message.Message {
	conf := homeassistant.NewSensorConfig(dev.entityConfig(mac, "humidity"))

	conf.SetDeviceClass("humidity")
	conf.SetValueTemplate("{{ value_json.humidity | round(1) }}")
	conf.SetStateClass("measurement")
	conf.SetUnitOfMeasurement("%")
	conf.SetName("Humidity")
	conf.SetStateTopic(dev.stateTopic(mac))

	return dev.configToMessage(
		homeassistant.AutoDiscoveryTopic(homeassistant.EntityTypeSensor, dev.Id(mac), "humidity"),
		&conf,
	)
}

func (dev ProtonHT) dewPointConfig(mac string) message.Message {
	conf := homeassistant.NewSensorConfig(dev.entityConfig(mac, "dew_point"))

	conf.SetDeviceClass("temperature")
	conf.SetValueTemplate("{{ value_json.dew_point | round(1) }}")
	conf.SetStateClass("measurement")
	conf.SetUnitOfMeasurement("°C")
	conf.SetName("Dew Point")
	conf.SetStateTopic(dev.stateTopic(mac))

	return dev.configToMessage(
		homeassistant.AutoDiscoveryTopic(homeassistant.EntityTypeSensor, dev.Id(mac), "dew_point"),
		&conf,
	)
}

func (dev ProtonHT) absoluteHumidityConfig(mac string) message.Message {
	conf := homeassistant.NewSensorConfig(dev.entityConfig(mac, "absolute_humidity"))

	conf.SetDeviceClass("water")
	conf.SetValueTemplate("{{ value_json.absolute_humidity | round(1) }}")
	conf.SetStateClass("measurement")
	conf.SetUnitOfMeasurement("mg/m³")
	conf.SetName("Absolute Humidity")
	conf.SetStateTopic(dev.stateTopic(mac))

	return dev.configToMessage(
		homeassistant.AutoDiscoveryTopic(homeassistant.EntityTypeSensor, dev.Id(mac), "absolute_humidity"),
		&conf,
	)
}

func (dev ProtonHT) voltageConfig(mac string) message.Message {
	conf := homeassistant.NewSensorConfig(dev.entityConfig(mac, "battery_voltage"))

	conf.SetDeviceClass("voltage")
	conf.SetValueTemplate("{{ value_json.battery_voltage | round(2) }}")
	conf.SetStateClass("measurement")
	conf.SetUnitOfMeasurement("V")
	conf.SetName("Battery Voltage")
	conf.SetEntityCategory("diagnostic")
	conf.SetStateTopic(dev.stateTopic(mac))

	return dev.configToMessage(
		homeassistant.AutoDiscoveryTopic(homeassistant.EntityTypeSensor, dev.Id(mac), "battery_voltage"),
		&conf,
	)
}

func (dev ProtonHT) currentConfig(mac string) message.Message {
	conf := homeassistant.NewSensorConfig(dev.entityConfig(mac, "battery_current"))

	conf.SetDeviceClass("current")
	conf.SetValueTemplate("{{ value_json.battery_current | round(2) }}")
	conf.SetStateClass("measurement")
	conf.SetUnitOfMeasurement("mA")
	conf.SetName("Battery Current")
	conf.SetEntityCategory("diagnostic")
	conf.SetStateTopic(dev.stateTopic(mac))

	return dev.configToMessage(
		homeassistant.AutoDiscoveryTopic(homeassistant.EntityTypeSensor, dev.Id(mac), "battery_current"),
		&conf,
	)
}

func (dev ProtonHT) levelConfig(mac string) message.Message {
	conf := homeassistant.NewSensorConfig(dev.entityConfig(mac, "battery_level"))

	conf.SetDeviceClass("battery")
	conf.SetValueTemplate("{{ value_json.battery_level | round(2) }}")
	conf.SetStateClass("measurement")
	conf.SetUnitOfMeasurement("%")
	conf.SetName("Battery Level")
	conf.SetEntityCategory("diagnostic")
	conf.SetStateTopic(dev.stateTopic(mac))

	return dev.configToMessage(
		homeassistant.AutoDiscoveryTopic(homeassistant.EntityTypeSensor, dev.Id(mac), "battery_level"),
		&conf,
	)
}

func (dev ProtonHT) Id(mac string) string {
	return fmt.Sprintf("protonht-%s", mac)
}

func (dev ProtonHT) configToMessage(topic string, config interface{}) message.Message {
	msg, err := message.Json(topic, config, true, 0)
	if err != nil {
		panic(err)
	}

	return msg
}

func (dev ProtonHT) Configuration(mac string) []message.Message {
	return []message.Message{
		dev.temperatureConfig(mac),
		dev.humidityConfig(mac),
		dev.absoluteHumidityConfig(mac),
		dev.dewPointConfig(mac),
		dev.voltageConfig(mac),
		dev.currentConfig(mac),
		dev.levelConfig(mac),
	}
}

func (dev ProtonHT) Process(packet packet.Packet) []message.Message {
	reader := bytes.NewReader(packet.Payload())
	payload := payload{}

	err := binary.Read(reader, binary.LittleEndian, &(payload.Temperature))
	if err != nil {
		return dev.offline(packet.Mac())
	}
	err = binary.Read(reader, binary.LittleEndian, &(payload.Humidity))
	if err != nil {
		return dev.offline(packet.Mac())
	}
	err = binary.Read(reader, binary.LittleEndian, &(payload.Voltage))
	if err != nil {
		return dev.offline(packet.Mac())
	}
	err = binary.Read(reader, binary.LittleEndian, &(payload.Current))
	if err != nil {
		return dev.offline(packet.Mac())
	}

	payload.AbsoluteHumidity = dev.absoluteHumidity(payload.Temperature, payload.Humidity)
	payload.DewPoint = dev.dewPoint(payload.Temperature, payload.Humidity)
	payload.Level = dev.level(payload.Voltage)

	availabilityPayload := "online"
	if packet.Timestamp().Sub(dev.lastMessage(packet.Mac())) >= 3*time.Minute {
		availabilityPayload = "offline"
	}

	stateMessage, err := message.Json(dev.stateTopic(packet.Mac()), &payload, false, 0)
	if err != nil {
		return dev.offline(packet.Mac())
	}

	dev.logMessage(packet.Mac(), packet.Timestamp())

	return []message.Message{
		message.NewMessage(dev.availabilityTopic(packet.Mac()), []byte(availabilityPayload), true, 0),
		stateMessage,
	}
}

func (dev ProtonHT) lastMessage(mac string) time.Time {
	timestamp, found := dev.messageTimestamps[mac]
	if !found {
		return time.Now()
	}

	return timestamp
}

func (dev ProtonHT) logMessage(mac string, time time.Time) {
	dev.messageTimestamps[mac] = time
}

func (dev ProtonHT) offline(mac string) []message.Message {
	return []message.Message{
		message.NewMessage(dev.availabilityTopic(mac), []byte("offline"), true, 0),
	}
}

func (dev ProtonHT) stateTopic(mac string) string {
	return fmt.Sprintf("protons/protonht-%s/state", mac)
}

func (dev ProtonHT) availabilityTopic(mac string) string {
	return fmt.Sprintf("protons/protonht-%s/status", mac)
}

func (dev ProtonHT) dewPoint(temperature float32, humidity float32) float32 {
	alpha := float32(math.Log(float64(humidity/100.0))) + (17.625*temperature)/(243.04+temperature)
	return (243.04 * alpha) / (17.624 - alpha)
}

func (dev ProtonHT) absoluteHumidity(temperature float32, humidity float32) float32 {
	PSat := 6.112 * float32(math.Pow(math.E, float64((17.67*temperature)/(temperature+243.5))))
	P := PSat * (humidity / 100.0)
	return ((P * 2.1674) / (273.15 + temperature)) * 1000.0 * 1000.0
}

func (dev ProtonHT) level(voltage float32) float32 {
	return voltage*100.0 - 320.0
}
