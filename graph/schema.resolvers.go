package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	model2 "github.com/lgunko/beauty-backend/EmployeeService/graph/model"
	"github.com/lgunko/beauty-organisation-service/graph/generated"
	"github.com/lgunko/beauty-organisation-service/graph/model"
	"github.com/lgunko/beauty-organisation-service/repository"
	model1 "github.com/lgunko/beauty-reuse/graph/model"
	"github.com/lgunko/beauty-reuse/headers"
	"github.com/lgunko/beauty-reuse/loggingutil"
	"github.com/lgunko/beauty-reuse/orgbasedrepository"
)

func (r *mutationResolver) RegisterOrg(ctx context.Context, input model.InitialOrgInput) (*model.OrgOutput, error) {
	ctx, output, err := repository.CreateOrUpdate(ctx, r.GetDatabase(), nil, input)
	if err != nil {
		loggingutil.GetLoggerFilledFromContext(ctx).Error(err)
		return nil, err
	}
	email := ctx.Value(headers.Email).(string)
	ctx, _, err = orgbasedrepository.New(output.ID, r.GetDatabase(), "Employee").Create(ctx, model2.EmployeeInput{Name: input.CreatorName, Surname: input.CreatorSurname, Email: email})
	if err != nil {
		loggingutil.GetLoggerFilledFromContext(ctx).Error(err)
		return nil, err
	}
	output.Role = model1.RoleManager

	return output, nil
}

func (r *queryResolver) AllowedOrgList(ctx context.Context) ([]*model.OrgOutput, error) {
	email := ctx.Value(headers.Email).(string)
	var result []*model.OrgOutput
	ctx, err := repository.FindAllOrgs(ctx, &result, r.GetDatabase(), email)
	if err != nil {
		loggingutil.GetLoggerFilledFromContext(ctx).Error(err)
		return nil, err
	}
	for _, org := range result {
		ctx, role, err := repository.GetUsersRole(ctx, r.GetDatabase(), org.ID, email)
		if err != nil {
			loggingutil.GetLoggerFilledFromContext(ctx).Error(err)
			return nil, err
		}
		org.Role = *role
	}

	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
