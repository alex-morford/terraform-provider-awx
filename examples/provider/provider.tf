terraform {
  required_providers {
    awx = {
      source = "TravisStratton/awx"
    }
  }
}

provider "awx" {
  endpoint = "http://localhost:8078"
  token    = "awxtoken"
}
