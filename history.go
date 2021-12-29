package cidirectinvesting

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type HistoryEntry struct {
	Date     string          `json:"date"`
	Total    float64         `json:"total"`
	Accounts map[int]float64 `json:"-"`
}

func (e *HistoryEntry) UnmarshalJSON(data []byte) error {
	var vs map[string]interface{}
	if err := json.Unmarshal(data, &vs); err != nil {
		return err
	}

	e.Accounts = map[int]float64{}
	for k, v := range vs {
		switch k {
		case "date":
			e.Date = v.(string)
		case "total":
			e.Total = v.(float64)
		default:
			accountID, err := strconv.Atoi(k)
			if err != nil {
				return err
			}
			e.Accounts[accountID] = v.(float64)
		}
	}

	return nil
}

type HistoryRequest struct {
	StartDate time.Time
	EndDate   time.Time
	AccountID *int
}

type HistoryResponse []HistoryEntry

func (c *Client) History(request HistoryRequest) (*HistoryResponse, error) {
	req, err := http.NewRequest("GET", "/investment_accounts/history", nil)
	if err != nil {
		return nil, err
	}

	values := req.URL.Query()
	if request.AccountID != nil {
		values.Set("id", strconv.Itoa(*request.AccountID))
	}
	values.Set("start_date", request.StartDate.Format(time.RFC3339))
	values.Set("end_date", request.EndDate.Format(time.RFC3339))
	values.Set("basic", "true")         // TODO what values are allowed here? what does this control?
	values.Set("selected_value", "all") // TODO what values are allowed here? what does this control?
	req.URL.RawQuery = values.Encode()

	// TODO decode multi-account responses

	var response = &HistoryResponse{}
	if err := c.doWithProcessing(req, response, []int{http.StatusOK}); err != nil {
		return nil, err
	}

	return response, nil
}
