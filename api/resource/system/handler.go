package system

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type API struct {
	buildTime string
}

func New(buildTime string) API {
	return API{
		buildTime: buildTime,
	}
}

// HandleAbout godoc
//
//	@Summary		Gives system info about api
//	@Description	api info
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	aboutDTO
//	@Router			/api/system/about [get]
func (api *API) HandleAbout(w http.ResponseWriter, r *http.Request) {
	dto := aboutDTO{
		Product:       "Books Api",
		Author:        "VPavliashvili",
		Version:       "1.0",
		BuildDatetime: api.buildTime,
	}

	json, _ := json.Marshal(dto)

	fmt.Fprint(w, string(json[:]))
}

// HandleHealth godoc
//
//	@Summary		Gives system health status
//	@Description	api health
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	healthDTO
//	@Router			/api/system/health [get]
func (api *API) HandleHealth(w http.ResponseWriter, r *http.Request) {
	dto := healthDTO{
		Healthy: true, // TO DO
	}

	json, _ := json.Marshal(dto)

	fmt.Fprint(w, string(json[:]))
}
