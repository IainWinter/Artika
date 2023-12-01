// https://github.com/MicahParks/keyfunc/blob/master/examples/json/main.go

// I replaced the keys with the current ones from "https://www.googleapis.com/oauth2/v3/certs"
// and the jwt from my account

package main

import (
	"encoding/json"
	"log"

	"github.com/golang-jwt/jwt/v5"

	"github.com/MicahParks/keyfunc/v2"
)

func main() {
	// Get the JWKS as JSON.
	jwksJSON := json.RawMessage(`{
		"keys": [
		  {
			"e": "AQAB",
			"n": "vZgPf9nruMYY71q5pgThDwmk6Z3DD7cwN-Z52__b4xHeY95wOeKpjSliaI8K1PpeBbm4NykHm6UmfB_pCw5P2owpHZ8JEF2FCeDFKcOtZOzolYVgKZY8Sunqxcr3Sm0n73jbGcPgqu5PpjnOR4WkZCnpmDEZ34KNQat_MYYNUZZE2RlbpppNHiLatdiLW-rWi9YCmpsE4EIdd-XKIyZpQZRKaAl-w72BboTD_Koq2CkAOZOab73Q_G5FVT0NrxEWqP6artVfg5Dc_VVPnvtsC9yMe8lNgU3c3a-mE-vzE9oxAjr0s8Ek0Ih_sv-CbWL8xHiI7MOygIPG_aQqvMhPaQ",
			"kid": "0e72da1df501ca6f756bf103fd7c720294772506",
			"alg": "RS256",
			"use": "sig",
			"kty": "RSA"
		  },
		  {
			"kty": "RSA",
			"alg": "RS256",
			"e": "AQAB",
			"kid": "e4adfb436b9e197e2e1106af2c842284e4986aff",
			"n": "psply8S991RswM0JQJwv51fooFFvZUtYdL8avyKObshyzj7oJuJD8vkf5DKJJF1XOGi6Wv2D-U4b3htgrVXeOjAvaKTYtrQVUG_Txwjebdm2EvBJ4R6UaOULjavcSkb8VzW4l4AmP_yWoidkHq8n6vfHt9alDAONILi7jPDzRC7NvnHQ_x0hkRVh_OAmOJCpkgb0gx9-U8zSBSmowQmvw15AZ1I0buYZSSugY7jwNS2U716oujAiqtRkC7kg4gPouW_SxMleeo8PyRsHpYCfBME66m-P8Zr9Fh1Qgmqg4cWdy_6wUuNc1cbVY_7w1BpHZtZCNeQ56AHUgUFmo2LAQQ",
			"use": "sig"
		  }
		]
	  }`)

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.NewJSON(jwksJSON)
	if err != nil {
		log.Fatalf("Failed to create JWKS from JSON.\nError: %s", err.Error())
	}

	// Get a JWT to parse.
	jwtB64 := "out your JWT here"

	// Parse the JWT.
	token, err := jwt.Parse(jwtB64, jwks.Keyfunc)
	if err != nil {
		log.Fatalf("Failed to parse the JWT.\nError: %s", err.Error())
	}

	// Check if the token is valid.
	if !token.Valid {
		log.Fatalf("The token is not valid.")
	}
	log.Println("The token is valid.")
}
