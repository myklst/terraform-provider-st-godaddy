terraform-provider-st-godaddy
=============================

A Terraform Provider for GoDaddy domain management.

## Prerequisites

First you'll need to apply for API access to GoDaddy. You can do that on
this [API admin page](https://developer.godaddy.com/getstarted).

Once you've done that, make note of the API key, your
username to fill into our `provider` block.

Supported Versions
------------------

| Terraform version | minimum provider version |maxmimum provider version
| ---- |--------------------------| ----|
| >= 1.3.x	| 0.1.0	                   | latest |

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 1.3.x
- [Go](https://golang.org/doc/install) 1.19 (to build the provider plugin)

Local Installation
------------------

1. Run make file `make install-local-custom-provider` to install the provider under ~/.terraform.d/plugins.

2. The provider source should be change to the path that configured in the *Makefile*:

    ```
    terraform {
      required_providers {
        st-godaddy = {
          source = "example.local/myklst/st-godaddy"
        }
      }
    }

    provider "st-godaddy" {
      baseurl = "XXX"
      key     = "XXX"
      secret  = "XXXX"
    }
    ```

Notes
-----
1. Changing purchase years will not affect past purchases. Will only affect future purchases
2. OTE environment behaviour -> Some country's IP cannot perform purchase and renew domain.
   All other API calls are allowed. A VPN is required to perform purchase and renew.
3. OTE name server behaviour is broken. Dont create ns_binding resource for OTE.
4. Domains purchased on OTE have the same lifecycle as production.
5. The same domain, once purchased on OTE, can still be purchased on production.
6. OTE is pre-funded with $10,000 credit.
