package health

type Status struct {
	Status string `json:"status"`
}

func OK() Status { return Status{Status: "ok"} }
