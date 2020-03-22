package models

import (
	"encoding/json"
	"io"
)

// Device contains the information about a device
type Device struct {
	ID       int64  `json:"id"`
	Hostname string `json:"hostname"`
	IsTurnOn bool   `json:"isTurnOn"`
}

// Bind transform the given `requestBody` int a `Device`
func (d *Device) Bind(responseBody io.ReadCloser) error {
	decoder := json.NewDecoder(responseBody)
	if err := decoder.Decode(d); err != nil {
		return err
	}

	return nil
}
