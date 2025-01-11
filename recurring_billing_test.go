package authorizenet

import (
	"fmt"
	"testing"
	"time"
)

var newSubscriptionId string
var newSecondSubscriptionId string

func TestCreateSubscription(t *testing.T) {

	nextMonthDate := time.Now().AddDate(0, 1, 0)
	expiration := fmt.Sprintf("%s/%s", nextMonthDate.Format("01"), nextMonthDate.Format("06"))

	subscription := Subscription{
		Name:   "New Subscription",
		Amount: RandomNumber(5, 99) + ".00",
		//TrialAmount: "0.00",
		PaymentSchedule: &PaymentSchedule{
			StartDate:        CurrentDate(),
			TotalOccurrences: "9999",
			//TrialOccurrences: "0",
			Interval: IntervalMonthly(),
		},
		Payment: &Payment{
			CreditCard: CreditCard{
				CardNumber:     "4007000000027",
				ExpirationDate: expiration,
			},
		},
		BillTo: &BillTo{
			FirstName: "Hunter",
			LastName:  "Long",
		},
	}

	res, err := subscription.Charge(client)
	if err != nil {
		t.Fail()
		return
	}

	if res.Approved() {
		newSecondSubscriptionId = res.SubscriptionID
		newSecondCustomerProfileId = res.CustomerProfileId()
		t.Log("New Subscription: ", res.SubscriptionID)
		t.Log("New Customer Profile ID: ", res.CustomerProfileId())
		t.Log("New Payment Profile ID: ", res.CustomerPaymentProfileId())
	} else {
		t.Log(res.ErrorMessage(), "\n")
	}

}

func TestGetSubscription(t *testing.T) {

	sub := SetSubscription{
		Id: newSubscriptionId,
	}

	subscriptionInfo, err := sub.Info(client)
	if err != nil {
		t.Fail()
		return
	}

	t.Log("Subscription Name: #", subscriptionInfo.Subscription.Name, "\n")
	t.Log("Subscription Status: ", subscriptionInfo.Subscription.Status, "\n")

}

func TestGetSubscriptionStatus(t *testing.T) {

	sub := SetSubscription{
		Id: newSubscriptionId,
	}

	subscriptionInfo, err := sub.Status(client)
	if err != nil {
		t.Fail()
		return
	}

	if subscriptionInfo.Active() {
		t.Log("Subscription ID", newSubscriptionId, " has status: ", subscriptionInfo.Status)
	} else {
		t.Log("Subscription ID", newSubscriptionId, "has status: ", subscriptionInfo.Status)
		t.Fail()
	}

}

func TestUpdateSubscription(t *testing.T) {
	nextMonthDate := time.Now().AddDate(0, 1, 0)
	expiration := fmt.Sprintf("%s/%s", nextMonthDate.Format("01"), nextMonthDate.Format("06"))

	subscription := Subscription{
		Payment: &Payment{
			CreditCard: CreditCard{
				CardNumber:     "5424000000000015",
				ExpirationDate: expiration,
			},
		},
		SubscriptionId: newSubscriptionId,
	}

	res, err := subscription.Update(client)
	if err != nil {
		t.Fail()
		return
	}

	if res.Approved() {
		t.Log("Updated Subscription \n")
	} else {
		t.Log(res.ErrorMessage(), "\n")
		t.Fail()
	}

}

func TestGetActiveSubscriptionList(t *testing.T) {

	subscriptionList, err := client.SubscriptionList("subscriptionActive")
	if err != nil {
		t.Fail()
		return
	}

	count := subscriptionList.Count()
	t.Log("Amount of Active Subscriptions: ", count)

	if count == 0 {
		t.Fail()
	}

}

func TestGetExpiringSubscriptionList(t *testing.T) {

	subscriptionList, err := client.SubscriptionList("subscriptionExpiringThisMonth")
	if err != nil {
		t.Fail()
		return
	}

	t.Log("Amount of Subscriptions Expiring This Month: ", subscriptionList.Count())

}

func TestGetCardExpiringSubscriptionList(t *testing.T) {

	subscriptionList, err := client.SubscriptionList("cardExpiringThisMonth")
	if err != nil {
		t.Fail()
		return
	}

	t.Log("Amount of Cards Expiring This Month: ", subscriptionList.Count())

}
