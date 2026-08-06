package requests_mongo_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// PatchRequest updates selected fields of an existing request using $set.
// Patched fields: "Title", "Status", "Origin", "Docs", "WorkTypes", "Geography", "Deadline", "Contract".
func (r *Repository) PatchRequest(ctx context.Context, id uuid.UUID, request domain.Request) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if id == uuid.Nil {
		return fmt.Errorf("`ID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}
	if request.Version < 1 {
		return fmt.Errorf("`Version` can't be uninitialized: %w", core_errors.ErrInvalidArgument)
	}

	model := domainToModel(request)

	now, versionInc := time.Now(), model.Version+1

	update := make([]bson.E, 0, 12)
	update = append(update, bson.E{Key: "v", Value: versionInc}, bson.E{Key: "updated_at", Value: now})

	update = append(update, collectPatchFields(&model)...)

	origin, err := parseOriginValue(model.Origin)
	if err != nil {
		return fmt.Errorf("parse origin: %w", err)
	}
	update = append(update, bson.E{Key: "origin", Value: origin})

	model.CalcRequiredEmptyFields()
	update = append(update, bson.E{Key: "required_empty_fields", Value: model.RequiredEmptyFields})

	filter := bson.D{{Key: "_id", Value: id.String()}, {Key: "v", Value: request.Version}}

	if err := defaultPatch(ctx, r.pool.Requests(), filter, update); err != nil {
		return err
	}

	return nil
}

func collectPatchFields(model *RequestModel) []bson.E {
	update := make([]bson.E, 0, 9)

	if model.Title != "" {
		update = append(update, bson.E{Key: "title", Value: model.Title})
	}

	if model.Docs.TechTaskDocID != nil || model.Docs.ProjectDocID != nil || model.Docs.SpecificationDocID != nil {
		update = append(update, bson.E{Key: "docs", Value: model.Docs})
	}

	update = append(update, bson.E{Key: "work_types", Value: model.WorkTypes})
	update = append(update, bson.E{Key: "geo_desc", Value: model.Geography})
	update = append(update, bson.E{Key: "deadline", Value: model.Deadline})
	update = append(update, bson.E{Key: "contract", Value: model.ContractDocID})
	update = append(update, bson.E{Key: "client", Value: model.ClientID})

	return update
}

func parseOriginValue(origins []RequestModelOrigin) ([]RequestModelOrigin, error) {
	if len(origins) == 0 {
		return nil, nil
	}

	for i, origin := range origins {
		switch origin.Type {
		case domain.RequestOriginTypePhone:
			phone, err := core_validation.ValidatePhoneNumber(origin.Value)
			if err != nil {
				return nil, fmt.Errorf("`Value` to phone: %w", err)
			}
			origin.Value = phone.E164
		case domain.RequestOriginTypeEmail:
			email, err := core_validation.ValidateStringEmail(&origin.Value)
			if err != nil {
				return nil, fmt.Errorf("`Value` to email: %w", err)
			}
			origin.Value = email.String()
		}

		origins[i] = origin
	}

	return origins, nil
}
