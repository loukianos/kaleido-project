// Command demo drives the API through the full loan lifecycle: deploy a
// contract, then originate, transfer, and repay a loan to settlement, plus a
// second loan that is marked defaulted. It expects the local stack (API,
// Postgres, Besu) to be running; see the README.
package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "demo failed:", err)
		os.Exit(1)
	}
}

func run() error {
	apiURL := cmp.Or(os.Getenv("API_URL"), "http://localhost:8080")
	lender := cmp.Or(os.Getenv("LENDER_ADDRESS"), "0x627306090abaB3A6e1400e9345bC60c78a8BEf57")
	newLender := cmp.Or(os.Getenv("NEW_LENDER_ADDRESS"), "0x1111111111111111111111111111111111111111")

	api := client{baseURL: apiURL, http: &http.Client{Timeout: 1 * time.Minute}}

	fmt.Println("Checking readiness at", apiURL)
	if err := api.get("/ready", nil); err != nil {
		return err
	}

	// Deploy is skipped when a contract already exists (e.g. on demo re-runs):
	// the base version supports a single contract per chain.
	if err := api.get("/contracts/active", nil); err != nil {
		fmt.Println("Deploying active LoanNote contract")
		if err := api.post("/admin/contracts/deploy", map[string]any{}, nil); err != nil {
			return err
		}
	} else {
		fmt.Println("Reusing existing LoanNote contract")
	}

	fmt.Println("Originating loan note")
	var loan struct {
		ID            int64 `json:"id"`
		TotalDueMinor int64 `json:"total_due_minor"`
	}
	err := api.post("/loans", map[string]any{
		"borrower_ref":    fmt.Sprintf("demo-%d", time.Now().Unix()),
		"lender_address":  lender,
		"principal_minor": 10000,
		"apr_bps":         0,
		"term_days":       30,
	}, &loan)
	if err != nil {
		return err
	}

	fmt.Printf("Transferring loan %d\n", loan.ID)
	if err := api.post(fmt.Sprintf("/loans/%d/transfer", loan.ID), map[string]any{"to_address": newLender}, nil); err != nil {
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

	fmt.Println("Originating loan note for default flow")
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

	fmt.Printf("Marking loan %d defaulted\n", defaultLoan.ID)
	if err := api.post(fmt.Sprintf("/loans/%d/default", defaultLoan.ID), nil, nil); err != nil {
		return err
	}

	fmt.Println("Reading defaulted loan state")
	if err := api.get(fmt.Sprintf("/loans/%d", defaultLoan.ID), nil); err != nil {
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

// do sends the request, pretty-prints the JSON response, and decodes it into
// out when out is non-nil.
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
