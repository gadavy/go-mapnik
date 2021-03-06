#include <mapnik/version.hpp>
#include <mapnik/map.hpp>
#include <mapnik/layer.hpp>
#include <mapnik/datasource.hpp>
#include <mapnik/color.hpp>
#include <mapnik/image_util.hpp>
#include <mapnik/agg_renderer.hpp>
#include <mapnik/load_map.hpp>
#include <mapnik/datasource_cache.hpp>
#include <mapnik/projection.hpp>
#include <mapnik/font_engine_freetype.hpp>

#if MAPNIK_VERSION >= 300000
#include <mapnik/image.hpp>
#define mapnik_image_type mapnik::image_rgba8
#else
#include <mapnik/graphics.hpp>
#define mapnik_image_type mapnik::image_32
#endif

#include "mapnik_c_api.h"

#include <stdlib.h>

#ifdef __cplusplus
extern "C"
{
#endif

    int mapnik_register_font(const char *path, char **err)
    {
        try
        {
            mapnik::freetype_engine::register_font(path);

            return 0;
        }
        catch (std::exception const &ex)
        {
            *err = strdup(ex.what());

            return -1;
        }
    }

    int mapnik_register_datasource(const char *path, char **err)
    {
        try
        {
#if MAPNIK_VERSION >= 200200
        mapnik::datasource_cache::instance().register_datasources(path);
#else
        mapnik::datasource_cache::instance()->register_datasources(path);
#endif
        }
        catch (std::exception const &ex)
        {
            *err = strdup(ex.what());

            return -1;
        }

        return 0;
    }

    struct _mapnik_map_t
    {
        mapnik::Map *m;
        std::string *err;
    };

    mapnik_map_t *mapnik_map(unsigned width, unsigned height)
    {
        mapnik_map_t *map = new mapnik_map_t;
        map->m = new mapnik::Map(width, height);
        map->err = NULL;

        return map;
    }

    void mapnik_map_free(mapnik_map_t *m)
    {
        if (m)
        {
            if (m->m)
                delete m->m;
            if (m->err)
                delete m->err;
            delete m;
        }
    }

    inline void mapnik_map_reset_last_error(mapnik_map_t *m)
    {
        if (m && m->err)
        {
            delete m->err;
            m->err = NULL;
        }
    }

    void mapnik_map_set_max_connections(mapnik_map_t *m, int num_threads)
    {
        unsigned int i;
        char *tmp = (char *)malloc(20);

        for (i = 0; i < m->m->layer_count(); i++)
        {
#if MAPNIK_VERSION >= 300000
            mapnik::layer &l = m->m->get_layer(i);
#else
            mapnik::layer &l = m->m->getLayer(i);
#endif
            mapnik::parameters params = l.datasource()->params();

            if (params.find("max_size") == params.end())
            {
                sprintf(tmp, "%i", num_threads + 2);
                params["max_size"] = std::string(tmp);
            }

#if MAPNIK_VERSION >= 200200
        l.set_datasource(mapnik::datasource_cache::instance().create(params));
#else
        l.set_datasource(mapnik::datasource_cache::instance()->create(params));
#endif
        }

        free(tmp);
    }

    const char *mapnik_map_get_srs(mapnik_map_t *m)
    {
        if (m && m->m)
            return m->m->srs().c_str();

        return NULL;
    }

    int mapnik_map_set_srs(mapnik_map_t *m, const char *srs)
    {
        if (m)
        {
            m->m->set_srs(srs);

            return 0;
        }

        return -1;
    }

    int mapnik_map_load(mapnik_map_t *m, const char *stylesheet)
    {
        mapnik_map_reset_last_error(m);

        if (m && m->m)
        {
            try
            {
                mapnik::load_map(*m->m, stylesheet);
            }
            catch (std::exception const &ex)
            {
                m->err = new std::string(ex.what());

                return -1;
            }

            return 0;
        }

        return -1;
    }

    int mapnik_map_load_string(mapnik_map_t *m, const char *stylesheet_string)
    {
        mapnik_map_reset_last_error(m);

        if (m && m->m)
        {
            try
            {
                mapnik::load_map_string(*m->m, stylesheet_string);
            }
            catch (std::exception const &ex)
            {
                m->err = new std::string(ex.what());

                return -1;
            }

            return 0;
        }

        return -1;
    }

    int mapnik_map_zoom_all(mapnik_map_t *m)
    {
        mapnik_map_reset_last_error(m);

        if (m && m->m)
        {
            try
            {
                m->m->zoom_all();
            }
            catch (std::exception const &ex)
            {
                m->err = new std::string(ex.what());

                return -1;
            }

            return 0;
        }

        return -1;
    }

    int mapnik_map_render_to_file(mapnik_map_t *m, const char *filepath)
    {
        mapnik_map_reset_last_error(m);

        if (m && m->m)
        {
            try
            {
                mapnik_image_type buf(m->m->width(), m->m->height());
                mapnik::agg_renderer<mapnik_image_type> ren(*m->m, buf);
                ren.apply();

                mapnik::save_to_file(buf, filepath);
            }
            catch (std::exception const &ex)
            {
                m->err = new std::string(ex.what());

                return -1;
            }

            return 0;
        }

        return -1;
    }

    void mapnik_map_resize(mapnik_map_t *m, unsigned int width, unsigned int height)
    {
        if (m && m->m)
        {
            m->m->resize(width, height);
        }
    }

    MAPNIKCAPICALL void mapnik_map_set_buffer_size(mapnik_map_t *m, int buffer_size)
    {
        m->m->set_buffer_size(buffer_size);
    }

    const char *mapnik_map_last_error(mapnik_map_t *m)
    {
        if (m && m->err)
        {
            return m->err->c_str();
        }

        return NULL;
    }

    struct _mapnik_projection_t
    {
        mapnik::projection *p;
    };

    mapnik_projection_t *mapnik_map_projection(mapnik_map_t *m)
    {
        mapnik_projection_t *proj = new mapnik_projection_t;
        if (m && m->m)
            proj->p = new mapnik::projection(m->m->srs());
        else
            proj->p = NULL;

        return proj;
    }

    void mapnik_projection_free(mapnik_projection_t *p)
    {
        if (p)
        {
            if (p->p)
                delete p->p;
            delete p;
        }
    }

    mapnik_coord_t mapnik_projection_forward(mapnik_projection_t *p, mapnik_coord_t c)
    {
        if (p && p->p)
        {
            p->p->forward(c.x, c.y);
        }

        return c;
    }

    struct _mapnik_bbox_t
    {
        mapnik::box2d<double> b;
    };

    mapnik_bbox_t *mapnik_bbox(double minx, double miny, double maxx, double maxy)
    {
        mapnik_bbox_t *b = new mapnik_bbox_t;
        b->b = mapnik::box2d<double>(minx, miny, maxx, maxy);

        return b;
    }

    void mapnik_bbox_free(mapnik_bbox_t *b)
    {
        if (b)
            delete b;
    }

    void mapnik_map_zoom_to_box(mapnik_map_t *m, mapnik_bbox_t *b)
    {
        if (m && m->m && b)
        {
            m->m->zoom_to_box(b->b);
        }
    }

    struct _mapnik_image_t
    {
        mapnik_image_type *i;
        std::string *err;
    };

    inline void mapnik_image_reset_last_error(mapnik_image_t *i)
    {
        if (i && i->err)
        {
            delete i->err;
            i->err = NULL;
        }
    }

    void mapnik_image_free(mapnik_image_t *i)
    {
        if (i)
        {
            if (i->i)
                delete i->i;
            delete i;
        }
    }

    const char *mapnik_image_last_error(mapnik_image_t *i)
    {
        if (i && i->err)
        {
            return i->err->c_str();
        }

        return NULL;
    }

    mapnik_image_t *mapnik_map_render_to_image(mapnik_map_t *m, double scale, double scale_factor)
    {
        mapnik_map_reset_last_error(m);

        mapnik_image_type *im = new mapnik_image_type(m->m->width(), m->m->height());
        if (m && m->m)
        {
            try
            {
                mapnik::agg_renderer<mapnik_image_type> ren(*m->m, *im, scale_factor);
                if (scale > 0.0)
                {
                    ren.apply(scale);
                }
                else
                {
                    ren.apply();
                }
            }
            catch (std::exception const &ex)
            {
                delete im;
                m->err = new std::string(ex.what());
                return NULL;
            }
        }

        mapnik_image_t *i = new mapnik_image_t;
        i->i = im;
        i->err = NULL;

        return i;
    }

    void mapnik_image_blob_free(mapnik_image_blob_t *b)
    {
        if (b)
        {
            if (b->ptr)
                delete[] b->ptr;

            delete b;
        }
    }

    mapnik_image_blob_t *mapnik_image_to_blob(mapnik_image_t *i, const char *format)
    {
        mapnik_image_reset_last_error(i);

        mapnik_image_blob_t *blob = new mapnik_image_blob_t;
        blob->ptr = NULL;
        blob->len = 0;

        if (i && i->i)
        {
            try
            {
                std::string s = save_to_string(*(i->i), format);
                blob->len = s.length();
                blob->ptr = new char[blob->len];

                memcpy(blob->ptr, s.c_str(), blob->len);
            }
            catch (std::exception const &ex)
            {
                i->err = new std::string(ex.what());
                delete blob;

                return NULL;
            }
        }

        return blob;
    }

    const uint8_t *mapnik_image_to_raw(mapnik_image_t *i, size_t *size)
    {
        if (i && i->i)
        {
            *size = i->i->width() * i->i->height() * 4;

            return (uint8_t *)i->i->data();
        }

        return NULL;
    }

    mapnik_image_t *mapnik_image_from_raw(const uint8_t *raw, int width, int height)
    {
        mapnik_image_t *img = new mapnik_image_t;
        img->i = new mapnik_image_type(width, height);
        memcpy(img->i->data(), raw, width * height * 4);
        img->err = NULL;

        return img;
    }

    const char *mapnik_version_string()
    {
#if MAPNIK_VERSION >= 200100
        return MAPNIK_VERSION_STRING;
#else
#define MAPNIK_C_API_STRINGIFY(n) #n
        return "ABI " MAPNIK_C_API_STRINGIFY(MAPNIK_VERSION);
#endif
    }

    void mapnik_logging_set_severity(int level)
    {
        mapnik::logger::severity_type severity;

        switch (level)
        {
        case MAPNIK_DEBUG:
            severity = mapnik::logger::debug;
            break;
        case MAPNIK_WARN:
            severity = mapnik::logger::warn;
            break;
        case MAPNIK_ERROR:
            severity = mapnik::logger::error;
            break;
        default:
            severity = mapnik::logger::none;
            break;
        }

        mapnik::logger::instance().set_severity(severity);
    }

#ifdef __cplusplus
}
#endif
