package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	needforheatmanualserver "github.com/energietransitie/needforheat-manual-server"
	custommiddleware "github.com/energietransitie/needforheat-manual-server/middleware"
	"github.com/energietransitie/needforheat-manual-server/parser"
	"github.com/energietransitie/needforheat-manual-server/wfs/dirfs"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conf, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	//Parser parses every manual so it can be served
	parsedFS := dirfs.New("./parsed")

	localDirParser := parser.New(parsedFS)

	err = localDirParser.Parse(conf.Source)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("generated folder structure to be served")

	server := needforheatmanualserver.NewServer(parsedFS, needforheatmanualserver.ServerOptions{
		FallbackLanguage: conf.FallbackLanguage,
	})

	r := chi.NewRouter()

	//CleanPathRedirect is the router, it will make sure everything redirects to the right page
	r.Use(middleware.Timeout(time.Second * 30))
	r.Use(middleware.Heartbeat("/healthcheck"))
	r.Use(custommiddleware.CleanPathRedirect)
	r.Use(middleware.Logger)

	r.Mount("/", server)

	httpServer := &http.Server{
		Addr:        ":8080",
		Handler:     r,
		BaseContext: returnContextFn(ctx),
	}

	err = listenAndServe(ctx, httpServer)
	if err != nil {
		log.Println(err)
	}
}

// Return a function for BaseContext that always returns context ctx.
func returnContextFn(ctx context.Context) func(net.Listener) context.Context {
	return func(_ net.Listener) context.Context {
		return ctx
	}
}

// Start HTTP server and gracefully shutdown if the context is cancelled.
func listenAndServe(ctx context.Context, server *http.Server) error {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return server.ListenAndServe()
	})

	log.Println("serving manuals on", server.Addr)

	g.Go(func() error {
		<-gCtx.Done()
		return server.Shutdown(context.Background())
	})

	return g.Wait()
}
