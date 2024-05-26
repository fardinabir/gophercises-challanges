package handler

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.RequestURI
		if value, ok := pathsToUrls[path]; ok {
			// Perform the redirection
			http.Redirect(w, r, value, http.StatusFound) // StatusFound is HTTP 302
			return
		} else {
			fallback.ServeHTTP(w, r)
			return
		}
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	// Convert the slice of structs into a map
	pathToURL := buildMap(parsedYaml)
	return MapHandler(pathToURL, fallback), nil
}

func parseYaml(yml []byte) ([]map[string]string, error) {
	ymlsMap := []map[string]string{}

	err := yaml.Unmarshal(yml, &ymlsMap)
	if err != nil {
		return nil, err
	}
	return ymlsMap, nil
}

func buildMap(ymlMap []map[string]string) map[string]string {
	pathToUrl := make(map[string]string)
	for _, mp := range ymlMap {
		pathToUrl[mp["path"]] = mp["url"]
	}
	return pathToUrl
}
