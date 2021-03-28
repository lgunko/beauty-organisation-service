package test

import (
	"context"
	"github.com/lgunko/beauty-organisation-service/client"
	"github.com/lgunko/beauty-organisation-service/graph/model"
	"github.com/stretchr/testify/assert"
)

func (suite *OrganizationServiceTestSuite) TestRegisterOrgWithoutHeaders() {
	assert.Empty(suite.T(), suite.hook.AllEntries())
	org, gqlErrors, err := client.RegisterOrg(context.Background(), "http://localhost:8080/api/query", model.InitialOrgInput{
		CreatorName:    "Leonid",
		CreatorSurname: "Gunko",
		OrgInput: model.OrgInput{
			Name:    "TestOrgName",
			City:    "TestOrgCity",
			Address: "TestOrgAddress",
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
		CreatorName:    "Leonid",
		CreatorSurname: "Gunko",
		OrgInput: model.OrgInput{
			Name:    "TestOrgName",
			City:    "TestOrgCity",
			Address: "TestOrgAddress",
			LogoURL: nil,
		},
	})

	// assertions
	assert.Equal(suite.T(), model.OrgOutput{
		ID:      org.ID,
		Name:    "TestOrgName",
		City:    "TestOrgCity",
		Address: "TestOrgAddress",
		LogoURL: nil,
	}, *org)
	assert.Nil(suite.T(), gqlErrors)
	assert.Nil(suite.T(), err)
	assert.Len(suite.T(), suite.hook.AllEntries(), 0)
}
