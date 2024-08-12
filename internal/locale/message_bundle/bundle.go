package locale

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/JulioZittei/wsrs-ama-go/internal/middlewares"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle
var localizer *i18n.Localizer

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	filePath1 := filepath.Join(currentDir, "messages.en.json")
	filePath2 := filepath.Join(currentDir, "messages.pt-BR.json")
	bundle = i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile(filePath1)
	bundle.MustLoadMessageFile(filePath2)
}

func GetMessage(ctx context.Context, key string, args ...string) (string, error) {
	arguments := make(map[string]string)

	for i, arg := range args {
		arguments[fmt.Sprintf("Arg%d", i+1)] = arg
	}

	lang, ok := ctx.Value(middlewares.LangKey).(string)
	if !ok {
		lang = "en"
	}

	localizeConfig := i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: arguments,
	}

	localizer = i18n.NewLocalizer(bundle, lang)
	localizedMessage, err := localizer.Localize(&localizeConfig)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return localizedMessage, nil
}
