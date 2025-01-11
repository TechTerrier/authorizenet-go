package authorizenet

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestCancelSubscription(t *testing.T) {

	sub := SetSubscription{
		Id: newSubscriptionId,
	}

	subscriptionInfo, err := sub.Cancel(client)
	if err != nil {
		t.Fail()
		return
	}

	if subscriptionInfo.Ok() {
		t.Log("Subscription ID has been canceled: ", sub.Id, "\n")
		t.Log(subscriptionInfo.ErrorMessage(), "\n")
	} else {
		t.Log(subscriptionInfo.ErrorMessage())
		t.Fail()
	}

}

func TestGetInactiveSubscriptionList(t *testing.T) {
	subscriptionList, err := client.SubscriptionList("subscriptionInactive")
	if err != nil {
		t.Fail()
		return
	}

	count := subscriptionList.Count()
	t.Log("Amount of Inactive Subscriptions: ", count)

	if count == 0 {
		t.Fail()
	}
}

func TestCancelSecondSubscription(t *testing.T) {

	sub := SetSubscription{
		Id: newSecondSubscriptionId,
	}

	subscriptionInfo, err := sub.Cancel(client)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if subscriptionInfo.Ok() {
		t.Log("Second Subscription ID has been canceled: ", sub.Id, "\n")
		t.Log(subscriptionInfo.ErrorMessage(), "\n")
	} else {
		t.Log(subscriptionInfo.ErrorMessage())
		t.Fail()
	}

}

func TestDeleteCustomerShippingProfile(t *testing.T) {

	// Allow enough time for all the profiles to propagate.
	time.Sleep(60 * time.Second)

	customer := Customer{
		ID:         newCustomerProfileId,
		ShippingID: newCustomerShippingId,
	}

	res, err := customer.DeleteShippingProfile(client)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if res.Ok() {
		t.Log("Shipping Profile was Deleted")
	} else {
		t.Log(res.ErrorMessage())
		t.Fail()
	}
}

func TestDeleteCustomerPaymentProfile(t *testing.T) {
	customer := Customer{
		ID:        newCustomerProfileId,
		PaymentID: newCustomerPaymentId,
	}

	res, err := customer.DeletePaymentProfile(client)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if res.Ok() {
		t.Log("Payment Profile was Deleted")
	} else {
		t.Log(res.ErrorMessage())
		t.Fail()
	}
}

func TestDeleteCustomerProfile(t *testing.T) {

	customer := Customer{
		ID: newCustomerProfileId,
	}

	t.Log("TestDeleteCustomerProfile - ID", newCustomerProfileId)

	res, err := customer.DeleteProfile(client)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if res.Ok() {
		t.Log("Customer was Deleted")
	} else {
		t.Log(res.ErrorMessage())
		t.Fail()
	}

}

func TestDeleteSecondCustomerProfile(t *testing.T) {
	customer := Customer{
		ID: newSecondCustomerProfileId,
	}

	res, err := customer.DeleteProfile(client)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if res.Ok() {
		t.Log("Second Customer was Deleted")
	} else {
		t.Log(res.ErrorMessage())
		t.Fail()
	}

}

func TestDeclineTransaction(t *testing.T) {
	nextMonthDate := time.Now().AddDate(0, 1, 0)
	expiration := fmt.Sprintf("%s/%s", nextMonthDate.Format("01"), nextMonthDate.Format("06"))
	randomTransId := RandomNumber(1000, 999999999)

	newTransaction := NewTransaction{
		RefTransId: randomTransId,
		Amount:     "1201.00",
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
		t.Log(err)
		t.Fail()
	}

	if res.Approved() {
		t.Log("DECLINED the previous transasction that was on Hold. ID #", oldTransaction.RefId)
		t.Log(res.TransactionID())
	} else {
		t.Log(res.ErrorMessage())
		t.Fail()
	}
}

func RandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomNumber(min, max int) string {
	num := rand.Intn(max-min) + min
	return strconv.Itoa(num)
}
