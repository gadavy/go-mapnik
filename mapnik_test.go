package mapnik_test

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gadavy/go-mapnik"
)

func init() {
	if err := mapnik.RegisterDatasources(mapnik.ConfigPlugins()); err != nil {
		log.Fatalf("register datasources: %v", err)
	}

	if err := mapnik.RegisterFonts(mapnik.ConfigFonts()); err != nil {
		log.Fatalf("register fonts: %v", err)
	}
}

func TestRender(t *testing.T) {
	m := mapnik.NewMap(800, 600)
	if err := m.Load("test/map.xml"); err != nil {
		t.Fatal(err)
	}

	if err := m.ZoomAll(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		opts    mapnik.RenderOpts
		wantErr bool
		err     string
	}{
		{
			name: "Pass_PNG256",
			opts: mapnik.RenderOpts{Format: mapnik.Png256},
		},
		{
			name: "Pass_JPG",
			opts: mapnik.RenderOpts{Format: mapnik.Jpeg100},
		},
		{
			name: "Pass_JPG",
			opts: mapnik.RenderOpts{Format: mapnik.Jpeg80},
		},
		{
			name:    "Error_InvalidFormat",
			opts:    mapnik.RenderOpts{Format: "invalid_format"},
			wantErr: true,
			err:     "mapnik: unknown file type: invalid_format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := m.Render(tt.opts)
			if tt.wantErr {
				require.EqualError(t, err, tt.err)
			} else {
				require.NoError(t, err)

				img, _, err := image.Decode(bytes.NewBuffer(data))
				if err != nil {
					t.Fatal(err)
				}

				if img.Bounds().Dx() != 800 || img.Bounds().Dy() != 600 {
					t.Error("unexpected size of output image: ", img.Bounds())
				}
			}
		})
	}
}

func TestConfigFonts(t *testing.T) {

}
