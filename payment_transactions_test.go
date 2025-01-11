package authorizenet

import (
	"fmt"
	"testing"
	"time"
)

var previousAuth string
var previousCharged string

func TestChargeCard(t *testing.T) {
	nextMonthDate := time.Now().AddDate(0, 1, 0)
	expiration := fmt.Sprintf("%s/%s", nextMonthDate.Format("01"), nextMonthDate.Format("06"))

	newTransaction := NewTransaction{
		Amount: "15.90",
		CreditCard: CreditCard{
			CardNumber:     "4007000000027",
			ExpirationDate: expiration,
		},
	}

	res, err := newTransaction.Charge(client)
	if err != nil {
		t.Fail()
		return
	}
	if res.Approved() {
		previousCharged = res.TransactionID()
		t.Log("#", res.TransactionID(), "Transaction was CHARGED $", newTransaction.Amount, "\n")
		t.Log("AVS Result Code: ", res.AVS().avsResultCode+"\n")
		t.Log("AVS ACVV Result Code: ", res.AVS().cavvResultCode+"\n")
		t.Log("AVS CVV Result Code: ", res.AVS().cvvResultCode+"\n")
	} else {
		t.Log(res)
		t.Log("Response code: ", res.Response.ResponseCode)
		t.Log(res.ErrorMessage(), "\n")
		t.Fail()
	}
}

func TestAVSDeclinedChargeCard(t *testing.T) {
	nextMonthDate := time.Now().AddDate(0, 1, 0)
	expiration := fmt.Sprintf("%s/%s", nextMonthDate.Format("01"), nextMonthDate.Format("06"))

	newTransaction := NewTransaction{
		Amount: RandomNumber(5, 99) + ".75",
		CreditCard: CreditCard{
			CardNumber:     "5424000000000015",
			ExpirationDate: expiration,
		},
		BillTo: &BillTo{
			FirstName:   RandomString(7),
			LastName:    RandomString(9),
			Address:     "1111 white ct",
			City:        "los angeles",
			State:       "CA",
			Zip:         "46205", // Must use this ZIP code for AVS status "N"
			Country:     "USA",
			PhoneNumber: "8885555555",
		},
		ShipTo: &Address{
			FirstName: RandomString(7),
			LastName:  RandomString(9),
			Address:   "1111 white ct",
			City:      "los angeles",
			State:     "CA",
			Zip:       "46205", // Must use this ZIP code for AVS status "N"
			Country:   "USA",
		},
	}
	res, err := newTransaction.Charge(client)
	if err != nil {
		t.Fail()
		return
	}

	t.Log(res)

	if res.AVS().avsResultCode == "N" {
		t.Log("#", res.TransactionID(), "AVS Transaction was DECLINED due to AVS Code. $", newTransaction.Amount, "\n")
		t.Log("AVS Result Text: ", res.AVS().Text(), "\n")
		t.Log("AVS Result Code: ", res.AVS().avsResultCode, "\n")
		t.Log("AVS ACVV Result Code: ", res.AVS().cavvResultCode, "\n")
		t.Log("AVS CVV Result Code: ", res.AVS().cvvResultCode, "\n")
	} else {
		t.Log("AVS Status Code: ", res.AVS().avsResultCode, "\n")
		t.Log(res.ErrorMessage(), "\n")
		t.Fail()
	}
}

func TestAVSChargeCard(t *testing.T) {

	nextMonthDate := time.Now().AddDate(0, 1, 0)
	expiration := fmt.Sprintf("%s/%s", nextMonthDate.Format("01"), nextMonthDate.Format("06"))

	newTransaction := NewTransaction{
		Amount: RandomNumber(5, 99) + ".51",
		CreditCard: CreditCard{
			CardNumber:     "4012888818888",
			ExpirationDate: expiration,
		},
		BillTo: &BillTo{
			FirstName:   RandomString(7),
			LastName:    RandomString(9),
			Address:     "1111 green ct",
			City:        "los angeles",
			State:       "CA",
			Zip:         "46203",
			Country:     "USA",
			PhoneNumber: "8885555555",
		},
		ShipTo: &Address{
			FirstName: RandomString(7),
			LastName:  RandomString(9),
			Address:   "1111 green ct",
			City:      "los angeles",
			State:     "CA",
			Zip:       "46203",
			Country:   "USA",
		},
	}
	res, err := newTransaction.Charge(client)
	if err != nil {
		t.Fail()
		return
	}

	if res.Held() {
		t.Log("Transaction is being Held for Review", "\n")
	}

	if res.AVS().avsResultCode == "E" {
		t.Log("#", res.TransactionID(), "AVS Transaction was CHARGED is now on HOLD$", newTransaction.Amount, "\n")
		t.Log("AVS Result Text: ", res.AVS().Text(), "\n")
		t.Log("AVS Result Code: ", res.AVS().avsResultCode, "\n")
		t.Log("AVS ACVV Result Code: ", res.AVS().cavvResultCode, "\n")
		t.Log("AVS CVV Result Code: ", res.AVS().cvvResultCode, "\n")
	} else {
		t.Log("AVS Status Code: ", res.AVS().avsResultCode, "\n")
		t.Log(res.ErrorMessage(), "\n")
		t.Fail()
	}
}

func TestDeclinedChargeCard(t *testing.T) {
	nextMonthDate := time.Now().AddDate(0, 1, 0)
	expiration := fmt.Sprintf("%s/%s", nextMonthDate.Format("01"), nextMonthDate.Format("06"))

	newTransaction := NewTransaction{
		Amount: RandomNumber(5, 99) + ".90",
		CreditCard: CreditCard{
			CardNumber:     "4007000000027",
			ExpirationDate: expiration,
		},
		BillTo: &BillTo{
			FirstName:   "Declined",
			LastName:    "User",
			Address:     "1337 Yolo Ln.",
			City:        "Beverly Hills",
			State:       "CA",
			Country:     "USA",
			Zip:         "46282",
			PhoneNumber: "8885555555",
		},
	}
	res, err := newTransaction.Charge(client)
	if err != nil {
		t.Fail()
		return
	}

	if res.Approved() {
		t.Fail()
	} else {
		t.Log("#", res.TransactionID(), "Transaction was DECLINED!!!", "\n")
		t.Log(res.Message(), "\n")
		t.Log("AVS Result Text: ", res.AVS().Text(), "\n")
		t.Log("AVS Result Code: ", res.AVS().avsResultCode, "\n")
		t.Log("AVS ACVV Result Code: ", res.AVS().cavvResultCode, "\n")
		t.Log("AVS CVV Result Code: ", res.AVS().cvvResultCode, "\n")
	}
}

func TestAuthOnlyCard(t *testing.T) {
	nextMonthDate := time.Now().AddDate(0, 1, 0)
	expiration := fmt.Sprintf("%s/%s", nextMonthDate.Format("01"), nextMonthDate.Format("06"))

	newTransaction := NewTransaction{
		Amount: "100.00",
		CreditCard: CreditCard{
			CardNumber:     "4012888818888",
			ExpirationDate: expiration,
		},
	}
	res, err := newTransaction.AuthOnly(client)
	if err != nil {
		t.Fail()
		return
	}

	if res.Approved() {
		previousAuth = res.TransactionID()
		t.Log("#", res.TransactionID(), "Transaction was AUTHORIZED $", newTransaction.Amount, "\n")
	} else {
		t.Log(res.ErrorMessage(), "\n")
	}
}

func TestCaptureAuth(t *testing.T) {
	oldTransaction := PreviousTransaction{
		Amount: "49.99",
		RefId:  previousAuth,
	}
	res, err := oldTransaction.Capture(client)
	if err != nil {
		t.Fail()
		return
	}
	if res.Approved() {
		t.Log("#", res.TransactionID(), "Transaction was CAPTURED $", oldTransaction.Amount, "\n")
	} else {
		t.Log(res.ErrorMessage(), "\n")
	}
}

func TestChargeCardChannel(t *testing.T) {
	nextMonthDate := time.Now().AddDate(0, 1, 0)
	expiration := fmt.Sprintf("%s/%s", nextMonthDate.Format("01"), nextMonthDate.Format("06"))

	newTransaction := NewTransaction{
		Amount: "38.00",
		CreditCard: CreditCard{
			CardNumber:     "4012888818888",
			ExpirationDate: expiration,
		},
		AuthCode: "RANDOMAUTHCODE",
	}
	res, err := newTransaction.Charge(client)
	if err != nil {
		t.Fail()
		return
	}

	if res.Approved() {
		previousAuth = res.TransactionID()
		t.Log("#", res.TransactionID(), "Transaction was Charged Through Channel (AuthCode) $", newTransaction.Amount, "\n")
	} else {
		t.Log(res.ErrorMessage(), "\n")
	}
}

func TestRefundCard(t *testing.T) {

	// Note: In case transaction is less than 24 hours, it must be voided instead of refunded.
	// From https://support.authorize.net/knowledgebase/Knowledgearticle/?code=000001244:
	// Please note that a void is used to cancel funding of a transaction that is pending settlement, and a refund is used to send money back to a customer.
	// Refunds are not applicable to unsettled transactions because funding has not been transferred for these transactions.

	nextMonthDate := time.Now().AddDate(0, 1, 0)
	expiration := fmt.Sprintf("%s/%s", nextMonthDate.Format("01"), nextMonthDate.Format("06"))
	randomTransId := RandomNumber(1000, 999999999)

	newTransaction := NewTransaction{
		Amount: fmt.Sprintf("%s.%s", RandomNumber(5, 200), RandomNumber(10, 99)),
		CreditCard: CreditCard{
			CardNumber:     "4012888818888",
			ExpirationDate: expiration,
		},
		RefTransId: randomTransId,
	}

	chargeRes, err := newTransaction.Charge(client)
	if err != nil {
		t.Log("Failed to perform the charge.")
		t.Fail()
		return
	}

	if !chargeRes.Approved() {
		t.Log("Transaction was declined.")
		t.Fail()
		return
	}

	newTransaction.RefTransId = chargeRes.TransactionID()

	res, err := newTransaction.Refund(client)
	if err != nil {
		t.Fail()
		return
	}

	resErrors := res.Response.Errors

	if res.Approved() {
		t.Log("#", res.TransactionID(), "Transaction was REFUNDED $", newTransaction.Amount, "\n")
	} else if len(resErrors) > 0 && resErrors[0].ErrorCode == "54" {
		t.Log("Transaction has not settled yet to be refunded.")
	} else {
		t.Log(res.ErrorMessage(), "\n")
		t.Fail()
	}
}

func TestVoidCard(t *testing.T) {
	newTransaction := PreviousTransaction{
		RefId: previousCharged,
	}
	res, err := newTransaction.Void(client)
	if err != nil {
		t.Fail()
		return
	}
	if res.Approved() {
		t.Log("#", res.TransactionID(), "Transaction was VOIDED $", newTransaction.Amount, "\n")
	} else {
		t.Log(res.ErrorMessage(), "\n")
	}
}

func TestChargeCustomerProfile(t *testing.T) {

	oldProfileId := "1810921101"
	oldPaymentId := "1805617738"

	customer := Customer{
		ID:        oldProfileId,
		PaymentID: oldPaymentId,
	}

	newTransaction := NewTransaction{
		Amount: "35.00",
	}

	res, err := newTransaction.ChargeProfile(customer, client)
	if err != nil {
		t.Fail()
		return
	}

	if res.Approved() {
		t.Log("#", res.TransactionID(), "Customer was Charged $", newTransaction.Amount, "\n")
	} else {
		t.Log(res.ErrorMessage(), "\n")
	}
}
