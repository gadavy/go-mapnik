//go:build cgo
// +build cgo

package mapnik

// #include <stdlib.h>
// #include "mapnik_c_api.h"
import "C"
import (
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"unsafe"
)

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
	return "Mapnik " + C.GoString(C.mapnik_version_string())
}

func SetLogLevel(level LogLevel) {
	var ll C.int

	switch level {
	case LogLevelNone:
		ll = C.MAPNIK_NONE
	case LogLevelDebug:
		ll = C.MAPNIK_DEBUG
	case LogLevelWarn:
		ll = C.MAPNIK_WARN
	case LogLevelError:
		ll = C.MAPNIK_ERROR
	default:
		ll = C.MAPNIK_WARN
	}

	C.mapnik_logging_set_severity(C.int(ll))
}

func RegisterDatasources(path string) error {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range fileInfos {
		cs := C.CString(filepath.Join(path, file.Name()))
		defer C.free(unsafe.Pointer(cs))

		if C.mapnik_register_datasource(cs) == 0 {
			e := C.GoString(C.mapnik_register_last_error())
			if e != "" {
				return errors.New("registering datasources: " + e)
			}

			return errors.New("error while registering datasources")
		}
	}

	return nil
}

func RegisterFonts(path string) error {
	walk := func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		if !isFontFile(path) {
			return nil
		}

		cs := C.CString(path)
		defer C.free(unsafe.Pointer(cs))

		ce := C.CString("")
		defer C.free(unsafe.Pointer(ce))

		if C.mapnik_register_font(cs, &ce) != 0 {
			if ce != "" {
				return errors.New("registering fonts: " + ce)
			}

			return errors.New("error while registering fonts")
		}

		log.Printf("[mapnik] font %s registered", info.Name())

		return nil
	}

	return filepath.Walk(path, walk)
}

func isFontFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".ttf" ||
		ext == ".otf" ||
		ext == ".woff" ||
		ext == ".ttc" ||
		ext == ".pfa" ||
		ext == ".pfb" ||
		ext == ".dfont"
}

type Map struct {
	m *C.struct__mapnik_map_t
}

func NewMap(width, height uint32) *Map {
	return &Map{C.mapnik_map(C.uint(width), C.uint(height))}
}

func (m *Map) lastError() error {
	return errors.New("mapnik: " + C.GoString(C.mapnik_map_last_error(m.m)))
}

// Load initializes the map by loading its stylesheet from stylesheetFile
func (m *Map) Load(stylesheetFile string) error {
	cs := C.CString(stylesheetFile)
	defer C.free(unsafe.Pointer(cs))

	if C.mapnik_map_load(m.m, cs) != 0 {
		return m.lastError()
	}

	return nil
}

// LoadString initializes the map not from a file but from a stylesheet
// provided as a string.
func (m *Map) LoadString(stylesheet string) error {
	cs := C.CString(stylesheet)
	defer C.free(unsafe.Pointer(cs))

	if C.mapnik_map_load_string(m.m, cs) != 0 {
		return m.lastError()
	}

	return nil
}

func (m *Map) Resize(width, height uint32) {
	C.mapnik_map_resize(m.m, C.uint(width), C.uint(height))
}

func (m *Map) Free() {
	C.mapnik_map_free(m.m)
	m.m = nil
}

func (m *Map) SRS() string {
	return C.GoString(C.mapnik_map_get_srs(m.m))
}

func (m *Map) SetSRS(srs string) {
	cs := C.CString(srs)
	defer C.free(unsafe.Pointer(cs))

	C.mapnik_map_set_srs(m.m, cs)
}

func (m *Map) ZoomAll() error {
	if C.mapnik_map_zoom_all(m.m) != 0 {
		return m.lastError()
	}

	return nil
}

func (m *Map) Zoom(minx, miny, maxx, maxy float64) {
	bbox := C.mapnik_bbox(C.double(minx), C.double(miny), C.double(maxx), C.double(maxy))
	defer C.mapnik_bbox_free(bbox)

	C.mapnik_map_zoom_to_box(m.m, bbox)
}

func (m *Map) RenderToFile(path string) error {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))

	if C.mapnik_map_render_to_file(m.m, cs) != 0 {
		return m.lastError()
	}

	return nil
}

func (m *Map) SetBufferSize(s int) {
	C.mapnik_map_set_buffer_size(m.m, C.int(s))
}

type RenderOpts struct {
	Scale       float64
	ScaleFactor float64
	Format      ImageFormat // Format for the image ('jpeg80', 'png256', etc.)
}

// Render returns the map as an encoded image.
func (m *Map) Render(opts RenderOpts) ([]byte, error) {
	scaleFactor := opts.ScaleFactor
	if scaleFactor == 0.0 {
		scaleFactor = 1.0
	}

	i := C.mapnik_map_render_to_image(m.m, C.double(opts.Scale), C.double(scaleFactor))
	if i == nil {
		return nil, m.lastError()
	}

	defer C.mapnik_image_free(i)

	if opts.Format == "raw" {
		size := 0
		raw := C.mapnik_image_to_raw(i, (*C.size_t)(unsafe.Pointer(&size)))
		return C.GoBytes(unsafe.Pointer(raw), C.int(size)), nil
	}

	var format *C.char

	if opts.Format != "" {
		format = C.CString(opts.Format.String())
	} else {
		format = C.CString("png256")
	}

	b := C.mapnik_image_to_blob(i, format)
	if b == nil {
		return nil, errors.New("mapnik: " + C.GoString(C.mapnik_image_last_error(i)))
	}

	C.free(unsafe.Pointer(format))
	defer C.mapnik_image_blob_free(b)

	return C.GoBytes(unsafe.Pointer(b.ptr), C.int(b.len)), nil
}

func (m *Map) Projection() *Projection {
	return &Projection{
		p: C.mapnik_map_projection(m.m),
	}
}

// Projection from one reference system to the other.
type Projection struct {
	p *C.struct__mapnik_projection_t
}

func (p *Projection) Free() {
	C.mapnik_projection_free(p.p)
	p.p = nil
}

func (p Projection) Forward(x, y float64) (_x, _y float64) {
	c := C.mapnik_coord_t{C.double(x), C.double(y)}
	c = C.mapnik_projection_forward(p.p, c)

	return float64(c.x), float64(c.y)
}
