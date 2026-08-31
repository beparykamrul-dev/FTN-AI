package devices

type Device struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Kind       string            `json:"kind"`
	Address    string            `json:"address"`
	Status     string            `json:"status"`
	Interfaces []DeviceInterface `json:"interfaces,omitempty"`
}

type DeviceInterface struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
