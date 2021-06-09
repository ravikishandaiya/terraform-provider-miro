# Terraform Miro Provider

This terraform provider allows to perform Create ,Read ,Update, Delete and Import Miro User.


## Requirements
tested on following versions of tools.
* [Go](https://golang.org/doc/install) 1.16 <br>
* [Terraform](https://www.terraform.io/downloads.html) 0.15.x <br/>
* [Miro](https://miro.com/login/) Create account & create team (follow instructions mentioned in Miro API docs)


## Setup Miro Account
    :heavy_exclamation_mark: [IMPORTANT] : This Provider is tested on Windows 10, So all examples are accordingly.
    :heavy_exclamation_mark: [IMPORTANT] : This provider can be successfully tested on any Miro Free account. <br><br>


1. Create a Miro account. (https://miro.com/signup/)
2. And Sign in to the Miro.
3. Follow the link create app. (https://developers.miro.com/docs/getting-started#section-step-1-get-developer-team)
4. Agree the T&C and click on Create APP.
5. Set App name and OAuth scopes and Copy API token.


## Initialise Miro Provider in local machine 
1. Clone the repository to your local machine.<br>
2. Add the API token to respective fields in `main.tf` <br>
3. Add the team_id to respective fields in `main.tf` <br>
4. Run the following command :
 ```golang
go mod init terraform-provider-miro
go mod vendor
```
5. Vendor will create a vendor directory that contains all the provider's dependencies. <br>

## Installation
1. Run the following command to create a vendor subdirectory which will comprise of  all provider dependencies. <br>
```
~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
``` 
Command: 
```bash
mkdir -p ~/.terraform.d/plugins/hashicorp.com/edu/docusign/0.2.0/[OS_ARCH]
```
For eg. `mkdir -p ~/.terraform.d/plugins/hashicorp.com/edu/miro/0.1.0/windows_amd64`<br>

2. Run `go build -o terraform-provider-docusign.exe`. This will save the binary (`.exe`) file in the main/root directory. <br>
3. Run this command to move this binary file to appropriate location.
 ```
  move terraform-provider-docusign.exe %APPDATA%\terraform.d\plugins\hashicorp.com\edu\miro\0.1.0\windows_amd64
 ``` 
 Otherwise you can manually move the file from current directory to destination directory.<br>

 [OR]

1. Download required binaries <br>
2. move binary `~/.terraform.d/plugins/[architecture name]/`

## Working with Terraform
`terraform init` is used to initialize a working directory containing Terraform configuration files.
`terraform plan` command will create an execution plan, which will preview you the changes that Terraform plans to make to your infrastructure.
`terraform apply` will execute that plan and will make changes in infrastructure accordingly.

#### Create User
1. Add the user email type in the respective field in `main.tf`
2. Initialize the terraform provider `terraform init`
3. Check the changes applicable using `terraform plan` and apply using `terraform apply`
4. You will see that a user has been successfully invited and an invitation mail has been sent to the user.
5. Set the account using invitation.

#### Update the user
Update the data of the user in the `main.tf` file (in this case only role can be updated.)
And apply using `terraform apply`

#### Read the User Data
Add data and output blocks in the `main.tf` file and run `terraform apply` to read the user data.

#### Delete the user
Delete the resource block of the particular user from `main.tf` file and run `terraform apply`.

#### Import a User Data
1. Write manually a resource configuration block for the User in `main.tf`, to which the imported object will be mapped.
2. Run the command `terraform import miro_user.user1 [EMAIL_ID]`
3. Check for the attributes in the `.tfstate` file and fill them accordingly in resource block.


### Testing the Provider
1. Navigate to the test file directory.
2. Run command `go test`. This command will give combined test result for the execution or errors if any failure occur.
3. If you want to see test result of each test function individually while running test in a single go, run command `go test -v`
4. To check test cover run `go test -cover`


## Example Usage
```terraform
terraform {
  required_providers {
    miro = {
      version = "0.1.0"
      source  = "hashicorp.com/edu/miro"
    }
  }
}

provider "miro" {
  refresh token = ""
  accountid = ""
}

resource "miro_user" "user1" {
   email      = "[EMAIL_ID]"
   role       = "[ROLE]"
}

data "miro_user" "user2" {
  id = "[EMAIL_ID]"
}

output "user2" {
  value = data.miro_user.user2
}
```

