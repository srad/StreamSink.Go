package workers

import (
	"github.com/srad/streamsink/utils"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	_, filestring, _, _ = runtime.Caller(0)
	basefolder          = filepath.Dir(filestring)
	file                = filepath.Join(basefolder, "..", "assets", "test.mp4")
	outputPath          = filepath.Join(basefolder, "..", "assets")
)

func TestGetFrameCount(t *testing.T) {
	count, err := utils.GetFrameCount(file)
	if err != nil {
		t.Errorf("error computing framecount: %v", err)
	}

	if count != 1445 {
		t.Errorf("Unexpected frame count")
	}
}

func TestGetVideoInfo(t *testing.T) {
	info, err := utils.GetVideoInfo(file)
	if err != nil {
		t.Errorf("error when getting video duration: %v", err)
	}

	if info.BitRate != 699297 {
		t.Errorf("BitRate wrong: %d", info.BitRate)
	}
	if info.Size != 5251725 {
		t.Errorf("Size wrong: %d", info.Size)
	}
	if info.Fps != 24.03846153846154 {
		t.Errorf("Fps wrong: %f", info.Fps)
	}
	if info.Duration != 60.08 {
		t.Errorf("Duration wrong: %f", info.Duration)
	}
	if info.Height != 360 {
		t.Errorf("Height wrong: %d", info.Height)
	}
	if info.Width != 640 {
		t.Errorf("Width wrong: %d", info.Width)
	}
}