terraform {
  required_providers {
    st-godaddy = {
      source  = "myklst/st-godaddy"
      version = "~> 0.1"
    }
  }
}

provider "st-godaddy" {
  baseurl = "https://api.godaddy.com"
  key     = "XXX"
  secret  = "XXX"
}
