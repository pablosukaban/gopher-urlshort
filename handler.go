package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v3"
)

type URLShortPath struct {
	OrigPath    string `yaml:"path"`
	RedirectUrl string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// paths := parseMap(pathsToUrls)

	return func(w http.ResponseWriter, r *http.Request) {
		rPath := r.URL.Path

		if p, ok := pathsToUrls[rPath]; ok {
			http.Redirect(w, r, p, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}

	pathsToUrls := makeMap(paths)

	return MapHandler(pathsToUrls, fallback), nil
}

func parseYaml(yml []byte) ([]URLShortPath, error) {
	var paths []URLShortPath
	err := yaml.Unmarshal(yml, &paths)
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func makeMap(paths []URLShortPath) map[string]string {
	var pathsToUrls = make(map[string]string, len(paths))
	for _, p := range paths {
		pathsToUrls[p.OrigPath] = p.RedirectUrl
	}

	return pathsToUrls
}
