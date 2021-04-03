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
	"flag"
	"fmt"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v2"
)

func getFolders(ctx context.Context, req *cloudresourcemanager.FoldersListCall, foldersArray []string) []string {
	retArray := make([]string, 0)

	for _, folderItem := range foldersArray {
		if err := req.Parent(folderItem).Pages(ctx, func(page *cloudresourcemanager.ListFoldersResponse) error {
			for _, folder := range page.Folders {
				retArray = append(retArray, folder.Name)
			}
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}

	return retArray
}

/*
   Example command lines:
       piyush -set -user="user:test@gmail.com" -role="roles/resourcemanager.folderEditor" -folder="folders/345573146175" -org="organizations/27464139858"
       piyush -get -org="organizations/27464139858"
       piyush -projects
*/
func main() {
	// get flag needs the org
	getArg := flag.Bool("get", false, "Get folder IAM policies")
	// set flag needs the user,role,folder,org
	setArg := flag.Bool("set", false, "Set folder IAM policy")
	overwriteArg := flag.Bool("overwrite", false, "Overwrite IAM policy")

	// Command line arguments
	user := flag.String("user", "", "user:[user]")
	role := flag.String("role", "", "roles/[role]")
	folder := flag.String("folder", "", "folders/[folder ID]")
	org := flag.String("org", "", "organizations/[org ID]")
	flag.Parse()

	// Get context
	ctx := context.Background()

	// Get default client
	c, err := google.DefaultClient(ctx, cloudresourcemanager.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	// Get cloud resource manager service
	cloudresourcemanagerService, err := cloudresourcemanager.New(c)
	if err != nil {
		log.Fatal(err)
	}

	if *getArg {
		// Get the top level folders in the hierarchy
		req := cloudresourcemanagerService.Folders.List()

		initArray := make([]string, 0)
		initArray = append(initArray, *org)

		foldersArray := getFolders(ctx, req, initArray)
		if len(foldersArray) == 0 {
			fmt.Println("The specified organization does not have any folders")
			return
		}

		// Get the remaining folders in the hierarchy
		temp := foldersArray
		for {
			temp2 := getFolders(ctx, req, temp)
			if len(temp2) == 0 {
				break
			}

			foldersArray = append(foldersArray, temp2...)
			temp = temp2
		}

		// Get the IAM policies for each folder
		for _, folderItem := range foldersArray {
			iamRequestPolicy := cloudresourcemanager.GetIamPolicyRequest{}
			reqIAM := cloudresourcemanagerService.Folders.GetIamPolicy(folderItem, &iamRequestPolicy)

			iamPolicy, err := reqIAM.Do()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Bindings for " + folderItem + ":")
			for _, binding := range iamPolicy.Bindings {
				fmt.Println("- members:")
				for _, member := range binding.Members {
					fmt.Println("  - " + member)
				}
				fmt.Println("  role: " + binding.Role)
			}
			fmt.Println("")
		}
	}

	if *setArg {
		if *overwriteArg {
			// Set the IAM policy of a folder
			setIamPolicyRequest := cloudresourcemanager.SetIamPolicyRequest{}

			binding := cloudresourcemanager.Binding{}
			binding.Members = append(binding.Members, *user)
			binding.Role = *role

			policy := cloudresourcemanager.Policy{}
			policy.Version = 1
			policy.Bindings = append(policy.Bindings, &binding)

			setIamPolicyRequest.Policy = &policy

			reqSetIAM := cloudresourcemanagerService.Folders.SetIamPolicy(*folder, &setIamPolicyRequest)
			_, err = reqSetIAM.Do()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			iamRequestPolicy := cloudresourcemanager.GetIamPolicyRequest{}
			reqIAM := cloudresourcemanagerService.Folders.GetIamPolicy(*folder, &iamRequestPolicy)

			policy, err := reqIAM.Do()
			if err != nil {
				log.Fatal(err)
			}

			binding := cloudresourcemanager.Binding{}
			binding.Members = append(binding.Members, *user)
			binding.Role = *role

			policy.Bindings = append(policy.Bindings, &binding)

			setIamPolicyRequest := cloudresourcemanager.SetIamPolicyRequest{}
			setIamPolicyRequest.Policy = policy

			reqSetIAM := cloudresourcemanagerService.Folders.SetIamPolicy(*folder, &setIamPolicyRequest)
			_, err = reqSetIAM.Do()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
