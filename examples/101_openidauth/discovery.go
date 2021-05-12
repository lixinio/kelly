package main

import (
	"net/http"

	"github.com/lixinio/kelly"
)

type WellKnown struct {
	// URL using the https scheme with no query or fragment component that the OP asserts as its IssuerURL Identifier.
	// If IssuerURL discovery is supported , this value MUST be identical to the issuer value returned
	// by WebFinger. This also MUST be identical to the iss Claim value in ID Tokens issued from this IssuerURL.
	//
	// required: true
	// example: https://playground.ory.sh/ory-hydra/public/
	Issuer string `json:"issuer"`

	// URL of the OP's OAuth 2.0 Authorization Endpoint.
	//
	// required: true
	// example: https://playground.ory.sh/ory-hydra/public/oauth2/auth
	AuthEndpoint string `json:"authorization_endpoint"`

	// URL of the OP's OAuth 2.0 Token Endpoint
	//
	// required: true
	// example: https://playground.ory.sh/ory-hydra/public/oauth2/token
	TokenEndpoint string `json:"token_endpoint"`

	// URL of the OP's JSON Web Key Set [JWK] document. This contains the signing key(s) the RP uses to validate
	// signatures from the OP. The JWK Set MAY also contain the Server's encryption key(s), which are used by RPs
	// to encrypt requests to the Server. When both signing and encryption keys are made available, a use (Key Use)
	// parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key's intended usage.
	// Although some algorithms allow the same key to be used for both signatures and encryption, doing so is
	// NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of
	// keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate.
	//
	// required: true
	JWKsURI string `json:"jwks_uri"`

	// JSON array containing a list of the Subject Identifier types that this OP supports. Valid types include
	// pairwise and public.
	//
	// required: true
	SubjectTypes []string `json:"subject_types_supported"`

	// JSON array containing a list of the OAuth 2.0 response_type values that this OP supports. Dynamic OpenID
	// Providers MUST support the code, id_token, and the token id_token Response Type values.
	//
	// required: true
	ResponseTypes []string `json:"response_types_supported"`

	// JSON array containing a list of the Claim Names of the Claims that the OpenID Provider MAY be able to supply
	// values for. Note that for privacy or other reasons, this might not be an exhaustive list.
	ClaimsSupported []string `json:"claims_supported"`

	// JSON array containing a list of the OAuth 2.0 response_mode values that this OP supports.
	ResponseModes []string `json:"response_modes_supported"`

	// URL of the OP's UserInfo Endpoint.
	UserInfoEndpoint string `json:"userinfo_endpoint"`

	// SON array containing a list of the OAuth 2.0 [RFC6749] scope values that this server supports. The server MUST
	// support the openid scope value. Servers MAY choose not to advertise some supported scope values even when this parameter is used
	ScopesSupported []string `json:"scopes_supported"`

	// JSON array containing a list of Client Authentication methods supported by this Token Endpoint. The options are
	// client_secret_post, client_secret_basic, client_secret_jwt, and private_key_jwt, as described in Section 9 of OpenID Connect Core 1.0
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`

	// JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for the ID Token
	// to encode the Claims in a JWT.
	//
	// required: true
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`

	// URL of the authorization server's OAuth 2.0 revocation endpoint.
	RevocationEndpoint string `json:"revocation_endpoint"`
}

var responseTypes = []string{
	"code",
	"token",
	"id_token",
	"code token",
	"code id_token",
	"token id_token",
	"code token id_token",
}

var claims = []string{
	"aud",
	"auth_time",
	"email",
	"exp",
	"iat",
	"iss",
	"name",
	"picture",
	"sub",
}

func (h *Handler) handleDiscovery(ac *kelly.AnnotationContext) kelly.HandlerFunc {
	return func(c *kelly.Context) {
		wk := &WellKnown{
			Issuer:                            h.issuer,
			AuthEndpoint:                      h.issuer + AuthPath,
			TokenEndpoint:                     h.issuer + TokenPath,
			JWKsURI:                           h.issuer + JWKPath,
			RevocationEndpoint:                h.issuer + RevocationPath,
			SubjectTypes:                      []string{"public"},
			ResponseModes:                     []string{"query", "fragment"},
			ResponseTypes:                     responseTypes,
			ClaimsSupported:                   claims,
			ScopesSupported:                   []string{"openid", "email", "profile"},
			UserInfoEndpoint:                  h.issuer + UserInfoPath,
			TokenEndpointAuthMethodsSupported: []string{"client_secret_post", "client_secret_basic"},
			IDTokenSigningAlgValuesSupported:  []string{"RS256"},
		}

		c.WriteJSON(http.StatusOK, wk)
	}
}
