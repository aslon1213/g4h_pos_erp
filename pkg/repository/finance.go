package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TransactionType string
type InitiatorType string

const (
	TransactionTypeCredit TransactionType = "credit" // credit means - income - when money is gained or received into an account
	TransactionTypeDebit  TransactionType = "debit"  // debit means - outcome - when money is lost, spent, or withdrawn from an account
)

type TransactionBase struct {
	Amount      float64         `json:"amount" bson:"amount"`
	Description string          `json:"description" bson:"description"`
	Type        TransactionType `json:"type" bson:"type"`
}

func NewTransactionBase(amount float64, description string, typeOfTransaction TransactionType) *TransactionBase {
	return &TransactionBase{
		Amount:      amount,
		Description: description,
		Type:        typeOfTransaction,
	}
}

const (
	InitiatorTypeSalary    InitiatorType = "salary"
	InitiatorTypeRent      InitiatorType = "rent"
	InitiatorTypeUtilities InitiatorType = "utilities"
	InitiatorTypeOther     InitiatorType = "other"
	InitiatorTypeSale      InitiatorType = "sale"
	InitiatorTypeSupplier  InitiatorType = "supplier"
)

type Transaction struct {
	TransactionBase
	Type      InitiatorType `json:"type" bson:"type"`
	ID        string        `json:"id" bson:"_id"`
	CreatedAt time.Time     `json:"date" bson:"created_at"`
	UpdatedAt time.Time     `json:"date" bson:"updated_at"`
}

func NewTransaction(transactionBase *TransactionBase, typeOfTransaction InitiatorType) *Transaction {
	return &Transaction{
		TransactionBase: *transactionBase,
		Type:            typeOfTransaction,
		ID:              uuid.New().String(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

type Finance struct {
	Transactions  []Transaction `json:"transactions" bson:"transactions"`
	Balance       float64       `json:"balance" bson:"balance"`
	TotalIncome   float64       `json:"total_income" bson:"total_income"`
	TotalExpenses float64       `json:"total_expenses" bson:"total_expenses"`
}

func (f Finance) CreateTransaction(transaction Transaction, mongoCollection *mongo.Collection) error {

	_, err := mongoCollection.InsertOne(context.Background(), transaction)
	if err != nil {
		return err
	}

	return nil
}
