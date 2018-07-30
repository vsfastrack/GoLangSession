package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

//*********************For the handler interface type************************
// type messageHandler struct {
// 	message string
// }

// func (m *messageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, m.message)
// }

// func main() {
// 	mux := http.NewServeMux()
// 	mh1 := &messageHandler{"Welcome to Go Web Development"}
// 	mux.Handle("/welcome", mh1)
// 	mh2 := &messageHandler{"net/http is awesome"}
// 	mux.Handle("/message", mh2)
// 	log.Println("Listening...")
// 	http.ListenAndServe(":8080", mux)
// }

//*********************For the handler func type************************
// func messageHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to Go Web Development")
// }
// func main() {
// 	mux := http.NewServeMux()
// 	// Convert the messageHandler function to a HandlerFunc type
// 	mh := http.HandlerFunc(messageHandler)
// 	mux.Handle("/welcome", mh)
// 	log.Println("Listening...")
// 	http.ListenAndServe(":8080", mux)
// }

//****************Handler logic into closure********************

//Handler logic into a Closure
// func messageHandler(message string) http.Handler {
// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// fmt.Fprintf(w, message)
// })
// }
// func main() {
// mux := http.NewServeMux()
// mux.Handle("/welcome", messageHandler("Welcome to Go Web Development"))
// mux.Handle("/message", messageHandler("net/http is awesome"))
// log.Println("Listening...")
// http.ListenAndServe(":8080", mux)
// }

//*************Server using server struct*******************
func messageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Go Web Development")
}
func main() {
	http.HandleFunc("/welcome", messageHandler)
	server := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
