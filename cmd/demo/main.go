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
	"net/url"
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
	keycloakURL := getenv("KEYCLOAK_URL", "http://localhost:8081")
	lender := getenv("LENDER_ADDRESS", "0x627306090abaB3A6e1400e9345bC60c78a8BEf57")
	newLender := getenv("NEW_LENDER_ADDRESS", "0x1111111111111111111111111111111111111111")

	httpClient := &http.Client{Timeout: 1 * time.Minute}
	public := client{baseURL: apiURL, http: httpClient}

	fmt.Println("Checking readiness at", apiURL)
	if err := public.get("/ready", nil); err != nil {
		return err
	}

	// Every actor authenticates against the seeded Keycloak realm: the servicer via client credentials, lenders via password grants.
	fmt.Println("Fetching tokens from", keycloakURL)
	tokenURL := keycloakURL + "/realms/" + getenv("KEYCLOAK_REALM", "loan-notes") + "/protocol/openid-connect/token"
	servicerClientSecret, err := requiredEnv("SERVICER_CLIENT_SECRET")
	if err != nil {
		return err
	}
	alicePassword, err := requiredEnv("ALICE_PASSWORD")
	if err != nil {
		return err
	}
	bobPassword, err := requiredEnv("BOB_PASSWORD")
	if err != nil {
		return err
	}
	servicerToken, err := fetchToken(httpClient, tokenURL, url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {getenv("SERVICER_CLIENT_ID", "servicer")},
		"client_secret": {servicerClientSecret},
	})
	if err != nil {
		return fmt.Errorf("servicer login: %w", err)
	}
	aliceToken, err := fetchToken(httpClient, tokenURL, url.Values{
		"grant_type": {"password"},
		"client_id":  {getenv("LENDER_CLIENT_ID", "loan-notes-app")},
		"username":   {"alice"},
		"password":   {alicePassword},
	})
	if err != nil {
		return fmt.Errorf("alice login: %w", err)
	}
	bobToken, err := fetchToken(httpClient, tokenURL, url.Values{
		"grant_type": {"password"},
		"client_id":  {getenv("LENDER_CLIENT_ID", "loan-notes-app")},
		"username":   {"bob"},
		"password":   {bobPassword},
	})
	if err != nil {
		return fmt.Errorf("bob login: %w", err)
	}

	servicer := client{baseURL: apiURL, http: httpClient, token: servicerToken}
	alice := client{baseURL: apiURL, http: httpClient, token: aliceToken}
	bob := client{baseURL: apiURL, http: httpClient, token: bobToken}

	// Onboarding is the explicit provisioning step: it creates the lender identity and custodial wallet, and returns the subject + address the lender hands to the servicer.
	fmt.Println("Onboarding alice as a custodial lender")
	var aliceAccount struct {
		Subject string `json:"subject"`
		Address string `json:"address"`
	}
	if err := alice.post("/lenders/onboard", nil, &aliceAccount); err != nil {
		return err
	}
	fmt.Println("Onboarding bob as a custodial lender")
	var bobAccount struct {
		Subject string `json:"subject"`
		Address string `json:"address"`
	}
	if err := bob.post("/lenders/onboard", nil, &bobAccount); err != nil {
		return err
	}
	aliceSub, bobSub := aliceAccount.Subject, bobAccount.Subject

	// The API refuses anonymous callers outright.
	fmt.Println("Verifying anonymous requests are rejected")
	if err := public.get("/loans", nil); err == nil {
		return errors.New("anonymous loan list should have been rejected")
	} else {
		fmt.Println("Rejected as designed:", err)
	}

	// The platform signer address doubles as the custody address for warehouse originations.
	fmt.Println("Reading service info")
	var info struct {
		SignerAddress string `json:"signer_address"`
	}
	if err := public.get("/", &info); err != nil {
		return err
	}
	// The lender must be a distinct party or the refused-transfer demonstration below would succeed.
	if strings.EqualFold(lender, info.SignerAddress) {
		return fmt.Errorf("LENDER_ADDRESS %s is the platform signer; set it to a different address", lender)
	}

	// Deploy is skipped when a contract already exists (e.g. on demo re-runs)
	if err := servicer.get("/contracts/active", nil); err != nil {
		fmt.Println("Deploying active LoanNote contract")
		if err := servicer.post("/admin/contracts/deploy", map[string]any{}, nil); err != nil {
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
	err = servicer.post("/loans", map[string]any{
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
	if err := servicer.post(fmt.Sprintf("/loans/%d/transfer", loan.ID), map[string]any{"to_address": lender}, nil); err != nil {
		return err
	}

	fmt.Printf("Repaying loan %d\n", loan.ID)
	err = servicer.post(fmt.Sprintf("/loans/%d/repayments", loan.ID), map[string]any{
		"amount_minor": loan.TotalDueMinor,
		"external_ref": fmt.Sprintf("demo-final-%d", loan.ID),
	}, nil)
	if err != nil {
		return err
	}

	fmt.Println("Reading final loan state")
	if err := servicer.get(fmt.Sprintf("/loans/%d", loan.ID), nil); err != nil {
		return err
	}

	// Lender-owned flow: the note is minted straight to the lender, so the platform provably cannot move it.
	fmt.Println("Originating lender-owned loan note for default flow")
	var defaultLoan struct {
		ID int64 `json:"id"`
	}
	err = servicer.post("/loans", map[string]any{
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
	if err := servicer.post(fmt.Sprintf("/loans/%d/transfer", defaultLoan.ID), map[string]any{"to_address": newLender}, nil); err == nil {
		return errors.New("transfer of a lender-owned note should have been refused")
	} else {
		fmt.Println("Refused as designed:", err)
	}

	fmt.Printf("Marking loan %d defaulted\n", defaultLoan.ID)
	if err := servicer.post(fmt.Sprintf("/loans/%d/default", defaultLoan.ID), nil, nil); err != nil {
		return err
	}

	fmt.Println("Reading defaulted loan state")
	if err := servicer.get(fmt.Sprintf("/loans/%d", defaultLoan.ID), nil); err != nil {
		return err
	}

	// Provisioning strictly precedes participation: naming a lender who never onboarded is refused.
	fmt.Println("Attempting origination to a never-onboarded lender (expected to be refused)")
	err = servicer.post("/loans", map[string]any{
		"borrower_ref":    fmt.Sprintf("demo-unknown-%d", time.Now().Unix()),
		"lender_subject":  "never-onboarded",
		"principal_minor": 1000,
		"apr_bps":         0,
		"term_days":       30,
	}, nil)
	if err == nil {
		return errors.New("origination to a never-onboarded lender should have been refused")
	}
	fmt.Println("Refused as designed:", err)

	// Custodial lenders: the note is minted to the wallet provisioned at onboarding, and transfers are signed by the owning lender's key under their own token.
	fmt.Println("Originating loan note for custodial lender alice")
	var custodialLoan struct {
		ID            int64  `json:"id"`
		LenderAddress string `json:"lender_address"`
		LenderSubject string `json:"lender_subject"`
		TotalDueMinor int64  `json:"total_due_minor"`
	}
	err = servicer.post("/loans", map[string]any{
		"borrower_ref":    fmt.Sprintf("demo-custodial-%d", time.Now().Unix()),
		"lender_subject":  aliceSub,
		"principal_minor": 20000,
		"apr_bps":         300,
		"term_days":       60,
	}, &custodialLoan)
	if err != nil {
		return err
	}
	if custodialLoan.LenderSubject != aliceSub {
		return fmt.Errorf("custodial loan lender subject is %q, want %q", custodialLoan.LenderSubject, aliceSub)
	}
	if !strings.EqualFold(custodialLoan.LenderAddress, aliceAccount.Address) {
		return fmt.Errorf("note minted to %s, want alice's onboarded wallet %s", custodialLoan.LenderAddress, aliceAccount.Address)
	}
	if strings.EqualFold(custodialLoan.LenderAddress, info.SignerAddress) {
		return errors.New("alice's custodial address must differ from the platform signer")
	}

	// Another lender can't move alice's note any more than the platform can.
	fmt.Printf("Attempting transfer of alice's loan %d as bob (expected to be refused)\n", custodialLoan.ID)
	if err := bob.post(fmt.Sprintf("/loans/%d/transfer", custodialLoan.ID), map[string]any{"to_subject": bobSub}, nil); err == nil {
		return errors.New("bob transferring alice's note should have been refused")
	} else {
		fmt.Println("Refused as designed:", err)
	}

	fmt.Printf("Transferring loan %d from alice to bob (alice's token, alice's custodial key)\n", custodialLoan.ID)
	var custodialTransfer struct {
		LenderSubject string `json:"lender_subject"`
		LenderAddress string `json:"lender_address"`
	}
	if err := alice.post(fmt.Sprintf("/loans/%d/transfer", custodialLoan.ID), map[string]any{"to_subject": bobSub}, &custodialTransfer); err != nil {
		return err
	}
	if custodialTransfer.LenderSubject != bobSub {
		return fmt.Errorf("transferred loan lender subject is %q, want %q", custodialTransfer.LenderSubject, bobSub)
	}

	fmt.Printf("Transferring loan %d back to alice (bob's token, bob's custodial key)\n", custodialLoan.ID)
	if err := bob.post(fmt.Sprintf("/loans/%d/transfer", custodialLoan.ID), map[string]any{"to_subject": aliceSub}, nil); err != nil {
		return err
	}

	// Lenders see only their own loans.
	fmt.Println("Listing loans as alice (scoped to her holdings)")
	var aliceLoans struct {
		Loans []struct {
			ID int64 `json:"id"`
		} `json:"loans"`
	}
	if err := alice.get("/loans?limit=100", &aliceLoans); err != nil {
		return err
	}
	sawOwn := false
	for _, item := range aliceLoans.Loans {
		if item.ID == custodialLoan.ID {
			sawOwn = true
		}
		if item.ID == loan.ID || item.ID == defaultLoan.ID {
			return fmt.Errorf("alice's loan list leaked loan %d", item.ID)
		}
	}
	if !sawOwn {
		return fmt.Errorf("alice's loan list is missing her loan %d", custodialLoan.ID)
	}

	fmt.Printf("Repaying loan %d\n", custodialLoan.ID)
	err = servicer.post(fmt.Sprintf("/loans/%d/repayments", custodialLoan.ID), map[string]any{
		"amount_minor": custodialLoan.TotalDueMinor,
		"external_ref": fmt.Sprintf("demo-custodial-final-%d", custodialLoan.ID),
	}, nil)
	if err != nil {
		return err
	}

	// Each contract instance is its own loan series; originations select one by contract_id.
	fmt.Println("Deploying a second loan series (non-default contract)")
	var series struct {
		ID int64 `json:"id"`
	}
	if err := servicer.post("/admin/contracts/deploy", map[string]any{}, &series); err != nil {
		return err
	}

	fmt.Printf("Originating lender-owned loan note on series %d\n", series.ID)
	var seriesLoan struct {
		ID         int64 `json:"id"`
		ContractID int64 `json:"contract_id"`
	}
	err = servicer.post("/loans", map[string]any{
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

	// Pooled signing + failure tolerance: a burst larger than the pool spreads across the keys, and overflow that finds every lock busy comes back 202 for the reconciler to drain.
	const concurrentMints = 6
	fmt.Printf("Originating %d loans concurrently across the servicer key pool\n", concurrentMints)
	quiet := servicer
	quiet.quiet = true
	type mintResult struct {
		id       int64
		accepted bool
		err      error
	}
	results := make(chan mintResult, concurrentMints)
	for i := range concurrentMints {
		go func(i int) {
			var out struct {
				ID int64 `json:"id"`
			}
			status, err := quiet.postStatus("/loans", map[string]any{
				"borrower_ref":    fmt.Sprintf("demo-pool-%d-%d", time.Now().Unix(), i),
				"lender_address":  lender,
				"principal_minor": 1000 + i,
				"apr_bps":         100,
				"term_days":       30,
				"external_ref":    fmt.Sprintf("demo-pool-%d-%d", time.Now().Unix(), i),
			}, &out)
			results <- mintResult{id: out.ID, accepted: status == http.StatusAccepted, err: err}
		}(i)
	}
	var mintIDs []int64
	acceptedCount := 0
	for range concurrentMints {
		result := <-results
		if result.err != nil {
			return fmt.Errorf("concurrent origination: %w", result.err)
		}
		if result.accepted {
			acceptedCount++
		}
		mintIDs = append(mintIDs, result.id)
	}
	fmt.Printf("%d mints completed synchronously, %d accepted as pending for the reconciler\n", concurrentMints-acceptedCount, acceptedCount)

	// Every loan converges to active with no client retries: pending ones are re-driven by the reconciler.
	fmt.Println("Polling until all loans are active")
	mintSigners := map[string]bool{}
	deadline := time.Now().Add(2 * time.Minute)
	for _, id := range mintIDs {
		for {
			var minted struct {
				Status            string `json:"status"`
				MintSignerAddress string `json:"mint_signer_address"`
			}
			if err := quiet.get(fmt.Sprintf("/loans/%d", id), &minted); err != nil {
				return err
			}
			if minted.Status == "active" {
				mintSigners[minted.MintSignerAddress] = true
				break
			}
			if minted.Status != "originating" {
				return fmt.Errorf("loan %d ended up %s, want active", id, minted.Status)
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("loan %d still %s after deadline", id, minted.Status)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
	fmt.Printf("%d concurrent mints converged to active, signed by %d distinct pool keys\n", concurrentMints, len(mintSigners))
	if len(mintSigners) < 2 {
		return errors.New("concurrent mints should have spread across at least two pool keys")
	}

	fmt.Println("Listing contracts")
	if err := servicer.get("/contracts", nil); err != nil {
		return err
	}

	fmt.Println("Reading all loans")
	if err := servicer.get("/loans?limit=100", nil); err != nil {
		return err
	}

	fmt.Println("Demo complete")
	return nil
}

// fetchToken performs an OAuth token request and returns the access token.
func fetchToken(httpClient *http.Client, tokenURL string, form url.Values) (string, error) {
	resp, err := httpClient.PostForm(tokenURL, form)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint: %s: %s", resp.Status, bytes.TrimSpace(body))
	}
	var payload struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}
	if payload.AccessToken == "" {
		return "", errors.New("token response missing access_token")
	}
	return payload.AccessToken, nil
}

type client struct {
	baseURL string
	http    *http.Client
	token   string
	// quiet suppresses response printing, for concurrent calls whose output would interleave.
	quiet bool
}

func (c client) get(path string, out any) error {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c client) post(path string, body any, out any) error {
	_, err := c.postStatus(path, body, out)
	return err
}

func (c client) postStatus(path string, body any, out any) (int, error) {
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return 0, err
		}
		reader = bytes.NewReader(encoded)
	}
	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, reader)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.doStatus(req, out)
}

func (c client) do(req *http.Request, out any) error {
	_, err := c.doStatus(req, out)
	return err
}

func (c client) doStatus(req *http.Request, out any) (int, error) {
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, fmt.Errorf("%s %s: %s: %s", req.Method, req.URL.Path, resp.Status, bytes.TrimSpace(body))
	}

	if !c.quiet {
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, body, "", "  "); err != nil {
			fmt.Println(string(body))
		} else {
			fmt.Println(pretty.String())
		}
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return resp.StatusCode, fmt.Errorf("decode %s response: %w", req.URL.Path, err)
		}
	}
	return resp.StatusCode, nil
}

func getenv(key, fallback string) string {
	return cmp.Or(os.Getenv(key), fallback)
}

func requiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s must be set", key)
	}
	return value, nil
}
