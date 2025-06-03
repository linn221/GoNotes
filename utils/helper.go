package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"net"
	"net/http"
	"strings"
)

func NewTrue() *bool {
	b := true
	return &b
}

func NewFalse() *bool {
	b := false
	return &b
}

// returns slice removing duplicate elements
func UniqueSlice[T comparable](slice []T) []T {
	inResult := make(map[T]struct{})
	var result []T
	for _, elm := range slice {
		if _, ok := inResult[elm]; !ok {
			// if not exists in map, append it, otherwise do nothing
			inResult[elm] = struct{}{}
			result = append(result, elm)
		}
	}
	return result
}

// func UniqueLength[T comparable](slice []T) int {

// }

func UniqueSliceWithDuplicateCount[T comparable](slice []T) ([]T, int) {
	inResult := make(map[T]struct{}, len(slice))
	var duplicates int
	unqSlice := make([]T, 0, len(slice))
	for _, elm := range slice {
		_, found := inResult[elm]
		if found {
			duplicates++
		} else {
			inResult[elm] = struct{}{}
			unqSlice = append(unqSlice, elm)
		}
	}

	return unqSlice, duplicates
}

// turn slice into map for faster look up
func SliceToMap[T any, K comparable](slice []T, getKey func(T) K) map[K]T {
	m := make(map[K]T, len(slice))
	for _, v := range slice {
		k := getKey(v)
		m[k] = v
	}

	return m
}

func NewExistsChecker[T comparable](slice []T) func(T) bool {
	m := make(map[T]struct{}, len(slice))
	for _, elm := range slice {
		m[elm] = struct{}{}
	}
	return func(t T) bool {
		_, found := m[t]
		return found
	}
}
func NewDuplicateChecker[T comparable](count int) func(T) bool {
	m := make(map[T]struct{}, count)
	return func(t T) bool {
		_, ok := m[t]
		if ok {
			return true
		} else {
			m[t] = struct{}{}
			return false
		}
	}

}
func SanitizeStr(s string) string {
	return strings.TrimSpace(s)
}

func HashString(s string) string {

	h := sha1.New()
	h.Write([]byte(s))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

// ai generated one
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (may contain multiple IPs)
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		// Return the first IP in the list
		return strings.Split(fwd, ",")[0]
	}

	// Fallback to X-Real-IP
	if realIP := r.Header.Get("X-Real-Ip"); realIP != "" {
		return realIP
	}

	// Fallback to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // return as is if parsing fails
	}
	return host
}
