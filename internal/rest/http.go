package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/saikrir/auth-svc/internal/models"
)

type LDAPAuthenticator interface {
	Authenticate(string, string) error
}

type TokenSvc interface {
	NewToken(string) (string, int64, error)
	VerifyToken(string) error
}

type Server struct {
	router           *http.ServeMux
	server           *http.Server
	ldapAutheticator LDAPAuthenticator
	tokenSvc         TokenSvc
	rootContext      string
}

func NewServer(rootContext string, port int, tokenSvc TokenSvc, ldapAuthenticator LDAPAuthenticator) *Server {
	svr := &Server{rootContext: rootContext, tokenSvc: tokenSvc, ldapAutheticator: ldapAuthenticator}

	svr.router = http.NewServeMux()
	rootCtxMux := http.NewServeMux()
	rootCtxMux.Handle(fmt.Sprintf("%s/", rootContext), http.StripPrefix(rootContext, svr.router))
	svr.mapRoutes()

	svr.server = &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%d", port),
		Handler:           rootCtxMux,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	return svr
}

func (s *Server) mapRoutes() {
	s.router.HandleFunc("POST /authentications", s.authenticate)
	s.router.HandleFunc("POST /tokens", s.validate)
}

func (s *Server) authenticate(resp http.ResponseWriter, r *http.Request) {

	log.Println("will authenticate ")
	var (
		credential models.AuthenticationRequest
		err        error
	)

	if err = json.NewDecoder(r.Body).Decode(&credential); err != nil {
		log.Println("failed to decode incoming payload ", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	if err = validator.New().Struct(credential); err != nil {
		log.Println("request did no pass validation", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.ldapAutheticator.Authenticate(credential.AccountName, credential.AccountPassword); err != nil {
		log.Println("authentication failed ", err)
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Println("Authentication succeeded for ", credential.AccountName)

	var (
		token string
		ttl   int64
	)
	if token, ttl, err = s.tokenSvc.NewToken(credential.AccountName); err != nil {
		log.Println("failed to generate token ", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := models.AuthenticationResponse{
		Token:     token,
		ExpriesAt: ttl,
	}

	if err = json.NewEncoder(resp).Encode(response); err != nil {
		log.Println("failed to encode response ", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) validate(resp http.ResponseWriter, r *http.Request) {
	var (
		authzRequest models.AuthorizationRequest
		err          error
	)

	if err = json.NewDecoder(r.Body).Decode(&authzRequest); err != nil {
		log.Println("failed to decode request ", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = validator.New().Struct(authzRequest); err != nil {
		log.Println("request validation failed", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.tokenSvc.VerifyToken(authzRequest.Token); err != nil {
		log.Println("Token validation failed", err)
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	return
}
func (s *Server) Serve() error {

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Println("Server will stop ", err)
		}
	}()

	// These 3 lines, capture CTRL+C and pass control to the app
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	s.server.Shutdown(ctx)
	log.Println("Server has shutdown")
	return nil
}
