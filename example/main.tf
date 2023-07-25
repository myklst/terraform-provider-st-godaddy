terraform {
  required_providers {
    godaddy = {
      source = "n3integration/godaddy"
      version = "2.2.0"
    }
  }
}


provider "godaddy" {
  baseurl = "https://api.ote-godaddy.com"
  key = "3mM44UdB63ixBA_tSx4tP52257DiFPnjutMU"
  secret = "P4KptM8XDHWhNoDUMjQiX4"

  admin_email = "XXX"
  admin_fax = "XXX"
  admin_jobtitle = "XXX"
  admin_namelast = "XXX"
  admin_namefirst = "XXX"
  admin_namemiddle = "XXX"
  admin_organization = "XXX"
  admin_phone = "XXX"
  admin_address = "XXX"
  admin_city = "XXX"
  admin_country = "XXX"
  admin_state = "XXX"
  admin_postcode = "XXX"

  billing_email = "XXX"
  billing_fax = "XXX"
  billing_jobtitle = "XXX"
  billing_namelast = "XXX"
  billing_namefirst = "XXX"
  billing_namemiddle = "XXX"
  billing_organization = "XXX"
  billing_phone = "XXX"
  billing_address = "XXX"
  billing_city = "XXX"
  billing_country = "XXX"
  billing_state = "XXX"
  billing_postcode = "XXX"

  reg_email = "XXX"
  reg_fax = "XXX"
  reg_jobtitle = "XXX"
  reg_namelast = "XXX"
  reg_namefirst = "XXX"
  reg_namemiddle = "XXX"
  reg_organization = "XXX"
  reg_phone = "XXX"
  reg_address = "XXX"
  reg_city = "XXX"
  reg_country = "XXX"
  reg_state = "XXX"
  reg_postcode = "XXX"

  tech_email = "XXX"
  tech_fax = "XXX"
  tech_jobtitle = "XXX"
  tech_namelast = "XXX"
  tech_namefirst = "XXX"
  tech_namemiddle = "XXX"
  tech_organization = "XXX"
  tech_phone = "XXX"
  tech_address = "XXX"
  tech_city = "XXX"
  tech_country = "XXX"
  tech_state = "XXX"
  tech_postcode = "XXX"

}

resource "godaddy_domain_record" "gd-fancy-domain" {
  domain   = "hohojiang.com"

}

