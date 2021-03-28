package client

import (
	"context"
	"encoding/json"
	"github.com/lgunko/beauty-organisation-service/graph/model"
	"github.com/lgunko/beauty-reuse/restyRequest"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func GetAllowedOrgList(ctx context.Context, orgServiceURI string) ([]model.OrgOutput, []*gqlerror.Error, error) {
	req := `
    	query {
			AllowedOrgList{
    			id,
    			name,
    			city,
    			address,
    			logoUrl,
				role
  			}
		}
	`
	type Data struct {
		AllowedOrgList []model.OrgOutput `json:"AllowedOrgList"`
	}
	var respData Data
	rawJson, gqlErrors, err := sendQueryRequest(ctx, orgServiceURI, req)
	json.Unmarshal(rawJson, &respData)
	return respData.AllowedOrgList, gqlErrors, err
}

func RegisterOrg(ctx context.Context, orgServiceURI string, initialOrgInput model.InitialOrgInput) (*model.OrgOutput, []*gqlerror.Error, error) {
	req := `
    	mutation($input: InitialOrgInput!){
			RegisterOrg(input: $input){
    			id,
    			name,
    			city,
    			address,
    			logoUrl,
				role
  			}
		}
	`
	type Data struct {
		RegisterOrg *model.OrgOutput `json:"RegisterOrg"`
	}
	var respData Data
	rawJson, gqlErrors, err := sendMutationRequest(ctx, orgServiceURI, req, initialOrgInput)
	json.Unmarshal(rawJson, &respData)
	return respData.RegisterOrg, gqlErrors, err
}

type graphqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type graphqlResponse struct {
	Data   json.RawMessage   `json:"data"`
	Errors []*gqlerror.Error `json:"errors"`
}

func sendQueryRequest(ctx context.Context, uri, query string) (json.RawMessage, []*gqlerror.Error, error) {

	req := graphqlRequest{
		query,
		map[string]interface{}{},
	}

	// run it and capture the response
	jsonStr, _ := json.Marshal(req)

	ctx, resp, err := restyRequest.New(ctx).SetBody(string(jsonStr)).SetHeader("Content-Type", "application/json").Post(uri)
	if err != nil {
		return nil, nil, err
	}
	response := graphqlResponse{}
	json.Unmarshal(resp.Body(), &response)

	return response.Data, response.Errors, nil
}

func sendMutationRequest(ctx context.Context, uri, query string, initialOrgInput model.InitialOrgInput) (json.RawMessage, []*gqlerror.Error, error) {

	req := graphqlRequest{
		query,
		map[string]interface{}{
			"input": initialOrgInput,
		},
	}

	// run it and capture the response
	jsonStr, _ := json.Marshal(req)

	ctx, resp, err := restyRequest.New(ctx).SetBody(string(jsonStr)).SetHeader("Content-Type", "application/json").Post(uri)
	if err != nil {
		return nil, nil, err
	}
	response := graphqlResponse{}
	json.Unmarshal(resp.Body(), &response)

	return response.Data, response.Errors, nil
}
