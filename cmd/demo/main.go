package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "demo failed:", err)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}

	apiURL := getenv("API_URL", "http://localhost:8080")
	lender := getenv("LENDER_ADDRESS", "0x627306090abaB3A6e1400e9345bC60c78a8BEf57")
	newLender := getenv("NEW_LENDER_ADDRESS", "0x1111111111111111111111111111111111111111")

	api := client{baseURL: apiURL, http: &http.Client{Timeout: 1 * time.Minute}}

	fmt.Println("Checking readiness at", apiURL)
	if err := api.get("/ready", nil); err != nil {
		return err
	}

	// The platform signer address doubles as the custody address for warehouse originations.
	fmt.Println("Reading service info")
	var info struct {
		SignerAddress string `json:"signer_address"`
	}
	if err := api.get("/", &info); err != nil {
		return err
	}
	// The lender must be a distinct party or the refused-transfer demonstration below would succeed.
	if strings.EqualFold(lender, info.SignerAddress) {
		return fmt.Errorf("LENDER_ADDRESS %s is the platform signer; set it to a different address", lender)
	}

	// Deploy is skipped when a contract already exists (e.g. on demo re-runs)
	if err := api.get("/contracts/active", nil); err != nil {
		fmt.Println("Deploying active LoanNote contract")
		if err := api.post("/admin/contracts/deploy", map[string]any{}, nil); err != nil {
			return err
		}
	} else {
		fmt.Println("Reusing existing LoanNote contract")
	}

	// Warehouse flow: the note is originated into platform custody, then sold to a lender with an owner-signed transfer.
	fmt.Println("Originating warehouse loan note (lender = platform signer)")
	var loan struct {
		ID            int64 `json:"id"`
		TotalDueMinor int64 `json:"total_due_minor"`
	}
	err := api.post("/loans", map[string]any{
		"borrower_ref":    fmt.Sprintf("demo-%d", time.Now().Unix()),
		"lender_address":  info.SignerAddress,
		"principal_minor": 10000,
		"apr_bps":         0,
		"term_days":       30,
	}, &loan)
	if err != nil {
		return err
	}

	fmt.Printf("Transferring warehoused loan %d to lender\n", loan.ID)
	if err := api.post(fmt.Sprintf("/loans/%d/transfer", loan.ID), map[string]any{"to_address": lender}, nil); err != nil {
		return err
	}

	fmt.Printf("Repaying loan %d\n", loan.ID)
	err = api.post(fmt.Sprintf("/loans/%d/repayments", loan.ID), map[string]any{
		"amount_minor": loan.TotalDueMinor,
		"external_ref": fmt.Sprintf("demo-final-%d", loan.ID),
	}, nil)
	if err != nil {
		return err
	}

	fmt.Println("Reading final loan state")
	if err := api.get(fmt.Sprintf("/loans/%d", loan.ID), nil); err != nil {
		return err
	}

	// Lender-owned flow: the note is minted straight to the lender, so the platform provably cannot move it.
	fmt.Println("Originating lender-owned loan note for default flow")
	var defaultLoan struct {
		ID int64 `json:"id"`
	}
	err = api.post("/loans", map[string]any{
		"borrower_ref":    fmt.Sprintf("demo-default-%d", time.Now().Unix()),
		"lender_address":  lender,
		"principal_minor": 7500,
		"apr_bps":         500,
		"term_days":       30,
	}, &defaultLoan)
	if err != nil {
		return err
	}

	fmt.Printf("Attempting platform transfer of lender-owned loan %d (expected to be refused)\n", defaultLoan.ID)
	if err := api.post(fmt.Sprintf("/loans/%d/transfer", defaultLoan.ID), map[string]any{"to_address": newLender}, nil); err == nil {
		return errors.New("transfer of a lender-owned note should have been refused")
	} else {
		fmt.Println("Refused as designed:", err)
	}

	fmt.Printf("Marking loan %d defaulted\n", defaultLoan.ID)
	if err := api.post(fmt.Sprintf("/loans/%d/default", defaultLoan.ID), nil, nil); err != nil {
		return err
	}

	fmt.Println("Reading defaulted loan state")
	if err := api.get(fmt.Sprintf("/loans/%d", defaultLoan.ID), nil); err != nil {
		return err
	}

	// Each contract instance is its own loan series; originations select one by contract_id.
	fmt.Println("Deploying a second loan series (non-default contract)")
	var series struct {
		ID int64 `json:"id"`
	}
	if err := api.post("/admin/contracts/deploy", map[string]any{}, &series); err != nil {
		return err
	}

	fmt.Printf("Originating lender-owned loan note on series %d\n", series.ID)
	var seriesLoan struct {
		ID         int64 `json:"id"`
		ContractID int64 `json:"contract_id"`
	}
	err = api.post("/loans", map[string]any{
		"borrower_ref":    fmt.Sprintf("demo-series-%d", time.Now().Unix()),
		"lender_address":  lender,
		"principal_minor": 5000,
		"apr_bps":         250,
		"term_days":       90,
		"contract_id":     series.ID,
	}, &seriesLoan)
	if err != nil {
		return err
	}
	if seriesLoan.ContractID != series.ID {
		return fmt.Errorf("series loan minted on contract %d, want %d", seriesLoan.ContractID, series.ID)
	}

	fmt.Println("Listing contracts")
	if err := api.get("/contracts", nil); err != nil {
		return err
	}

	fmt.Println("Reading all loans")
	if err := api.get("/loans?limit=100", nil); err != nil {
		return err
	}

	fmt.Println("Demo complete")
	return nil
}

type client struct {
	baseURL string
	http    *http.Client
}

func (c client) get(path string, out any) error {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c client) post(path string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(encoded)
	}
	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, out)
}

func (c client) do(req *http.Request, out any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s %s: %s: %s", req.Method, req.URL.Path, resp.Status, bytes.TrimSpace(body))
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, body, "", "  "); err != nil {
		fmt.Println(string(body))
	} else {
		fmt.Println(pretty.String())
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("decode %s response: %w", req.URL.Path, err)
		}
	}
	return nil
}

func getenv(key, fallback string) string {
	return cmp.Or(os.Getenv(key), fallback)
}
