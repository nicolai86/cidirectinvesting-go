package cidirectinvesting

import (
	"net/http"
	"time"
)

type Currency string

var (
	CurrencyCAD Currency = "cad"
)

func toRequestParameter(currency Currency) string {
	return string(currency)
}

type DashboardRequest struct {
	Currencies map[Currency]struct{}
	StartDate  time.Time
	EndDate    time.Time
}

type Account struct {
	Id                          int         `json:"id"`
	UserId                      int         `json:"user_id"`
	AccountNumber               string      `json:"account_number"`
	Currency                    string      `json:"currency"`
	Status                      int         `json:"status"`
	Archived                    bool        `json:"archived"`
	JoinedAccounts              interface{} `json:"joined_accounts"`
	SupportsDeposits            bool        `json:"supports_deposits"`
	SupportsWithdrawals         bool        `json:"supports_withdrawals"`
	MaskedAccountNumber         string      `json:"masked_account_number"`
	Name                        string      `json:"name"`
	DefaultName                 string      `json:"default_name"`
	AccountName                 string      `json:"account_name"`
	AccountType                 string      `json:"account_type"`
	AccountTypeFriendly         string      `json:"account_type_friendly"`
	CustodianName               string      `json:"custodian_name"`
	Custodian                   string      `json:"custodian"`
	CustodianAlias              string      `json:"custodian_alias"`
	CustodianId                 int         `json:"custodian_id"`
	ClientName                  string      `json:"client_name"`
	PortfolioModelName          string      `json:"portfolio_model_name"`
	CreatedAt                   time.Time   `json:"created_at"`
	UpdatedAt                   time.Time   `json:"updated_at"`
	InceptionDate               string      `json:"inception_date"`
	EndingBalanceDate           string      `json:"ending_balance_date"`
	BeginningBalanceDate        string      `json:"beginning_balance_date"`
	NetContributions            float64     `json:"net_contributions"`
	TotalContributions          float64     `json:"total_contributions"`
	Fees                        float64     `json:"fees"`
	TotalWithdrawals            float64     `json:"total_withdrawals"`
	Withdrawals                 float64     `json:"withdrawals"`
	Deposits                    float64     `json:"deposits"`
	BeginningBalance            float64     `json:"beginning_balance"`
	StartingBalance             float64     `json:"starting_balance"`
	InitialValue                float64     `json:"initial_value"`
	Balance                     float64     `json:"balance"`
	EndingBalance               float64     `json:"ending_balance"`
	TotalValue                  float64     `json:"total_value"`
	MarketValue                 float64     `json:"market_value"`
	MarketValueEndingBalance    float64     `json:"market_value_ending_balance"`
	MarketValueBeginningBalance float64     `json:"market_value_beginning_balance"`
	CashValue                   float64     `json:"cash_value"`
	CashValueEndingBalance      float64     `json:"cash_value_ending_balance"`
	CashValueBeginningBalance   float64     `json:"cash_value_beginning_balance"`
	WithdrawalsValue            float64     `json:"withdrawals_value"`
	TradingBalance              float64     `json:"trading_balance"`
	Gains                       float64     `json:"gains"`
	RateOfReturn                float64     `json:"rate_of_return"`
	AnnualizedXirr              float64     `json:"annualized_xirr"`
}

type User struct {
	Id        int      `json:"id"`
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Roles     []string `json:"roles"`
}

type DashboardResponse struct {
	BeginningBalance  float64   `json:"beginning_balance"`
	EndingBalance     float64   `json:"ending_balance"`
	Balance           float64   `json:"balance"`
	EndingBalanceDate string    `json:"ending_balance_date"`
	RateOfReturn      float64   `json:"rate_of_return"`
	InceptionDate     string    `json:"inception_date"`
	Deposits          float64   `json:"deposits"`
	Withdrawals       float64   `json:"withdrawals"`
	NetContributions  float64   `json:"net_contributions"`
	Gains             float64   `json:"gains"`
	Fees              float64   `json:"fees"`
	Accounts          []Account `json:"accounts"`
	CanCreateHisa     struct {
		Message string `json:"message"`
	} `json:"can_create_hisa"`
	User          User          `json:"user"`
	TaskReminders []interface{} `json:"task_reminders"`
}

func (c *Client) Dashboard(request DashboardRequest) (*DashboardResponse, error) {
	req, err := http.NewRequest("GET", "/investment_accounts/dashboard", nil)
	if err != nil {
		return nil, err
	}

	values := req.URL.Query()
	for currency, _ := range request.Currencies {
		values.Set(toRequestParameter(currency), "true")
	}
	values.Set("start_date", request.StartDate.Format("2006-01-02"))
	values.Set("end_date", request.EndDate.Format("2006-01-02"))
	values.Set("selected_value", "all") // TODO what values are allowed here? what does this control?
	req.URL.RawQuery = values.Encode()

	var response = &DashboardResponse{}
	if err := c.doWithProcessing(req, response, []int{http.StatusOK}); err != nil {
		return nil, err
	}

	return response, nil
}
