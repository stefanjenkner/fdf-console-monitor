package serialmonitor

type capture struct {
	elapsedTime      uint16
	distance         uint16
	time500m         uint16
	strokesPerMinute uint8
	watts            uint16
	cals             uint16
	level            uint8
}
