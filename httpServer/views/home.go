package views

import (
	"fmt"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request, url_matches []string) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Welcome!")
}
