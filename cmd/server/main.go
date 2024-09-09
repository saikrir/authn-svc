package main

import (
	"fmt"
	"log"
	"os"

	"github.com/saikrir/auth-svc/internal/ldap"
	"github.com/saikrir/auth-svc/internal/rest"
	"github.com/saikrir/auth-svc/internal/token"
)

const SecretEnvVar = "secret"
const ValidityInHours = 6
const RootContext = "/v1/auth"
const AuthPort = 9999

func Run() error {
	ldapAuth := ldap.NewLdapAuth("ldap.skrao.net", "ou=ServiceAccounts,dc=skrao,dc=net")
	tokenSecret := os.Getenv(SecretEnvVar)
	if len(tokenSecret) == 0 {
		log.Printf("failed to locate env variable [%s] \n", SecretEnvVar)
		return fmt.Errorf("failed to locate env variable [%s]", SecretEnvVar)
	}

	tokenSvc := token.NewJWTManager(tokenSecret, ValidityInHours)

	server := rest.NewServer(RootContext, AuthPort, tokenSvc, ldapAuth)

	log.Println("start server on ", AuthPort)

	return server.Serve()
}

func main() {
	if err := Run(); err != nil {
		log.Fatalln("failed to start server ", err)
	}
}
