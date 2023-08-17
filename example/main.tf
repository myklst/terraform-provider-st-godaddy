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

  billing_email = "john.doe@test-domain.com"
  billing_fax = "+48.111111111"
  billing_jobtitle = "XXX"
  billing_namelast = "Doe"
  billing_namefirst = "John"
  billing_namemiddle = "XXX"
  billing_organization = "Corporation Inc."
  billing_phone = "+48.111111111"
  billing_address = "Street Ave. 666"
  billing_city = "New City"
  billing_country = "PL"
  billing_state = "state of art"
  billing_postcode = "11-111"

  reg_email = "john.doe@test-domain.com"
  reg_fax = "+48.111111111"
  reg_jobtitle = "XXX"
  reg_namelast = "Doe"
  reg_namefirst = "John"
  reg_namemiddle = "XXX"
  reg_organization = "Corporation Inc."
  reg_phone = "+48.111111111"
  reg_address = "Street Ave. 666"
  reg_city = "New City"
  reg_country = "PL"
  reg_state = "state of art"
  reg_postcode = "11-111"

  tech_email = "john.doe@test-domain.com"
  tech_fax = "+48.111111111"
  tech_jobtitle = "XXX"
  tech_namelast = "Doe"
  tech_namefirst = "John"
  tech_namemiddle = "XXX"
  tech_organization = "Corporation Inc."
  tech_phone = "+48.111111111"
  tech_address = "Street Ave. 666"
  tech_city = "New City"
  tech_country = "PL"
  tech_state = "state of art"
  tech_postcode = "11-111"

}


resource "godaddy_domain" "gd-fancy-domain" {
  domain   = "test-domain.com"

}


resource "godaddy_domain_record" "gd-fancy-domain-record" {

  domain   = "test-domain.com"

  record {
    name = "www"
    type = "CNAME"
    data = "fancy.github.io"
    ttl = 3600
  }

  // specify any A records associated with the domain
  addresses   = ["192.168.1.2", "192.168.1.3"]

  nameservers = ["ns7.domains.com", "ns6.domains.com"]

}

