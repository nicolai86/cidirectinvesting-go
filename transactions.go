package cidirectinvesting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Transaction struct {
	Id                   int         `json:"id"`
	ProcessDate          string      `json:"process_date"`
	AccountIdentifier    string      `json:"account_identifier"`
	TransactionId        string      `json:"transaction_id"`
	SubTransactionId     string      `json:"sub_transaction_id"`
	Commission           string      `json:"commission"`
	SecurityType         string      `json:"security_type"`
	ExchangeRate         string      `json:"exchange_rate"`
	CreatedAt            string      `json:"created_at"`
	UpdatedAt            string      `json:"updated_at"`
	TransferId           interface{} `json:"transfer_id"`
	EffectiveDate        string      `json:"effective_date"`
	TransactionDesc      string      `json:"transaction_desc"`
	TransactionType      string      `json:"transaction_type"`
	HumanTransactionType string      `json:"human_transaction_type"`
	Quantity             string      `json:"quantity"`
	Price                string      `json:"price"`
	PriceCad             string      `json:"price_cad"`
	Value                string      `json:"value"`
	ValueCad             string      `json:"value_cad"`
	HoldingListed        bool        `json:"holding_listed"`
	HasMatchedTransfer   bool        `json:"has_matched_transfer"`
	CurrencyCode         string      `json:"currency_code"`
	Symbol               *string     `json:"symbol"`
	Cusip                *string     `json:"cusip"`
	FullCount            int         `json:"full_count"`
	FundsUrl             string      `json:"funds_url,omitempty"`
}

type TransactionResponse struct {
	Metadata struct {
		TransactionTypes []struct {
			Category *string `json:"category"`
		} `json:"transaction_types"`
		User struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
		InvestmentAccount struct {
			Id            int    `json:"id"`
			Name          string `json:"name"`
			InceptionDate string `json:"inception_date"`
		} `json:"investment_account"`
	} `json:"metadata"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Pagination struct {
		Next               int         `json:"next"`
		PageAfterNext      interface{} `json:"page_after_next"`
		Previous           interface{} `json:"previous"`
		PageBeforePrevious interface{} `json:"page_before_previous"`
	} `json:"pagination"`
	Data []Transaction `json:"data"`
}

type Filter struct {
	Quantity     *Range   `json:"quantity"`
	Price        *Range   `json:"price"`
	Value        *Range   `json:"value"`
	Symbol       *string  `json:"symbol,omitempty"` // ie VSB
	Description  *string  `json:"description,omitempty"`
	Transactions []string `json:"transaction,omitempty"` // Adjustment, Buy, Cancel Buy, Cancel Sell, Cancel transfer, Deposit, Dividend, Fees, Interest, Other, Reorganization, Sell, Tax, Transfer, Withdrawal, Withholding
}

type Range struct {
	From int `json:"from,omitempty"`
	To   int `json:"to,omitempty"`
}

type TransactionRequest struct {
	AccountID int
	StartDate time.Time
	EndDate   time.Time
	Page      int
	Limit     int
	Filter    *Filter
	OrderBy   *string
	Reverse   bool
}

func (c *Client) Transactions(request TransactionRequest) (*TransactionResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/investment_accounts/%d/transactions", request.AccountID), nil)
	if err != nil {
		return nil, err
	}

	values := req.URL.Query()
	if request.OrderBy != nil {
		values.Set("order_by", *request.OrderBy)
	} else {
		values.Set("order_by", "description")
	}

	values.Set("reverse", fmt.Sprintf("%t", request.Reverse))
	values.Set("start_date", request.StartDate.Format("2006-01-02"))
	values.Set("end_date", request.EndDate.Format("2006-01-02"))
	values.Set("limit", strconv.Itoa(request.Limit))
	values.Set("page", strconv.Itoa(request.Page))
	values.Set("selected", "custom")
	if request.Filter != nil {
		buf := &bytes.Buffer{}
		json.NewEncoder(buf).Encode(request.Filter)
		filters := buf.String()
		values.Set("filters", filters[0:len(filters)-1]) // skip \n
	}
	req.URL.RawQuery = values.Encode()

	var response = &TransactionResponse{}
	if err := c.doWithProcessing(req, response, []int{http.StatusOK}); err != nil {
		return nil, err
	}

	return response, nil
}
