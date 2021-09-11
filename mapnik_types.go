package mapnik

import (
	"errors"
)

var (
	_ = Version
	_ = SetLogLevel
	_ = RegisterDatasources
	_ = RegisterFonts
	_ = NewMap
	_ = new(Map).SetMaxConnections
	_ = new(Map).Load
	_ = new(Map).LoadString
	_ = new(Map).Resize
	_ = new(Map).Free
	_ = new(Map).SRS
	_ = new(Map).SetSRS
	_ = new(Map).ZoomAll
	_ = new(Map).Zoom
	_ = new(Map).RenderToFile
	_ = new(Map).SetBufferSize
	_ = new(Map).Render
	_ = new(Map).Projection
	_ = new(Projection).Free
	_ = new(Projection).Forward
)

var ErrMapnikError = errors.New("mapnik")

type LogLevel int

const (
	LogLevelNone  LogLevel = 0
	LogLevelDebug LogLevel = 1
	LogLevelWarn  LogLevel = 2
	LogLevelError LogLevel = 3
)

type ImageFormat string

const (
	Png256 ImageFormat = "png256"

	Jpeg100 ImageFormat = "jpeg100"
	Jpeg80  ImageFormat = "jpeg80"
)

func (f ImageFormat) String() string {
	return string(f)
}

type RenderOpts struct {
	Scale       float64
	ScaleFactor float64
	Format      ImageFormat // Format for the image ('jpeg80', 'png256', etc.)
}
