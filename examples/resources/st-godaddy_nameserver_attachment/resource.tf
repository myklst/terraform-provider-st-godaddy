resource "st-godaddy_nameserver_attachment" "some_custom_nameservers" {
  domain = "sige-test11.com"

  nameservers = [
    "ns7.alidns.com",
    "ns8.alidns.com"
  ]

  depends_on = [st-godaddy_domain.domain-com]
}
