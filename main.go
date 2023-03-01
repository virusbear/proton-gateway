package main

import (
	"bytes"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"os"
	"proton-gateway/config"
	"proton-gateway/device"
	"proton-gateway/gateway"
	"proton-gateway/message"
	"proton-gateway/packet"
	"sync"
	"time"
)

func main() {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("error opening config.yaml: %v", err)
	}
	log.Infof("configuration opened")

	conf, err := config.Load(file)
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}
	_ = file.Close()
	log.Infof("configuration loaded")

	devices := make(map[string]device.Device)

	log.Infof("creating new mqtt client")
	options := mqtt.NewClientOptions()
	options.SetAutoReconnect(true)
	options.AddBroker(fmt.Sprintf("tcp://%s:%d", conf.Mqtt.Host, conf.Mqtt.Port))
	client := mqtt.NewClient(options)

	log.Infof("connecting to mqtt server")
	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("error connecting to mqtt broker: %v", token.Error())
	}
	log.Infof("connected to mqtt server")

	log.Infof("building devices and announcing configuration")
	for _, deviceConfig := range conf.Devices {
		dev := device.GetDeviceByType(deviceConfig.Type)
		if dev != nil {
			devices[deviceConfig.Mac] = dev

			log.Infof("announcing configuration for device: %s", deviceConfig.Mac)
			for _, msg := range dev.Configuration(deviceConfig.Mac) {
				client.Subscribe(msg.Topic(), msg.Qos(), func(client mqtt.Client, m mqtt.Message) {
					if bytes.Compare(m.Payload(), msg.Payload()) != 0 {
						client.Publish(msg.Topic(), msg.Qos(), msg.Retain(), msg.Payload()).Wait()
					}
				})
				client.Publish(msg.Topic(), msg.Qos(), msg.Retain(), msg.Payload()).Wait()
			}
		} else {
			log.Warnf("unknown device type %s. Packets for %s won't be handled", deviceConfig.Type, deviceConfig.Mac)
		}
	}
	log.Infof("configuration announced")

	log.Infof("opening gateway")
	gw, err := gateway.OpenGateway(conf.Serial.Port, int(conf.Serial.BaudRate))
	if err != nil {
		log.Fatalf("error opening connection to gateway: %v", err)
	}

	packets := make(chan packet.Packet)
	messages := make(chan message.Message)

	wg := sync.WaitGroup{}

	go func() {
		wg.Add(1)

		log.Infof("starting gateway connection")
		err := gw.Start(func(p packet.Packet) {
			packets <- p
		})

		log.Fatalf("gateway encountered error: %v", err)
	}()

	go func() {
		wg.Add(1)

		log.Infof("listening for incoming packets")
		for p := range packets {
			dev, found := devices[p.Mac()]
			if !found {
				log.Warnf("received packet from unknown device: %s", p.Mac())
				continue
			}

			for _, msg := range dev.Process(p) {
				messages <- msg
			}
		}

		wg.Done()
	}()

	go func() {
		wg.Add(1)

		log.Infof("starting message handler")
		for msg := range messages {
			client.Publish(msg.Topic(), msg.Qos(), msg.Retain(), msg.Payload()).Wait()
		}

		wg.Done()
	}()

	time.Sleep(100 * time.Millisecond)
	wg.Wait()

	log.Errorf("finishing execution")
}

//ec94cb6bd6f00c
