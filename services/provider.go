package services

import (
	"encoding/json"
	"io"

	"net/http"

	"github.com/kjuvi/rubus-api/models"
)

func request(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetDevice adds a new Rubus `Device` into the database
func GetDevice(port string) (device *models.Device, err error) {
	url := "http://rubus_provider:1080/device/" + port

	res, err := request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(res.Body).Decode(&device); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return device, nil
}

// GetAllDevices add all the Rubus `Device` into the database
func GetAllDevices() (devices *[]models.Device, err error) {
	url := "http://rubus_provider:1080/device"

	res, err := request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(res.Body).Decode(&devices); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return devices, nil
}

// PowerDeviceOn boots the `Device` on the given `port`
func PowerDeviceOn(port string) error {
	url := "http://rubus_provider:1080/device/" + port + "/on"

	_, err := request("POST", url, nil)
	if err != nil {
		return err
	}

	return nil
}

// PowerDeviceOff shuts down the `Device` on the given `port`
func PowerDeviceOff(port string) error {
	url := "http://rubus_provider:1080/device/" + port + "/off"

	_, err := request("POST", url, nil)
	if err != nil {
		return err
	}

	return nil
}
