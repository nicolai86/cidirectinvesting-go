package cidirectinvesting

import (
	"fmt"
	"net/http"
	"time"
)

type AccountDetailRequest struct {
	AccountID int
}

type Position struct {
	InvestmentAccountId int     `json:"investment_account_id"`
	InsertDate          string  `json:"insert_date"`
	CurrencyCode        string  `json:"currency_code"`
	Symbol              string  `json:"symbol"`
	Cusip               *string `json:"cusip"`
	SecurityType        string  `json:"security_type"`
	ExchangeRate        float64 `json:"exchange_rate"`
	PriceCad            float64 `json:"price_cad"`
	Quantity            float64 `json:"quantity"`
	Value               float64 `json:"value"`
	ValueCad            float64 `json:"value_cad"`
	Description         string  `json:"description"`
	Mer                 float64 `json:"mer"`
	AssetClass          string  `json:"asset_class"`
	Allocation          float64 `json:"allocation"`
	FundsUrl            string  `json:"funds_url"`
}

type AssetClassAllocation struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Allocation  float64 `json:"allocation"`
}

type AccountDetailResponse struct {
	User                        User                   `json:"user"`
	Id                          int                    `json:"id"`
	UserId                      int                    `json:"user_id"`
	AccountNumber               string                 `json:"account_number"`
	Currency                    string                 `json:"currency"`
	Status                      int                    `json:"status"`
	Archived                    bool                   `json:"archived"`
	JoinedAccounts              interface{}            `json:"joined_accounts"`
	SupportsDeposits            bool                   `json:"supports_deposits"`
	SupportsWithdrawals         bool                   `json:"supports_withdrawals"`
	MaskedAccountNumber         string                 `json:"masked_account_number"`
	Name                        string                 `json:"name"`
	DefaultName                 string                 `json:"default_name"`
	AccountName                 string                 `json:"account_name"`
	AccountType                 string                 `json:"account_type"`
	AccountTypeFriendly         string                 `json:"account_type_friendly"`
	CustodianName               string                 `json:"custodian_name"`
	Custodian                   string                 `json:"custodian"`
	CustodianAlias              string                 `json:"custodian_alias"`
	CustodianId                 int                    `json:"custodian_id"`
	ClientName                  string                 `json:"client_name"`
	PortfolioModelName          string                 `json:"portfolio_model_name"`
	CreatedAt                   time.Time              `json:"created_at"`
	UpdatedAt                   time.Time              `json:"updated_at"`
	InceptionDate               string                 `json:"inception_date"`
	EndingBalanceDate           string                 `json:"ending_balance_date"`
	BeginningBalanceDate        string                 `json:"beginning_balance_date"`
	NetContributions            float64                `json:"net_contributions"`
	TotalContributions          float64                `json:"total_contributions"`
	Fees                        float64                `json:"fees"`
	TotalWithdrawals            float64                `json:"total_withdrawals"`
	Withdrawals                 float64                `json:"withdrawals"`
	Deposits                    float64                `json:"deposits"`
	BeginningBalance            float64                `json:"beginning_balance"`
	StartingBalance             float64                `json:"starting_balance"`
	InitialValue                float64                `json:"initial_value"`
	Balance                     float64                `json:"balance"`
	EndingBalance               float64                `json:"ending_balance"`
	TotalValue                  float64                `json:"total_value"`
	MarketValue                 float64                `json:"market_value"`
	MarketValueEndingBalance    float64                `json:"market_value_ending_balance"`
	MarketValueBeginningBalance float64                `json:"market_value_beginning_balance"`
	CashValue                   float64                `json:"cash_value"`
	CashValueEndingBalance      float64                `json:"cash_value_ending_balance"`
	CashValueBeginningBalance   float64                `json:"cash_value_beginning_balance"`
	WithdrawalsValue            float64                `json:"withdrawals_value"`
	TradingBalance              float64                `json:"trading_balance"`
	Gains                       float64                `json:"gains"`
	RateOfReturn                float64                `json:"rate_of_return"`
	AssetClassAllocation        []AssetClassAllocation `json:"asset_class_allocation"`
	Positions                   []Position             `json:"positions"`
	Mer                         float64                `json:"mer"`
	FundedStatus                string                 `json:"funded_status"`
}

func (c *Client) AccountDetails(request AccountDetailRequest) (*AccountDetailResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/investment_accounts/%d", request.AccountID), nil)
	if err != nil {
		return nil, err
	}

	values := req.URL.Query()
	values.Set("basic", "true") // no idea what this does
	req.URL.RawQuery = values.Encode()

	var response = &AccountDetailResponse{}
	if err := c.doWithProcessing(req, response, []int{http.StatusOK}); err != nil {
		return nil, err
	}

	return response, nil
}
