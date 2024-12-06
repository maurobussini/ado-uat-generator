package sdks

import (
	"encoding/base64"
	"strconv"
	"strings"

	"zenprogramming.it/ado-uat-generator/utils"
)

//const WORKITEM_TYPE_USER_ACCEPTANCE_TESTS = "User Acceptance Tests"

type AzureDevOpsServiceConfiguration struct {
	TenantName  string
	ProjectName string
	ApiVersion  string
	ApiKey      string
	ApiUsername string
}

func (config AzureDevOpsServiceConfiguration) generateUrl(
	relativeUrl string) string {

	var baseUrl = "https://dev.azure.com/" +
		config.TenantName +
		"/" +
		config.ProjectName +
		"/"

	var fullUrl = baseUrl + relativeUrl

	if strings.Contains(fullUrl, "?") {
		fullUrl = fullUrl + "&" + config.ApiVersion
	} else {
		fullUrl = fullUrl + "?" + config.ApiVersion
	}

	return fullUrl
}

func (config AzureDevOpsServiceConfiguration) getBasicAuth() string {

	auth := config.ApiUsername +
		":" +
		config.ApiKey

	b64 := base64.StdEncoding.EncodeToString([]byte(auth))

	return "Basic " + b64
}

// Creates settings used for invoke Azure DevOps APIs
func CreateSettings(tenantName string,
	projectName string,
	apiKey string) AzureDevOpsServiceConfiguration {

	return AzureDevOpsServiceConfiguration{
		TenantName:  tenantName,
		ProjectName: projectName,
		ApiVersion:  "api-version=7.1",
		ApiKey:      apiKey,
		ApiUsername: "fakeuser",
	}
}

func GetWorkItem(config AzureDevOpsServiceConfiguration,
	id int) (
	WorkItemDetailsResponse, error) {

	var url = config.generateUrl("_apis/wit/workitems/" +
		strconv.Itoa(id) +
		"?$expand=relations")

	return utils.HttpGet[WorkItemDetailsResponse](
		url,
		config.getBasicAuth())
}

func CreateWorkItem(config AzureDevOpsServiceConfiguration,
	workItemType string,
	request *[]IWorkItemFieldRequest) (
	WorkItemDetailsResponse, error) {

	var url = config.generateUrl("_apis/wit/workitems/" +
		"$" + workItemType)

	return utils.HttpPost[[]IWorkItemFieldRequest, WorkItemDetailsResponse](
		url,
		config.getBasicAuth(),
		request)
}

func UpdateWorkItem(config AzureDevOpsServiceConfiguration,
	id int,
	request *[]IWorkItemFieldRequest) (
	WorkItemDetailsResponse, error) {

	var url = config.generateUrl("_apis/wit/workitems/" +
		strconv.Itoa(id))

	return utils.HttpPatch[[]IWorkItemFieldRequest, WorkItemDetailsResponse](
		url,
		config.getBasicAuth(),
		request)
}
