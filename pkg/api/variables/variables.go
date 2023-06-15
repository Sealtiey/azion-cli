package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/utils"

	sdk "github.com/aziontech/azionapi-go-sdk/variables"
)

type Client struct {
	apiClient *sdk.APIClient
}

type VariablesResponse interface {
	GetUuid() string
	GetKey() string
	GetValue() string
	GetSecret() bool
	GetLastEditor() string
}

func NewClient(c *http.Client, url string, token string) *Client {
	conf := sdk.NewConfiguration()
	conf.AddDefaultHeader("Authorization", "Token "+token)
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}
	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

func (c *Client) List(ctx context.Context) ([]VariablesResponse, error) {
	resp, httpResp, err := c.apiClient.VariablesApi.ApiVariablesList(ctx).Execute()

	if err != nil {
		fmt.Println(err.Error())
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	var result []VariablesResponse

	for i := range resp {
		result = append(result, &resp[i])
	}

	return result, nil
}
