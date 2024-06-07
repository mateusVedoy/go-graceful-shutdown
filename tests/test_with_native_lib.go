package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func With_native_libs() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 3)
		w.Write([]byte("Hello"))
	})

	server := &http.Server{Addr: ":8080"}

	go func() {
		if er := server.ListenAndServe(); er != nil && http.ErrServerClosed != er {
			panic(er)
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	fmt.Printf("\nstopped\n")

}
