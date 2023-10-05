resource "st-godaddy_domain" "domain-com" {
  domain = "sige-test11.com"
  purchase_years = 1
  min_days_remaining = 90
  contact = local.contact
}
