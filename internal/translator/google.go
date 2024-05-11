package translator

import (
	"context"
	"fmt"

	googleTranslate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type GoogleClient struct {
	client    *googleTranslate.TranslationClient
	projectID string
}

type GoogleCredentials struct {
	APIKeyJSON string
	ProjectID  string
}

func NewGoogleClient(ctx context.Context, gc *GoogleCredentials) (*GoogleClient, error) {
	credentials, err := google.CredentialsFromJSON(ctx, []byte(gc.APIKeyJSON), googleTranslate.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to create credentials: %w", err)
	}

	tc, err := googleTranslate.NewTranslationClient(ctx, option.WithCredentials(credentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create translation client: %w", err)
	}

	return &GoogleClient{
		client:    tc,
		projectID: gc.ProjectID,
	}, nil
}

func (gc *GoogleClient) Close() error {
	return gc.client.Close()
}

func (gc *GoogleClient) TranslateEnToRu(ctx context.Context, text string) (string, error) {
	req := gc.request(text, "en", "ru")
	resp, err := gc.client.TranslateText(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to translate text: %w", err)
	}
	return resp.Translations[0].TranslatedText, nil
}

func (gc *GoogleClient) request(text, source, target string) *translatepb.TranslateTextRequest {
	return &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", gc.projectID),
		MimeType:           "text/plain",
		SourceLanguageCode: source,
		TargetLanguageCode: target,
		Contents:           []string{text},
	}
}
