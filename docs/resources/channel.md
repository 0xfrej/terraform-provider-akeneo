---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "akeneo_channel Resource - terraform-provider-akeneo"
subcategory: ""
description: |-
  Akeneo channel resource
---

# akeneo_channel (Resource)

Akeneo channel resource



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `category_tree` (String) Category tree assigned to the channel
- `code` (String) Channel code
- `currencies` (List of String) Currencies assigned to the channel
- `locales` (List of String) Locales assigned to the channel

### Optional

- `conversion_units` (Map of List of String) Converion units assigned to the chennel
- `labels` (Map of String) Label definition per locale
