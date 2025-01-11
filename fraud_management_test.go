package authorizenet

import (
	"fmt"
	"testing"
	"time"
)

func TestGetUnsettledTransactions(t *testing.T) {
	transactions, err := client.UnsettledBatchList()
	if err != nil {
		t.Fail()
		return
	}

	t.Log("Count Unsettled: ", transactions.Count())
	t.Log(transactions.List(client))
}

// NOTE FOR TESTS BELOW: Under AFDS, go to Amount Filter and change the filter action to "Authorize and hold for review."
// Lower limit: 0, Upper limit: 1000

func TestApproveTransaction(t *testing.T) {

	nextMonthDate := time.Now().AddDate(0, 1, 0)
	expiration := fmt.Sprintf("%s/%s", nextMonthDate.Format("01"), nextMonthDate.Format("06"))
	randomTransId := RandomNumber(1000, 999999999)

	newTransaction := NewTransaction{
		RefTransId: randomTransId,
		Amount:     "1200.00",
		CreditCard: CreditCard{
			CardNumber:     "4007000000027",
			ExpirationDate: expiration,
		},
	}

	resTrans, err := newTransaction.Charge(client)
	if err != nil {
		t.Log("Failed to perform charge.")
		t.Fail()
		return
	}

	if resTrans.Held() == false {
		t.Log("The charge failed to go in held status.")
		t.Fail()
		return
	}

	oldTransaction := PreviousTransaction{
		Amount: "1200.00",
		RefId:  resTrans.TransactionID(),
	}

	res, err := oldTransaction.Approve(client)
	if err != nil {
		t.Fail()
		return
	}

	if res.Approved() {
		t.Log(res.ErrorMessage())
	} else {
		t.Log(res.ErrorMessage())
		t.Fail()
	}
}

func TestDeclineTransaction2(t *testing.T) {

	nextMonthDate := time.Now().AddDate(0, 1, 0)
	expiration := fmt.Sprintf("%s/%s", nextMonthDate.Format("01"), nextMonthDate.Format("06"))
	randomTransId := RandomNumber(1000, 999999999)

	newTransaction := NewTransaction{
		RefTransId: randomTransId,
		Amount:     "1205.00",
		CreditCard: CreditCard{
			CardNumber:     "4007000000027",
			ExpirationDate: expiration,
		},
	}

	resTrans, err := newTransaction.Charge(client)
	if err != nil {
		t.Log("Failed to perform charge.")
		t.Fail()
		return
	}

	if resTrans.Held() == false {
		t.Log("The charge failed to go in held status.")
		t.Fail()
		return
	}

	oldTransaction := PreviousTransaction{
		Amount: "1200.00",
		RefId:  resTrans.TransactionID(),
	}

	res, err := oldTransaction.Decline(client)
	if err != nil {
		t.Fail()
		return
	}

	if res.Approved() {
		t.Log(res.ErrorMessage())
	} else {
		t.Log(res.ErrorMessage())
	}
}
