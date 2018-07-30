package main

import (
	"net/http"
)

// func main() {
// 	mux := http.NewServeMux()
// 	fs := http.FileServer(http.Dir("public"))
// 	mux.Handle("/", fs)
// 	log.Println("==========Server started at port:3000============")
// 	http.ListenAndServe(":3000", fs)
// }

//*************For Http middleware
func main() {
	// To serve a directory on disk (/public) under an alternate URL
	// path (/public/), use StripPrefix to modify the request
	// URL's path before the FileServer sees it:
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
}
