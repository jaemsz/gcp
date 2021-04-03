package main

// BEFORE RUNNING:
// ---------------
// 1. If not already done, enable the Cloud Resource Manager API
//    and check the quota for your project at
//    https://console.developers.google.com/apis/api/cloudresourcemanager
// 2. This sample uses Application Default Credentials for authentication.
//    If not already done, install the gcloud CLI from
//    https://cloud.google.com/sdk/ and run
//    `gcloud beta auth application-default login`.
//    For more information, see
//    https://developers.google.com/identity/protocols/application-default-credentials
// 3. Install and update the Go dependencies by running `go get -u` in the
//    project directory.

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v3"
)

func test() {
}

func main() {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, cloudresourcemanager.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	cloudresourcemanagerService, err := cloudresourcemanager.New(c)
	if err != nil {
		log.Fatal(err)
	}

	// Search for projects.
	// NOTE:  Could not call List() because for some reason I do not have permission
	reqProj := cloudresourcemanagerService.Projects.Search()
	projSearch, err := reqProj.Do()
	if err != nil {
		log.Fatal(err)
	}

	// key is a string and the value is an array of strings
	projMap := make(map[string][]string)

	for _, project := range projSearch.Projects {
		if project.State == "ACTIVE" {
			// Add only active projects to the map
			projMap[project.Parent] = append(projMap[project.Parent], project.Name)
		}
	}

	for key, value := range projMap {
		fmt.Println("Parent: " + key)
		for _, proj := range value {
			fmt.Println("  Project: " + proj)
		}
		fmt.Println("")
	}
}
