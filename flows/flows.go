package flows

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"zenprogramming.it/ado-uat-generator/sdks"
)

const ADO_UAT_GENERATOR_VERSION = "1.0.0"
const WORKITEM_TYPE_USER_STORY_TYPE = "User Story"
const WORKITEM_TYPE_USER_ACCEPTANCE_TESTS = "User Acceptance Tests"
const UAT_TITLE = "[UAT] Automatic (" + ADO_UAT_GENERATOR_VERSION + ")"
const UAT_RELATION_TYPE_FORWARD = "Microsoft.VSTS.Common.TestedBy-Forward"
const UAT_RELATION_TYPE_REVERSE = "Microsoft.VSTS.Common.TestedBy-Reverse"

// const FORCE_ITERATION_PATH = "Zero\\Dev\\Sprint 134"
// const RELATION_TYPE = "System.LinkTypes.Related"
//const RELATION_TYPE = "Microsoft.VSTS.Common.TestedBy-Forward"

func GetUserStory(config sdks.AzureDevOpsServiceConfiguration, id int) (sdks.WorkItemDetailsResponse, error) {

	var result sdks.WorkItemDetailsResponse

	// Get workitem with provided id
	data, err := sdks.GetWorkItem(config, id)

	// Check if recovery of workitem has success
	if err != nil {
		return result, errors.New("Error during recovery of workitem with id '" + strconv.Itoa(id) + "'.")
	}

	// Check if workitem is "User Story"
	if data.Fields.WorkItemType != WORKITEM_TYPE_USER_STORY_TYPE {
		return result, errors.New("Workitem with id '" +
			strconv.Itoa(id) +
			"' should be of type '" +
			WORKITEM_TYPE_USER_STORY_TYPE +
			"', but is instead of type '" +
			data.Fields.WorkItemType +
			"'")
	}

	return data, nil
}

func RenderWorkItemTitle(workItem sdks.WorkItemDetailsResponse) {

	fmt.Println("****************************************************************")
	fmt.Printf("Processing workitem [%v] '%v'", workItem.Id, workItem.Fields.Title)
	fmt.Println()
}

func RenderExistingUserAcceptanceTest(uatWorkItem sdks.WorkItemDetailsResponse) {

	fmt.Println("> UAT Already existing")
	fmt.Printf("  Id    : %v", uatWorkItem.Id)
	fmt.Println()
	fmt.Printf("  Title : %v", uatWorkItem.Fields.Title)
	fmt.Println()
}

func RenderUpdatedUserAcceptanceTest(uatWorkItem sdks.WorkItemDetailsResponse) {

	fmt.Println("> UAT Existing Updated")
	fmt.Printf("  Id    : %v", uatWorkItem.Id)
	fmt.Println()
	fmt.Printf("  Title : %v", uatWorkItem.Fields.Title)
	fmt.Println()
}

func RenderCreatedUserAcceptanceTest(uatWorkItem sdks.WorkItemDetailsResponse) {

	fmt.Println("> UAT Created")
	fmt.Printf("  Id    : %v", uatWorkItem.Id)
	fmt.Println()
	fmt.Printf("  Title : %v", uatWorkItem.Fields.Title)
	fmt.Println()
}

func appendExecutionToDescription(currentDescription string,
	executionDate time.Time,
	isSuccess bool) string {

	var successString = renderSuccessString(isSuccess)
	var executionString = fmt.Sprintf("<p>History: %v => %v</p>", executionDate.Format(time.RFC3339), successString)

	if currentDescription == "" {
		return executionString
	} else {
		return currentDescription + executionString
	}
}

func renderSuccessString(isSuccess bool) string {
	if isSuccess {
		return "SUCCESS"
	} else {
		return "FAIL"
	}
}

func getUatFields(title string,
	teamProject string,
	areaPath string,
	iterationPath string,
	uATResults string,
	uATEndDate string,
	description string,
	userStoryUrl string) []sdks.IWorkItemFieldRequest {

	return []sdks.IWorkItemFieldRequest{
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.Title",
			Value: title,
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.TeamProject",
			Value: teamProject,
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.AreaPath",
			Value: areaPath,
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.IterationPath",
			Value: iterationPath,
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.State",
			Value: "Closed",
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/Custom.UATResults",
			Value: uATResults,
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/Custom.UATEndDate",
			Value: uATEndDate,
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.Description",
			Value: description,
		},
		sdks.WorkItemAddComplexFieldRequest{
			Op:   "add",
			Path: "/relations/-",
			Value: sdks.WorkItemAddRelationFieldRequest{
				Rel: UAT_RELATION_TYPE_FORWARD,
				Url: userStoryUrl,
			},
		},
	}
}

func CreateAttachedUserAcceptanceTests(config sdks.AzureDevOpsServiceConfiguration,
	userStory sdks.WorkItemDetailsResponse,
	executionDate time.Time,
	isSuccess bool) (sdks.WorkItemDetailsResponse, error) {

	var successString = renderSuccessString(isSuccess)
	var description = appendExecutionToDescription("", executionDate, isSuccess)

	var fields = getUatFields(
		UAT_TITLE,
		userStory.Fields.TeamProject,
		userStory.Fields.AreaPath,
		userStory.Fields.IterationPath,
		successString,
		executionDate.Format(time.RFC3339),
		description,
		userStory.Url)

	return sdks.CreateWorkItem(config, WORKITEM_TYPE_USER_ACCEPTANCE_TESTS, &fields)
}

func UpdateExistingUserAcceptanceTests(settings sdks.AzureDevOpsServiceConfiguration,
	existingUat sdks.WorkItemDetailsResponse,
	userStory sdks.WorkItemDetailsResponse,
	executionDate time.Time,
	isSuccess bool) (sdks.WorkItemDetailsResponse, error) {

	var successString = renderSuccessString(isSuccess)
	var description = appendExecutionToDescription(
		existingUat.Fields.Description,
		executionDate,
		isSuccess)

	var fields = getUatFields(
		UAT_TITLE,
		userStory.Fields.TeamProject,
		userStory.Fields.AreaPath,
		userStory.Fields.IterationPath,
		successString,
		executionDate.Format(time.RFC3339),
		description,
		userStory.Url)

	return sdks.UpdateWorkItem(settings, existingUat.Id, &fields)
}

func GetAttachedUserAcceptanceTests(config sdks.AzureDevOpsServiceConfiguration,
	targetWorkItem sdks.WorkItemDetailsResponse) (sdks.WorkItemDetailsResponse, error) {

	var defaultValue sdks.WorkItemDetailsResponse

	// Iterate all relations of provided workitem
	for k := 0; k < len(targetWorkItem.Relations); k++ {

		// Check if relation is of type "Related". Parent, Child, etc.
		// should be excluded because are not of conventional type
		if targetWorkItem.Relations[k].Rel != UAT_RELATION_TYPE_REVERSE {
			continue
		}

		userAcceptanceTestsId := getWorkItemIbByRelationUrl(targetWorkItem.Relations[k].Url)

		uatWorkItem, err := sdks.GetWorkItem(config, userAcceptanceTestsId)
		if err != nil {
			return defaultValue, err
		}

		// If current "related" is not a User Acceptance Tests, skip
		if uatWorkItem.Fields.WorkItemType != WORKITEM_TYPE_USER_ACCEPTANCE_TESTS {
			continue
		}

		// Otherwise
		return uatWorkItem, nil
	}

	// If a user acceptance tests workitem was not found
	// attached to target workitem , simply return empty element
	return defaultValue, nil
}

func getWorkItemIbByRelationUrl(relationUrl string) int {

	// Split relation url using "/"
	chunks := strings.Split(relationUrl, "/")

	if len(chunks) <= 0 {
		panic("Provided relation url '" + relationUrl + "' is not a valid url.")
	}

	idAsString := chunks[len(chunks)-1]
	id, _ := strconv.Atoi(idAsString)

	return id
}
