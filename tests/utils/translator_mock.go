package utils

import "context"

type TranslatorMock struct {
	EnToRu map[string]string
}

func (t *TranslatorMock) TranslateEnToRu(_ context.Context, text string) (string, error) {
	return t.EnToRu[text], nil
}

func (t *TranslatorMock) Close() error {
	return nil
}
