---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "awx_inventory_source Data Source - awx"
subcategory: ""
description: |-
  Get inventory_source datasource
---

# awx_inventory_source (Data Source)

Get inventory_source datasource

## Example Usage

```terraform
data "awx_inventory_source" "example" {
  id = "1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Inventory Source ID.

### Read-Only

- `credential` (Number) Inventory source credential ID.
- `description` (String) InventorySource description.
- `enabled_value` (String) This field is ignored unless an Enabled Variable is set. If the enabled variable matches this value, the host will be enabled on import.
- `enabled_var` (String) Retrieve the enabled state from the given dict of host variables. The enabled variable may be specified using dot notation, e.g: 'foo.bar'
- `execution_environment` (Number) The ID of the execution environment this inventory source.
- `host_filter` (String) Regular expression where only matching host names will be imported. The filter is applied as a post-processing step after any inventory plugin filters are applied.
- `inventory` (Number) Inventory ID for the inventory source to be attached to.
- `name` (String) Inventory Source name.
- `overwrite` (Boolean) If checked, any hosts and groups that were previously present on the external source but are now removed will be removed from the inventory. Hosts and groups that were not managed by the inventory source will be promoted to the next manually created group or if there is no manually created group to promote them into, they will be left in the `all` default group for the inventory. When not checked, local child hosts and groups not found on the external source will remain untouched by the inventory update process.
- `overwrite_vars` (Boolean) If checked, all variables for child groups and hosts will be removed and replaced by those found on the external source. When not checked, a merge will be performed, combining local variables with those found on the external source.
- `scm_branch` (String) Branch to use on inventory sync. Project default used if blank. Only allowed if project allow_override field is set to true.
- `source` (String) Type of SCM resource. Options: `scm`, `ec2`, `gce`, `azure_rm`, `vmware`, `satellite6`, `openstack`, `rhv`, `controller`, `insights`, `terraform`, `openshift_virtualization`.
- `source_path` (String) (Inventory file) - The inventory file to be synced by this source.
- `source_project` (Number) The ID of the source project.
- `source_vars` (String) Default value is `"---"`
- `update_cache_timeout` (Number) Time in seconds to consider an inventory sync to be current. During job runs and callbacks the task system will evaluate the timestamp of the latest sync. If it is older than Cache Timeout, it is not considered current, and a new inventory sync will be performed.
- `update_on_launch` (Boolean) Each time a job runs using this inventory, refresh the inventory from the selected source before executing job tasks.
- `verbosity` (Number) Control the level of output Ansible will produce for inventory source update jobs. `0 - Warning`, `1 - Info`, `2 - Debug`