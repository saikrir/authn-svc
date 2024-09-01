package ldap

type Credential struct {
	AccountName, AccountPassword string
}

type LdapAuth struct {
	Host, BindDN    string
	BindCredentials Credential
}

func NewLdapAuth(host, bindDN string, bindCredentials Credential) *LdapAuth {
	return &LdapAuth{
		Host:            host,
		BindDN:          bindDN,
		BindCredentials: bindCredentials,
	}
}
