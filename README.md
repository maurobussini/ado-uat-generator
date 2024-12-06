# Azure DevOps User Acceptance Tests Generator

This is a simple program writteng in "Go Lang" that generates automatically `User Acceptance Tests` workitems, applying a relation with reference `User Story`, using results contained in a `uat-result.json` file.

## Usage
```sh

# Run application agains your Azure DevOps subscription
go run . -t TENANT_ID -p PROJECT_NAME -k API_KEY -u UAT_RESULTS_JSON_FILE
```
Given the sample project on Azure DevOps services `https://dev.azure.com/foo/bar/`:

- `TENANT_ID`: Unique identifier of your subscription (ex. `foo`)
- `PROJECT_ID`: Unique identifier of the project on the tenant (ex. `bar`)
- `API_KEY`: Personal Access Token (PAT) generated in Azure DevOps > `User Settings` > `Personal Access Tokens`
- `UAT_RESULTS_JSON_FILE`: File with tests results (same folder of the executable)

## UAT Results file sample
```json
[
    {
        "workItem": 1234,
        "executionDate": "2024-12-06T13:22:00Z",
        "isSuccess": true
    },
    {
        "workItem": 5678,
        "executionDate": "2024-12-05T09:33:12Z",
        "isSuccess": false
    }
]
```
