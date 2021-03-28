package test

import (
	"context"
	"github.com/lgunko/beauty-reuse/constants"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OrganizationServiceTestSuite struct {
	HookTestSuite
}

func TestOrganizationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationServiceTestSuite))
}

func (suite *OrganizationServiceTestSuite) SetupTest() {
	suite.cleanUp()
}

func (suite *OrganizationServiceTestSuite) cleanUp() {
	//delete everything on db level
	suite.database.Collection(constants.OrgCollection).Drop(context.Background())
	suite.database.Collection(constants.EmployeeCollection).Drop(context.Background())
	suite.hook.Reset()
}
