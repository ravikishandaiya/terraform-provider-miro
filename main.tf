terraform {
    required_providers {
        miro = {
            version = "0.1.0"
            source  = "hashicorp.com/edu/miro"
        }
    }
}

provider "miro" {
    token   = ""
    team_id = ""
}

resource "miro_user" "admin" {
    email      = ""
    role       = ""
}

output "admin" {
    value = miro_user.admin
}
