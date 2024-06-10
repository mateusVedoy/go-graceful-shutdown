package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	RenderChi "github.com/go-chi/render"
	RenderPkg "github.com/unrolled/render"
)

var render *RenderPkg.Render
var wg sync.WaitGroup

func With_go_chi() {

	route := &Route{}
	first_ctx := context.Background()

	wg.Add(1)
	go route.Start(first_ctx)
	finish(first_ctx, route)
	wg.Wait()

}

type Route struct {
	server *http.Server
}

func (R *Route) Start(ctx context.Context) {
	contentType := middleware.AllowContentType("application/json")
	render = RenderPkg.New()
	route := chi.NewRouter()
	route.Use(middleware.RequestID)
	route.Use(middleware.RealIP)
	route.Use(middleware.Recoverer)
	route.Use(contentType)
	route.Use(RenderChi.SetContentType(RenderChi.ContentTypeJSON))
	route.Use(middleware.Timeout(60 * time.Second))

	route.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 5)
		render.JSON(w, 200, "hello")
	})

	R.server = &http.Server{
		Addr:    ":8080",
		Handler: route,
	}

	if er := R.server.ListenAndServe(); er != nil && http.ErrServerClosed != er {
		panic(er)
	}
}

func (R *Route) Close(ctx context.Context) error {
	return R.server.Shutdown(ctx)
}

func finish(ctx context.Context, router *Route) {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-stop

	timeout, shutdown := context.WithTimeout(ctx, time.Second*5)

	defer shutdown()

	if err := router.Close(timeout); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 5)
	fmt.Println("stopped")
	os.Exit(1)
}
