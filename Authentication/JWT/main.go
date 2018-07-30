package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// using asymmetric crypto/RSA keys
// location of the files used for signing and verification
const (
	privKeyPath = "app.rsa"     // openssl genrsa -out app.rsa 1024
	pubKeyPath  = "app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

// verify key and sign key
var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

//struct User for parsing login credentials
type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// read the key files before starting http handlers
func init() {
	var err error
	reader := rand.Reader
	bitsizze := 2048
	signKey, err = rsa.GenerateKey(reader, bitsizze)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}
	verifyKey = &signKey.PublicKey
}

// reads the login credentials, checks them and creates JWT the token
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	//decode into User struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error in request body")
		return
	}
	// validate user credentials
	if user.UserName != "shijuvar" && user.Password != "pass" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "Wrong info")
		return
	}
	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	// set our claims
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["userInfo"] = "userInfo"
	token.Claims = claims
	tokenString, err := t.SignedString(signKey)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Sorry, error while Signing Token!")
		log.Printf("Token Signing error: %v\n", err)
		return
	}
	response := Token{tokenString}
	jsonResponse(response, w)
}

// only accessible with a valid token
func authHandler(w http.ResponseWriter, r *http.Request) {
	// validate the tok
	jwtToken := r.Header.Get("Authorization")
	jwtToken = strings.Replace(jwtToken, "Bearer", "", -1)
	jwtToken = strings.TrimSpace(jwtToken)
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			fmt.Println("===============Authentication Error===============")
			fmt.Println(err.Error())
			fmt.Println("===============Authentication Error===============") // something was wrong during the validation
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "Token Expired, get a new one.")
				return
			default:
				fmt.Println("===============Auth Validation Error===============")
				fmt.Println(err.Error())
				fmt.Println("===============Auth Validation  Error===============")
				w.WriteHeader(http.StatusUnauthorized)
				response := Response{"Token Validation Error"}
				jsonResponse(response, w)
				log.Printf("ValidationError error: %+v\n", vErr.Errors)
				return
			}
		default: // something else went wrong
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error while Parsing Token!")
			log.Printf("Token parse error: %v\n", err.Error())
			return
		}
	}
	if token.Valid {
		response := Response{"Authorized to the system"}
		jsonResponse(response, w)
	} else {
		response := Response{"Invalid token"}
		jsonResponse(response, w)
	}
}

type Response struct {
	Text string `json:"text"`
}
type Token struct {
	Token string `json:"token"`
}

func jsonResponse(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

//Entry point of the program
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/auth", authHandler).Methods("POST")
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
