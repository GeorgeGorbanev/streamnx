package translator

import (
	"context"
)

type Translator interface {
	TranslateEnToRu(context.Context, string) (string, error)
	Close() error
}
