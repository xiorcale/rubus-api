package services

import (
	"encoding/json"
	"io"

	"net/http"

	"github.com/kjuvi/rubus-api/models"
)

func request(method, url string, body io.Reader) (*http.Response, *models.JSONError) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, models.NewInternalServerError()
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, models.NewInternalServerError()
	}

	return res, nil
}

// GetDevice adds a new Rubus `Device` into the database
func GetDevice(port string) (*models.Device, *models.JSONError) {
	url := "http://rubus_provider:1080/device/" + port

	res, jsonErr := request("GET", url, nil)
	if jsonErr != nil {
		return nil, jsonErr
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, &models.JSONError{
			Status: http.StatusNotFound,
			Error:  "device not found",
		}
	}

	device := models.Device{}
	if err := json.NewDecoder(res.Body).Decode(&device); err != nil {
		return nil, models.NewInternalServerError()
	}
	defer res.Body.Close()

	return &device, nil
}

// GetAllDevices add all the Rubus `Device` into the database
func GetAllDevices() (*[]models.Device, *models.JSONError) {
	url := "http://rubus_provider:1080/device"

	res, jsonErr := request("GET", url, nil)
	if jsonErr != nil {
		return nil, jsonErr
	}

	devices := []models.Device{}
	if err := json.NewDecoder(res.Body).Decode(&devices); err != nil {
		return nil, models.NewInternalServerError()
	}
	defer res.Body.Close()

	return &devices, nil
}

// PowerDeviceOn boots the `Device` on the given `port`
func PowerDeviceOn(port string) *models.JSONError {
	url := "http://rubus_provider:1080/device/" + port + "/on"

	res, jsonErr := request("POST", url, nil)
	if jsonErr != nil {
		return jsonErr
	}

	if res.StatusCode == http.StatusNotFound {
		return &models.JSONError{
			Status: http.StatusNotFound,
			Error:  "device not found",
		}
	}

	return nil
}

// PowerDeviceOff shuts down the `Device` on the given `port`
func PowerDeviceOff(port string) *models.JSONError {
	url := "http://rubus_provider:1080/device/" + port + "/off"

	res, jsonErr := request("POST", url, nil)
	if jsonErr != nil {
		return jsonErr
	}

	if res.StatusCode == http.StatusNotFound {
		return &models.JSONError{
			Status: http.StatusNotFound,
			Error:  "device not found",
		}
	}

	return nil
}
