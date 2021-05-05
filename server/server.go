package server

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/yura-under-review/port-api/models"
)

const (
	// MemoryLimit = 200 << 10 // 200 KB
	MemoryLimit = 5 << 10
)

type Server struct {
	addr             string
	rootPageTemplate string
	rootPageRendered []byte
	s                *http.Server
	repo             PortsRepository
	sinkBatchSize    int
}

type PortsRepository interface {
	UpsertPorts(context.Context, []*models.PortInfo) error
}

func New(addr, rootPageTemplate string, repo PortsRepository, sinkBatchSize int) *Server {
	return &Server{
		addr:             addr,
		rootPageTemplate: rootPageTemplate,
		repo:             repo,
		sinkBatchSize:    sinkBatchSize,
	}
}

func (srv *Server) Run(ctx context.Context, wg *sync.WaitGroup) error {

	// TODO: implement root page templating to setup host:port

	var err error
	srv.rootPageRendered, err = ioutil.ReadFile(srv.rootPageTemplate)
	if err != nil {
		log.Errorf("failed to read root page file: %v", err)
		return err
	}

	r := mux.NewRouter()

	r.HandleFunc("/", srv.rootHandler)
	r.HandleFunc("/upload", srv.uploadHandler)

	srv.s = &http.Server{
		Addr:    srv.addr,
		Handler: r,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := srv.s.ListenAndServe(); err != http.ErrServerClosed {
			log.Errorf("server failed: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()

		if err := srv.Close(); err != nil {
			log.Errorf("failed to close http server: %v", err)
		}
	}()

	log.Info("http server runs")

	return nil
}

func (srv *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	if srv.s != nil {
		err = srv.s.Shutdown(ctx)
	}

	return err
}

func (srv *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("new ROOT request [host: %s, url: %s]", r.Host, r.URL)

	_, err := w.Write(srv.rootPageRendered)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (srv *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {

	log.Infof("new UPLOAD request [host: %s, url: %s]", r.Host, r.URL)

	if err := r.ParseMultipartForm(MemoryLimit); err != nil {
		log.Errorf("failed to parse form: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	f, h, err := r.FormFile("file")
	if err != nil {
		log.Errorf("failed to read form file: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debugf("request file: [name: %s, size: %d]", h.Filename, h.Size)

	if err := srv.sinkData(r.Context(), f); err != nil {
		log.Errorf("failed to sink ports: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (srv *Server) sinkData(ctx context.Context, r io.Reader) error {

	parser := NewFileParser(r)

	idx := 1
	var ports []*models.PortInfo
	lastSink := false

	for {
		p, err := parser.Read()
		if err != nil {
			lastSink = true
		}

		if p != nil {
			ports = append(ports, p)
		}

		if (idx%srv.sinkBatchSize == 0) || lastSink {
			if err := srv.repo.UpsertPorts(ctx, ports); err != nil {
				return fmt.Errorf("failed to sink data: %w", err)
			}

			ports = ports[:0]
		}

		if lastSink {
			break
		}
		idx++
	}

	return nil
}
