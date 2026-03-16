package collector

import (
	"fmt"
	"net/http"
)

func StartHTTPServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Login page")
	})

	handler := HTTPCollector(mux)

	fmt.Println("HTTP Server running on :8080")

	http.ListenAndServe(":8080", handler)
}