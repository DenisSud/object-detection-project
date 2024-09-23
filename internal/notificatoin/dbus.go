package notification

import (
	"github.com/godbus/dbus/v5"
)

type Notifier struct {
	conn *dbus.Conn
}

func NewNotifier() (*Notifier, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	return &Notifier{conn: conn}, nil
}

func (n *Notifier) Notify(objects []Object) error {
	obj := n.conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0, "ObjectDetector", uint32(0),
		"", "Object Detected", "Detected: "+formatObjects(objects), []string{},
		map[string]dbus.Variant{}, int32(5000))
	return call.Err
}

func formatObjects(objects []Object) string {
	// Format objects into a string for notification
	// Implementation details omitted for brevity
	return ""
}
