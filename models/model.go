package models

type Credential struct {
	AccountName     string `json:"accountName"`
	AccountPassword string `json:"accountPassword"`
}

type AuthenticationRequest struct {
	Credential
}

type AuthenticationResponse struct {
	Token     string `json:"token"`
	ExpriesAt int64  `json:"expires_at"`
}

type AuthorizationRequest struct {
	Token string `json:"token"`
}
