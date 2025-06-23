package test

import (
	"aslon1213/magazin_pos/pkg/app"
	"aslon1213/magazin_pos/pkg/configs"
	"aslon1213/magazin_pos/test/client"
	"context"
	"net/http"
	"testing"
	"time"

	models "aslon1213/magazin_pos/pkg/repository"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestMain(m *testing.M) {

	app := app.NewApp()
	app.DB.Database("magazin").Collection("finance").DeleteMany(context.Background(), bson.M{})
	app.DB.Database("magazin").Collection("suppliers").DeleteMany(context.Background(), bson.M{})
	app.DB.Database("magazin").Collection("transactions").DeleteMany(context.Background(), bson.M{})

	go app.Run()
	m.Run()
	// delete all records from magazin database
}

// TestCreateBranches tests the creation of branches in the finance module
func TestCreateBranches(t *testing.T) {
	configs, err := configs.LoadConfig("../../")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	// Test data for creating branches
	branches := []models.NewBranchFinanceInput{
		{
			BranchName: "Branch A",
			Details:    map[string]interface{}{"location": "Location A"},
		},
		{
			BranchName: "Branch B",
			Details:    map[string]interface{}{"location": "Location B"},
		},
		{
			BranchName: "Branch C",
			Details:    map[string]interface{}{"location": "Location C"},
		},
	}

	client := client.NewClient("localhost", configs.Server.Port, "admin", "admin")

	for _, branch := range branches {
		resp, output, err := client.CreateBranch(branch)
		assert.Equal(t, resp.StatusCode, http.StatusCreated, "Expected status code 201, but got %d", resp.StatusCode)
		if err != nil {
			t.Fatal(err)
		}
		assert.Nil(t, output.Error, "Expected no error, but got one")
	}
}

func TestCreateBranchDuplicate(t *testing.T) {
	configs, err := configs.LoadConfig("../../")
	if err != nil {
		t.Fatal(err)
	}

	client := client.NewClient("localhost", configs.Server.Port, "admin", "admin")

	branch := models.NewBranchFinanceInput{
		BranchName: "Branch A",
		Details:    map[string]interface{}{"location": "Location A"},
	}
	resp, output, err := client.CreateBranch(branch)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "Expected status code 400, but got %d", resp.StatusCode)
	assert.NotNil(t, output.Error, "Expected error, but got none")
}

func TestGetBranches(t *testing.T) {
	configs, err := configs.LoadConfig("../../")
	if err != nil {
		t.Fatal(err)
	}

	client := client.NewClient("localhost", configs.Server.Port, "admin", "admin")

	_, output, err := client.GetAllBranches()
	if err != nil {
		t.Fatal(err)
	}

	log.Info().Interface("output", output).Msg("Output")

	assert.Nil(t, output.Error, "Expected no error, but got one")
}

func TestGetBranchByID(t *testing.T) {
	configs, err := configs.LoadConfig("../../")
	if err != nil {
		t.Fatal(err)
	}

	client := client.NewClient("localhost", configs.Server.Port, "admin", "admin")
	// get all branches first
	resp, output, err := client.GetAllBranches()
	if err != nil {
		t.Fatal(err)
	}

	financeBranches := output.Data

	if len(financeBranches) == 0 {
		t.Fatal("No branches found")
	}

	firstBranch := financeBranches[0]
	branchID := firstBranch.BranchID

	resp, output_, err := client.GetBranchByID(branchID)
	if err != nil {
		log.Error().Err(err).Interface("Response", output_.Data).Msg("Failed to decode response")
		t.Fatal(err)
	}
	log.Info().Interface("resp", resp).Msg("Response")
	assert.Nil(t, output.Error, "Expected no error, but got one")
}

func TestGetBranchByName(t *testing.T) {
	configs, err := configs.LoadConfig("../../")
	if err != nil {
		t.Fatal(err)
	}

	client := client.NewClient("localhost", configs.Server.Port, "admin", "admin")
	resp, output, err := client.GetAllBranches()
	if err != nil {
		t.Fatal(err)
	}

	financeBranches := output.Data

	branchName := financeBranches[0].BranchName

	resp, output_, err := client.GetBranchByName(branchName)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output.Error, "Expected no error, but got one")
	assert.Equal(t, output_.Data.BranchName, branchName, "Expected branch name to be %s, but got %s", branchName, output_.Data.BranchName)

}

func TestCreateSupplier(t *testing.T) {
	configs, err := configs.LoadConfig("../../")
	if err != nil {
		t.Fatal(err)
	}

	client := client.NewClient("localhost", configs.Server.Port, "admin", "admin")
	newSupplier := models.SupplierBase{
		Name:    "New Supplier",
		Email:   "new@example.com",
		Phone:   "0987654321",
		Address: "New Address",
		INN:     "987654321",
		Notes:   "New Notes",
		Branch:  "Branch A",
	}

	resp, output, err := client.CreateSupplier(newSupplier)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201, but got %d", resp.StatusCode)
	assert.Nil(t, output.Error, "Expected no error, but got one")
	assert.Equal(t, output.Data.Name, newSupplier.Name, "Expected supplier name to be %s, but got %s", newSupplier.Name, output.Data.Name)
}

func TestGetSuppliers(t *testing.T) {
	configs, err := configs.LoadConfig("../../")
	if err != nil {
		t.Fatal(err)
	}

	client := client.NewClient("localhost", configs.Server.Port, "admin", "admin")
	resp, output, err := client.GetAllSuppliers()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output.Error, "Expected no error, but got one")
	assert.NotEmpty(t, output.Data, "Expected to find suppliers, but got none")
}

func TestGetSupplierByID(t *testing.T) {
	configs, err := configs.LoadConfig("../../")
	if err != nil {
		t.Fatal(err)
	}

	client := client.NewClient("localhost", configs.Server.Port, "admin", "admin")
	resp, output, err := client.GetAllSuppliers()
	if err != nil {
		t.Fatal(err)
	}

	suppliers := output.Data
	if len(suppliers) == 0 {
		t.Fatal("No suppliers found")
	}

	firstSupplier := suppliers[0]
	supplierID := firstSupplier.ID

	resp, output_, err := client.GetSupplierByID(supplierID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output_.Error, "Expected no error, but got one")
	assert.Equal(t, output_.Data.ID, supplierID, "Expected supplier ID to be %s, but got %s", supplierID, output_.Data.ID)
}

func TestNewSupplierTransaction(t *testing.T) {
	configs, err := configs.LoadConfig("../../")
	if err != nil {
		t.Fatal(err)
	}

	client := client.NewClient("localhost", configs.Server.Port, "admin", "admin")
	_, output, err := client.GetAllSuppliers()
	if err != nil {
		t.Fatal(err)
	}

	suppliers := output.Data
	supplierID := ""
	Branch := ""

	// choose oe supplier which has branch
	for _, supplier := range suppliers {
		if supplier.Branch != "" {
			supplierID = supplier.ID
			Branch = supplier.Branch
			break
		}
	}

	transactions := []models.TransactionBase{
		{
			Amount:        10000000,
			Description:   "Test Transaction 1",        // 10 million
			Type:          models.TransactionTypeDebit, // we get the money from the supplier
			PaymentMethod: models.PaymentMethodBank,    // incrase debt section
		},
		{
			Amount:        20000000,
			Description:   "Test Transaction 2",        // 20 million
			Type:          models.TransactionTypeDebit, // we get the money from the supplier
			PaymentMethod: models.PaymentMethodBank,    // increase the debt section
		},
		{
			Amount:        15000000,
			Description:   "Test Transaction 3",        // 15 million
			Type:          models.TransactionTypeDebit, // we get the money from the supplier
			PaymentMethod: models.PaymentMethodBank,    // increase the debt section
		},
		{
			Amount:        5000000,
			Description:   "Test Transaction 3",         // 5 million
			Type:          models.TransactionTypeCredit, // we give the money to the supplier
			PaymentMethod: models.PaymentMethodCash,     // decrease the balance.cash by 5 million
		},
		{
			Amount:        2500000,
			Description:   "Test Transaction 3",         // 2.5 million
			Type:          models.TransactionTypeCredit, // we get the money from the supplier
			PaymentMethod: models.PaymentMethodBank,     // decrease the balance.bank by 2.5 million
		},
		{
			Amount:        3000000,
			Description:   "Test Transaction 3",         // 3 million
			Type:          models.TransactionTypeCredit, // we get the money from the supplier
			PaymentMethod: models.PaymentMethodBank,     // decrease the balance.bank by 3 million
		},
		{
			Amount:        4000000,
			Description:   "Test Transaction 3",         // 4 million
			Type:          models.TransactionTypeCredit, // we get the money from the supplier
			PaymentMethod: models.PaymentMethodCash,     // decrease the balance.cash by 4 million
		},
		{
			Amount:        1000000,
			Description:   "Test Transaction 3",         // 1 million
			Type:          models.TransactionTypeCredit, // we get the money from the supplier
			PaymentMethod: models.OnlineTransfer,        // decrease the balance.mobile_apps by 1 million
		},
	}
	total_income := 15500000
	total_expenses := 45000000
	balance := total_income - total_expenses

	for _, transaction := range transactions {
		resp, output, err := client.NewSupplierTransaction(Branch, supplierID, transaction)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201, but got %d", resp.StatusCode)
		assert.Nil(t, output.Error, "Expected no error, but got one")
	}

	// fetch finance and supplier from database and check the balance and total income and total expenses

	// get finance by branch name
	resp, branch_response, err := client.GetBranchByID(Branch)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, branch_response.Error, "Expected no error, but got one")

	// check the balance and total income and total expenses
	assert.Equal(t, int32(-balance), branch_response.Data.Finance.Debt, "Expected balance to be %f, but got %f", -balance, branch_response.Data.Finance.Debt)
	assert.Equal(t, int32(-5500000), branch_response.Data.Finance.Balance.Bank, "Expected balance to be %f, but got %f", 5500000, branch_response.Data.Finance.Balance.Bank)
	assert.Equal(t, int32(-9000000), branch_response.Data.Finance.Balance.Cash, "Expected balance to be %f, but got %f", 10000000, branch_response.Data.Finance.Balance.Cash)
	assert.Equal(t, int32(-1000000), branch_response.Data.Finance.Balance.MobileApps, "Expected balance to be %f, but got %f", 1000000, branch_response.Data.Finance.Balance.MobileApps)
	// assert.Equal(t, branch_response.Data.Finance.TotalIncome, total_income, "Expected total income to be %f, but got %f", total_income, branch_response.Data.Finance.TotalIncome)
	// assert.Equal(t, branch_response.Data.Finance.TotalExpenses, total_expenses, "Expected total expenses to be %f, but got %f", total_expenses, branch_response.Data.Finance.TotalExpenses)

	// get supplier by id
	resp, supplier_response, err := client.GetSupplierByID(supplierID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, supplier_response.Error, "Expected no error, but got one")
	assert.Equal(t, supplier_response.Data.ID, supplierID, "Expected supplier ID to be %s, but got %s", supplierID, supplier_response.Data.ID)
	assert.NotEqual(t, supplier_response.Data.FinancialData.Balance, 0, "Expected balance to be %f, but got %f", balance, supplier_response.Data.FinancialData.Balance)

}

// test get transactions, delete transactions of sales, create sales transaction
func TestSalesTransaction(t *testing.T) {
	configs, err := configs.LoadConfig("../../")
	if err != nil {
		t.Fatal(err)
	}

	client := client.NewClient("localhost", configs.Server.Port, "admin", "admin")

	// get all branches
	resp, output, err := client.GetAllBranches()
	if err != nil {
		t.Fatal(err)
	}

	branches := output.Data

	branchID := ""

	for _, branch := range branches {
		if branch.BranchName == "Branch A" {
			branchID = branch.BranchID
			break
		}
	}

	// create sales transaction
	transactions := []models.TransactionBase{{
		Amount:        10000000, // 10 million
		Description:   "Test Transaction 1",
		Type:          models.TransactionTypeDebit, // transaction are always positive to balance
		PaymentMethod: models.PaymentMethodBank,    // increase balance.bank by 10 million
	}, {
		Amount:        20000000, // 20 million
		Description:   "Test Transaction 2",
		Type:          models.TransactionTypeDebit,
		PaymentMethod: models.PaymentMethodBank, // increase balance.bank by 20 million
	}, {
		Amount:        30000000, // 30 million
		Description:   "Test Transaction 3",
		Type:          models.TransactionTypeCredit,
		PaymentMethod: models.PaymentMethodBank, // increase balance.bank by 30 million
	}, {
		Amount:        40000000, // 40 million
		Description:   "Test Transaction 4",
		Type:          models.TransactionTypeDebit,
		PaymentMethod: models.PaymentMethodCash, // increase balance.cash by 40 million
	}, {
		Amount:        50000000, // 50 million
		Description:   "Test Transaction 5",
		Type:          models.TransactionTypeCredit,
		PaymentMethod: models.PaymentMethodCash, // increase balance.cash by 50 million
	}, {
		Amount:        60000000, // 60 million
		Description:   "Test Transaction 6",
		Type:          models.TransactionTypeCredit,
		PaymentMethod: models.PaymentMethodCash, // increase balance.cash by 60 million
	}, {
		Amount:        500000, // 500 thousand
		Description:   "Test Transaction 7",
		Type:          models.TransactionTypeCredit,
		PaymentMethod: models.OnlineMobileAppPayment, // increase balance.mobile_apps by 500 thousand
	}, {
		Amount:        205000, // 205 thousand
		Description:   "Test Transaction 8",
		Type:          models.TransactionTypeCredit,
		PaymentMethod: models.OnlineMobileAppPayment, // increase balance.mobile_apps by 205 thousand
	}, {
		Amount:        1000000, // 1 million
		Description:   "Test Transaction 9",
		Type:          models.TransactionTypeCredit,
		PaymentMethod: models.OnlineTransfer, // increase balance.mobile_apps by 1 million
	}, {
		Amount:        1000000, // 1 million
		Description:   "Test Transaction 10",
		Type:          models.TransactionTypeCredit,
		PaymentMethod: models.OnlineTransfer, // increase balance.mobile_apps by 1 million
	}, {
		Amount:        1000000, // 1 million
		Description:   "Test Transaction 11",
		Type:          models.TransactionTypeCredit,
		PaymentMethod: models.PaymentMethodTerminal, // increase balance.terminal by 1 million
	}, {
		Amount:        1000000, // 1 million
		Description:   "Test Transaction 12",
		Type:          models.TransactionTypeCredit,
		PaymentMethod: models.PaymentMethodTerminal, // increase balance.terminal by 1 million
	},
	}

	for _, transaction := range transactions {
		resp, output, err := client.CreateSalesTransaction(branchID, transaction)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201, but got %d", resp.StatusCode)
		assert.Equal(t, output.Data.Amount, transaction.Amount, "Expected amount to be %f, but got %f", transaction.Amount, output.Data.Amount)
		assert.Nil(t, output.Error, "Expected no error, but got one")
	}
	// get transactions
	// get branch finance info
	cash := -9000000
	bank := -5500000
	terminal := 0
	mobile_apps := -1000000
	for _, transaction := range transactions {
		if transaction.PaymentMethod == models.PaymentMethodCash {
			cash += int(transaction.Amount)
		}
		if transaction.PaymentMethod == models.PaymentMethodBank {
			bank += int(transaction.Amount)
		}
		if transaction.PaymentMethod == models.PaymentMethodTerminal {
			terminal += int(transaction.Amount)
		}
		if transaction.PaymentMethod == models.OnlineMobileAppPayment {
			mobile_apps += int(transaction.Amount)
		}
		if transaction.PaymentMethod == models.OnlineTransfer {
			mobile_apps += int(transaction.Amount)
		}
	}
	balance := models.Balance{
		Cash:       int32(cash),
		Bank:       int32(bank),
		Terminal:   int32(terminal),
		MobileApps: int32(mobile_apps),
	}
	resp, output_branch, err := client.GetBranchByID(branchID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output_branch.Error, "Expected no error, but got one")
	assert.Equal(t, output_branch.Data.Finance.Balance, balance, "Expected balance to be %f, but got %f", balance, output_branch.Data.Finance.Balance)

	// get transactions
	resp, output_transactions, err := client.GetTransactions(branchID, "", 0, 0, "", "", models.InitiatorTypeSales, time.Time{}, time.Time{}, 1, 100)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output_transactions.Error, "Expected no error, but got one")
	assert.Equal(t, len(output_transactions.Data), 12, "Expected 12 transactions, but got %d", len(output_transactions.Data))
}
