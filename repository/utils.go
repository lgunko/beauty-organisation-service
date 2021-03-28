package repository

import (
	"context"
	"fmt"
	"github.com/lgunko/beauty-reuse/constants"
	"github.com/lgunko/beauty-reuse/errorsutil"
	model1 "github.com/lgunko/beauty-reuse/graph/model"
	"github.com/lgunko/beauty-reuse/orgbasedrepository"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

func GetUsersRole(ctx context.Context, db mongo.Database, orgID, email string) (context.Context, *model1.Role, error) {
	var usersRole *model1.Role
	var err error
	wg := sync.WaitGroup{}
	wg.Add(len(model1.AllRole))
	for _, role := range model1.AllRole {
		go func(role model1.Role) {
			ctx, count, err := orgbasedrepository.New(orgID, db, constants.EmployeeCollection).CountDocuments(ctx, orgbasedrepository.GetFilterFromEmail(email), orgbasedrepository.GetFilterFromRole(role))
			if err != nil {
				ctx, err = errorsutil.GetInternalServerError(ctx, err)
				wg.Done()
				return
			}
			switch count {
			case 0:

			case 1:
				usersRole = &role
			default:
				ctx, err = errorsutil.GetInternalServerError(ctx, fmt.Errorf("duplicated employee"))
				wg.Done()
				return
			}
			wg.Done()
		}(role)
	}
	wg.Wait()
	if err != nil {
		return ctx, nil, err
	}
	return ctx, usersRole, nil
}
