terraform-provider-st-godaddy
===============================

A Terraform Provider for Godaddy domain management.

## Prerequisites

First you'll need to apply for API access to Godaddy. You can do that on
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

-	[Terraform](https://www.terraform.io/downloads.html) 1.3.x
-	[Go](https://golang.org/doc/install) 1.19 (to build the provider plugin)

Local Installation
------------------

1. Run make file `make install-local-custom-provider` to install the provider under ~/.terraform.d/plugins.

2. The provider source should be change to the path that configured in the *Makefile*:

    ```
    terraform {
      required_providers {
        st-alicloud = {
          source = "example.local/myklst/st-godaddy"
        }
      }
    }

    provider "st-godaddy" {
        baseurl   = "XXX"
        key    = "XXX"
        secret     = "XXXX"

    }
    ```
