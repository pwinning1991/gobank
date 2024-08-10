package main

import (
	"encoding/json"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewApiServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{listenAddr: listenAddr,
		store: store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAith(makeHTTPHandleFunc(s.handleGetAccountByID)))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))
	log.Println("JSON Api server Listening on " + s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, router))

}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetAccount(w, r)
	}
	if r.Method == http.MethodPost {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("Unsupported method: %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		id, err := getID(r)
		if err != nil {
			return err
		}
		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}
		return WriteJson(w, http.StatusOK, account)
	}

	if r.Method == http.MethodDelete {
		return s.handleDeleteAccount(w, r)

	}
	return fmt.Errorf("Unsupported method: %s", r.Method)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}
	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, tokenString)
	return WriteJson(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, map[string]int{"delete": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJson(w, http.StatusOK, transferReq)
}

func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"ExpiresAt":     15000,
		"accountNumber": account.Number,
	}

	secrect := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secrect))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func withJWTAith(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")
		tokenString := r.Header.Get("Authorization")
		_, err := validateJWT(tokenString)
		if err != nil {
			WriteJson(w, http.StatusUnauthorized, ApiError{Error: "invalid token"})
			return
		}
		handlerFunc(w, r)
	}
}

func validateJWT(token string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return id, fmt.Errorf("invalid id: %s", idstr)
	}
	return id, err

}
