package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/srad/streamsink/app"
	"github.com/srad/streamsink/services"
	"log"
	"net/http"
)

// IsRecording godoc
// @Summary     Return if server is current recording
// @Description Return if server is current recording.
// @Tags        recorder
// @Produce     plain
// @Success     200 {} bool "true if currently recording, otherwise false."
// @Success     500
// @Router      /recorder [get]
func IsRecording(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, services.IsRecording())
}

// StopRecorder godoc
// @Summary     StopRecorder server recording
// @Tags        recorder
// @Success     200
// @Router      /recorder/pause [post]
func StopRecorder(c *gin.Context) {
	appG := app.Gin{C: c}

	log.Println("Pausing recorder")
	go services.StopRecorder()

	appG.Response(http.StatusOK, nil)
}

// StartRecorder godoc
// @Summary     StartRecorder server recording
// @Tags        recorder
// @Success     200
// @Router      /recorder/resume [post]
func StartRecorder(c *gin.Context) {
	appG := app.Gin{C: c}

	log.Println("Resuming recorder")
	services.StartRecorder()
	appG.Response(http.StatusOK, nil)
}