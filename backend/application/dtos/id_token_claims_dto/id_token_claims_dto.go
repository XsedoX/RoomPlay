package id_token_claims_dto

type IdTokenClaimsDto struct {
	Subject    string `json:"sub" validate:"required"`
	GivenName  string `json:"given_name" validate:"required"`
	FamilyName string `json:"family_name" validate:"required"`
}
