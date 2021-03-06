// Code generated by "dbusutil-gen -type Manager manager.go"; DO NOT EDIT.

package accounts

func (v *Manager) setPropAllowGuest(value bool) (changed bool) {
	if v.AllowGuest != value {
		v.AllowGuest = value
		v.emitPropChangedAllowGuest(value)
		return true
	}
	return false
}

func (v *Manager) emitPropChangedAllowGuest(value bool) error {
	return v.service.EmitPropertyChanged(v, "AllowGuest", value)
}
