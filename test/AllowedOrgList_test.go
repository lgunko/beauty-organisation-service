package test

import (
	"context"
	model2 "github.com/lgunko/beauty-backend/EmployeeService/graph/model"
	model5 "github.com/lgunko/beauty-backend/ReusableLib/graph/model"
	"github.com/lgunko/beauty-organisation-service/client"
	model4 "github.com/lgunko/beauty-organisation-service/graph/model"
	"github.com/lgunko/beauty-organisation-service/repository"
	"github.com/lgunko/beauty-reuse/constants"
	"github.com/lgunko/beauty-reuse/orgbasedrepository"
	"github.com/stretchr/testify/assert"
)

func (suite *OrganizationServiceTestSuite) TestAllowedOrgListWithoutHeaders() {
	assert.Empty(suite.T(), suite.hook.AllEntries())
	allowedOrgList, gqlErrors, err := client.GetAllowedOrgList(context.Background(), "http://localhost:8080/api/query")
	suite.testWithoutHeader(allowedOrgList, gqlErrors, err)
}

func (suite *OrganizationServiceTestSuite) TestNoRoleInOrgAllowedOrgsWithHeaders() {
	// set up
	ctx, _, err := orgbasedrepository.New(testingOrgID, suite.database, constants.EmployeeCollection).Create(context.Background(), model2.EmployeeInput{
		Name:        "TestName",
		Surname:     "TestSurname",
		Email:       testingEmail,
		MobilePhone: "",
		Role:        "",
		FavourList:  nil,
		AvatarURL:   nil,
	})

	// call
	allowedOrgs, gqlErrors, err := client.GetAllowedOrgList(setHeaders(ctx), "http://localhost:8080/api/query")

	// assertions
	assert.Empty(suite.T(), allowedOrgs)
	assert.Nil(suite.T(), gqlErrors)
	assert.Nil(suite.T(), err)
	assert.Len(suite.T(), suite.hook.AllEntries(), 0)
}

func (suite *OrganizationServiceTestSuite) TestMasterRoleInOrgAllowedOrgsWithHeaders() {
	// set up
	ctx, orgOutput, err := repository.CreateOrUpdate(context.Background(), suite.database, nil, model4.OrgInput{
		Name:    "TestOrgName",
		City:    "TestOrgCity",
		Address: "TestOrgAddress",
		LogoURL: nil,
	})
	ctx, _, err = orgbasedrepository.New(orgOutput.ID, suite.database, constants.EmployeeCollection).Create(context.Background(), model2.EmployeeInput{
		Name:        "TestName",
		Surname:     "TestSurname",
		Email:       testingEmail,
		MobilePhone: "",
		Role:        model5.RoleMaster,
		FavourList:  nil,
		AvatarURL:   nil,
	})

	// call
	allowedOrgs, gqlErrors, err := client.GetAllowedOrgList(setHeaders(ctx), "http://localhost:8080/api/query")

	// assertions
	assert.Len(suite.T(), allowedOrgs, 1)
	assert.Contains(suite.T(), allowedOrgs, model4.OrgOutput{
		ID:      orgOutput.ID,
		Name:    "TestOrgName",
		City:    "TestOrgCity",
		Address: "TestOrgAddress",
		LogoURL: nil,
	})
	assert.Nil(suite.T(), gqlErrors)
	assert.Nil(suite.T(), err)
	assert.Len(suite.T(), suite.hook.AllEntries(), 0)
}

func (suite *OrganizationServiceTestSuite) TestManagerRoleInOrgAllowedOrgsWithHeaders() {
	// set up
	ctx, orgOutput, err := repository.CreateOrUpdate(context.Background(), suite.database, nil, model4.OrgInput{
		Name:    "TestOrgName",
		City:    "TestOrgCity",
		Address: "TestOrgAddress",
		LogoURL: nil,
	})
	ctx, _, err = orgbasedrepository.New(orgOutput.ID, suite.database, constants.EmployeeCollection).Create(context.Background(), model2.EmployeeInput{
		Name:        "TestName",
		Surname:     "TestSurname",
		Email:       testingEmail,
		MobilePhone: "",
		Role:        model5.RoleManager,
		FavourList:  nil,
		AvatarURL:   nil,
	})

	// call
	allowedOrgs, gqlErrors, err := client.GetAllowedOrgList(setHeaders(ctx), "http://localhost:8080/api/query")

	// assertions
	assert.Len(suite.T(), allowedOrgs, 1)
	assert.Contains(suite.T(), allowedOrgs, model4.OrgOutput{
		ID:      orgOutput.ID,
		Name:    "TestOrgName",
		City:    "TestOrgCity",
		Address: "TestOrgAddress",
		LogoURL: nil,
	})
	assert.Nil(suite.T(), gqlErrors)
	assert.Nil(suite.T(), err)
	assert.Len(suite.T(), suite.hook.AllEntries(), 0)
}
