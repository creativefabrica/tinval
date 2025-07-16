package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/creativefabrica/tinval"
	"github.com/creativefabrica/tinval/abn"
	"github.com/creativefabrica/tinval/euvat"
	"github.com/creativefabrica/tinval/ukvat"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	tin, err := tinval.Parse("NL822010690B01")
	if err != nil {
		logger.Error("invalid tax id number", "error", err)
		os.Exit(1)

		return
	}

	logger.Info("parsed tax id number", "country_code", tin.CountryCode, "number", tin.Number)

	tinval.MustParse("NL822010690B01")
	retries := 3

	httpClient := &http.Client{}
	validator := tinval.NewValidator(
		tinval.WithEUVATClient(
			euvat.NewClient(
				euvat.WithHTTPClient(httpClient),
				euvat.WithRetries(retries),
			),
		),
		tinval.WithUKVATClient(
			ukvat.NewClient(
				ukvat.ClientCredentials{
					Secret: os.Getenv("UKVAT_API_CLIENT_SECRET"),
					ID:     os.Getenv("UKVAT_API_CLIENT_ID"),
				},
				ukvat.WithBaseURL(ukvat.TestServiceBaseURL),
				ukvat.WithHTTPClient(httpClient),
			),
		),
		tinval.WithANBClient(
			abn.NewClient(
				os.Getenv("ABN_API_AUTH_GUID"),
				abn.WithHTTPClient(httpClient),
			),
		),
	)

	tins := []string{
		"GB146295999727",
		"NL822010690B01",
		"NL822010690B02",
		"GB123456789",
		"AU51824753556",
		"AU41824753556",
	}

	for _, tin := range tins {
		err = validator.Validate(context.Background(), tin)
		if err != nil {
			logger.Error("tax id number is invalid", "error", err, "tin", tin)

			continue
		}

		logger.Info("tax id number is valid", "tin", tin)
	}
}
