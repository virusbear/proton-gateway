package homeassistant

type DeviceConfig struct {
	ConfigurationUrl *string   `json:"configuration_url,omitempty"`
	Connections      *[]string `json:"connections,omitempty"`
	HwVersion        *string   `json:"hw_version,omitempty"`
	Identifiers      *[]string `json:"identifiers,omitempty"`
	Manufacturer     *string   `json:"manufacturer,omitempty"`
	Model            *string   `json:"model,omitempty"`
	Name             *string   `json:"name,omitempty"`
	SuggestedArea    *string   `json:"suggested_area,omitempty"`
	SwVersion        *string   `json:"sw_version,omitempty"`
	ViaDevice        *string   `json:"via_device,omitempty"`
}

func NewDeviceConfig() *DeviceConfig {
	return &DeviceConfig{}
}

func (conf *DeviceConfig) SetManufacturer(manufacturer string) {
	conf.Manufacturer = &manufacturer
}

func (conf *DeviceConfig) SetModel(model string) {
	conf.Model = &model
}

func (conf *DeviceConfig) SetName(name string) {
	conf.Name = &name
}

func (conf *DeviceConfig) SetSoftwareVersion(version string) {
	conf.SwVersion = &version
}

func (conf *DeviceConfig) AddConnection(connectionType string, connection string) {
	if conf.Connections == nil {
		empty := make([]string, 0)
		conf.Connections = &empty
	}

	*conf.Connections = append(*conf.Connections, connectionType, connection)
}

func (conf *DeviceConfig) AddIdentifier(id string) {
	if conf.Identifiers == nil {
		empty := make([]string, 0)
		conf.Identifiers = &empty
	}

	*conf.Identifiers = append(*conf.Identifiers, id)
}
