package cuid

import "github.com/lucsky/cuid"

// New generates a unique CUID2 string (Prisma-compatible)
func New() string {
	return cuid.New()
}