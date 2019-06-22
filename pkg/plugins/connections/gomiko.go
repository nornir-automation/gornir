package connections

import (
	"github.com/Ali-aqrabawi/gomiko/pkg/types"
	"gomiko/pkg"
)

type GomikoConn struct {
	Host       string
	Username   string
	Password   string
	DeviceType string
	device     types.Device
}

func (c *GomikoConn) Open() error {
	c.device = gomiko.NewDevice(c.Host, c.Username, c.Password, c.DeviceType)
	return c.device.Connect()
}

func (c *GomikoConn) Close() {
	c.device.Disconnect()
}

func (c *GomikoConn) Send(cmd string) (RemoteCommandResults, error) {
	var results RemoteCommandResults
	stdout, err := c.device.SendCommand(cmd)
	if err != nil {
		return results, err
	}
	stderr := "" // gomiko can't generate stderr as it uses shell session.
	results.Stdout = stdout
	results.Stderr = stderr
	return results, nil
}
