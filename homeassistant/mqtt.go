package homeassistant

import "fmt"

func AutoDiscoveryTopic(entityType EntityType, deviceId string, sensorName string) string {
	return fmt.Sprintf("homeassistant/%s/%s/%s/config", entityType, deviceId, sensorName)
}
