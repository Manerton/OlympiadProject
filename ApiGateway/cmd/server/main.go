package main

import (
	"flag"
	"log"
	"main/internal/config"
	handlers "main/internal/handlers"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

func getConfigPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "../../config-yaml/config.yaml")
}

func main() {
	cfgPath := flag.String("config", getConfigPath(), "path to config file")
	addr := "0.0.0.0:6611"
	flag.Parse()

	cfg := config.GetConfig(*cfgPath)
	mux := http.NewServeMux()

	for _, route := range cfg.HTTPServer.Routes {
		var handler http.Handler

		if route.Aggregate {
			handler = handlers.NewAggregateHandler(route, cfg.JwtTemp.Key)
		} else {
			handler = handlers.NewProxyHandler(route, cfg.JwtTemp.Key)
		}

		// http.StripPrefix возвращает http.Handler, а mux.Handle принимает тот же тип
		mux.Handle(route.Prefix, http.StripPrefix(strings.TrimSuffix(route.Prefix, "/"), handler))
	}

	// логируем все запросы
	loggedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})

	log.Printf("API Gateway listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, loggedMux))
}
