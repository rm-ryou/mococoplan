package route

import (
	"io"
	"net/http"
)

func Setup() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello, World!\n")
	})

	return mux
}
