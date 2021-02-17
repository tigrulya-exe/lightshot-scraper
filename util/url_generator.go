package util

import (
	"math/rand"
	"time"
)

type urlGenerator struct {
	randGenerator rand.Rand
	alphabet      []byte
}

type UrlGenerator interface {
	Generate(int) string
}

func NewUrlGenerator() UrlGenerator {
	return &urlGenerator{
		alphabet: []byte("abcdefghijklmnopqrstuvwxyz01234567890"),
		randGenerator: *rand.New(
			rand.NewSource(
				time.Now().Unix(),
			),
		),
	}
}

func (g *urlGenerator) Generate(size int) string {
	buf := make([]byte, size)
	alphabetLen := len(g.alphabet)
	for i := range buf {
		buf[i] = g.alphabet[g.randGenerator.Intn(alphabetLen)]
	}
	return string(buf)
}
