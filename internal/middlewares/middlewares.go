package middlewares

import (
	"context"
	"net/http"
	"strings"
)

type ctxKeyLanguage string

const LangKey ctxKeyLanguage = "language"

var languageSuports = []string{"en", "pt-BR"}

func LanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var lang string

		if al := r.Header.Get("Accept-Language"); al != "" {
			lang = parseAcceptLanguage(al)
		}

		ctx := context.WithValue(r.Context(), LangKey, lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseAcceptLanguage(al string) string {
	if strings.Contains(al, ",") {
		parts := strings.Split(al, ",")

		for _, l := range parts {
			if contains(languageSuports, l) {
				return l
			}
		}
	}

	if contains(languageSuports, al) {
		return al
	}

	return "en"
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
