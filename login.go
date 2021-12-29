package cidirectinvesting

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type SessionResponse struct {
	Id               int      `json:"id"`
	Email            string   `json:"email"`
	Name             string   `json:"name"`
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	Roles            []string `json:"roles"`
	IntercomUserHash string   `json:"intercom_user_hash"`
	Archived         bool     `json:"archived"`
	AdvisoryFee      float64  `json:"advisory_fee"`
	Preferences      struct {
		Dashboard struct {
			ChartRange string `json:"chart_range"`
		} `json:"dashboard"`
	} `json:"preferences"`
	HasAccounts struct {
		Id         int `json:"id"`
		UserId     int `json:"user_id"`
		Properties struct {
			Balance    float64 `json:"balance"`
			ClientName string  `json:"client_name"`
		} `json:"properties"`
		CreatedAt              time.Time   `json:"created_at"`
		UpdatedAt              time.Time   `json:"updated_at"`
		Status                 string      `json:"status"`
		Payment                interface{} `json:"payment"`
		LastBalanced           interface{} `json:"last_balanced"`
		Archived               bool        `json:"archived"`
		Currency               string      `json:"currency"`
		AccountType            string      `json:"account_type"`
		EtfReady               bool        `json:"etf_ready"`
		CustodianAccountType   string      `json:"custodian_account_type"`
		PortfolioModelId       interface{} `json:"portfolio_model_id"`
		DoNotTrade             bool        `json:"do_not_trade"`
		DoNotTradeReason       interface{} `json:"do_not_trade_reason"`
		DoNotTradeChangedBy    interface{} `json:"do_not_trade_changed_by"`
		DoNotTradeChangedAt    interface{} `json:"do_not_trade_changed_at"`
		AutomateProfileChanges bool        `json:"automate_profile_changes"`
		TransferOut            interface{} `json:"transfer_out"`
		TransferOutChangedBy   interface{} `json:"transfer_out_changed_by"`
		TransferOutChangedAt   interface{} `json:"transfer_out_changed_at"`
	} `json:"has_accounts"`
	HasActiveAccounts          bool   `json:"has_active_accounts"`
	HasTransfers               bool   `json:"has_transfers"`
	HasActiveTransfers         bool   `json:"has_active_transfers"`
	HasPlans                   bool   `json:"has_plans"`
	HasApplication             bool   `json:"has_application"`
	IsArchivable               bool   `json:"is_archivable"`
	EmailConfirmed             bool   `json:"email_confirmed"`
	HasExternalAdvisor         bool   `json:"has_external_advisor"`
	HideBeneficiaryDesignation bool   `json:"hide_beneficiary_designation"`
	TwoFactorEnabled           bool   `json:"two_factor_enabled"`
	PromptTwoFactorSignup      bool   `json:"prompt_two_factor_signup"`
	Language                   string `json:"language"`
	IntercomUserHashIos        string `json:"intercom_user_hash_ios"`
	IntercomUserHashAndroid    string `json:"intercom_user_hash_android"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Client) Login() (*SessionResponse, error) {
	requestBody := &bytes.Buffer{}
	if err := json.NewEncoder(requestBody).Encode(loginRequest{
		Email:    c.apiKey,
		Password: c.apiSecret,
	}); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "/sessions", requestBody)
	if err != nil {
		return nil, err
	}

	var response = &SessionResponse{}
	if err := c.doWithProcessing(req, response, []int{http.StatusOK}); err != nil {
		return nil, err
	}
	return response, nil
}
