package api

import (
	"io"
	"log/slog"
	"net/http"
)

func newTestHandler(opts Options) http.Handler {
	if opts.Loans == nil {
		opts.Loans = &fakeLoansService{}
	}
	if opts.Contracts == nil {
		opts.Contracts = &fakeContractsService{}
	}
	return New("test", slog.New(slog.NewTextHandler(io.Discard, nil)), opts)
}
