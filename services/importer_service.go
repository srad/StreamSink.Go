package services

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/srad/streamsink/conf"
	"github.com/srad/streamsink/database"
	"github.com/srad/streamsink/helpers"
	"os"
	"path/filepath"
)

var (
	ctx, cancelImport     = context.WithCancel(context.Background())
	importing             = false
	importSize        int = 0
	importProgress    int = 0
)

func StartImport() {
	go runImport()
}

func StopImport() {
	cancelImport()
}

func IsImporting() bool {
	return importing
}

func GetImportProgress() (int, int) {
	return importProgress, importSize
}

func runImport() {
	importing = true

	if err := ImportChannels(ctx); err != nil {
		log.Errorln(err)
	}

	importing = false
}

// ImportChannels Imports folders and videos found on disk.
//
// 1. Import all folders as channels found in the recording path.
// 2. If the folder contains the channel.json backup file, then reconstruct the channel information from this file.
// 3. Then search on each folder for media files to import as recordings.
// 4. If the recordings do not contain previews, their creation will be scheduled.
func ImportChannels(context.Context) error {
	cfg := conf.Read()

	log.Infoln("------------------------------------------------------------------------------------------")
	log.Infof("Scanning file system for media: %s", cfg.RecordingsAbsolutePath)
	log.Infoln("------------------------------------------------------------------------------------------")

	recordingFolder, err := os.Open(cfg.RecordingsAbsolutePath)
	if err != nil {
		return fmt.Errorf("failed opening directory '%s': %s", cfg.RecordingsAbsolutePath, err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Errorf("error closing folder %s", file.Name())
		}
	}(recordingFolder)

	// ---------------------------------------------------------------------------------
	// Traverse folders (channels)
	// ---------------------------------------------------------------------------------
	channelFolders, _ := recordingFolder.Readdirnames(0)
	importSize = len(channelFolders)
	importProgress = 0
	for _, name := range channelFolders {
		importProgress++
		channelName := database.ChannelName(name)
		log.Infof("Import/%s (%d/%d)] Scanning folder", channelName, importProgress, importSize)
		// Is no directory, skip
		if dir, err := os.Stat(channelName.AbsoluteChannelPath()); err != nil || !dir.IsDir() {
			continue
		}

		log.Infof("[Import/%s (%d/%d)] Reading folder", channelName, importProgress, importSize)

		newChannel, err4 := database.CreateChannel(channelName, channelName.String(), "https://"+channelName.String())
		if err4 != nil {
			log.Errorf("[Import/%s (%d/%d)] Error adding %s", channelName, importProgress, importSize, err4)
		}

		// ---------------------------------------------------------------------------------
		// Import individual files
		// ---------------------------------------------------------------------------------
		log.Infof("[Import/%s (%d/%d)] Import individual files ...", channelName, importProgress, len(channelFolders))
		files, err2 := os.ReadDir(channelName.AbsoluteChannelPath())
		if err2 != nil {
			log.Errorf("[Import/%s] Error reading: %s", channelName, err2)
			continue
		}

		// ---------------------------------------------------------------------------------
		// Traverse all mp4 files and add to models if not existent
		// ---------------------------------------------------------------------------------
		var j = 0
		log.Infof("[Import/%s (%d/%d)] Traverse all mp4 files and add to models if not existent (files: %d) ...", channelName, importProgress, importSize, len(files))
		for _, file := range files {
			j++
			mp4File := !file.IsDir() && filepath.Ext(file.Name()) == ".mp4"
			if !mp4File {
				continue
			}

			log.Infof("[Import/%s (%d/%d) (%d/%d)] Checking file: %s", channelName, importProgress, importSize, j, len(files), file.Name())

			recording := database.Recording{ChannelId: newChannel.ChannelId, ChannelName: channelName, Filename: database.RecordingFileName(file.Name())}

			video := &helpers.Video{FilePath: channelName.AbsoluteChannelFilePath(database.RecordingFileName(file.Name()))}

			if _, errVideoInfo := video.GetVideoInfo(); errVideoInfo != nil {
				log.Errorf("[Import/%s] File '%s' seems corrupted, deleting: %s", channelName, file.Name(), errVideoInfo)
				if errDestroy := recording.Destroy(); errDestroy != nil {
					log.Errorf("[Import/%s] Error deleting: %s: %s", channelName, file.Name(), errDestroy)
				} else {
					log.Infof("[Import/%s] Deleted: %s", channelName, file.Name())
				}
				continue
			}

			// File seems ok, try to add.
			newRecording, errAdd := database.AddIfNotExists(newChannel.ChannelId, newChannel.ChannelName, database.RecordingFileName(file.Name()))
			if errAdd != nil {
				log.Errorf("[Import/%s] Error: %s", channelName, errAdd)
				continue
			}

			// ---------------------------------------------------------------------------------
			// Not new record inserted and therefore not automatically new previews generated.
			// So check if the files exist and if not generate them.
			// Create preview if any not existent
			// ---------------------------------------------------------------------------------
			if database.PreviewsExist(newRecording.ChannelName, newRecording.Filename) {
				log.Infof("[Import/%s] Preview files exist", channelName)
				newRecording.RecordingId.AddPreviews()
			} else {
				log.Infof("[Import/%s] Adding job for: %s", channelName, file.Name())
				EnqueuePreviewJob(newRecording.RecordingId)
			}
		}
	}

	return nil
}