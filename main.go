package main

import (
	"flag"

	"zenprogramming.it/ado-uat-generator/flows"
	"zenprogramming.it/ado-uat-generator/sdks"
	"zenprogramming.it/ado-uat-generator/utils"
)

func main() {

	var tenantName string
	var projectName string
	var apiKey string
	var resultsFile string

	// Get tenant id, project name, API key and UAT Results file from commandline
	flag.StringVar(&tenantName, "t", "MY_TENANT", "Azure DevOps tenant name")
	flag.StringVar(&projectName, "p", "MY_PROJECT", "Azure DevOps project name")
	flag.StringVar(&apiKey, "k", "MY_API_KEY", "Azure DevOps API Key (PAT)")
	flag.StringVar(&resultsFile, "f", "my-uat-results.json", "JSON file with UAT results")
	flag.Parse()

	// Create settings for Azure DevOps API
	settings := sdks.CreateSettings(
		tenantName,
		projectName,
		apiKey)

	// Read content of UAT results file
	resultsData, err := utils.ReadUatResults(resultsFile)
	if err != nil {
		panic(err)
	}

	// Iterate all test results on source JSON
	for i := 0; i < len(resultsData); i++ {

		userStory, err := flows.GetUserStory(settings, resultsData[i].WorkItemId)
		if err != nil {
			panic(err)
		}

		flows.RenderWorkItemTitle(userStory)

		existingUat, err := flows.GetAttachedUserAcceptanceTests(settings, userStory)
		if err != nil {
			panic(err)
		}

		// Check if UAT does exists
		if existingUat.Id != 0 {

			// Render existing UAT
			flows.RenderUserAcceptanceTestAlreadyAvailable(existingUat)

			// Update existing UAT with execution
			updatedUat, err := flows.UpdateExistingUserAcceptanceTests(
				settings,
				existingUat,
				resultsData[i].ExecutionDate,
				resultsData[i].IsSuccess)

			if err != nil {
				panic(err)
			}

			flows.RenderUserAcceptanceTestUpdated(updatedUat)

			continue
		}

		newUat, err := flows.CreateAttachedUserAcceptanceTests(
			settings,
			userStory,
			resultsData[i].ExecutionDate,
			resultsData[i].IsSuccess)

		if err != nil {
			panic(err)
		}

		flows.RenderUserAcceptanceTestCreated(newUat)
	}
}
