# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Output variable definitions


output "domain_name" {
  description = "Domain name to maintain"
  value       = st-godaddy_domain.domain-com.domain
}


