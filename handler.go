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
	paths := parseMap(pathsToUrls)

	return func(w http.ResponseWriter, r *http.Request) {
		for _, p := range paths {
			if p.OrigPath == r.URL.Path {
				http.Redirect(w, r, p.RedirectUrl, 301)
			}
		}

		fallback.ServeHTTP(w, r)
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
	var paths []URLShortPath

	err := yaml.Unmarshal(yml, &paths)
	if err != nil {
		return nil, err
	}

	var f = func(w http.ResponseWriter, r *http.Request) {
		for _, p := range paths {
			// fmt.Println(p.OrigPath, r.URL.Path, p.OrigPath == r.URL.Path)

			if p.OrigPath == r.URL.Path {
				http.Redirect(w, r, p.RedirectUrl, 301)
			}
		}

		fallback.ServeHTTP(w, r)
	}

	return f, nil
}

func parseMap(m map[string]string) []URLShortPath {
	res := make([]URLShortPath, len(m))
	for k, v := range m {
		res = append(res, URLShortPath{OrigPath: k, RedirectUrl: v})
	}

	return res
}
