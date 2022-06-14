package libqatapult

type Device interface {
	// GetCliArgs returns arguments to be passed to the QEMU
	// command line for this device.
	GetCliArgs() ([]string, error)
}

type FilesProvider interface {
	Device

	// GetFiles provides any files needed by the given device to
	// function properly.
	GetFiles() []File
}

// DeviceGroup groups individual devices together.
type DeviceGroup struct {
	devices []Device
}

func (g *DeviceGroup) GetFiles() []File {
	var files []File
	for _, dev := range g.devices {
		if dev == nil {
			continue
		}

		if p, ok := dev.(FilesProvider); ok {
			files = append(files, p.GetFiles()...)
		}
	}
	return files
}

func (g *DeviceGroup) GetCliArgs() (out []string, err error) {
	for _, dev := range g.devices {
		if dev == nil {
			continue
		}
		args, err := dev.GetCliArgs()
		if err != nil {
			return nil, err
		}
		out = append(out, args...)
	}
	return
}

func NewDeviceGroup(devices ...Device) *DeviceGroup {
	return &DeviceGroup{devices: devices}
}
