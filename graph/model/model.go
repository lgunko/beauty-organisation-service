package model

import "github.com/lgunko/beauty-reuse/graph/model"

type IOrgInput interface {
	ToOutput(id string) OrgOutput
	GetOrgInput() OrgInput
}

type OrgInput struct {
	Name    string  `json:"name"`
	City    string  `json:"city"`
	Address string  `json:"address"`
	LogoURL *string `json:"logoUrl"`
}

type InitialOrgInput struct {
	CreatorName    string `json:"creatorName"`
	CreatorSurname string `json:"creatorSurname"`
	OrgInput
}

type OrgOutput struct {
	ID      string     `json:"id" bson:"_id"`
	Name    string     `json:"name"`
	City    string     `json:"city"`
	Address string     `json:"address"`
	LogoURL *string    `json:"logoUrl"`
	Role    model.Role `json:"role"`
}

func (input OrgInput) GetOrgInput() OrgInput {
	return input
}

func (input OrgInput) ToOutput(id string) OrgOutput {
	return OrgOutput{
		ID:      id,
		Name:    input.Name,
		City:    input.City,
		Address: input.Address,
		LogoURL: input.LogoURL,
	}
}

func (input InitialOrgInput) GetOrgInput() OrgInput {
	return input.OrgInput
}

func (input InitialOrgInput) ToOutput(id string) OrgOutput {
	return input.OrgInput.ToOutput(id)
}

func (eo OrgOutput) GetID() string {
	return eo.ID
}
