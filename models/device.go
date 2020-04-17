package models

import (
	"net/http"
	"strings"

	"github.com/astaxie/beego/orm"
)

// Device contains the information about a device
type Device struct {
	ID       int64  `json:"id" orm:"pk"`
	Hostname string `json:"hostname"`
	IsTurnOn bool   `json:"isTurnOn"`
	Owner    *int64 `json:"owner" orm:"null"`
}

func init() {
	orm.RegisterModel(new(Device))
}

// AddDevice inserts a new `Device` into the database
func AddDevice(d *Device) *JSONError {
	o := orm.NewOrm()

	if _, err := o.Insert(d); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return &JSONError{
				Status: http.StatusConflict,
				Error:  "id already exists",
			}
		}
		return NewInternalServerError()
	}

	return nil
}

// AddDeviceMulti insters multiple `Device` into the database
func AddDeviceMulti(devices *[]Device) *JSONError {
	o := orm.NewOrm()

	if _, err := o.InsertMulti(100, devices); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return &JSONError{
				Status: http.StatusConflict,
				Error:  "id already exists",
			}
		}
		return NewInternalServerError()
	}

	return nil
}

// GetDevice returns the `Device` with the given `deviceID` from the database
func GetDevice(deviceID int64) (*Device, *JSONError) {
	o := orm.NewOrm()

	device := Device{ID: deviceID}
	if err := o.Read(&device); err != nil {
		if err == orm.ErrNoRows {
			return nil, &JSONError{
				Status: http.StatusNotFound,
				Error:  "Device does not exists",
			}
		}
		return nil, NewInternalServerError()
	}

	return &device, nil
}

// GetAllDevices returns all the `Device` from the database
func GetAllDevices() ([]*Device, *JSONError) {
	o := orm.NewOrm()

	devices := []*Device{}
	if _, err := o.QueryTable("device").All(&devices); err != nil {
		return nil, NewInternalServerError()
	}

	return devices, nil
}

// DeleteDevice removes the given Rubus `Device` from the database
func DeleteDevice(uid int64) *JSONError {
	o := orm.NewOrm()

	device := Device{ID: uid}
	uid, err := o.Delete(&device)
	if uid == 0 {
		return &JSONError{
			Status: http.StatusNotFound,
			Error:  "device does not exists",
		}
	}
	if err != nil {
		return NewInternalServerError()
	}

	return nil
}

// SwitchDevicePower inverse the `isTurnOn` field of the `Device` with the given
// `deviceID`
func SwitchDevicePower(device *Device) (*Device, *JSONError) {
	o := orm.NewOrm()

	device.IsTurnOn = !device.IsTurnOn

	if _, err := o.Update(device); err != nil {
		return nil, NewInternalServerError()
	}

	return device, nil
}

// AcquireDevice sets the `User` parameter as the owner of the `Device`
func AcquireDevice(device *Device, uid int64) *JSONError {
	o := orm.NewOrm()

	device.Owner = &uid

	if _, err := o.Update(device); err != nil {
		return NewInternalServerError()
	}

	return nil
}

// ReleaseDevice sets the `owner` of the `Device` as nil, if it is the current
// `owner` which is modifying it.
func ReleaseDevice(device *Device) *JSONError {
	o := orm.NewOrm()

	device.Owner = nil

	if _, err := o.Update(device); err != nil {
		return NewInternalServerError()
	}

	return nil
}
