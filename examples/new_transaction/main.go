package main

import (
	"fmt"
	"github.com/techterrier/authorizenet-go"
	"os"
)

var newTransactionId string

func main() {
	apiName := os.Getenv("API_NAME")
	apiKey := os.Getenv("API_KEY")

	if apiName == "" || apiKey == "" {
		panic("API_KEY or API_NAME environment variable not set.")
	}
	
	client := *authorizenet.New(apiName, apiKey, true)
	ChargeCustomer(client)
	VoidTransaction(client)
}

func ChargeCustomer(client authorizenet.Client) {
	newTransaction := authorizenet.NewTransaction{
		Amount: "13.75",
		CreditCard: authorizenet.CreditCard{
			CardNumber:     "4012888818888",
			ExpirationDate: "08/25",
			CardCode:       "393",
		},
		BillTo: &authorizenet.BillTo{
			FirstName:   "Timmy",
			LastName:    "Jimmy",
			Address:     "1111 green ct",
			City:        "los angeles",
			State:       "CA",
			Zip:         "43534",
			Country:     "USA",
			PhoneNumber: "8885555555",
		},
	}
	res, err := newTransaction.Charge(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.Approved() {
		newTransactionId = res.TransactionID()
		fmt.Println("Transaction was Approved! #", res.TransactionID())
	}

}

func VoidTransaction(client authorizenet.Client) {

	newTransaction := authorizenet.PreviousTransaction{
		RefId: newTransactionId,
	}
	res, err := newTransaction.Void(client)

	if err != nil {
		fmt.Println(err)
		return
	}

	if res.Approved() {
		fmt.Println("Transaction was Voided!")
	}
}
