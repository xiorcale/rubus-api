package models

import (
	"net/http"

	"github.com/go-pg/pg/v9"
)

// Device contains the information about a device
type Device struct {
	ID       int64  `json:"id" pg:",pk"`
	Hostname string `json:"hostname"`
	IsTurnOn bool   `json:"isTurnOn"`
	Owner    *int64 `json:"owner" orm:"null"`
}

// AddDevice inserts a new `Device` into the database
func AddDevice(db *pg.DB, d *Device) *JSONError {
	if err := db.Insert(d); err != nil {
		if pgErr, ok := err.(pg.Error); ok && pgErr.IntegrityViolation() {
			return &JSONError{
				Status: http.StatusConflict,
				Error:  "id already exists.",
			}
		}
		return NewInternalServerError()
	}

	return nil
}

// AddDeviceMulti insters multiple `Device` into the database
func AddDeviceMulti(db *pg.DB, devices *[]Device) *JSONError {
	for _, d := range *devices {
		if err := AddDevice(db, &d); err != nil {
			return err
		}
	}
	return nil
}

// GetDevice returns the `Device` with the given `deviceID` from the database
func GetDevice(db *pg.DB, deviceID int64) (*Device, *JSONError) {
	device := &Device{ID: deviceID}
	if err := db.Select(device); err != nil {
		if err == pg.ErrNoRows {
			return nil, &JSONError{
				Status: http.StatusNotFound,
				Error:  "Device does not exist.",
			}
		}
		return nil, NewInternalServerError()
	}

	return device, nil
}

// GetAllDevices returns all the `Device` from the database
func GetAllDevices(db *pg.DB) (devices []*Device, jsonErr *JSONError) {
	if err := db.Model(devices).Select(); err != nil {
		return nil, NewInternalServerError()
	}

	return devices, nil
}

// DeleteDevice removes the given Rubus `Device` from the database
func DeleteDevice(db *pg.DB, uid int64) *JSONError {
	device := &Device{ID: uid}
	if err := db.Delete(device); err != nil {
		if err == pg.ErrNoRows {
			return &JSONError{
				Status: http.StatusNotFound,
				Error:  "device does not exist.",
			}
		}
		return NewInternalServerError()
	}

	return nil
}

// SwitchDevicePower inverse the `isTurnOn` field of the `Device` with the given
// `deviceID`
func SwitchDevicePower(db *pg.DB, device *Device) (*Device, *JSONError) {
	device.IsTurnOn = !device.IsTurnOn

	if err := db.Update(device); err != nil {
		return nil, NewInternalServerError()
	}

	return device, nil
}

// AcquireDevice sets the `User` parameter as the owner of the `Device`
func AcquireDevice(db *pg.DB, device *Device, uid int64) *JSONError {
	device.Owner = &uid

	if err := db.Update(device); err != nil {
		return NewInternalServerError()
	}

	return nil
}

// ReleaseDevice sets the `owner` of the `Device` as nil, if it is the current
// `owner` which is modifying it.
func ReleaseDevice(db *pg.DB, device *Device) *JSONError {
	device.Owner = nil

	if err := db.Update(device); err != nil {
		return NewInternalServerError()
	}

	return nil
}
