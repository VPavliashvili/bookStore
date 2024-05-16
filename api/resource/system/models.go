package system

type aboutDTO struct {
	Product       string `json:"product"`
	Author        string `json:"author"`
	Version       string `json:"version"`
	BuildDatetime string `json:"buildDatetime"`
}

type healthDTO struct {
	Healthy bool `json:"healhty"`
}
