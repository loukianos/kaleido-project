package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"kaleido-project/db/sqlc"
	"kaleido-project/internal/contracts"
)

func TestDeployContract(t *testing.T) {
	service := &fakeContractsService{
		deployContract: db.Contract{
			ID:           1,
			ChainID:      1337,
			Address:      "0x123",
			DeployTxHash: db.Ptr("0xabc"),
			BaseUri:      "http://localhost:8080/loans/",
			Active:       true,
		},
	}
	handler := newTestHandler(Options{
		Contracts: service,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/admin/contracts/deploy", strings.NewReader(`{"base_uri":"https://example.test/loans/"}`))
	handler.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, "https://example.test/loans/", service.deployBaseURI)

	var body contractResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "0x123", body.Address)
	require.Equal(t, "0xabc", body.DeployTxHash)
	require.True(t, body.Active)
}

func TestDeployContractLockBusy(t *testing.T) {
	handler := newTestHandler(Options{
		Contracts: &fakeContractsService{deployErr: db.ErrLockBusy},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/admin/contracts/deploy", strings.NewReader(`{}`))
	handler.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
}

func TestDeployContractAlreadyDeployed(t *testing.T) {
	handler := newTestHandler(Options{
		Contracts: &fakeContractsService{deployErr: contracts.ErrContractAlreadyDeployed},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/admin/contracts/deploy", strings.NewReader(`{}`))
	handler.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestActiveContract(t *testing.T) {
	handler := newTestHandler(Options{
		Contracts: &fakeContractsService{
			activeContract: db.Contract{
				ID:      2,
				ChainID: 1337,
				Address: "0x456",
				BaseUri: "http://localhost:8080/loans/",
				Active:  true,
			},
		},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/contracts/active", nil)
	handler.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)

	var body contractResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "0x456", body.Address)
}

func TestActiveContractNotFound(t *testing.T) {
	handler := newTestHandler(Options{
		Contracts: &fakeContractsService{activeErr: pgx.ErrNoRows},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/contracts/active", nil)
	handler.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

type fakeContractsService struct {
	deployBaseURI  string
	deployContract db.Contract
	deployErr      error
	activeContract db.Contract
	activeErr      error
}

func (f *fakeContractsService) Deploy(_ context.Context, baseURI string) (db.Contract, error) {
	f.deployBaseURI = baseURI
	if f.deployErr != nil {
		return db.Contract{}, f.deployErr
	}
	return f.deployContract, nil
}

func (f *fakeContractsService) ActiveContract(context.Context) (db.Contract, error) {
	if f.activeErr != nil {
		return db.Contract{}, f.activeErr
	}
	if f.activeContract.ID == 0 {
		return db.Contract{}, errors.New("missing fake active contract")
	}
	return f.activeContract, nil
}
