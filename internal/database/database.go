package database

import (
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

func CreateClient(jwt string) *supabase.Client {
	// TODO: Create database client

	//API_URL and API_KEY variables in .env file

	API_URL := os.Getenv("API_URL")
	API_KEY := os.Getenv("API_KEY")
	client, err := supabase.NewClient(API_URL, API_KEY, &supabase.ClientOptions{})
	if err != nil {
		panic(err)
	}
	//client.Functions.TokenAuth(jwt)
	jwt = strings.Replace(jwt, "Bearer ", "", 1)

	session := types.Session{
		AccessToken: jwt,
	}
	client.UpdateAuthSession(session)
	// client.Auth.WithToken(jwt)

	return client
}

func FetchUserID(client *supabase.Client) (uuid.UUID, error) {
	clientdetails, err := client.Auth.GetUser()
	if err != nil {
		return uuid.Nil, err
	}
	return clientdetails.ID, nil
}
