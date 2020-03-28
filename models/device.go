package models

import (
	"errors"
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
func AddDevice(d *Device) error {
	o := orm.NewOrm()

	if _, err := o.Insert(d); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("ID already exists")
		}
		return errors.New("Internal Server Error")
	}

	return nil
}

// AddDeviceMulti insters multiple `Device` into the database
func AddDeviceMulti(devices *[]Device) error {
	o := orm.NewOrm()

	if _, err := o.InsertMulti(100, devices); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("ID already exists")
		}
		return errors.New("Internal Server Error")
	}

	return nil
}

// GetDevice returns the `Device` with the given `deviceID` from the database
func GetDevice(deviceID int64) (d *Device, err error) {
	o := orm.NewOrm()

	device := Device{ID: deviceID}
	if err = o.Read(&device); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.New("Device does not exists")
		}
		return nil, errors.New("Internal Server Error")
	}

	return &device, nil
}

// GetAllDevices returns all the `Device` from the database
func GetAllDevices() (devices []*Device, err error) {
	o := orm.NewOrm()

	if _, err = o.QueryTable("device").All(&devices); err != nil {
		return nil, errors.New("Internal Server Error")
	}

	return devices, nil
}

// SwitchDevicePower inverse the `isTurnOn` field of the `Device` with the given
// `deviceID`
func SwitchDevicePower(deviceID int64) (device *Device, err error) {
	o := orm.NewOrm()

	device, err = GetDevice(deviceID)
	if err != nil {
		return nil, err
	}

	device.IsTurnOn = !device.IsTurnOn

	if _, err = o.Update(device); err != nil {
		return nil, errors.New("Internal Server Error")
	}

	return device, nil
}

// AcquireDevice sets the `User` parameter as the owner of the `Device`
func AcquireDevice(device *Device, uid int64) error {
	o := orm.NewOrm()

	device.Owner = &uid

	if _, err := o.Update(device); err != nil {
		return errors.New("Internal Server Error")
	}

	return nil
}

// ReleaseDevice sets the `owner` of the `Device` as nil, if it is the current
// `owner` which is modifying it.
func ReleaseDevice(device *Device) error {
	o := orm.NewOrm()

	device.Owner = nil

	if _, err := o.Update(device); err != nil {
		return errors.New("Internal Server Error")
	}

	return nil
}
