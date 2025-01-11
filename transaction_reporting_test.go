package authorizenet

import (
	"testing"
)

func TestGetSettledBatchList(t *testing.T) {

	list := Range{
		Start: LastWeek(),
		End:   Now(),
	}

	batches, err := list.SettledBatch(client)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	batchList := batches.List()

	for _, v := range batchList {
		t.Log("Batch ID: ", v.BatchID, "\n")
		t.Log("Payment Method: ", v.PaymentMethod, "\n")
		t.Log("State: ", v.SettlementState, "\n")
	}

}

func TestGetTransactionList(t *testing.T) {

	list := Range{
		BatchId: "6933560",
	}

	batches, err := list.Transactions(client)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	batchList := batches.List()

	for _, v := range batchList {
		t.Log("Transaction ID: ", v.TransID, "\n")
		t.Log("Amount: ", v.Amount, "\n")
		t.Log("Account: ", v.AccountNumber, "\n")
	}

}

func TestGetTransactionDetails(t *testing.T) {

	newTransaction := PreviousTransaction{
		RefId: "60019493304",
	}
	res, err := newTransaction.Info(client)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log("Transaction Status: ", res.TransactionStatus, "\n")
}

func TestGetUnSettledBatchList(t *testing.T) {

	batches, err := client.UnSettledBatch()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	batchList := batches.List()

	for _, v := range batchList {
		t.Log("Status: ", v.TransactionStatus, "\n")
		t.Log("Amount: ", v.Amount, "\n")
		t.Log("Transaction ID: #", v.TransID, "\n")
	}

}

func TestGetBatchStatistics(t *testing.T) {

	// Prepare to get a batch ID.
	searchList := Range{
		Start: LastWeek(),
		End:   Now(),
	}
	batches, err := searchList.SettledBatch(client)
	if err != nil {
		t.Log("Failed to get batches.")
		t.Log(err)
		t.Fail()
	}

	batchList := batches.List()

	if len(batchList) == 0 {
		t.Log("There are no batches to get statistics on.")
		return
	}

	list := Range{
		BatchId: batchList[0].BatchID, // Use the first batch ID provided.
	}

	batch, err := list.Statistics(client)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log("Refund Count: ", batch.RefundCount, "\n")
	t.Log("Charge Count: ", batch.ChargeCount, "\n")
	t.Log("Void Count: ", batch.VoidCount, "\n")
	t.Log("Charge Amount: ", batch.ChargeAmount, "\n")
	t.Log("Refund Amount: ", batch.RefundAmount, "\n")

}

func TestGetMerchantDetails(t *testing.T) {

	info, err := client.GetMerchantDetails()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log("Test Mode: ", info.IsTestMode, "\n")
	t.Log("Gateway ID: ", info.GatewayID, "\n")
}
