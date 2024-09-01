package ldap

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

type Credential struct {
	AccountName, AccountPassword string
}

type LdapAuth struct {
	Host, SearchBaseDN string
}

func NewLdapAuth(host, searchBaseDN string) *LdapAuth {
	return &LdapAuth{
		Host:         host,
		SearchBaseDN: searchBaseDN,
	}
}

func (l *LdapAuth) openConnection() (*ldap.Conn, error) {
	ldapUrl := fmt.Sprintf("ldap://%s:389", l.Host)
	var (
		conn *ldap.Conn
		err  error
	)
	if conn, err = ldap.DialURL(ldapUrl); err != nil {
		log.Println("failed to dail ldap url ", ldapUrl)
		return nil, err
	}

	return conn, nil
}

func (l *LdapAuth) Authenticate(userCreds Credential) error {
	var (
		conn          *ldap.Conn
		err           error
		searchResults *ldap.SearchResult
	)
	if conn, err = l.openConnection(); err != nil {
		log.Println("failed to obtain connection ", err)
		return err
	}
	defer conn.Close()

	userBindDN := fmt.Sprintf("uid=%s,%s", ldap.EscapeFilter(userCreds.AccountName), l.SearchBaseDN)

	if err = conn.Bind(userBindDN, userCreds.AccountPassword); err != nil {
		log.Println("authentication failed ", err)
		return err
	}

	searchFilter := fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(userCreds.AccountName))
	ldapQuery := ldap.NewSearchRequest(l.SearchBaseDN, ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false, searchFilter, []string{}, []ldap.Control{})

	if searchResults, err = conn.Search(ldapQuery); err != nil {
		log.Println("failed to peform search ", err)
		return err
	}

	if len(searchResults.Entries) == 0 {
		return fmt.Errorf("failed to find any matches for %s", userCreds.AccountName)
	}

	log.Printf("%s was successfully authenticated \n", searchResults.Entries[0].GetAttributeValue("cn"))
	return nil
}
