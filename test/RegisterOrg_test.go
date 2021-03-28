package test

import (
	"context"
	"github.com/lgunko/beauty-organisation-service/client"
	"github.com/lgunko/beauty-organisation-service/graph/model"
	"github.com/lgunko/beauty-reuse/test"
	"github.com/stretchr/testify/assert"
)

func (suite *OrganizationServiceTestSuite) TestRegisterOrgWithoutHeaders() {
	assert.Empty(suite.T(), suite.hook.AllEntries())
	org, gqlErrors, err := client.RegisterOrg(context.Background(), "http://localhost:8080/api/query", model.InitialOrgInput{
		CreatorName:    test.TestingEmployeeName,
		CreatorSurname: test.TestingEmployeeSurname,
		OrgInput: model.OrgInput{
			Name:    test.TestingOrgName,
			City:    test.TestingOrgCity,
			Address: test.TestingOrgAddress,
			LogoURL: nil,
		},
	})
	suite.testWithoutHeader(org, gqlErrors, err)
}

func (suite *OrganizationServiceTestSuite) TestRegisterOrgWithHeaders() {
	// set up
	assert.Empty(suite.T(), suite.hook.AllEntries())

	// call
	org, gqlErrors, err := client.RegisterOrg(setHeaders(context.Background()), "http://localhost:8080/api/query", model.InitialOrgInput{
		CreatorName:    test.TestingEmployeeName,
		CreatorSurname: test.TestingEmployeeSurname,
		OrgInput: model.OrgInput{
			Name:    test.TestingOrgName,
			City:    test.TestingOrgCity,
			Address: test.TestingOrgAddress,
			LogoURL: nil,
		},
	})

	// assertions
	assert.Equal(suite.T(), model.OrgOutput{
		ID:      org.ID,
		Name:    test.TestingOrgName,
		City:    test.TestingOrgCity,
		Address: test.TestingOrgAddress,
		LogoURL: nil,
	}, *org)
	assert.Nil(suite.T(), gqlErrors)
	assert.Nil(suite.T(), err)
	assert.Len(suite.T(), suite.hook.AllEntries(), 0)
	//employeeClient.Get...
}
