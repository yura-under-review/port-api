package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/yura-under-review/port-api/models"
)

type Server struct {
	addr             string
	rootPageTemplate string
	rootPageRendered []byte
	s                *http.Server
	repo             PortsRepository
}

type PortsRepository interface {
	UpsertPorts(context.Context, []models.Port) error
}

func New(addr, rootPageTemplate string, repo PortsRepository) *Server {
	return &Server{
		addr:             addr,
		rootPageTemplate: rootPageTemplate,
		repo:             repo,
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
	log.Infof("new ROOT request [host: %s, url: %s]",
		r.Host,
		r.URL,
	)

	_, err := w.Write(srv.rootPageRendered)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (srv *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {

	log.Infof("new UPLOAD request [host: %s, url: %s]",
		r.Host,
		r.URL,
	)

	if err := r.ParseMultipartForm(20 << 10); err != nil {
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

	fileBuffer := make([]byte, 20<<10)
	fileLen, err := f.Read(fileBuffer)
	if err != nil {
		log.Errorf("failed to read file to buffer: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// temporary printing file
	fmt.Printf("----- [%s] -----\n", h.Filename)
	fmt.Println(string(fileBuffer[:fileLen]))
	fmt.Println("----------")

	// TODO: deal with memory limitation
	// TODO: forward received file

	// ports, err := FileToPorts(fileBuffer[:fileLen])
	// if err != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// }
	// srv.repo.UpsertPorts()

	w.WriteHeader(http.StatusOK)
}
