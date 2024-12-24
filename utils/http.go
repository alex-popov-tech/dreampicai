package utils

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func AddAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Value:    accessToken,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "rt",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Value:    refreshToken,
	})
}

func CleanAllCookies(w http.ResponseWriter, r *http.Request) {
	// Get all cookies from the request
	cookies := r.Cookies()

	// For each cookie, create a new cookie with the same name but expired
	for _, cookie := range cookies {
		// Create a new cookie with the same name
		expiredCookie := &http.Cookie{
			Name:     cookie.Name,
			Value:    "",              // Empty the value
			Path:     "/",             // Cover all paths
			Domain:   cookie.Domain,   // Use same domain
			Expires:  time.Unix(0, 0), // Set to epoch time (expired)
			MaxAge:   -1,              // Delete cookie immediately
			Secure:   cookie.Secure,   // Maintain secure flag
			HttpOnly: cookie.HttpOnly, // Maintain HttpOnly flag
			SameSite: cookie.SameSite, // Maintain SameSite policy
		}

		// Set the expired cookie in response
		http.SetCookie(w, expiredCookie)
	}
}

func GetTokensFromQuery(values url.Values) (accessToken string, refreshToken string, err error) {
	accessToken = values.Get("access_token")
	refreshToken = values.Get("refresh_token")
	if accessToken == "" || refreshToken == "" {
		return "", "", fmt.Errorf(
			"Cannot get 'access_token' or 'refresh_token' from query: %v",
			values,
		)
	}
	return accessToken, refreshToken, nil
}

func MakeRoute(handler func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			log.Printf("Handler unhanded error %s %s\nError: %v", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Sorry, something went wrong\n%v", err.Error())))
		}
	})
}

// func UserAuthMiddleware(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		atCookie, err := r.Cookie("at")
// 		at := atCookie.Value
// 		rtCookie, err := r.Cookie("rt")
// 		rt := rtCookie.Value
// 		requestWithEmptyUserContext := r.WithContext(
// 			context.WithValue(r.Context(), model.UserContextKey, model.User{}),
// 		)
// 		if err != nil {
// 			log.Printf("UserAuthMiddleware, cookie not found: %v\n", err)
// 			handler.ServeHTTP(w, requestWithEmptyUserContext)
// 			return
// 		}
// 		if at == "" {
// 			log.Printf("UserAuthMiddleware, access_token cookie is empty: %v\n", err)
// 			handler.ServeHTTP(w, requestWithEmptyUserContext)
// 			return
// 		}
// 		if rt == "" {
// 			log.Printf("UserAuthMiddleware, refresh_token cookie is empty: %v\n", err)
// 			handler.ServeHTTP(w, requestWithEmptyUserContext)
// 			return
// 		}
// 		user, err := ParseSupabaseToken(at)
// 		if errors.Is(err, jwt.ErrTokenExpired) {
// 			authDetails, err := supabase.Client.Auth.RefreshUser(context.Background(), at, rt)
// 			if err != nil {
// 				log.Printf(
// 					"UserAuthMiddleware, access_token is expired and cannot be refreshed: %v\n",
// 					err,
// 				)
// 				handler.ServeHTTP(
// 					w,
// 					r.WithContext(context.WithValue(r.Context(), model.UserContextKey, user)),
// 				)
// 				return
// 			}
// 			atCookie := http.Cookie{
// 				Name:     "at",
// 				Path:     "/",
// 				HttpOnly: true,
// 				Secure:   true,
// 				Value:    authDetails.AccessToken,
// 			}
// 			rtCookie := http.Cookie{
// 				Name:     "rt",
// 				Path:     "/",
// 				HttpOnly: true,
// 				Secure:   true,
// 				Value:    authDetails.RefreshToken,
// 			}
// 			http.SetCookie(w, &atCookie)
// 			http.SetCookie(w, &rtCookie)
// 			handler.ServeHTTP(
// 				w,
// 				r.WithContext(
// 					context.WithValue(
// 						r.Context(),
// 						model.UserContextKey,
// 						model.User{
// 							Email:      authDetails.User.Email,
// 							ID:         authDetails.User.ID,
// 							IsLoggedIn: true,
// 						},
// 					),
// 				),
// 			)
// 			return
// 		}
// 		if err != nil {
// 			log.Printf("UserAuthMiddleware, jwt parsing error: %v\n", err)
// 			handler.ServeHTTP(
// 				w,
// 				r.WithContext(context.WithValue(r.Context(), model.UserContextKey, user)),
// 			)
// 			return
// 		}
// 		handler.ServeHTTP(
// 			w,
// 			r.WithContext(context.WithValue(r.Context(), model.UserContextKey, user)),
// 		)
// 	})
// }
//
// func AuthProtectedMiddleware(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if user, ok := r.Context().Value(model.UserContextKey).(model.User); ok && user.IsLoggedIn {
// 			handler.ServeHTTP(w, r)
// 			return
// 		}
// 		http.Redirect(w, r, fmt.Sprintf("/signin?redirect=%s", r.URL.Path), http.StatusFound)
// 	})
// }
