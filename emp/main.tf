terraform {
  required_providers {
    st-godaddy = {
      source = "myklst/st-godaddy"
      version = "2.2.0"
    }
  }
}


provider "st-godaddy" {

  baseurl = "https://api.godaddy.com"
  key = "gHpoo7Q1ghjT_MYPpWKMo3H7ZishACKkhqY"
  secret = "WqAFKLsqtkvNLEJmUPTXgY"

  /*
  admin_email = "john.doe@test-domain.com"
  admin_fax = "+48.111111111"
  admin_jobtitle = "XXX"
  admin_namelast = "Doe"
  admin_namefirst = "John"
  admin_namemiddle = "XXX"
  admin_organization = "Corporation Inc."
  admin_phone = "+48.111111111"
  admin_address = "Street Ave. 666"
  admin_city = "New City"
  admin_country = "PL"
  admin_state = "state of art"
  admin_postcode = "11-111"
  */
}


resource "st-godaddy_domain" "test-domain" {

  domain   = "test-domain2.online"

  auto_renew_years = 1

  min_days_remaining = 90

  contact = local.contact

}


