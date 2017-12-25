package main

import (
	"fmt"
	"net/http"
	"regexp"
	"./views"
	"strings"
)

var GET uint8 = 0x01
var POST uint8 = 0x02
var PUT uint8 = 0x04
var DELETE uint8 = 0x08

func convertMethod(check string) uint8 {
	check = strings.ToUpper(check)
	switch check {
		case "GET":
			return GET
		case "POST":
			return POST
		case "PUT":
			return PUT
		case "DELETE":
			return DELETE
		default:
			return 0x00
	}
}

type Route struct {
	Methods uint8
	Path *regexp.Regexp
	Func func(http.ResponseWriter, *http.Request, []string)
}

func (r *Route) pathMatches(path string) []string {
	return r.Path.FindStringSubmatch(path)
}

func (r *Route) process(w http.ResponseWriter, req *http.Request) bool {
	method := convertMethod(req.Method)
	if r.Methods & method == 0 {
		return false
	}
	if matches := r.pathMatches(req.URL.Path); matches != nil {
		r.Func(w, req, matches)
		return true
	}
	return false
}

type Routes struct {
	routes []Route
}

func (R *Routes) addRoute(m_flags uint8, path string, fn func(http.ResponseWriter, *http.Request, []string)) Route {
	if (GET + POST + PUT + DELETE) & m_flags == 0 {
		panic("Invalid route method.")
	}
	var r = regexp.MustCompile("^/" + path + "/*$")
	var new_route Route = Route{m_flags, r, fn}
	R.routes = append(R.routes, new_route)
	return new_route
}

func (R *Routes) handleRoute(w http.ResponseWriter, r *http.Request) {
	var found bool = false
	for _, route := range R.routes {
		if route.process(w, r) == true {
			found = true
			break
		}
	}
	if !found {
		fmt.Println("Request: " + r.Method + " " + r.URL.Path + " 404")
		http.NotFound(w, r)
		return
	}
	fmt.Println("Request: " + r.Method + " " + r.URL.Path + " 200")
}

func main() {
	var router Routes = Routes{}
	
	router.addRoute(GET, "(home|)", views.Home)
	router.addRoute(GET+PUT, "user/([^/]+)", views.User)

	http.HandleFunc("/", router.handleRoute)
	fmt.Println(http.ListenAndServe(":8080", nil))
}