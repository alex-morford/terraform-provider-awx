---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "awx Provider"
subcategory: ""
description: |-
  Warning: All v0 releases are considered alpha and subject to breaking changes at any time.
---

# awx Provider

**Warning**: All v0 releases are considered alpha and subject to breaking changes at any time.

## Example Usage

```terraform
terraform {
  required_providers {
    awx = {
      source = "TravisStratton/awx"
    }
  }
}

provider "awx" {
  endpoint = "http://localhost:8078"
  token    = "awxtoken"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `endpoint` (String)
- `token` (String)
