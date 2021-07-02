package mlog

import (
	"io"
	"path"

	"github.com/arthurkiller/rollingwriter"
)

type RollingWriter struct {
	config rollingwriter.Config
	Writer io.Writer
}

func NewRollingWriter(logPath, fileName string) (writer *RollingWriter, err error) {
	writer = &RollingWriter{}
	writer.config = rollingwriter.Config{
		TimeTagFormat:     "20060102T150405",
		LogPath:           path.Dir(logPath),
		FileName:          path.Base(fileName),
		RollingPolicy:     rollingwriter.VolumeRolling,
		RollingVolumeSize: "50MB",
		Compress:          false,
		WriterMode:        "lock",
		MaxRemain:         5,
	}
	writer.Writer, err = rollingwriter.NewWriterFromConfig(&writer.config)
	return writer, err
}

func (w *RollingWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
