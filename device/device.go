package device

import (
	"proton-gateway/message"
	"proton-gateway/packet"
)

type Device interface {
	Configuration(mac string) []message.Message
	Process(packet packet.Packet) []message.Message
}

var devices map[string]Device

func RegisterDevice(deviceType string, device Device) {
	devices[deviceType] = device
}

func GetDeviceByType(deviceType string) Device {
	device, _ := devices[deviceType]
	return device
}

func init() {
	devices = make(map[string]Device)
	RegisterDevice("ht", NewProtonHT())
}
