//go:build !cgo
// +build !cgo

package mapnik

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
	Jpeg80 ImageFormat = "jpeg80"
)

func (f ImageFormat) String() string {
	return string(f)
}

func Version() string {
	return "Mapnik disabled"
}

func SetLogLevel(level LogLevel) {}

func RegisterDatasources(path string) error {
	return nil
}

func RegisterFonts(path string) error {
	return nil
}

type Map struct{}

func NewMap(width, height uint32) *Map {
	return &Map{}
}

// Load initializes the map by loading its stylesheet from stylesheetFile
func (m *Map) Load(stylesheetFile string) error {
	return nil
}

// LoadString initializes the map not from a file but from a stylesheet
// provided as a string.
func (m *Map) LoadString(stylesheet string) error {
	return nil
}

func (m *Map) Resize(width, height uint32) {}

func (m *Map) Free() {}

func (m *Map) SRS() string {
	return ""
}

func (m *Map) SetSRS(srs string) {}

func (m *Map) ZoomAll() error {
	return nil
}

func (m *Map) Zoom(minx, miny, maxx, maxy float64) {}

func (m *Map) RenderToFile(path string) error {
	return nil
}

func (m *Map) SetBufferSize(s int) {}

type RenderOpts struct {
	Scale       float64
	ScaleFactor float64
	Format      ImageFormat // Format for the image ('jpeg80', 'png256', etc.)
}

// Render returns the map as an encoded image.
func (m *Map) Render(opts RenderOpts) ([]byte, error) {
	return []byte{}, nil
}

func (m *Map) Projection() *Projection {
	return &Projection{}
}

// Projection from one reference system to the other.
type Projection struct{}

func (p *Projection) Free() {}

func (p Projection) Forward(x, y float64) (_x, _y float64) {
	return x, y
}
