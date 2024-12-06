package flows

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"zenprogramming.it/ado-uat-generator/sdks"
)

func GetUserStory(config sdks.AzureDevOpsServiceConfiguration, id int) (sdks.WorkItemDetailsResponse, error) {

	var result sdks.WorkItemDetailsResponse
	const USER_STORY_TYPE = "User Story"

	// Get workitem with provided id
	data, err := sdks.GetWorkItem(config, id)

	// Check if recovery of workitem has success
	if err != nil {
		return result, errors.New("Error during recovery of workitem with id '" + strconv.Itoa(id) + "'.")
	}

	// Check if workitem is "User Story"
	if data.Fields.WorkItemType != USER_STORY_TYPE {
		return result, errors.New("Workitem with id '" +
			strconv.Itoa(id) +
			"' should be of type '" +
			USER_STORY_TYPE +
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

func RenderUserAcceptanceTestAlreadyAvailable(uatWorkItem sdks.WorkItemDetailsResponse) {

	fmt.Println("> UAT Already existing")
	fmt.Printf("  Id    : %v", uatWorkItem.Id)
	fmt.Println()
	fmt.Printf("  Title : %v", uatWorkItem.Fields.Title)
	fmt.Println()
}

func RenderUserAcceptanceTestUpdated(uatWorkItem sdks.WorkItemDetailsResponse) {

	fmt.Println("> UAT Existing Updated")
	fmt.Printf("  Id    : %v", uatWorkItem.Id)
	fmt.Println()
	fmt.Printf("  Title : %v", uatWorkItem.Fields.Title)
	fmt.Println()
}

func RenderUserAcceptanceTestCreated(uatWorkItem sdks.WorkItemDetailsResponse) {

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

func CreateAttachedUserAcceptanceTests(config sdks.AzureDevOpsServiceConfiguration,
	userStory sdks.WorkItemDetailsResponse,
	executionDate time.Time,
	isSuccess bool) (sdks.WorkItemDetailsResponse, error) {

	const WORKITEM_TYPE_USER_ACCEPTANCE_TESTS = "User Acceptance Tests"
	const UAT_TITLE = "Automatic UAT"

	const FORCE_ITERATION_PATH = "Zero\\Dev\\Sprint 134"

	var successString = renderSuccessString(isSuccess)
	var description = appendExecutionToDescription("", executionDate, isSuccess)

	var fields = []sdks.IWorkItemFieldRequest{
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.Title",
			Value: UAT_TITLE,
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.TeamProject",
			Value: userStory.Fields.TeamProject,
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.AreaPath",
			Value: userStory.Fields.AreaPath,
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.IterationPath",
			Value: FORCE_ITERATION_PATH, // TODO FOR DEVELOPMENT:  userStory.Fields.IterationPath,
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.State",
			Value: "Closed",
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/Custom.UATResults",
			Value: successString,
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/Custom.UATEndDate",
			Value: executionDate.Format(time.RFC3339),
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
				Rel: "System.LinkTypes.Related",
				Url: userStory.Url,
			},
		},
	}

	return sdks.CreateWorkItem(config, WORKITEM_TYPE_USER_ACCEPTANCE_TESTS, &fields)
}

func UpdateExistingUserAcceptanceTests(settings sdks.AzureDevOpsServiceConfiguration,
	existingUat sdks.WorkItemDetailsResponse,
	executionDate time.Time,
	isSuccess bool) (sdks.WorkItemDetailsResponse, error) {

	var successString = renderSuccessString(isSuccess)
	var description = appendExecutionToDescription(
		existingUat.Fields.Description,
		executionDate,
		isSuccess)

	var fields = []sdks.IWorkItemFieldRequest{
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.State",
			Value: "Closed",
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/Custom.UATResults",
			Value: successString,
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/Custom.UATEndDate",
			Value: executionDate.Format(time.RFC3339),
		},
		sdks.WorkItemAddPlainFieldRequest{
			Op:    "add",
			Path:  "/fields/System.Description",
			Value: description,
		},
	}

	return sdks.UpdateWorkItem(settings, existingUat.Id, &fields)
}

func GetAttachedUserAcceptanceTests(config sdks.AzureDevOpsServiceConfiguration,
	targetWorkItem sdks.WorkItemDetailsResponse) (sdks.WorkItemDetailsResponse, error) {

	const RELATION_TYPE = "System.LinkTypes.Related"
	const WORKITEM_TYPE_USER_ACCEPTANCE_TESTS = "User Acceptance Tests"

	var defaultValue sdks.WorkItemDetailsResponse

	// Iterate all relations of provided workitem
	for k := 0; k < len(targetWorkItem.Relations); k++ {

		// Check if relation is of type "Related". Parent, Child, etc.
		// should be excluded because are not of conventional type
		if targetWorkItem.Relations[k].Rel != RELATION_TYPE {
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
