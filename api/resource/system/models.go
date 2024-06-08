package system

type aboutDTO struct {
	Product       string `json:"product"`
	Author        string `json:"author"`
	Version       string `json:"version"`
	BuildDatetime string `json:"buildDatetime"`
}

type dependency struct {
	Name         string `json:"name"`
	HealthStatus struct {
		Healthy bool   `json:"healthy"`
		Err     string `json:"err"`
	} `json:"healthStatus"`
	Address string `json:"address"`
}

type healthDTO struct {
    Dependencies []dependency `json:"dependencies"`
}
