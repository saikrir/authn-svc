package main

import (
	"fmt"

	"github.com/saikrir/auth-svc/ldap"
	"github.com/saikrir/auth-svc/models"
)

func main() {

	ldapAuth := ldap.NewLdapAuth("ldap.skrao.net", "ou=ServiceAccounts,dc=skrao,dc=net")

	fmt.Println("Search ", ldapAuth.Authenticate(models.Credential{AccountName: "", AccountPassword: ""}))
}
