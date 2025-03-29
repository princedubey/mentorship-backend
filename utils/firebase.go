package utils

import (
	"context"
	"encoding/json"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var firebaseAuth *auth.Client

// InitFirebase initializes Firebase Admin SDK
func InitFirebase() {
	// Build Firebase config from environment variables
	firebaseConfig := map[string]interface{}{
		"type":                        os.Getenv("FIREBASE_TYPE"),
		"project_id":                  os.Getenv("FIREBASE_PROJECT_ID"),
		"private_key_id":              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		"private_key":                 os.Getenv("FIREBASE_PRIVATE_KEY"),
		"client_email":                os.Getenv("FIREBASE_CLIENT_EMAIL"),
		"client_id":                   os.Getenv("FIREBASE_CLIENT_ID"),
		"auth_uri":                    os.Getenv("FIREBASE_AUTH_URI"),
		"token_uri":                   os.Getenv("FIREBASE_TOKEN_URI"),
		"auth_provider_x509_cert_url": os.Getenv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
		"client_x509_cert_url":        os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
		"universe_domain":            os.Getenv("FIREBASE_UNIVERSE_DOMAIN"),
	}

	configJSON, err := json.Marshal(firebaseConfig)
	if err != nil {
		log.Fatal("Error creating Firebase config JSON:", err)
	}

	opt := option.WithCredentialsJSON(configJSON)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	firebaseAuth = auth
}

// VerifyFirebaseToken verifies the Firebase ID token
func VerifyFirebaseToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := firebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// GetUserByEmail gets Firebase user by email
func GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error) {
	return firebaseAuth.GetUserByEmail(ctx, email)
}

// GetUserByUID gets Firebase user by UID
func GetUserByUID(ctx context.Context, uid string) (*auth.UserRecord, error) {
	return firebaseAuth.GetUser(ctx, uid)
}
