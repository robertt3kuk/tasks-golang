package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Use(AllowOnly(NewIPWhitelist("127.0.0.1", "::1", "172.16.238.11")))
	r.Get("/generate-salt", func(w http.ResponseWriter, r *http.Request) {
		type Salt struct {
			S string
		}
		salt, _ := json.Marshal(Salt{S: generate()})
		w.Write(salt)
	})
	fmt.Println("started serving on port :3000")
	http.ListenAndServe(":3000", r)
}

func AllowOnly(allowedHandler func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return allowedHandler(next)
	}
}
func NewIPWhitelist(allowedIPs ...string) func(next http.Handler) http.Handler {
	allowedMap := make(map[string]bool)
	for _, ip := range allowedIPs {
		allowedMap[ip] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if !allowedMap[remoteIP] {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
func generate() string {
	// Создаем массив символов, из которых будет состоять случайная строка
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Создаем буфер для случайной строки нужной длины
	var b = make([]rune, 12)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
