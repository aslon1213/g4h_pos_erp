package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ActivityType string

const (
	ActivityTypeLogin             ActivityType = "login_success"
	ActivityTypeLoginFailed       ActivityType = "login_failed"
	ActivityTypeLogout            ActivityType = "logout_success"
	ActivityTypeLogoutFailed      ActivityType = "logout_failed"
	ActivityTypeRegister          ActivityType = "register_success"
	ActivityTypeRegisterFailed    ActivityType = "register_failed"
	ActivityTypeCreateTransaction ActivityType = "create_transaction"
	ActivityTypeCreateJournal     ActivityType = "create_journal"
	ActivityTypeCloseJournal      ActivityType = "close_journal"
	ActivityTypeCreateSupplier    ActivityType = "create_supplier"
	ActivityTypeCreateProduct     ActivityType = "create_product"
	ActivityTypeEditTransaction   ActivityType = "edit_transaction"
	ActivityTypeEditSupplier      ActivityType = "edit_supplier"
	ActivityTypeEditProduct       ActivityType = "edit_product"
	ActivityTypeDeleteTransaction ActivityType = "delete_transaction"
	ActivityTypeReopenJournal     ActivityType = "reopen_journal"
	ActivityTypeDeleteSupplier    ActivityType = "delete_supplier"
	ActivityTypeDeleteProduct     ActivityType = "delete_product"
	ActivityTypeCloseSalesSession ActivityType = "close_sales_session"
	ActivityTypeOpenSalesSession  ActivityType = "open_sales_session"
	ActivityTypeProductIncome     ActivityType = "product_income"
	ActivityTypeProductTransfer   ActivityType = "product_transfer"
	ActivityTypeCreateOperation   ActivityType = "create_operation"
	ActivityTypeEditOperation     ActivityType = "edit_operation"
	ActivityTypeDeleteOperation   ActivityType = "delete_operation"
	ActivityTypeCreateFinance     ActivityType = "create_finance"
	ActivityTypeEditFinance       ActivityType = "edit_finance"
	ActivityTypeDeleteFinance     ActivityType = "delete_finance"
)

type Activity struct {
	UserID string       `bson:"user_id"`
	Action ActivityType `bson:"action"`
	Data   interface{}  `bson:"data"`
	IP     string       `bson:"ip"`
	Date   time.Time    `bson:"date"`
	Status int          `bson:"status"`
}

func SetActionType(ctx *fiber.Ctx, action ActivityType) {
	log.Debug().Str("action", string(action)).Msg("SetActionType called")
	ctx.Locals(
		"action",
		action,
	)
}

func SetUser(ctx *fiber.Ctx, user string) {
	log.Debug().Str("user", user).Msg("SetUser called")
	ctx.Locals(
		"user",
		user,
	)
}

func SetData(ctx *fiber.Ctx, data interface{}) {
	log.Debug().Interface("data", data).Msg("SetData called")
	ctx.Locals(
		"data",
		data,
	)
}

func LogActivity(ctx *fiber.Ctx) {
	log.Debug().Msg("LogActivity called")
	ctx.Locals(
		"log_activity",
		true,
	)
}

func DontLogActivity(ctx *fiber.Ctx) {
	log.Debug().Msg("DontLogActivity called")
	ctx.Locals(
		"log_activity",
		false,
	)
}

func RecordActicity(ctx context.Context, activity Activity, activities *mongo.Collection) error {
	_, err := activities.InsertOne(ctx, activity)
	if err != nil {
		return err
	}
	return nil
}

func (a *Middlewares) RecordActicityMiddleware(c *fiber.Ctx) error {
	//

	log.Debug().Msg("RecordActicityMiddleware calling next functions")
	c_err := c.Next()

	log_activity, ok := c.Locals("log_activity").(bool)
	if !ok {
		log_activity = false
	}
	log.Debug().Msg("log_activity: " + strconv.FormatBool(log_activity))
	if log_activity {

		action, ok := c.Locals("action").(ActivityType)
		if !ok {
			log.Error().Msg("Action not found")
			return c_err
		}

		user, ok := c.Locals("user").(string)
		if !ok {
			log.Error().Msg("User not found")
			return c_err
		}

		data := c.Locals("data")

		log.Debug().Interface("data", data).Str("action", string(action)).Str("user", user).Msg("Activity data")
		activity := Activity{
			UserID: user,
			Action: action,
			Data:   data,
			IP:     c.IP(),
			Date:   time.Now(),
			Status: c.Response().StatusCode(),
		}

		err := RecordActicity(c.Context(), activity, a.ActivitiesCollection)
		log.Debug().Msg("Activity recorded")
		if err != nil {
			return c_err
		}
	}
	return c_err
}
