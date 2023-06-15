package gateway

import (
	"encoding/binary"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"proton-gateway/packet"
	"proton-gateway/utils"
	"time"
)

type PacketHandler func(packet.Packet)

type Gateway interface {
	Start(packets PacketHandler) error
}

var ErrOutOfSync = errors.New("gateway: communication out of sync")
var ErrComTimeout = errors.New("gateway: communication timeout")
var ErrInvalidResponse = errors.New("gateway: invalid response")
var syncDelay = 1 * time.Second

type Cmd uint8

const (
	CmdSynchronize  Cmd = 0x81
	CmdAwait        Cmd = 0x42
	CmdRead         Cmd = 0xc3
	CmdMessageCount Cmd = 0x24
	CmdReadMac      Cmd = 0xa5
)

const (
	maxSyncAttempts      = 16
	maxReconnectAttempts = 16
	syncMagic            = 0x0055ffaa
)

type ProtonGateway struct {
	config *serial.Config
	port   *serial.Port
}

func OpenGateway(port string, baudRate int) (Gateway, error) {
	gateway := ProtonGateway{
		config: &serial.Config{Name: port, Baud: baudRate},
	}

	com, err := serial.OpenPort(gateway.config)
	if err != nil {
		return nil, err
	}
	gateway.port = com

	err = gateway.synchronize()
	if err != nil {
		return nil, err
	}

	return &gateway, nil
}

func (gw ProtonGateway) Start(handler PacketHandler) error {
	if err := gw.ensureSynchronized(); err != nil {
		return err
	}

	mac, err := gw.mac()
	if err != nil {
		return err
	}
	log.Infof("Gateway mac address: %s", mac)

	for {
		if err := gw.ensureSynchronized(); err != nil {
			return err
		}

		messageCount, err := gw.messageCount()
		if err != nil {
			return err
		}
		for messageCount > 0 {
			received, err := gw.receivePacket()
			if err != nil {
				return err
			}
			handler(received)

			messageCount, err = gw.messageCount()
			if err != nil {
				return err
			}
		}

		if err := gw.await(); err != nil {
			return err
		}
	}
}

func (gw ProtonGateway) receivePacket() (packet.Packet, error) {
	if err := gw.synchronize(); err != nil {
		return nil, err
	}

	var result packet.Packet
	reader := func() error {
		var err error
		result, err = packet.Read(gw.port)
		return err
	}

	if err := gw.execute(CmdRead, reader); err != nil {
		return nil, err
	}

	return result, nil
}

func (gw ProtonGateway) mac() (string, error) {
	if err := gw.synchronize(); err != nil {
		return "", err
	}

	var mac *string
	reader := func() error {
		var err error
		mac, err = utils.ReadMac(gw.port)
		return err
	}
	err := gw.execute(CmdReadMac, reader)
	return *mac, err
}

func (gw ProtonGateway) await() error {
	if err := gw.synchronize(); err != nil {
		return err
	}

	var dummy uint8
	reader := func() error {
		return binary.Read(gw.port, binary.LittleEndian, &dummy)
	}
	if err := gw.execute(CmdAwait, reader); err != nil {
		if err == ErrComTimeout {
			return nil
		} else {
			return err
		}
	}

	if dummy != 0x00 {
		return ErrInvalidResponse
	}

	return nil
}

func (gw ProtonGateway) messageCount() (int, error) {
	if err := gw.synchronize(); err != nil {
		return 0, err
	}

	var messageCount uint8
	reader := func() error {
		return binary.Read(gw.port, binary.LittleEndian, &messageCount)
	}
	if err := gw.execute(CmdMessageCount, reader); err != nil {
		return 0, err
	}

	return int(messageCount), nil
}

func (gw ProtonGateway) ensureConnected() error {
	var err error
	for i := 0; i < maxReconnectAttempts; i++ {
		err = gw.reconnect()

		if err == nil {
			return nil
		} else {
			log.Errorf("gateway not connected. Attempt %d/%d", i, maxReconnectAttempts)
		}
	}

	return err
}

func (gw ProtonGateway) reconnect() error {
	if gw.port != nil {
		_ = gw.port.Close()
	}

	com, err := serial.OpenPort(gw.config)
	gw.port = com

	if err != nil {
		return err
	}

	return gw.ensureSynchronized()
}

func (gw ProtonGateway) ensureSynchronized() error {
	var err error
	for i := 0; i < maxSyncAttempts; i++ {
		err = gw.synchronize()
		if err == nil {
			return nil
		}
		if err != ErrOutOfSync {
			log.Errorf("error synchronizing gateway: %v", err)
			return err
		} else {
			log.Errorf("gateway not in sync. Attempt %d/%d", i, maxSyncAttempts)
		}
		time.Sleep(syncDelay)
	}

	return err
}

func (gw ProtonGateway) synchronize() error {
	if err := gw.port.Flush(); err != nil {
		return err
	}

	var response uint32
	reader := func() error {
		return binary.Read(gw.port, binary.BigEndian, &response)
	}
	if err := gw.execute(CmdSynchronize, reader); err != nil {
		return err
	}

	if int(response) != syncMagic {
		return ErrOutOfSync
	}

	return nil
}

func (gw ProtonGateway) execute(cmd Cmd, read func() error) error {
	if err := binary.Write(gw.port, binary.LittleEndian, cmd); err != nil {
		return err
	}

	return read()
}
