package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/srad/streamsink/app"
	"github.com/srad/streamsink/services"
	"net/http"
)

// TriggerImport godoc
// @Summary     Run once the import of mp4 files in the recordings folder, which are not yet in the system
// @Schemes
// @Description Return a list of channels
// @Tags        admin
// @Accept      json
// @Produce     json
// @Success     200
// @Failure     500 {}  http.StatusInternalServerError
// @Router      /admin/import [post]
func TriggerImport(c *gin.Context) {
	appG := app.Gin{C: c}

	services.StopImport()
	services.StartImport()

	appG.Response(http.StatusOK, nil)
}

// IsImporting godoc
// @Summary     Run once the import of mp4 files in the recordings folder, which are not yet in the system
// @Schemes
// @Description Return a list of channels
// @Tags        admin
// @Accept      json
// @Produce     json
// @Success     200 {bool} Importing flag
// @Router      /admin/importing [post]
func IsImporting(c *gin.Context) {
	appG := app.Gin{C: c}

	appG.Response(http.StatusOK, services.IsImporting())
}