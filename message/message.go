package message

import "encoding/json"

type Message interface {
	Topic() string
	Qos() byte
	Retain() bool
	Payload() []byte
}

type messageImpl struct {
	topic   string
	qos     byte
	retain  bool
	payload []byte
}

func NewMessage(topic string, payload []byte, retain bool, qos byte) Message {
	return messageImpl{
		topic:   topic,
		payload: payload,
		retain:  retain,
		qos:     qos,
	}
}

func Json(topic string, payload interface{}, retain bool, qos byte) (Message, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return NewMessage(topic, data, retain, qos), nil
}

func (msg messageImpl) Topic() string {
	return msg.topic
}

func (msg messageImpl) Qos() byte {
	return msg.qos
}

func (msg messageImpl) Retain() bool {
	return msg.retain
}

func (msg messageImpl) Payload() []byte {
	return msg.payload
}
