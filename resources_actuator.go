package micro

type Status struct {
	Value string `json:"status"`
}

type ApiInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

func ActuatorRoutes() func(r Router) {
	return func(r Router) {
		r.GET("/", Handle0(handleInfo))
		r.GET("/status", Handle0(handleStatus))
	}
}

func handleStatus(_ *Ctx) (interface{}, error) {
	return &Status{Value: "UP"}, nil
}

func handleInfo(_ *Ctx) (interface{}, error) {
	return ApiInfo{
		Name:    AppName,
		Version: AppVersion,
		Status:  "UP",
	}, nil
}
