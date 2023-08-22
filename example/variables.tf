


locals {
  contact = jsonencode(
    {
      "addressMailing": {
        "address1": "Street Ave. 666",
        "address2": "string",
        "city": "New City",
        "country": "US",
        "postalCode": "11-111",
        "state": "state of art"
      },
      "email": "john.doe@test-domain.com",
      "fax": "+48.111111111",
      "jobTitle": "XXX",
      "nameFirst": "John",
      "nameLast": "Doe",
      "nameMiddle": "XXX",
      "organization": "Corporation Inc.",
      "phone": "+48.111111111"
    }
  )
}