resource "st-godaddy_domain" "domain" {
  domain             = "example.com"
  purchase_years     = 1
  min_days_remaining = 90
  contact = jsonencode(
    {
      addressMailing = {
        address1   = "1501 India Street",
        address2   = "",
        city       = "San Diego",
        country    = "US",
        postalCode = "92101",
        state      = "California"
      },
      email        = "john.doe@test-domain.com",
      fax          = "+48.111111111",
      jobTitle     = "CEO",
      nameFirst    = "John",
      nameLast     = "Doe",
      nameMiddle   = " ",
      organization = "Corporation Inc.",
      phone        = "+48.111111111"
    }
  )
}

resource "st-godaddy_nameserver_attachment" "ns" {
  domain = st-godaddy_domain.domain.domain

  nameservers = [
    "ns7.alidns.com",
    "ns8.alidns.com"
  ]
}
