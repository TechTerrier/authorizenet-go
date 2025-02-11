package authorizenet

import (
	"encoding/json"
	"time"
)

type Range struct {
	Start   time.Time `json:"start,omitempty"`
	End     time.Time `json:"end,omitempty"`
	BatchId string    `json:"batchId,omitempty"`
	Sorting Sorting   `json:"sorting,omitempty"`
	Paging  Paging    `json:"paging,omitempty"`
}

func (r BatchListResponse) List() []BatchList {
	return r.BatchList
}

func (r Range) SettledBatch(c Client) (*BatchListResponse, error) {
	newRequest := GetSettledBatchListRequest{
		GetSettledBatchList: GetSettledBatchList{
			MerchantAuthentication: c.GetAuthentication(),
			IncludeStatistics:      "true",
			FirstSettlementDate:    r.Start,
			LastSettlementDate:     r.End,
		},
	}
	req, err := json.Marshal(newRequest)
	if err != nil {
		return nil, err
	}
	res, err := c.SendRequest(req)
	var dat BatchListResponse
	json.Unmarshal(res, &dat)
	return &dat, err
}

func (c Client) UnSettledBatch() (*UnsettledTransactionListResponse, error) {
	newRequest := GetUnsettledBatchTransactionListRequest{
		GetUnsettledTransactionList: GetUnsettledTransactionList{
			MerchantAuthentication: c.GetAuthentication(),
		},
	}
	req, err := json.Marshal(newRequest)
	if err != nil {
		return nil, err
	}
	res, err := c.SendRequest(req)
	var dat UnsettledTransactionListResponse
	err = json.Unmarshal(res, &dat)
	return &dat, err
}

func (r UnsettledTransactionListResponse) List() []Transaction {
	return r.Transactions
}

func (r *GetTransactionListResponse) List() []Transaction {
	return r.GetTransactionList.Transactions.Transaction
}

func (r *GetTransactionListResponse) Count() int {
	return r.GetTransactionList.TotalNumInResultSet
}

func (r Range) Transactions(c Client) (*GetTransactionListResponse, error) {
	newRequest := GetTransactionListRequest{
		GetTransactionList: GetTransactionList{
			MerchantAuthentication: c.GetAuthentication(),
			BatchID:                r.BatchId,
		},
	}
	req, err := json.Marshal(newRequest)
	if err != nil {
		return nil, err
	}
	res, err := c.SendRequest(req)
	var dat GetTransactionListResponse
	json.Unmarshal(res, &dat)
	return &dat, err
}

func (r Range) Statistics(c Client) (*Statistics, error) {
	newRequest := GetBatchStatisticsRequest{
		GetBatchStatistics: GetBatchStatistics{
			MerchantAuthentication: c.GetAuthentication(),
			BatchID:                r.BatchId,
		},
	}
	req, err := json.Marshal(newRequest)
	if err != nil {
		return nil, err
	}
	res, err := c.SendRequest(req)
	var dat BatchStatisticsResponse
	err = json.Unmarshal(res, &dat)
	return &dat.Batch.Statistics[0], err
}

func (c Client) GetMerchantDetails() (*MerchantDetailsResponse, error) {
	newRequest := GetMerchantDetailsRequest{
		GetMerchantDetailsReq: GetMerchantDetailsReq{
			MerchantAuthentication: c.GetAuthentication(),
		},
	}
	req, err := json.Marshal(newRequest)
	if err != nil {
		return nil, err
	}
	res, err := c.SendRequest(req)
	var dat MerchantDetailsResponse
	err = json.Unmarshal(res, &dat)
	return &dat, err
}

func (tranx PreviousTransaction) Info(c Client) (*FullTransaction, error) {
	newRequest := GetTransactionDetailsRequest{
		GetTransactionDetails: GetTransactionDetails{
			MerchantAuthentication: c.GetAuthentication(),
			TransID:                tranx.RefId,
		},
	}
	req, err := json.Marshal(newRequest)
	if err != nil {
		return nil, err
	}
	res, err := c.SendRequest(req)
	var dat TransactionDetailsResponse
	err = json.Unmarshal(res, &dat)
	return &dat.Transaction, err
}

type GetSettledBatchListRequest struct {
	GetSettledBatchList GetSettledBatchList `json:"getSettledBatchListRequest"`
}

type GetSettledBatchList struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	IncludeStatistics      string                 `json:"includeStatistics"`
	FirstSettlementDate    time.Time              `json:"firstSettlementDate"`
	LastSettlementDate     time.Time              `json:"lastSettlementDate"`
}

type BatchListResponse struct {
	MessagesResponse
	BatchList []BatchList `json:"batchList,omitempty"`
}

type BatchList struct {
	BatchID                      string    `json:"batchId"`
	SettlementTimeUTC            time.Time `json:"settlementTimeUTC"`
	SettlementTimeUTCSpecified   bool      `json:"settlementTimeUTCSpecified"`
	SettlementTimeLocal          string    `json:"settlementTimeLocal"`
	SettlementTimeLocalSpecified bool      `json:"settlementTimeLocalSpecified"`
	SettlementState              string    `json:"settlementState"`
	PaymentMethod                string    `json:"paymentMethod"`
}

type GetTransactionListRequest struct {
	GetTransactionList GetTransactionList `json:"getTransactionListRequest"`
}

type GetTransactionList struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	BatchID                string                 `json:"batchId,omitempty"`
	Sorting                *Sorting               `json:"sorting,omitempty"`
	Paging                 *Paging                `json:"paging,omitempty"`
}

type GetTransactionListResponse struct {
	GetTransactionList struct {
		MessagesResponse
		Transactions        Transactions `json:"transactions"`
		TotalNumInResultSet int          `json:"totalNumInResultSet"`
	} `json:"getTransactionListResponse"`
}

type Transactions struct {
	Transaction []Transaction `json:"transaction"`
}

type Transaction struct {
	TransID           string  `json:"transId"`
	SubmitTimeUTC     string  `json:"submitTimeUTC"`
	SubmitTimeLocal   string  `json:"submitTimeLocal"`
	TransactionStatus string  `json:"transactionStatus"`
	Invoice           string  `json:"invoice,omitempty"`
	FirstName         string  `json:"firstName,omitempty"`
	LastName          string  `json:"lastName,omitempty"`
	Amount            string  `json:"amount,omitempty"`
	AccountType       string  `json:"accountType,omitempty"`
	AccountNumber     string  `json:"accountNumber,omitempty"`
	SettleAmount      float64 `json:"settleAmount,omitempty"`
	Subscription      struct {
		ID     int `json:"id"`
		PayNum int `json:"payNum,omitempty"`
	} `json:"subscription,omitempty"`
	MarketType     string `json:"marketType,omitempty"`
	Product        string `json:"product,omitempty"`
	MobileDeviceID string `json:"mobileDeviceId,omitempty"`
}

type GetTransactionDetailsRequest struct {
	GetTransactionDetails GetTransactionDetails `json:"getTransactionDetailsRequest"`
}

type GetTransactionDetails struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	TransID                string                 `json:"transId"`
}

type TransactionDetailsResponse struct {
	Transaction FullTransaction `json:"transaction"`
	MessagesResponse
}

type FullTransaction struct {
	TransID                   string    `json:"transId"`
	SubmitTimeUTC             time.Time `json:"submitTimeUTC"`
	SubmitTimeLocal           string    `json:"submitTimeLocal"`
	TransactionType           string    `json:"transactionType"`
	TransactionStatus         string    `json:"transactionStatus"`
	ResponseCode              int       `json:"resCode"`
	ResponseReasonCode        int       `json:"resReasonCode"`
	ResponseReasonDescription string    `json:"resReasonDescription"`
	AVSResponse               string    `json:"AVSResponse"`
	Batch                     struct {
		BatchID                      string    `json:"batchId"`
		SettlementTimeUTC            time.Time `json:"settlementTimeUTC"`
		SettlementTimeUTCSpecified   bool      `json:"settlementTimeUTCSpecified"`
		SettlementTimeLocal          string    `json:"settlementTimeLocal"`
		SettlementTimeLocalSpecified bool      `json:"settlementTimeLocalSpecified"`
		SettlementState              string    `json:"settlementState"`
	} `json:"batch"`
	Order struct {
		InvoiceNumber string `json:"invoiceNumber"`
	} `json:"order"`
	RequestedAmountSpecified         bool    `json:"requestedAmountSpecified"`
	AuthAmount                       float64 `json:"authAmount"`
	SettleAmount                     float64 `json:"settleAmount"`
	PrepaidBalanceRemainingSpecified bool    `json:"prepaidBalanceRemainingSpecified"`
	TaxExempt                        bool    `json:"taxExempt"`
	TaxExemptSpecified               bool    `json:"taxExemptSpecified"`
	Payment                          struct {
		BankAccount struct {
			AccountType          int         `json:"accountType"`
			AccountTypeSpecified bool        `json:"accountTypeSpecified"`
			RoutingNumber        string      `json:"routingNumber"`
			AccountNumber        string      `json:"accountNumber"`
			NameOnAccount        string      `json:"nameOnAccount"`
			EcheckType           int         `json:"echeckType"`
			EcheckTypeSpecified  bool        `json:"echeckTypeSpecified"`
			BankName             interface{} `json:"bankName"`
		} `json:"bankAccount"`
	} `json:"payment"`
	RecurringBilling          bool `json:"recurringBilling"`
	RecurringBillingSpecified bool `json:"recurringBillingSpecified"`
	ReturnedItems             []struct {
		ID          string    `json:"id"`
		DateUTC     time.Time `json:"dateUTC"`
		DateLocal   string    `json:"dateLocal"`
		Code        string    `json:"code"`
		Description string    `json:"description"`
	} `json:"returnedItems"`
}

type GetUnsettledBatchTransactionListRequest struct {
	GetUnsettledTransactionList GetUnsettledTransactionList `json:"getUnsettledTransactionListRequest"`
}

type GetUnsettledTransactionList struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	Sorting                *Sorting               `json:"sorting,omitempty"`
	Paging                 *Paging                `json:"paging,omitempty"`
}

type UnsettledTransactionListResponse struct {
	Transactions []Transaction `json:"transactions"`
	MessagesResponse
}

type GetBatchStatisticsRequest struct {
	GetBatchStatistics GetBatchStatistics `json:"getBatchStatisticsRequest"`
}

type GetBatchStatistics struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	BatchID                string                 `json:"batchId"`
}

type BatchStatisticsResponse struct {
	Batch struct {
		BatchID                      string       `json:"batchId"`
		SettlementTimeUTC            time.Time    `json:"settlementTimeUTC"`
		SettlementTimeUTCSpecified   bool         `json:"settlementTimeUTCSpecified"`
		SettlementTimeLocal          string       `json:"settlementTimeLocal"`
		SettlementTimeLocalSpecified bool         `json:"settlementTimeLocalSpecified"`
		SettlementState              string       `json:"settlementState"`
		PaymentMethod                string       `json:"paymentMethod"`
		Statistics                   []Statistics `json:"statistics"`
	} `json:"batch"`
	MessagesResponse
}

type Statistics struct {
	AccountType                        string  `json:"accountType"`
	ChargeAmount                       float64 `json:"chargeAmount"`
	ChargeCount                        int     `json:"chargeCount"`
	RefundAmount                       float64 `json:"refundAmount"`
	RefundCount                        int     `json:"refundCount"`
	VoidCount                          int     `json:"voidCount"`
	DeclineCount                       int     `json:"declineCount"`
	ErrorCount                         int     `json:"errorCount"`
	ReturnedItemAmount                 int     `json:"returnedItemAmount"`
	ReturnedItemAmountSpecified        bool    `json:"returnedItemAmountSpecified"`
	ReturnedItemCount                  int     `json:"returnedItemCount"`
	ReturnedItemCountSpecified         bool    `json:"returnedItemCountSpecified"`
	ChargebackAmount                   int     `json:"chargebackAmount"`
	ChargebackAmountSpecified          bool    `json:"chargebackAmountSpecified"`
	ChargebackCount                    int     `json:"chargebackCount"`
	ChargebackCountSpecified           bool    `json:"chargebackCountSpecified"`
	CorrectionNoticeCount              int     `json:"correctionNoticeCount"`
	CorrectionNoticeCountSpecified     bool    `json:"correctionNoticeCountSpecified"`
	ChargeChargeBackAmount             int     `json:"chargeChargeBackAmount"`
	ChargeChargeBackAmountSpecified    bool    `json:"chargeChargeBackAmountSpecified"`
	ChargeChargeBackCount              int     `json:"chargeChargeBackCount"`
	ChargeChargeBackCountSpecified     bool    `json:"chargeChargeBackCountSpecified"`
	RefundChargeBackAmount             int     `json:"refundChargeBackAmount"`
	RefundChargeBackAmountSpecified    bool    `json:"refundChargeBackAmountSpecified"`
	RefundChargeBackCount              int     `json:"refundChargeBackCount"`
	RefundChargeBackCountSpecified     bool    `json:"refundChargeBackCountSpecified"`
	ChargeReturnedItemsAmount          float64 `json:"chargeReturnedItemsAmount"`
	ChargeReturnedItemsAmountSpecified bool    `json:"chargeReturnedItemsAmountSpecified"`
	ChargeReturnedItemsCount           int     `json:"chargeReturnedItemsCount"`
	ChargeReturnedItemsCountSpecified  bool    `json:"chargeReturnedItemsCountSpecified"`
	RefundReturnedItemsAmount          int     `json:"refundReturnedItemsAmount"`
	RefundReturnedItemsAmountSpecified bool    `json:"refundReturnedItemsAmountSpecified"`
	RefundReturnedItemsCount           int     `json:"refundReturnedItemsCount"`
	RefundReturnedItemsCountSpecified  bool    `json:"refundReturnedItemsCountSpecified"`
}

type GetMerchantDetailsRequest struct {
	GetMerchantDetailsReq GetMerchantDetailsReq `json:"getMerchantDetailsRequest"`
}

type GetMerchantDetailsReq struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
}

type MerchantDetailsResponse struct {
	IsTestMode bool `json:"isTestMode"`
	Processors []struct {
		Name string `json:"name"`
	} `json:"processors"`
	MerchantName   string   `json:"merchantName"`
	GatewayID      string   `json:"gatewayId"`
	MarketTypes    []string `json:"marketTypes"`
	ProductCodes   []string `json:"productCodes"`
	PaymentMethods []string `json:"paymentMethods"`
	Currencies     []string `json:"currencies"`
	MessagesResponse
}
