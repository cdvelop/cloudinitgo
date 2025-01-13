package cloudinitgo

import (
	"fmt"
	"net/http"
	"path/filepath"
)

func (s *cloudInit) createMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve YML files directly
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(s.Config.DataDir, "user-data"))
			return
		}
		http.ServeFile(w, r, filepath.Join(s.Config.DataDir, r.URL.Path))
	})
	return mux
}

func (s *cloudInit) Start() error {
	if err := s.CreateConfigFiles(); err != nil {
		return err
	}

	mux := s.createMux()
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: mux,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	return nil
}

func (s *cloudInit) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}
