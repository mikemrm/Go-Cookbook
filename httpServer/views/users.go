package views

import (
	"fmt"
	"net/http"
)

func User(w http.ResponseWriter, r *http.Request, url_matches []string){
	fmt.Fprintf(w, "Hello " + url_matches[1])
}
