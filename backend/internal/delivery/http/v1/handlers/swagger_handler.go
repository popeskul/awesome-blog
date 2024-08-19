package handlers

import (
	"net/http"
	"path/filepath"
)

func SwaggerHandler(staticPath string) http.Handler {
	swaggerUIPath := filepath.Join(staticPath, "swagger-ui")
	fs := http.FileServer(http.Dir(swaggerUIPath))
	return http.StripPrefix("/swagger", fs)
}
