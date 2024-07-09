resource "akeneo_attribute_group" "general" {
  code       = "general_data"
  sort_order = 10
  labels = {
    en_US = "General"
  }
}

resource "akeneo_attribute" "name" {
  code  = "name"
  type  = "pim_catalog_text"
  group = akeneo_attribute_group.general.code
  labels = {
    en_US = "Name"
  }
  localizable = true
}

resource "akeneo_attribute" "description" {
  code  = "description"
  type  = "pim_catalog_textarea"
  group = akeneo_attribute_group.general.code
  labels = {
    en_US = "Description"
  }
  localizable = true
  scopable    = true
}

resource "akeneo_attribute" "brand" {
  code  = "brand"
  type  = "pim_catalog_text"
  group = akeneo_attribute_group.general.code
  labels = {
    en_US = "Brand"
  }
}