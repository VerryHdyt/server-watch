package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	collector Collector
}

func (s Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		m, err := s.collector.Collect()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, m)
	})
	return logging(mux)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	addr := os.Getenv("SERVER_WATCH_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	s := Server{collector: ProcCollector{}}
	log.Printf("server-watch listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, s.routes()))
}
