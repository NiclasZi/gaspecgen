package server

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/NiclasZi/gaspecgen/db"
	"github.com/NiclasZi/gaspecgen/pkg/generator"
	"github.com/NiclasZi/gaspecgen/pkg/loader"
	"github.com/NiclasZi/gaspecgen/pkg/renderer"
	"github.com/Phillezi/common/utils/or"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	defaultMaxMemoryUploadBytes = 10 << 20 // 100mb
)

//go:embed embed/index.html
var indexHTML string

type Server struct {
	host       string
	httpServer *http.Server
	ctx        context.Context
	cancel     context.CancelFunc

	maxMemoryUploadBytes int64

	l *zap.Logger
}

func New(ctx context.Context, port int, hostOpt ...string) *Server {
	ctxx, cancel := context.WithCancel(ctx)

	s := &Server{
		host:   or.Or(or.Or(hostOpt...), "localhost"),
		ctx:    ctxx,
		cancel: cancel,
		l:      zap.L().Named("[SERVER]"),
	}

	router := mux.NewRouter()

	// API endpoints
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/query", s.handleStreamTransform).Methods("POST")

	router.PathPrefix("/").Handler(func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, indexHTML)
		}
	}())

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, port),
		Handler: router,
	}

	return s
}

func (s *Server) Addr() string {
	return fmt.Sprintf("http://%s", s.httpServer.Addr)
}

func (s *Server) Start() error {
	errCh := make(chan error, 1)

	// Start server
	go func() {
		s.l.Info("Starting server on", zap.String("address", s.httpServer.Addr))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			s.cancel()
		}
	}()

	// Wait for context cancel or server error
	select {
	case <-s.ctx.Done():
		s.l.Info("Context cancelled, shutting down server...")
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(ctxShutdown)
	case err := <-errCh:
		return err
	}
}

// Very large method...
// Divide into multiple functions
func (s *Server) handleStreamTransform(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := r.ParseMultipartForm(or.Or(s.maxMemoryUploadBytes, defaultMaxMemoryUploadBytes))
	if err != nil {
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}

	sqlFile, _, err := r.FormFile("sql_file")
	if err != nil {
		http.Error(w, "Missing sql_file", http.StatusBadRequest)
		return
	}
	defer sqlFile.Close()

	sqlBytes, err := io.ReadAll(sqlFile)
	if err != nil {
		http.Error(w, "Failed to read sql_file", http.StatusInternalServerError)
		return
	}

	var config map[string]any
	if cf, _, err := r.FormFile("config"); err == nil {
		defer cf.Close()
		// Try parsing as YAML or JSON
		data, _ := io.ReadAll(cf)
		if yaml.Unmarshal(data, &config) != nil {
			json.Unmarshal(data, &config)
		}
	}

	rend := renderer.NewGoTemplateRenderer()

	var query string
	var dataRows []map[string]string
	if vf, header, err := r.FormFile("values_file"); err == nil {
		defer vf.Close()

		loaderOpts := loader.LoadOpts{
			Sheet:      getString(config, "sheet-name-in", s.l),
			SheetIndex: getT[int](config, "sheet-index-in", s.l),
		}
		ld, err := loader.GetLoaderIO(header.Filename, vf, loaderOpts)
		if err != nil {
			http.Error(w, "Failed to get loader for values_file: "+err.Error(), http.StatusBadRequest)
			return
		}
		dataRows, err = ld.LoadIO(vf)
		if err != nil {
			http.Error(w, "Failed to load values_file: "+err.Error(), http.StatusBadRequest)
			return
		}
		q, err := rend.Render(string(sqlBytes), *renderer.FromMapArr(dataRows))
		if err != nil {
			http.Error(w, "Failed to render SQL query with the provided input data, error: "+err.Error(), http.StatusBadRequest)
			return
		}
		query = q
	} else {
		q, err := rend.Render(string(sqlBytes), renderer.QueryData{})
		if err != nil {
			http.Error(w, "Failed to render SQL query with the provided input data, error: "+err.Error(), http.StatusBadRequest)
			return
		}
		query = q
	}

	db, err := db.Get()
	if err != nil {
		http.Error(w, "Failed to connect to the database, error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := db.GetConnection().Prepare(query)
	if err != nil {
		http.Error(w, "Failed to prepare query, error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		http.Error(w, "Query execution failed, error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		http.Error(w, "Failed to get columns, error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var results []map[string]string
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			http.Error(w, "Failed to scan row, error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		rowMap := make(map[string]string)
		for i, col := range columns {
			var val string
			if b, ok := values[i].([]byte); ok {
				val = string(b)
			} else if values[i] != nil {
				val = fmt.Sprintf("%v", values[i])
			} else {
				val = ""
			}
			rowMap[col] = val
		}
		results = append(results, rowMap)
	}

	g, err := generator.GetGenerator(getString(config, "output", s.l), generator.GenerationOptions{
		SheetName: getString(config, "sheet", s.l),
	})
	if err != nil {
		http.Error(w, "Failed to get generator, error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Pipe for streaming
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()

		if err := g.GenerateIO(pw, results); err != nil {
			pw.CloseWithError(fmt.Errorf("file generation error: %w", err))
		}
	}()

	switch g.(type) {
	case *generator.XLSXGenerator:
		// For .xlsx
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	case *generator.CSVGenerator:
		// For .csv
		w.Header().Set("Content-Type", "text/csv")
	default:
		w.Header().Set("Content-Type", "text/plain")
	}

	w.WriteHeader(http.StatusOK)

	// Stream result back to client
	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := pr.Read(buf)
			if n > 0 {
				w.Write(buf[:n])
				w.(http.Flusher).Flush()
			}
			if err == io.EOF {
				return
			}
			if err != nil {
				s.l.Error("stream error", zap.Error(err))
				return
			}
		}
	}
}
