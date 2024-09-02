package models

type Credential struct {
	AccountName     string `json:"accountName" validate:"required"`
	AccountPassword string `json:"accountPassword" validate:"required"`
}

type AuthenticationRequest struct {
	Credential
}

type AuthenticationResponse struct {
	Token     string `json:"token"`
	ExpriesAt int64  `json:"expires_at"`
}

type AuthorizationRequest struct {
	Token string `json:"token" validate:"required"`
}
