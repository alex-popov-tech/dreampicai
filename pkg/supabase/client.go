package supabase

import (
	s "github.com/nedpals/supabase-go"
)

type UserCredentials = s.UserCredentials
type ProviderSignInOptions = s.ProviderSignInOptions

var Client *s.Client

func InitClient(url, key string) *s.Client {
	c := s.CreateClient(url, key)
	// init global supabase client for later usage across the app
	// if err != nil app won't start, so its safe to use global here
	Client = c
	return c
}
