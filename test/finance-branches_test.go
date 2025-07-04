package test

import (
	"context"
	"io"
	logging "log"
	"net/http"
	"testing"
	"time"

	"github.com/aslon1213/go-pos-erp/pkg/app"
	"github.com/aslon1213/go-pos-erp/pkg/configs"
	"github.com/aslon1213/go-pos-erp/pkg/utils"
	"github.com/aslon1213/go-pos-erp/test/client"
	"go.mongodb.org/mongo-driver/bson"

	models "github.com/aslon1213/go-pos-erp/pkg/repository"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logging.SetOutput(io.Discard)
	app := app.New()

	app.DB.Database("magazin").Collection("finance").DeleteMany(context.Background(), bson.M{})
	app.DB.Database("magazin").Collection("suppliers").DeleteMany(context.Background(), bson.M{})
	app.DB.Database("magazin").Collection("transactions").DeleteMany(context.Background(), bson.M{})
	app.DB.Database("magazin").Collection("journals").DeleteMany(context.Background(), bson.M{})

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
			BranchName: "Xonobod",
			Details:    map[string]interface{}{"location": "Location A"},
		},
		{
			BranchName: "Polevoy",
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
		BranchName: "Xonobod",
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

	suppliers := []models.SupplierBase{
		{
			Name:    "Supplier One",
			Email:   "supplier1@example.com",
			Phone:   "1234567890",
			Address: "Address One",
			INN:     "111111111",
			Notes:   "Notes for Supplier One",
			Branch:  "Xonobod",
		},
		{
			Name:    "Supplier Two",
			Email:   "supplier2@example.com",
			Phone:   "2345678901",
			Address: "Address Two",
			INN:     "222222222",
			Notes:   "Notes for Supplier Two",
			Branch:  "Xonobod",
		},
		{
			Name:    "Supplier Three",
			Email:   "supplier3@example.com",
			Phone:   "3456789012",
			Address: "Address Three",
			INN:     "333333333",
			Notes:   "Notes for Supplier Three",
			Branch:  "Xonobod",
		},
		{
			Name:    "Supplier Four",
			Email:   "supplier4@example.com",
			Phone:   "4567890123",
			Address: "Address Four",
			INN:     "444444444",
			Notes:   "Notes for Supplier Four",
			Branch:  "Polevoy",
		},
		{
			Name:    "Supplier Five",
			Email:   "supplier5@example.com",
			Phone:   "5678901234",
			Address: "Address Five",
			INN:     "555555555",
			Notes:   "Notes for Supplier Five",
			Branch:  "Polevoy",
		},
	}
	for _, newSupplier := range suppliers {
		resp, output, err := client.CreateSupplier(newSupplier)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201, but got %d", resp.StatusCode)
		assert.Nil(t, output.Error, "Expected no error, but got one")
		assert.Equal(t, output.Data.Name, newSupplier.Name, "Expected supplier name to be %s, but got %s", newSupplier.Name, output.Data.Name)
	}
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
		if branch.BranchName == "Xonobod" {
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

func TestOpenJournals(t *testing.T) {
	loc := utils.GetTimeZone()

	configs, err := configs.LoadConfig("../../")
	if err != nil {
		t.Fatal(err)
	}

	client := client.NewClient("localhost", configs.Server.Port, "admin", "admin")

	// number of days to create journals
	days := 10
	journals := []models.NewJournalEntryInput{}
	for i := 0; i < days; i++ {
		journals = append(journals, models.NewJournalEntryInput{
			BranchNameOrID: "Xonobod",
			Date:           time.Now().In(loc).AddDate(0, 0, -i),
		})
	}

	for _, journal := range journals {
		resp, output, err := client.OpenJournal(journal)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201, but got %d", resp.StatusCode)
		assert.Nil(t, output.Error, "Expected no error, but got one")
	}

	// create once again with the same date

	// create once again with the same date to check if it will be created again --- should not be created again
	resp, output, err := client.OpenJournal(journals[0])
	if err != nil {
		log.Error().Err(err).Msg("Failed to open journal")
		t.Fatal(err)
	}
	assert.NotEqual(t, http.StatusOK, resp.StatusCode, "Expected status code 400, but got %d", resp.StatusCode)
	assert.NotNil(t, output.Error, "Expected error, but got none")
}

func TestQueryJournal(t *testing.T) {
	// loc := utils.GetTimeZone()

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

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output.Error, "Expected no error, but got one")

	// choose the first branch
	branchID := output.Data[0].BranchID

	// get journal entries
	resp, journals_output, err := client.QueryJournalEntries(branchID, models.JournalQueryParams{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, journals_output.Error, "Expected no error, but got one")
	assert.Equal(t, len(journals_output.Data), 10, "Expected 1 journal, but got %d", len(journals_output.Data))

	// get journal by id
	resp, output_, err := client.GetJournalByID(journals_output.Data[0].ID.Hex())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output_.Error, "Expected no error, but got one")
	assert.Equal(t, journals_output.Data[0].ID, output_.Data.ID, "Expected journal ID to be %s, but got %s", journals_output.Data[0].ID, output_.Data.ID)
}

func TestNewOperation(t *testing.T) {

	// query journals
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

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output.Error, "Expected no error, but got one")

	// choose the first branch
	branch := output.Data[0]
	branchID := output.Data[0].BranchID

	// get journal entries
	resp, journals_output, err := client.QueryJournalEntries(branchID, models.JournalQueryParams{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, journals_output.Error, "Expected no error, but got one")
	assert.Equal(t, len(journals_output.Data), 10, "Expected 1 journal, but got %d", len(journals_output.Data))

	operations := []models.JournalOperationInput{
		{
			TransactionBase: models.TransactionBase{
				Amount:        1000000,
				Description:   "Test Operation 1 in journal " + journals_output.Data[0].ID.Hex(),
				Type:          models.TransactionTypeDebit,
				PaymentMethod: models.PaymentMethodBank,
			},
			SupplierTransaction: false,
		},
		{
			TransactionBase: models.TransactionBase{
				Amount:        1000000,
				Description:   "Test Operation 2 in journal " + journals_output.Data[0].ID.Hex(),
				Type:          models.TransactionTypeCredit,
				PaymentMethod: models.PaymentMethodBank,
			},
			SupplierTransaction: false,
		},
		{
			TransactionBase: models.TransactionBase{
				Amount:        1000000,
				Description:   "Test Operation 3 in journal " + journals_output.Data[0].ID.Hex(),
				Type:          models.TransactionTypeDebit,
				PaymentMethod: models.PaymentMethodCash,
			},
			SupplierTransaction: false,
		},
		{
			TransactionBase: models.TransactionBase{
				Amount:        1000000,
				Description:   "Test Operation 4 in journal " + journals_output.Data[0].ID.Hex(),
				Type:          models.TransactionTypeCredit,
				PaymentMethod: models.PaymentMethodCash,
			},
			SupplierTransaction: false,
		},
		{
			TransactionBase: models.TransactionBase{
				Amount:        1000000,
				Description:   "Test Operation 5 in journal " + journals_output.Data[0].ID.Hex(),
				Type:          models.TransactionTypeCredit,
				PaymentMethod: models.PaymentMethodCash,
			},
			SupplierTransaction: true,
			SupplierID:          branch.Suppliers[0],
		},
		{
			TransactionBase: models.TransactionBase{
				Amount:        1000000,
				Description:   "Test Operation 6 in journal " + journals_output.Data[0].ID.Hex(),
				Type:          models.TransactionTypeCredit,
				PaymentMethod: models.PaymentMethodCash,
			},
			SupplierTransaction: true,
			SupplierID:          branch.Suppliers[0],
		},
		{
			TransactionBase: models.TransactionBase{
				Amount:        1000000,
				Description:   "Test Operation 7 in journal " + journals_output.Data[0].ID.Hex(),
				Type:          models.TransactionTypeCredit,
				PaymentMethod: models.PaymentMethodCash,
			},
			SupplierTransaction: true,
			SupplierID:          branch.Suppliers[0],
		},
	}

	for _, operation := range operations {
		resp, output, err := client.NewJournalOperation(journals_output.Data[0].ID.Hex(), operation)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201, but got %d", resp.StatusCode)
		assert.Nil(t, output.Error, "Expected no error, but got one")
	}
}

func TestOperationsGotInserted(t *testing.T) {

	// query journals
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

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output.Error, "Expected no error, but got one")

	// choose the first branch
	// branch := output.Data[0]
	branchID := output.Data[0].BranchID

	// get journal entries
	resp, journals_output, err := client.QueryJournalEntries(branchID, models.JournalQueryParams{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, journals_output.Error, "Expected no error, but got one")
	// assert.Equal(t, len(journals_output.Data), 10, "Expected 1 journal, but got %d", len(journals_output.Data))

	journal := journals_output.Data[0]
	assert.Equal(t, journal.Branch.ID, branchID, "Expected branch ID to be %s, but got %s", branchID, journal.Branch.ID)
	assert.Equal(t, len(journal.Operations), 7, "Expected 7 operations, but got %d", len(journal.Operations))
}

func TestCloseJournal(t *testing.T) {
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

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output.Error, "Expected no error, but got one")

	branchID := output.Data[0].BranchID

	// get journals
	resp, journals_output, err := client.QueryJournalEntries(branchID, models.JournalQueryParams{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, journals_output.Error, "Expected no error, but got one")

	// find first open journal

	journal := models.Journal{}
	log.Info().Interface("journals", journals_output.Data[0].ID.Hex()).Msg("Journals")
	for _, journal_ := range journals_output.Data {
		if !journal_.Shift_is_closed {
			journal = journal_
			break
		}
	}
	log.Info().Interface("journal chosen", journal).Msg("Journal")

	input := models.CloseJournalEntryInput{
		CashLeft:       1000000,
		TerminalIncome: 1000000,
	}

	resp, output_, err := client.CloseJournal(journal.ID.Hex(), input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to close journal")
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output_.Error, "Expected no error, but got one")

	// fetch this journal now
	resp, output_2, err := client.GetJournalByID(journal.ID.Hex())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get journal")
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output_2.Error, "Expected no error, but got one")
	assert.Equal(t, output_2.Data.Shift_is_closed, true, "Expected shift to be closed, but got %t", output_2.Data.Shift_is_closed)
	assert.Equal(t, output_2.Data.Total, journal.Total+input.CashLeft+input.TerminalIncome)
	assert.Equal(t, output_2.Data.Cash_left, input.CashLeft)
	assert.Equal(t, output_2.Data.Terminal_income, input.TerminalIncome)
}

func TestReOpenJournal(t *testing.T) {
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

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output.Error, "Expected no error, but got one")

	branchID := output.Data[0].BranchID

	// get journals
	resp, journals_output, err := client.QueryJournalEntries(branchID, models.JournalQueryParams{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, journals_output.Error, "Expected no error, but got one")

	// find first closed journal
	journal := models.Journal{}
	for _, journal_ := range journals_output.Data {
		if journal_.Shift_is_closed {
			journal = journal_
			break
		}
	}
	log.Info().Interface("journal chosen", journal.ID.Hex()).Msg("Journal")

	resp, output_, err := client.ReOpenJournal(journal.ID.Hex())
	if err != nil {
		log.Error().Err(err).Msg("Failed to re-open journal")
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output_.Error, "Expected no error, but got one")

	// fetch this journal now
	resp, output_2, err := client.GetJournalByID(journal.ID.Hex())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get journal")
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
	assert.Nil(t, output_2.Error, "Expected no error, but got one")
	assert.Equal(t, false, output_2.Data.Shift_is_closed, "Expected shift to be open, but got %t", output_2.Data.Shift_is_closed)
	// assert.Equal(t, output_2.Data.Total, journal.Total)
	assert.Equal(t, uint32(0), output_2.Data.Cash_left)
	assert.Equal(t, uint32(0), output_2.Data.Terminal_income)
}
