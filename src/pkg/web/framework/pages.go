package framework

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func Run(handler func(w http.ResponseWriter, r *http.Request)) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	fmt.Print("running server...")

	err := http.ListenAndServe("0.0.0.0:80", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
