package main

import (
	"fmt"

	"github.com/saikrir/auth-svc/ldap"
)

func main() {

	ldapAuth := ldap.NewLdapAuth("ldap.skrao.net", "ou=ServiceAccounts,dc=skrao,dc=net")

	fmt.Println("Search ", ldapAuth.Authenticate(ldap.Credential{AccountName: "", AccountPassword: ""}))
}
