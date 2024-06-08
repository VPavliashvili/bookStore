package system

import (
	"booksapi/api/database"
	"booksapi/config"
	"encoding/json"
	"fmt"
	"net/http"
)

type API struct {
	compDate string
}

func New(compDate string) API {
	return API{
		compDate: compDate,
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
func (api API) HandleAbout(w http.ResponseWriter, r *http.Request) {
	dto := aboutDTO{
		Product:       "Books Api",
		Author:        "VPavliashvili",
		Version:       "1.0",
		BuildDatetime: api.compDate,
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
func (api API) HandleHealth(w http.ResponseWriter, r *http.Request) {
	dependencies := []dependency{
		{
			Name: func() string {
				db := config.GetAppsettings().Database
				return db.Db
			}(),
			HealthStatus: func() struct {
				Healthy bool   `json:"healthy"`
				Err     string `json:"err"`
			} {
				e := ""
				b := true
				err := database.Ping()
				if err != nil {
					b = false
					e = err.Error()
				}
				return struct {
					Healthy bool   `json:"healthy"`
					Err     string `json:"err"`
				}{
					Healthy: b,
					Err:     e,
				}
			}(),
			Address: func() string {
				db := config.GetAppsettings().Database
				return fmt.Sprintf("%s:%d/%s", db.Host, db.Port, db.Db)
			}(),
		},
	}

	dto := healthDTO{
		Dependencies: dependencies,
	}

	json, _ := json.Marshal(dto)

	fmt.Fprint(w, string(json[:]))
}
