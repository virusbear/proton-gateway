package homeassistant

type SensorConfig struct {
	EntityConfig

	DeviceClass            *string `json:"device_class,omitempty"`
	ExpireAfter            *int    `json:"expire_after,omitempty"`
	ForceUpdate            *bool   `json:"force_update,omitempty"`
	LastResetValueTemplate *string `json:"last_reset_value_template,omitempty"`
	StateClass             *string `json:"state_class,omitempty"`
	StateTopic             *string `json:"state_topic,omitempty"`
	UnitOfMeasurement      *string `json:"unit_of_measurement,omitempty"`
	ValueTemplate          *string `json:"value_template,omitempty"`
}

func NewSensorConfig(config *EntityConfig) *SensorConfig {
	return &SensorConfig{
		EntityConfig: *config,
	}
}

func (conf *SensorConfig) SetDeviceClass(class string) {
	conf.DeviceClass = &class
}

func (conf *SensorConfig) SetValueTemplate(template string) {
	conf.ValueTemplate = &template
}

func (conf *SensorConfig) SetStateClass(class string) {
	conf.StateClass = &class
}

func (conf *SensorConfig) SetUnitOfMeasurement(unit string) {
	conf.UnitOfMeasurement = &unit
}

func (conf *SensorConfig) SetEntityCategory(category string) {
	conf.EntityCategory = &category
}

func (conf *SensorConfig) SetStateTopic(topic string) {
	conf.StateTopic = &topic
}
