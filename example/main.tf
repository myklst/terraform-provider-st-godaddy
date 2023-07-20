terraform {
  required_providers {
    godaddy = {
      source = "n3integration/godaddy"
      version = "2.0.0"
    }
  }
}


provider "godaddy" {
  baseurl = "https://api.ote-godaddy.com"
  key = "3mM44UdB63ixBA_tSx4tP52257DiFPnjutMU"
  secret = "P4KptM8XDHWhNoDUMjQiX4"
}

resource "godaddy_domain_record" "gd-fancy-domain" {
  domain   = "hohojiang.com"

}

