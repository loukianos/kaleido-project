package main

import (
	"bytes"
	"cmp"
	"encoding/base64"
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
	servicerToken, err := fetchToken(httpClient, tokenURL, url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {getenv("SERVICER_CLIENT_ID", "servicer")},
		"client_secret": {getenv("SERVICER_CLIENT_SECRET", "servicer-secret")},
	})
	if err != nil {
		return fmt.Errorf("servicer login: %w", err)
	}
	aliceToken, err := fetchToken(httpClient, tokenURL, url.Values{
		"grant_type": {"password"},
		"client_id":  {getenv("LENDER_CLIENT_ID", "loan-notes-app")},
		"username":   {"alice"},
		"password":   {getenv("ALICE_PASSWORD", "alice-password")},
	})
	if err != nil {
		return fmt.Errorf("alice login: %w", err)
	}
	bobToken, err := fetchToken(httpClient, tokenURL, url.Values{
		"grant_type": {"password"},
		"client_id":  {getenv("LENDER_CLIENT_ID", "loan-notes-app")},
		"username":   {"bob"},
		"password":   {getenv("BOB_PASSWORD", "bob-password")},
	})
	if err != nil {
		return fmt.Errorf("bob login: %w", err)
	}

	// Custodial identities are keyed by the OIDC subject, so the servicer names lenders by the sub claim in their tokens.
	aliceSub, err := subjectOf(aliceToken)
	if err != nil {
		return err
	}
	bobSub, err := subjectOf(bobToken)
	if err != nil {
		return err
	}
	fmt.Printf("alice sub=%s bob sub=%s\n", aliceSub, bobSub)

	servicer := client{baseURL: apiURL, http: httpClient, token: servicerToken}
	alice := client{baseURL: apiURL, http: httpClient, token: aliceToken}
	bob := client{baseURL: apiURL, http: httpClient, token: bobToken}

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

	// Custodial lenders: the note is minted to a per-identity key provisioned on first sight, and transfers are signed by the owning lender's key under their own token.
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

// subjectOf reads the sub claim out of a JWT without verifying it; the demo only needs the identifier, the API does the verification.
func subjectOf(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errors.New("token is not a JWT")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("decode token payload: %w", err)
	}
	var claims struct {
		Sub string `json:"sub"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("parse token claims: %w", err)
	}
	if claims.Sub == "" {
		return "", errors.New("token missing sub claim")
	}
	return claims.Sub, nil
}

type client struct {
	baseURL string
	http    *http.Client
	token   string
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
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
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
