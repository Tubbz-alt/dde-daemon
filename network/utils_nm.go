package main

import nm "dbus/org/freedesktop/networkmanager"
import "dlib/dbus"

func nmNewDevice(devPath dbus.ObjectPath) (dev *nm.Device, err error) {
	dev, err = nm.NewDevice(NMDest, devPath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmNewAccessPoint(apPath dbus.ObjectPath) (ap *nm.AccessPoint, err error) {
	ap, err = nm.NewAccessPoint(NMDest, apPath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmNewActiveConnection(apath dbus.ObjectPath) (ac *nm.ActiveConnection, err error) {
	ac, err = nm.NewActiveConnection(NMDest, apath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmGetDevices() (devPaths []dbus.ObjectPath, err error) {
	devPaths, err = NMManager.GetDevices()
	if err != nil {
		Logger.Error(err)
	}
	return
}

func nmNewSettingsConnection(cpath dbus.ObjectPath) (conn *nm.SettingsConnection, err error) {
	conn, err = nm.NewSettingsConnection(NMDest, cpath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmGetDeviceInterface(devPath dbus.ObjectPath) (devInterface string) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	devInterface = dev.Interface.Get()
	return
}

func nmAddAndActivateConnection(data _ConnectionData, devPath dbus.ObjectPath) (cpath, apath dbus.ObjectPath, err error) {
	spath := dbus.ObjectPath("/")
	cpath, apath, err = NMManager.AddAndActivateConnection(data, devPath, spath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmActivateConnection(cpath, devPath dbus.ObjectPath) (apath dbus.ObjectPath, err error) {
	spath := dbus.ObjectPath("/")
	apath, err = NMManager.ActivateConnection(cpath, devPath, spath)
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmGetActiveConnections() (apath []dbus.ObjectPath) {
	apath = NMManager.ActiveConnections.Get()
	return
}

func nmGetActiveConnectionByUuid(uuid string) (apath dbus.ObjectPath, ok bool) {
	for _, apath = range nmGetActiveConnections() {
		if ac, err := nmNewActiveConnection(apath); err == nil {
			if ac.Uuid.Get() == uuid {
				ok = true
				return
			}
		}
	}
	ok = false
	return
}

func nmGetConnectionData(cpath dbus.ObjectPath) (data _ConnectionData, err error) {
	nmConn, err := nm.NewSettingsConnection(NMDest, cpath)
	if err != nil {
		Logger.Error(err)
		return
	}
	data, err = nmConn.GetSettings()
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmGetConnectionUuid(cpath dbus.ObjectPath) (uuid string) {
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	uuid = getSettingConnectionUuid(data)
	if len(uuid) == 0 {
		Logger.Error("get uuid of connection failed, uuid is empty")
	}
	return
}

func nmGetConnectionType(cpath dbus.ObjectPath) (ctype string) {
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	ctype = getSettingConnectionType(data)
	if len(ctype) == 0 {
		Logger.Error("get type of connection failed, type is empty")
	}
	return
}

func nmGetConnectionList() (connections []dbus.ObjectPath) {
	connections, err := NMSettings.ListConnections()
	if err != nil {
		Logger.Error(err)
		return
	}
	return
}

func nmGetConnectionById(id string) (cpath dbus.ObjectPath, ok bool) {
	for _, cpath = range nmGetConnectionList() {
		data, err := nmGetConnectionData(cpath)
		if err != nil {
			continue
		}
		if getSettingConnectionId(data) == id {
			ok = true
			return
		}
	}
	ok = false
	return
}

func nmGetConnectionByUuid(uuid string) (cpath dbus.ObjectPath, ok bool) {
	for _, cpath = range nmGetConnectionList() {
		data, err := nmGetConnectionData(cpath)
		if err != nil {
			continue
		}
		if getSettingConnectionUuid(data) == uuid {
			ok = true
			return
		}
	}
	ok = false
	return
}

func nmGetWirelessConnectionBySsid(ssid []byte) (cpath dbus.ObjectPath, ok bool) {
	for _, cpath = range nmGetConnectionList() {
		data, err := nmGetConnectionData(cpath)
		if err != nil {
			continue
		}
		if isSettingWirelessSsidExists(data) && string(getSettingWirelessSsid(data)) == string(ssid) {
			ok = true
			return
		}
	}
	ok = false
	return
}

func nmAddConnection(data _ConnectionData) {
	_, err := NMSettings.AddConnection(data)
	if err != nil {
		Logger.Error(err)
	}
	return
}
