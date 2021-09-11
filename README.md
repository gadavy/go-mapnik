# go-mapnik

Import
------------

```shell
CGO_ENABLED="0" go get gitlab.ozon.ru/osm/go-mapnik
```


Installation
------------

This package requires [Mapnik](http://mapnik.org/) (`libmapnik-dev` on Ubuntu/Debian, `mapnik` in Homebrew).
Make sure `mapnik-config` is in your `PATH`.

You need to set the `CGO_LDFLAGS` and `CGO_CXXFLAGS` environment variables for successful compilation and linking with Mapnik.
Refer to the Makefile how `mapnik-config` can be used to extract the required `CGO_LDFLAGS` and `CGO_CXXFLAGS` values.

Example:
```shell
export CGO_LDFLAGS=$(shell mapnik-config --libs)
export CGO_CXXFLAGS=$(shell mapnik-config --cxxflags --includes --dep-includes | tr '\n' ' ')
```
