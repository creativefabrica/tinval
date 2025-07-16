# TINVal

![Build](https://github.com/creativefabrica/tinval/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/creativefabrica/tinval)](https://goreportcard.com/report/github.com/creativefabrica/tinval)
[![GoDoc](https://godoc.org/github.com/creativefabrica/tinval?status.svg)](https://godoc.org/github.com/creativefabrica/tinval)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/creativefabrica/vat/master/LICENSE)

Package for parsing and validating Tax Identification numbers (TINs)

based on https://github.com/Teamwork/vat with some different design choices

## Installing

```shell
go get https://github.com/creativefabrica/tinval
```

```shell
import "github.com/creativefabrica/tinval"
```

## Usage

Parsing a Tax ID Number (TIN):

```go
tin, err := tinval.Parse("NL822010690B01")
if err != nil {
    fmt.Printf("Invalid tax id number: %s\n", err)
    return
}
fmt.Printf("Country Code: %s Number: %s\n", tin.CountryCode, tin.Number)
```

You can also use the `Must` variant if you want to `panic` on error; this is useful on tests:

```go
tin := tinval.MustParse("INVALID")
```

For validating that a TIN actually exists, different APIs are used depending of the type of TIN:

* AU ANBs are looked up using Australian government [ABN Search Service](https://abr.business.gov.au/abrxmlsearch).
* EU VAT numbers are looked up using the [VIES VAT validation API](http://ec.europa.eu/taxation_customs/vies/).
* UK VAT numbers are looked up
using the [UK GOV VAT validation API](https://developer.service.hmrc.gov.uk/api-documentation/docs/api/service/vat-registered-companies-api/2.0)
    * Requires [signing up for the UK API](https://developer.service.hmrc.gov.uk/api-documentation/docs/using-the-hub).

You can pass the clients implemented on the `abn`, `euvat` and `ukvat` packages as functional options to the vat Validator initializer:

```go
validator := tinval.NewValidator(
    tinval.WithEUVATClient(euvat.NewClient()),
    tinval.WithUKVATClient(ukvat.NewClient(
        ukvat.ClientCredentials{
            Secret: os.Getenv("UKVAT_API_CLIENT_SECRET"),
            ID:     os.Getenv("UKVAT_API_CLIENT_ID"),
        },
    )),
    tinval.WithANBClient(
        abn.NewClient(
            os.Getenv("ABN_API_AUTH_GUID"),
            abn.WithHTTPClient(httpClient),
        ),
    ),
)

err := validator.Validate(context.Background(), "GB146295999727")
if err != nil {
    return err
}
```

If you only need EU validation and/or UK validation for some reason, you can skip passing the unneeded clients.<br>
In this case the `Validate` function will only validate format using the `Parse` function.

[Full example](/example/main.go)

### Package usage: euvat

```go
httpClient := &http.Client{}
client := euvat.NewClient(
    // Use this option to provide a custom http client
    euvat.WithHTTPClient(httpClient),
    // Use this option to enable retries in case of rate limiting from the VIES API
    euvat.WithRetries(3),
)
```

### Package usage: ukvat

> [!IMPORTANT]
> For validating VAT numbers that begin with **GB** you will need to [sign up](https://developer.service.hmrc.gov.uk/api-documentation/docs/using-the-hub) to gain access to the UK government's VAT API.
> Once you have signed up and acquired a client ID and client secret you can provide them on the intitalizer

```go
httpClient := &http.Client{}
client := ukvat.NewClient(
    ukvat.ClientCredentials{
        Secret: os.Getenv("UKVAT_API_CLIENT_SECRET"),
        ID:     os.Getenv("UKVAT_API_CLIENT_ID"),
    },
    // Use this option to provide a custom http client
    ukvat.WithHTTPClient(httpClient),
)
```

> [!NOTE]
> The `ukvat.Client` struct will cache the auth token needed for the validation requests.
> To avoid getting `403` responses when validating VAT numbers, the client will refresh the token 2 minutes before it expires

If you need to hit the sandbox version of the UK VAT API you can use the following option:

```go
ukvat.WithBaseURL(ukvat.TestServiceBaseURL)
```

### Package usage: abn

> [!IMPORTANT]
> For validating Australian VAT numbers (or ABNs) that begin with **AU** you will need to [register](https://abr.business.gov.au/Tools/WebServicesRegister?AcceptLicenceTerms=Y) for an authentication GUID.

```go
httpClient := &http.Client{}
client := abn.NewClient(
    os.Getenv("ABN_API_AUTH_GUID"),
    // Use this option to provide a custom http client
    abn.WithHTTPClient(httpClient),
)
```

### Package usage: tinvaltest

You can use this package to provide a mock validation client to the tinvaltest.Validator.
This is useful in tests:

```go
validationClientMock := tinvaltest.NewMockValidationClient(gomock.NewController(t))
validator := tinval.NewValidator(
    tinval.WithUKVATClient(validationClientMock),
    tinval.WithEUVATClient(validationClientMock),
)
```