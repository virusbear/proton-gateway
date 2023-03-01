package homeassistant

type EntityConfig struct {
	Availability           *AvailabilityConfig `json:"availability,omitempty"`
	AvailabilityMode       *AvailabilityMode   `json:"availability_mode,omitempty"`
	AvailabilityTemplate   *string             `json:"availability_template,omitempty"`
	AvailabilityTopic      *string             `json:"availability_topic,omitempty"`
	Device                 *DeviceConfig       `json:"device,omitempty"`
	EnabledByDefault       *bool               `json:"enabled_by_default,omitempty"`
	Encoding               *string             `json:"encoding,omitempty"`
	EntityCategory         *string             `json:"entity_category,omitempty"`
	JsonAttributesTemplate *string             `json:"json_attributes_template,omitempty"`
	JsonAttributesTopic    *string             `json:"json_attributes_topic,omitempty"`
	Name                   *string             `json:"name,omitempty"`
	ObjectId               *string             `json:"object_id,omitempty"`
	Qos                    *int8               `json:"qos,omitempty"`
	UniqueId               *string             `json:"unique_id,omitempty"`
	Icon                   *string             `json:"icon,omitempty"`
	PayloadAvailable       *string             `json:"payload_available,omitempty"`
	PayloadNotAvailable    *string             `json:"payload_not_available,omitempty"`
}

func NewEntityConfig() *EntityConfig {
	return &EntityConfig{}
}

func (conf *EntityConfig) SetAvailabilityTopic(topic string) {
	conf.AvailabilityTopic = &topic
}

func (conf *EntityConfig) SetName(name string) {
	conf.Name = &name
}

func (conf *EntityConfig) SetObjectId(id string) {
	conf.ObjectId = &id
}

func (conf *EntityConfig) SetUniqueId(id string) {
	conf.UniqueId = &id
}

func (conf *EntityConfig) SetPayloadAvailable(payload string) {
	conf.PayloadAvailable = &payload
}

func (conf *EntityConfig) SetPayloadNotAvailable(payload string) {
	conf.PayloadNotAvailable = &payload
}

type AvailabilityConfig struct {
	PayloadAvailable    *string `json:"payload_available,omitempty"`
	PayloadNotAvailable *string `json:"payload_not_available,omitempty"`
	Topic               string  `json:"topic,omitempty"`
	ValueTemplate       *string `json:"value_template,omitempty"`
}

type AvailabilityMode string

const (
	AvailabilityModeAll    AvailabilityMode = "all"
	AvailabilityModeAny    AvailabilityMode = "any"
	AvailabilityModeLatest AvailabilityMode = "latest"
)

type EntityType string

const (
	EntityTypeSensor EntityType = "sensor"
)
