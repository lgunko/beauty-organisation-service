package repository

import (
	"context"
	"fmt"
	"github.com/lgunko/beauty-organisation-service/graph/model"
	"github.com/lgunko/beauty-reuse/constants"
	"github.com/lgunko/beauty-reuse/errorsutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	emailField    = "email"
	orgCollection = "Org"
)

type orgID struct {
	OrgID string `bson:"orgID"`
}

type orgToSave struct {
	ID     string `bson:"_id"`
	bson.M `bson:",inline"`
}

func toOrgToSave(input model.IOrgInput) orgToSave {
	bytes, _ := bson.Marshal(input.GetOrgInput())
	var raw bson.M
	bson.Unmarshal(bytes, &raw)
	return orgToSave{
		primitive.NewObjectID().Hex(),
		raw,
	}
}

func FindAllOrgs(ctx context.Context, result *[]*model.OrgOutput, db mongo.Database, email string) (context.Context, error) {
	globalFilter := bson.M{emailField: bson.M{"$eq": email}}
	cur, err := db.Collection("Employee").Find(ctx, globalFilter, options.Find().SetProjection(bson.M{constants.OrgIdField: 1}))
	if err != nil {
		return errorsutil.GetInternalServerError(ctx, err)
	}
	defer cur.Close(ctx)
	var currentEmployeeOrgIdList []orgID
	err = cur.All(ctx, &currentEmployeeOrgIdList)
	if err != nil {
		return errorsutil.GetInternalServerError(ctx, fmt.Errorf("can not decode mongo result : %s", err))
	}
	orgIdList := []string{}
	for _, elem := range currentEmployeeOrgIdList {
		orgIdList = append(orgIdList, elem.OrgID)
	}

	globalFilter = bson.M{constants.IdField: bson.M{"$in": orgIdList}}
	cur, err = db.Collection(orgCollection).Find(ctx, globalFilter, options.Find())
	if err != nil {
		return errorsutil.GetInternalServerError(ctx, err)
	}
	defer cur.Close(ctx)
	err = cur.All(ctx, result)
	if err != nil {
		return errorsutil.GetInternalServerError(ctx, fmt.Errorf("can not decode mongo result : %s", err))
	}
	return ctx, nil
}

func CreateOrUpdate(ctx context.Context, db mongo.Database, _id *string, org model.IOrgInput) (context.Context, *model.OrgOutput, error) {
	if _id != nil {
		return updateOrg(ctx, db, *_id, org)
	} else {
		return createOrg(ctx, db, org)
	}
}

func createOrg(ctx context.Context, db mongo.Database, org model.IOrgInput) (context.Context, *model.OrgOutput, error) {
	return createWithoutOrg(ctx, "Org", db, org)
}

func updateOrg(ctx context.Context, db mongo.Database, _id string, org model.IOrgInput) (context.Context, *model.OrgOutput, error) {
	return updateWithoutOrg(ctx, "Org", db, _id, org)
}

func createWithoutOrg(ctx context.Context, entityType string, db mongo.Database, entity model.IOrgInput) (context.Context, *model.OrgOutput, error) {
	mongoResult, err := db.Collection(entityType).InsertOne(ctx, toOrgToSave(entity))
	if err != nil {
		ctx, err := errorsutil.GetInternalServerError(ctx, err)
		return ctx, nil, err
	}
	output := entity.ToOutput(mongoResult.InsertedID.(string))
	return ctx, &output, nil
}

func updateWithoutOrg(ctx context.Context, entityType string, db mongo.Database, _id string, entity model.IOrgInput) (context.Context, *model.OrgOutput, error) {
	mongoResult, err := db.Collection(entityType).UpdateOne(ctx, bson.M{"_id": bson.M{"$eq": _id}}, bson.M{"$set": entity})
	if err != nil {
		ctx, err := errorsutil.GetInternalServerError(ctx, err)
		return ctx, nil, err
	}
	if mongoResult.MatchedCount != 1 || mongoResult.ModifiedCount != 1 {
		ctx, err := errorsutil.GetNotFoundError(ctx, nil)
		return ctx, nil, err
	}
	output := entity.ToOutput(_id)
	return ctx, &output, nil
}
