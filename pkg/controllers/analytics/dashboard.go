package analytics

import (
	"context"
	"time"

	journal_handlers "github.com/aslon1213/go-pos-erp/pkg/controllers/journals"
	models "github.com/aslon1213/go-pos-erp/pkg/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var branches []models.BranchFinance

type DashboardHandler struct {
	journals *mongo.Collection
	tracer   trace.Tracer
}

func New(db *mongo.Database) *DashboardHandler {
	tracer := otel.Tracer("DashboardHandler")
	ctx := context.Background()
	cursor, err := db.Collection("finance").Find(ctx, bson.M{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get branches")
	}

	err = cursor.All(ctx, &branches)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get branches")
	}

	log.Info().Msgf("Fetched %d branches", len(branches))
	return &DashboardHandler{
		journals: db.Collection("journals"),
		tracer:   tracer,
	}
}

func (d *DashboardHandler) MainPage(c *fiber.Ctx) error {
	body := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Analytics Dashboards</title>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <style>
            body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif; margin: 0; padding: 24px; background: #f3f4f6; }
            .container { max-width: 900px; margin: 0 auto; }
            .header { text-align: center; margin-bottom: 24px; }
            .cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 16px; }
            .card { background: white; padding: 20px; border-radius: 12px; box-shadow: 0 8px 25px rgba(0,0,0,0.05); }
            .card h2 { margin: 0 0 8px 0; font-size: 18px; color: #111827; }
            .card p { margin: 0 0 12px 0; color: #6b7280; }
            .card a { display: inline-block; color: #fff; background: #4f46e5; padding: 10px 14px; border-radius: 8px; text-decoration: none; font-weight: 600; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>Analytics Dashboards</h1>
                <p>Select a dashboard to view analytics</p>
            </div>
            <div class="cards">
                <div class="card">
                    <h2>General Analytics</h2>
                    <p>Branch comparison and detailed operational insights.</p>
                    <a href="/dashboard/general">Open General</a>
                </div>
                <div class="card">
                    <h2>Daily Analytics</h2>
                    <p>Global performance, payment method trends, and KPIs.</p>
                    <a href="/dashboard/journals">Open Daily</a>
                </div>
            </div>
        </div>
    </body>
    </html>`
	c.Set(
		"Content-Type",
		"text/html",
	)
	return c.SendString(body)
}

func (d *DashboardHandler) ServeDashBoardGeneral(c *fiber.Ctx) error {
	// General dashboard - branch comparison, growth, operations health
	from_date := c.Query(
		"from_date",
		time.Now().AddDate(0, 0, -30).Format("2006-01-02"),
	)
	to_date := c.Query(
		"to_date",
		time.Now().Local().Format("2006-01-02"),
	)
	branch_id := c.Query("branch_id", "")

	from_date_time, err := time.Parse("2006-01-02", from_date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid from date",
		})
	}
	to_date_time, err := time.Parse("2006-01-02", to_date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid to date",
		})
	}
	length := int(to_date_time.Sub(from_date_time).Hours() / 24)
	log.Info().Msgf("Length: %d", length)

	ctx, span := d.tracer.Start(c.Context(), "DashboardHandler.ServeDashBoardGeneral")
	defer span.End()

	var allJournals []models.Journal
	if branch_id != "" {
		journals, err := journal_handlers.QueryJournals(span, ctx, c, models.JournalQueryParams{
			BranchID: branch_id,
			Page:     1,
			PageSize: length,
		}, d.journals)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get journal data",
			})
		}
		allJournals = journals
	} else {
		for _, branch := range branches {
			journals, err := journal_handlers.QueryJournals(span, ctx, c, models.JournalQueryParams{
				BranchID: branch.BranchID,
				Page:     1,
				PageSize: length,
			}, d.journals)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to get journal data for branch %s", branch.BranchName)
				continue
			}
			allJournals = append(allJournals, journals...)
		}
	}

	c.Set("Content-Type", "text/html")
	RenderGeneralDashboard(allJournals, c.Response().BodyWriter())
	return nil
}

func (d *DashboardHandler) ServeDashBoardDays(c *fiber.Ctx) error {
	// this is where the days dashboard - for showing [Global Performance Analytics, Payment Method Trends, Moving Average, Financial KPIs]
	from_date := c.Query(
		"from_date",
		time.Now().AddDate(0, 0, -30).Format("2006-01-02"),
	)
	to_date := c.Query(
		"to_date",
		time.Now().Local().Format("2006-01-02"),
	)
	branch_id := c.Query("branch_id", "") // Optional branch filter

	from_date_time, err := time.Parse("2006-01-02", from_date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid from date",
		})
	}
	to_date_time, err := time.Parse("2006-01-02", to_date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid to date",
		})
	}

	log.Info().Msgf("From date: %s", from_date)
	log.Info().Msgf("To date: %s", to_date)
	log.Info().Msgf("Branch ID filter: %s", branch_id)

	ctx, span := d.tracer.Start(c.Context(), "DashboardHandler.ServeDashBoard")
	defer span.End()

	var allJournals []models.Journal
	length := int(to_date_time.Sub(from_date_time).Hours() / 24)
	// limit := length * len(branches)
	// log.Info().Msgf("Length: %d", length)
	// log.Info().Msgf("Limit: %d", limit)
	// log.Info().Msgf("From date: %s", from_date_time)
	// log.Info().Msgf("To date: %s", to_date_time)

	log.Info().Msgf("Length: %d", length)
	if branch_id != "" {
		// Query for specific branch
		journals, err := journal_handlers.QueryJournals(span, ctx, c, models.JournalQueryParams{
			BranchID: branch_id,
			Page:     1,
			PageSize: length,
		}, d.journals)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get journal data",
			})
		}
		allJournals = journals
	} else {
		// Query for all branches

		for _, branch := range branches {
			journals, err := journal_handlers.QueryJournals(span, ctx, c, models.JournalQueryParams{
				BranchID: branch.BranchID,
				Page:     1,
				PageSize: length,
			}, d.journals)

			if err != nil {
				log.Error().Err(err).Msgf("Failed to get journal data for branch %s", branch.BranchName)
				continue
			}
			allJournals = append(allJournals, journals...)
		}
	}

	log.Info().Msgf("Fetched %d journals", len(allJournals))

	c.Set(
		"Content-Type",
		"text/html",
	)
	RenderDaysDashboard(allJournals, c.Response().BodyWriter())

	return nil
}

// ServeDashBoardComparison builds multiple periods from query params and renders comparison
func (d *DashboardHandler) ServeDashBoardComparison(c *fiber.Ctx) error {
	// Query params: from1,to1, from2,to2, from3,to3, branch_id (optional)
	branchID := c.Query("branch_id", "")
	ranges := [][2]string{
		{c.Query("from1", ""), c.Query("to1", "")},
		{c.Query("from2", ""), c.Query("to2", "")},
		{c.Query("from3", ""), c.Query("to3", "")},
	}

	ctx, span := d.tracer.Start(c.Context(), "DashboardHandler.ServeDashBoardComparison")
	defer span.End()

	var periods []ComparisonPeriod
	for idx, r := range ranges {
		if r[0] == "" || r[1] == "" {
			continue
		}
		fromTime, err := time.Parse("2006-01-02", r[0])
		if err != nil {
			continue
		}
		toTime, err := time.Parse("2006-01-02", r[1])
		if err != nil {
			continue
		}

		var journals []models.Journal
		if branchID != "" {
			j, err := journal_handlers.QueryJournals(span, ctx, c, models.JournalQueryParams{
				BranchID: branchID,
				FromDate: fromTime,
				ToDate:   toTime,
				Page:     1,
				PageSize: 200,
			}, d.journals)
			if err == nil {
				journals = j
			}
		} else {
			for _, br := range branches {
				j, err := journal_handlers.QueryJournals(span, ctx, c, models.JournalQueryParams{
					BranchID: br.BranchID,
					FromDate: fromTime,
					ToDate:   toTime,
					Page:     1,
					PageSize: 200,
				}, d.journals)
				if err != nil {
					log.Error().Err(err).Msgf("Failed to get journal data for branch %s", br.BranchName)
					continue
				}
				journals = append(journals, j...)
			}
		}

		label := c.Query("label"+string(rune('1'+idx)), "")
		if label == "" {
			label = fromTime.Format("2006-01-02") + " to " + toTime.Format("2006-01-02")
		}
		periods = append(periods, ComparisonPeriod{Label: label, Journals: journals})
	}

	c.Set("Content-Type", "text/html")
	RenderComparisonDashboard(periods, c.Response().BodyWriter())
	return nil
}
