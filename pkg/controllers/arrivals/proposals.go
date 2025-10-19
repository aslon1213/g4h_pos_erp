package arrivals

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	models "github.com/aslon1213/g4h_pos_erp/pkg/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel"
)

var branches = []string{"xonobod", "polevoy"}

type ArrivalsHandlers struct {
	ctx                context.Context
	ArrivalsCollection *mongo.Collection
}

// GetImage handles GET /images/{image_name}
func (h *ArrivalsHandlers) GetImage(c *fiber.Ctx) error {
	tracer := otel.Tracer("arrivals-handlers")
	_, span := tracer.Start(h.ctx, "GetImage")
	defer span.End()

	imageName := c.Params("image_name")
	imagePath := filepath.Join("images", imageName)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Warn().Str("image_name", imageName).Msg("get_image.not_found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Image not found",
		})
	}

	return c.SendFile(imagePath)
}

// GetImageByArrivalsID handles GET /{arrivals_id}/image
func (h *ArrivalsHandlers) GetImageByArrivalsID(c *fiber.Ctx) error {
	tracer := otel.Tracer("arrivals-handlers")
	ctx, span := tracer.Start(h.ctx, "GetImageByArrivalsID")
	defer span.End()

	arrivalsID := c.Params("arrivals_id")
	objectID, err := bson.ObjectIDFromHex(arrivalsID)
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("get_image_by_arrivals_id.invalid_id")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid arrivals ID",
		})
	}

	var arrivals models.Arrivals
	err = h.ArrivalsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&arrivals)
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("get_image_by_arrivals_id.not_found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Arrivals not found",
		})
	}

	return c.Render("image_upload", fiber.Map{
		"arrivals_id": arrivalsID,
		"image":       arrivals.ImageFile,
		"url":         fmt.Sprintf("/arrivals/%s/image", arrivalsID),
	})
}

// UploadImage handles POST /{arrivals_id}/image
func (h *ArrivalsHandlers) UploadImage(c *fiber.Ctx) error {
	tracer := otel.Tracer("arrivals-handlers")
	ctx, span := tracer.Start(h.ctx, "UploadImage")
	defer span.End()

	arrivalsID := c.Params("arrivals_id")
	objectID, err := bson.ObjectIDFromHex(arrivalsID)
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("upload_image.invalid_id")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid arrivals ID",
		})
	}

	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("upload_image.no_file")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Check if arrivals exists
	var arrivals models.Arrivals
	err = h.ArrivalsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&arrivals)
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("upload_image.arrivals_not_found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Arrivals not found",
		})
	}

	// Delete old file if exists
	if arrivals.ImageFile != nil && *arrivals.ImageFile != "" {
		oldFilePath := filepath.Join(".", *arrivals.ImageFile)
		if _, err := os.Stat(oldFilePath); err == nil {
			if err := os.Remove(oldFilePath); err != nil {
				log.Warn().Err(err).Str("arrivals_id", arrivalsID).Msg("upload_image.delete_old_file_failed")
			}
		}
	}

	// Save new file
	newFileName := fmt.Sprintf("%s_%s", arrivalsID, file.Filename)
	newFilePath := filepath.Join("images", newFileName)

	if err := c.SaveFile(file, newFilePath); err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("upload_image.save_failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Update arrivals with new image path
	_, err = h.ArrivalsCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"image_file": newFilePath}},
	)
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("upload_image.update_failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update arrivals",
		})
	}

	log.Info().Str("arrivals_id", arrivalsID).Str("filename", file.Filename).Msg("upload_image.success")
	return c.Redirect("/arrivals/", http.StatusSeeOther)
}

// NewArrivals handles POST /{branch}/new
func (h *ArrivalsHandlers) NewArrivals(c *fiber.Ctx) error {
	tracer := otel.Tracer("arrivals-handlers")
	ctx, span := tracer.Start(h.ctx, "NewArrivals")
	defer span.End()

	branch := c.Params("branch")
	if !h.checkBranch(branch) {
		log.Error().Str("branch", branch).Msg("new_arrivals.invalid_branch")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid branch",
		})
	}

	var arrivalsNames []string
	if err := c.BodyParser(&arrivalsNames); err != nil {
		log.Error().Err(err).Str("branch", branch).Msg("new_arrivals.parse_error")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var arrivalsToInsert []interface{}
	now := time.Now()

	for _, name := range arrivalsNames {
		arrivals := models.Arrivals{
			ID:        bson.NewObjectID(),
			Name:      name,
			Date:      now,
			Branch:    branch,
			Fulfilled: false,
		}
		arrivalsToInsert = append(arrivalsToInsert, arrivals)
	}

	result, err := h.ArrivalsCollection.InsertMany(ctx, arrivalsToInsert)
	if err != nil {
		log.Error().Err(err).Str("branch", branch).Msg("new_arrivals.insert_failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create arrivals",
		})
	}

	var outputIDs []string
	for _, id := range result.InsertedIDs {
		outputIDs = append(outputIDs, id.(bson.ObjectID).Hex())
	}

	log.Info().Str("branch", branch).Int("count", len(arrivalsNames)).Msg("new_arrivals.success")
	return c.JSON(fiber.Map{
		"message":     "Arrivals created successfully",
		"arrivals_id": outputIDs,
	})
}

// GetArrivals handles GET ""
func (h *ArrivalsHandlers) GetArrivals(c *fiber.Ctx) error {
	tracer := otel.Tracer("arrivals-handlers")
	ctx, span := tracer.Start(h.ctx, "GetArrivals")
	defer span.End()

	filter := bson.M{}

	// Name filter
	if name := c.Query("name"); name != "" {
		filter["name"] = bson.M{"$regex": regexp.QuoteMeta(name), "$options": "i"}
	}

	// Branch filter
	if branch := c.Query("branch"); branch != "" {
		filter["branch"] = bson.M{"$regex": regexp.QuoteMeta(branch), "$options": "i"}
	}

	// Fulfilled filter
	if fulfilled := c.Query("fulfilled", "false"); fulfilled != "" {
		if fulfilled == "true" {
			filter["fulfilled"] = true
		} else if fulfilled == "false" {
			filter["fulfilled"] = false
		}
	}

	// Date filters
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		parsedDate, err := time.Parse("2006-01-02", dateFrom)
		if err != nil {
			log.Error().Str("date_from", dateFrom).Msg("get_arrivals.invalid_date_from")
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date_from format. Use YYYY-MM-DD.",
			})
		}
		filter["date"] = bson.M{"$gte": parsedDate}
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		parsedDate, err := time.Parse("2006-01-02", dateTo)
		if err != nil {
			log.Error().Str("date_to", dateTo).Msg("get_arrivals.invalid_date_to")
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date_to format. Use YYYY-MM-DD.",
			})
		}
		parsedDate = parsedDate.Add(24 * time.Hour) // Add one day
		if _, exists := filter["date"]; exists {
			filter["date"].(bson.M)["$lte"] = parsedDate
		} else {
			filter["date"] = bson.M{"$lte": parsedDate}
		}
	}

	cursor, err := h.ArrivalsCollection.Find(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("get_arrivals.find_failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch arrivals",
		})
	}
	defer cursor.Close(ctx)

	var arrivals []models.Arrivals
	if err := cursor.All(ctx, &arrivals); err != nil {
		log.Error().Err(err).Msg("get_arrivals.decode_failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode arrivals",
		})
	}

	// Format dates for response
	var response []fiber.Map
	for _, arrival := range arrivals {
		response = append(response, fiber.Map{
			"_id":       arrival.ID.Hex(),
			"name":      arrival.Name,
			"date":      arrival.Date.Format("02-01-2006"),
			"branch":    arrival.Branch,
			"fulfilled": arrival.Fulfilled,
		})
	}

	log.Info().Interface("filter", filter).Int("count", len(arrivals)).Msg("get_arrivals.success")
	return c.JSON(response)
}

// GetArrivals handles GET /{arrivals_id}
func (h *ArrivalsHandlers) GetArrivalsDetail(c *fiber.Ctx) error {
	tracer := otel.Tracer("arrivals-handlers")
	ctx, span := tracer.Start(h.ctx, "GetArrivalsDetail")
	defer span.End()

	arrivalsID := c.Params("arrivals_id")
	objectID, err := bson.ObjectIDFromHex(arrivalsID)
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("get_arrivals_detail.invalid_id")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid arrivals ID",
		})
	}

	var arrivals models.Arrivals
	err = h.ArrivalsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&arrivals)
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("get_arrivals_detail.not_found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Arrivals not found",
		})
	}

	log.Info().Str("arrivals_id", arrivalsID).Msg("get_arrivals_detail.success")
	return c.Render("arrivals", fiber.Map{
		"arrivals": fiber.Map{
			"_id":       arrivals.ID.Hex(),
			"name":      arrivals.Name,
			"date":      arrivals.Date.Format("02-01-2006"),
			"branch":    arrivals.Branch,
			"fulfilled": arrivals.Fulfilled,
		},
		"action": "view",
	})
}

// EditArrivals handles POST /{arrivals_id}/edit
func (h *ArrivalsHandlers) EditArrivals(c *fiber.Ctx) error {
	tracer := otel.Tracer("arrivals-handlers")
	ctx, span := tracer.Start(h.ctx, "EditArrivals")
	defer span.End()

	arrivalsID := c.Params("arrivals_id")
	objectID, err := bson.ObjectIDFromHex(arrivalsID)
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("edit_arrivals.invalid_id")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid arrivals ID",
		})
	}

	update := bson.M{"$set": bson.M{}}

	contentType := c.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("edit_arrivals.parse_json_error")
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON data",
			})
		}

		if name, ok := data["name"].(string); ok && name != "" {
			update["$set"].(bson.M)["name"] = name
		}
		if branch, ok := data["branch"].(string); ok && branch != "" {
			update["$set"].(bson.M)["branch"] = branch
		}
		if fulfilled, ok := data["fulfilled"]; ok {
			if fulfilledStr, ok := fulfilled.(string); ok {
				update["$set"].(bson.M)["fulfilled"] = strings.ToLower(fulfilledStr) == "true"
			} else if fulfilledBool, ok := fulfilled.(bool); ok {
				update["$set"].(bson.M)["fulfilled"] = fulfilledBool
			}
		}
	} else {
		if name := c.FormValue("name"); name != "" {
			update["$set"].(bson.M)["name"] = name
		}
		if branch := c.FormValue("branch"); branch != "" {
			update["$set"].(bson.M)["branch"] = branch
		}
		if fulfilled := c.FormValue("fulfilled"); fulfilled != "" {
			update["$set"].(bson.M)["fulfilled"] = strings.ToLower(fulfilled) == "true"
		}
	}

	_, err = h.ArrivalsCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("edit_arrivals.update_failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update arrivals",
		})
	}

	log.Info().Str("arrivals_id", arrivalsID).Interface("updates", update["$set"]).Msg("edit_arrivals.success")
	return c.Redirect("/arrivals", http.StatusFound)
}

// DeleteArrivals handles POST /{arrivals_id}/delete
func (h *ArrivalsHandlers) DeleteArrivals(c *fiber.Ctx) error {
	tracer := otel.Tracer("arrivals-handlers")
	ctx, span := tracer.Start(h.ctx, "DeleteArrivals")
	defer span.End()

	arrivalsID := c.Params("arrivals_id")
	objectID, err := bson.ObjectIDFromHex(arrivalsID)
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("delete_arrivals.invalid_id")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid arrivals ID",
		})
	}

	// Get the arrivals first to check for image file
	var arrivals models.Arrivals
	err = h.ArrivalsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&arrivals)
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("delete_arrivals.not_found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Arrivals not found",
		})
	}

	// Delete the arrivals
	_, err = h.ArrivalsCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("delete_arrivals.delete_failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete arrivals",
		})
	}

	// Delete the image file if it exists
	if arrivals.ImageFile != nil && *arrivals.ImageFile != "" {
		imageFilePath := filepath.Join(".", *arrivals.ImageFile)
		if _, err := os.Stat(imageFilePath); err == nil {
			if err := os.Remove(imageFilePath); err != nil {
				log.Warn().Err(err).Str("arrivals_id", arrivalsID).Msg("delete_arrivals.image_delete_failed")
			}
		}
	}

	log.Info().Str("arrivals_id", arrivalsID).Msg("delete_arrivals.success")
	return c.Redirect("/arrivals", http.StatusFound)
}

// FulfillArrivals handles GET /{branch}/fulfill
func (h *ArrivalsHandlers) FulfillArrivals(c *fiber.Ctx) error {
	tracer := otel.Tracer("arrivals-handlers")
	ctx, span := tracer.Start(h.ctx, "FulfillArrivals")
	defer span.End()

	branch := c.Params("branch")
	if !h.checkBranch(branch) {
		log.Error().Str("branch", branch).Msg("fulfill_arrivals.invalid_branch")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid branch",
		})
	}

	idsParam := c.Query("ids", "")
	if idsParam == "" {
		log.Error().Str("branch", branch).Msg("fulfill_arrivals.no_ids")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "No arrivals IDs provided",
		})
	}

	arrivalsIDs := strings.Split(idsParam, ",")
	var objectIDs []bson.ObjectID

	for _, id := range arrivalsIDs {
		objectID, err := bson.ObjectIDFromHex(strings.TrimSpace(id))
		if err != nil {
			log.Error().Err(err).Str("id", id).Msg("fulfill_arrivals.invalid_id")
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}

	if len(objectIDs) == 0 {
		log.Error().Str("branch", branch).Msg("fulfill_arrivals.no_valid_ids")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "No valid arrivals IDs provided",
		})
	}

	result, err := h.ArrivalsCollection.UpdateMany(
		ctx,
		bson.M{
			"_id":    bson.M{"$in": objectIDs},
			"branch": branch,
		},
		bson.M{"$set": bson.M{"fulfilled": true}},
	)
	if err != nil {
		log.Error().Err(err).Str("branch", branch).Msg("fulfill_arrivals.update_failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fulfill arrivals",
		})
	}

	log.Info().
		Str("branch", branch).
		Strs("arrivals_ids", arrivalsIDs).
		Int64("matched_count", result.MatchedCount).
		Int64("modified_count", result.ModifiedCount).
		Msg("fulfill_arrivals.success")

	return c.JSON(fiber.Map{
		"acknowledged":        true,
		"total_number_of_ids": len(arrivalsIDs),
		"matched_count":       result.MatchedCount,
		"modified_count":      result.ModifiedCount,
		"upserted_id":         nil,
	})
}

// GeneratePDF handles GET /pdf/pdf
func (h *ArrivalsHandlers) GeneratePDF(c *fiber.Ctx) error {
	tracer := otel.Tracer("arrivals-handlers")
	ctx, span := tracer.Start(h.ctx, "GeneratePDF")
	defer span.End()

	filter := bson.M{"fulfilled": false}
	cursor, err := h.ArrivalsCollection.Find(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("generate_pdf.find_failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch arrivals",
		})
	}
	defer cursor.Close(ctx)

	var arrivals []models.Arrivals
	if err := cursor.All(ctx, &arrivals); err != nil {
		log.Error().Err(err).Msg("generate_pdf.decode_failed")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode arrivals",
		})
	}

	// Format arrivals for PDF
	var formattedArrivals []fiber.Map
	baseURL := os.Getenv("ARRIVALS_BASE_URL")
	for _, arrival := range arrivals {
		imageFile := ""
		if arrival.ImageFile != nil && *arrival.ImageFile != "" && baseURL != "" {
			imageFile = baseURL + *arrival.ImageFile
		}

		formattedArrivals = append(formattedArrivals, fiber.Map{
			"_id":        arrival.ID.Hex(),
			"name":       arrival.Name,
			"date":       arrival.Date.Format("02-01-2006"),
			"branch":     arrival.Branch,
			"fulfilled":  arrival.Fulfilled,
			"image_file": imageFile,
		})
	}

	// Render PDF template (this would need to be implemented with a PDF library)
	// For now, returning JSON response
	log.Info().Int("arrivals_count", len(arrivals)).Msg("generate_pdf.success")

	return c.JSON(fiber.Map{
		"message":  "PDF generation not implemented yet",
		"arrivals": formattedArrivals,
	})
}

// Helper function to check branch validity
func (h *ArrivalsHandlers) checkBranch(branch string) bool {

	// Implement branch validation logic here
	// This is a placeholder implementation
	validBranches := branches
	for _, validBranch := range validBranches {
		if branch == validBranch {
			return true
		}
	}
	return false
}

// ##########################################################################################################################################
// ##########################################################################################################################################
// ##########################################################################################################################################
// NewArrivalsTemplate handles GET /{branch}/new
// func (h *ArrivalsHandlers) NewArrivalsTemplate(c *fiber.Ctx) error {
// 	tracer := otel.Tracer("arrivals-handlers")
// 	_, span := tracer.Start(h.ctx, "NewArrivalsTemplate")
// 	defer span.End()

// 	branch := c.Params("branch")
// 	if branch == "" {
// 		log.Error().Str("branch", branch).Msg("new_arrivals_template.invalid_branch")
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid branch",
// 		})
// 	}

// 	// Check branch validity (implement check_branch logic)
// 	if !h.checkBranch(branch) {
// 		log.Error().Str("branch", branch).Msg("new_arrivals_template.invalid_branch_check")
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid branch",
// 		})
// 	}

// 	return c.Render("new_arrivals", fiber.Map{
// 		"branch": branch,
// 	})
// }

// // DeleteArrivalsTemplate handles GET /{arrivals_id}/delete
// func (h *ArrivalsHandlers) DeleteArrivalsTemplate(c *fiber.Ctx) error {
// 	tracer := otel.Tracer("arrivals-handlers")
// 	ctx, span := tracer.Start(h.ctx, "DeleteArrivalsTemplate")
// 	defer span.End()

// 	arrivalsID := c.Params("arrivals_id")
// 	objectID, err := bson.ObjectIDFromHex(arrivalsID)
// 	if err != nil {
// 		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("delete_arrivals_template.invalid_id")
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid arrivals ID",
// 		})
// 	}

// 	var arrivals models.Arrivals
// 	err = h.ArrivalsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&arrivals)
// 	if err != nil {
// 		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("delete_arrivals_template.not_found")
// 		return c.Status(http.StatusNotFound).JSON(fiber.Map{
// 			"error": "Arrivals not found",
// 		})
// 	}

// 	log.Info().Str("arrivals_id", arrivalsID).Msg("delete_arrivals_template.success")
// 	return c.Render("arrivals", fiber.Map{
// 		"arrivals": fiber.Map{
// 			"_id": arrivals.ID.Hex(),
// 		},
// 		"action": "delete",
// 	})
// }

// // EditArrivalsTemplate handles GET /{arrivals_id}/edit
// func (h *ArrivalsHandlers) EditArrivalsTemplate(c *fiber.Ctx) error {
// 	tracer := otel.Tracer("arrivals-handlers")
// 	ctx, span := tracer.Start(h.ctx, "EditArrivalsTemplate")
// 	defer span.End()

// 	arrivalsID := c.Params("arrivals_id")
// 	objectID, err := bson.ObjectIDFromHex(arrivalsID)
// 	if err != nil {
// 		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("edit_arrivals_template.invalid_id")
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid arrivals ID",
// 		})
// 	}

// 	var arrivals models.Arrivals
// 	err = h.ArrivalsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&arrivals)
// 	if err != nil {
// 		log.Error().Err(err).Str("arrivals_id", arrivalsID).Msg("edit_arrivals_template.not_found")
// 		return c.Status(http.StatusNotFound).JSON(fiber.Map{
// 			"error": "Arrivals not found",
// 		})
// 	}

// 	log.Info().Str("arrivals_id", arrivalsID).Msg("edit_arrivals_template.success")
// 	return c.Render("arrivals", fiber.Map{
// 		"arrivals": fiber.Map{
// 			"_id":       arrivals.ID.Hex(),
// 			"name":      arrivals.Name,
// 			"date":      arrivals.Date.Format("2006-01-02"),
// 			"branch":    arrivals.Branch,
// 			"fulfilled": arrivals.Fulfilled,
// 		},
// 		"action": "edit",
// 	})
// }
