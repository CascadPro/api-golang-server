package core_i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

var (
	bundle *i18n.Bundle

	matcher = language.NewMatcher([]language.Tag{
		language.Russian,
		language.English,
	})
)

func Init() error {
	bundle = i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, locale := range []string{"en", "ru"} {
		if _, err := bundle.LoadMessageFileFS(
			localeFS,
			"locales/"+locale+".json",
		); err != nil {
			return fmt.Errorf("load %s locale: %w", locale, err)
		}
	}

	return nil
}

func Resolve(acceptLanguage string) language.Tag {
	if strings.TrimSpace(acceptLanguage) == "" {
		return language.English
	}

	tags, _, err := language.ParseAcceptLanguage(acceptLanguage)
	if err != nil || len(tags) == 0 {
		return language.English
	}

	tag, _ := language.MatchStrings(matcher, tags[0].String())

	return tag
}

func Translate(locale language.Tag, code, fallback string) string {
	if bundle == nil {
		return fallback
	}

	localizer := i18n.NewLocalizer(bundle, locale.String(), "en")

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: code,
	})

	if err != nil || message == "" {
		return fallback
	}

	return message
}
