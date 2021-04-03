# Piyush
## Description
This golang application gets and sets IAM policies on folders in an organization on Google Cloud Platform.
## Command line:
* piyush -set -user="user:test@gmail.com" -role="roles/resourcemanager.folderEditor" -folder="folders/1098214482154" -org="organizations/27464139858"
* piyush -set -overwrite -user="user:test@gmail.com" -role="roles/resourcemanager.folderEditor" -folder="folders/1098214482154" -org="organizations/27464139858"
* piyush -get -org="organizations/27464139858"


# Piyush2
## Description
This golang application enumerates projects in an organization on Google Cloud Platform and displays them with their parent node.
## Command line:
piyush2
