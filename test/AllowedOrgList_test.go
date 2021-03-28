package test

import (
	"context"
	modelEmployee "github.com/lgunko/beauty-employee-service/graph/model"
	"github.com/lgunko/beauty-organisation-service/client"
	"github.com/lgunko/beauty-organisation-service/graph/model"
	"github.com/lgunko/beauty-organisation-service/repository"
	"github.com/lgunko/beauty-reuse/constants"
	modelReuse "github.com/lgunko/beauty-reuse/graph/model"
	"github.com/lgunko/beauty-reuse/orgbasedrepository"
	"github.com/lgunko/beauty-reuse/test"
	"github.com/stretchr/testify/assert"
)

func (suite *OrganizationServiceTestSuite) TestAllowedOrgListWithoutHeaders() {
	assert.Empty(suite.T(), suite.hook.AllEntries())
	allowedOrgList, gqlErrors, err := client.GetAllowedOrgList(context.Background(), "http://localhost:8080/api/query")
	suite.testWithoutHeader(allowedOrgList, gqlErrors, err)
}

func (suite *OrganizationServiceTestSuite) TestNoRoleInOrgAllowedOrgsWithHeaders() {
	suite.execTestAllowedOrgsWithRole("", false)
}

func (suite *OrganizationServiceTestSuite) TestMasterRoleInOrgAllowedOrgsWithHeaders() {
	suite.execTestAllowedOrgsWithRole(modelReuse.RoleMaster, true)
}

func (suite *OrganizationServiceTestSuite) TestAdministratorRoleInOrgAllowedOrgsWithHeaders() {
	suite.execTestAllowedOrgsWithRole(modelReuse.RoleAdministrator, true)
}

func (suite *OrganizationServiceTestSuite) TestManagerRoleInOrgAllowedOrgsWithHeaders() {
	suite.execTestAllowedOrgsWithRole(modelReuse.RoleManager, true)
}

func (suite *OrganizationServiceTestSuite) execTestAllowedOrgsWithRole(role modelReuse.Role, success bool) {
	// set up
	ctx, orgOutput, err := repository.CreateOrUpdate(context.Background(), suite.database, nil, model.OrgInput{
		Name:    test.TestingOrgName,
		City:    test.TestingOrgCity,
		Address: test.TestingOrgAddress,
		LogoURL: nil,
	})
	ctx, _, err = orgbasedrepository.New(orgOutput.ID, suite.database, constants.EmployeeCollection).Create(context.Background(), modelEmployee.EmployeeInput{
		Name:        test.TestingEmployeeName,
		Surname:     test.TestingEmployeeSurname,
		Email:       testingEmail,
		MobilePhone: test.TestingEmployeeMobilePhone,
		Role:        role,
		FavourList:  nil,
		AvatarURL:   nil,
	})

	// call
	allowedOrgs, gqlErrors, err := client.GetAllowedOrgList(setHeaders(ctx), "http://localhost:8080/api/query")

	// assertions
	if success {
		assert.Len(suite.T(), allowedOrgs, 1)
		assert.Contains(suite.T(), allowedOrgs, model.OrgOutput{
			ID:      orgOutput.ID,
			Name:    test.TestingOrgName,
			City:    test.TestingOrgCity,
			Address: test.TestingOrgAddress,
			LogoURL: nil,
		})
		assert.Nil(suite.T(), gqlErrors)
		assert.Nil(suite.T(), err)
		assert.Len(suite.T(), suite.hook.AllEntries(), 0)
	} else {
		assert.Empty(suite.T(), allowedOrgs)
		assert.Nil(suite.T(), gqlErrors)
		assert.Nil(suite.T(), err)
		assert.Len(suite.T(), suite.hook.AllEntries(), 0)
	}
}
