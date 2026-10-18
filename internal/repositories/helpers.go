package repositories

import "github.com/zaem/backend-go/pkg/cuid"

func cuidNew() string {
	return cuid.New()
}