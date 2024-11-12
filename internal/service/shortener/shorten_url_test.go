package shortener

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestShortenURL(t *testing.T) {
	s := &Service{}
	src := rand.NewSource(time.Now().UnixNano())

	t.Run("ValidLength", func(t *testing.T) {
		length := 10
		shortURL := s.ShortenURL(length, src)

		if len(shortURL) != length {
			t.Errorf("Expected length %d, got %d", length, len(shortURL))
		}
	})

	t.Run("ValidCharacters", func(t *testing.T) {
		length := 20
		shortURL := s.ShortenURL(length, src)
		for _, char := range shortURL {
			if !strings.ContainsRune(letterBytes, char) {
				t.Errorf("Invalid character found: %c", char)
			}
		}
	})

	t.Run("Uniqueness", func(t *testing.T) {
		length := 15
		generated := make(map[string]bool)
		for i := 0; i < 100; i++ {
			shortURL := s.ShortenURL(length, src)
			if generated[shortURL] {
				t.Errorf("Duplicate URL found: %s", shortURL)
			}
			generated[shortURL] = true
		}
	})
}
