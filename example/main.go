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
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	tin, err := tinval.Parse("NL822010690B01", "NL")
	if err != nil {
		logger.ErrorContext(ctx, "invalid tax id number", slog.Any("error", err))
		os.Exit(1)

		return
	}

	logger.InfoContext(
		ctx, "parsed tax id number",
		slog.String("country_code", tin.CountryCode),
		slog.String("number", tin.Number),
	)

	tinval.MustParse("NL822010690B01", "NL")
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

	tins := []struct {
		value   string
		country string
	}{
		{
			value:   "GB146295999727",
			country: "GB",
		},
		{
			value:   "NL822010690B01",
			country: "NL",
		},
		{
			value:   "NL822010690B02",
			country: "NL",
		},
		{
			value:   "GB123456789",
			country: "GB",
		},
		{
			value:   "51824753556",
			country: "AU",
		},
		{
			value:   "41824753556",
			country: "AU",
		},
		{
			value:   "EL123456789",
			country: "GR",
		},
		{
			value:   "XI123456789",
			country: "GB",
		},
	}

	for _, tin := range tins {
		err = validator.Validate(context.Background(), tin.value, tin.country)
		if err != nil {
			logger.ErrorContext(
				ctx,
				"tax id number is invalid",
				slog.Any("error", err),
				slog.Any("tin", tin),
			)

			continue
		}

		logger.InfoContext(ctx, "tax id number is valid", slog.Any("tin", tin))
	}
}
