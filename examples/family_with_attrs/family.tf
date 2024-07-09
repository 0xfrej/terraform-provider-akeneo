resource "akeneo_family" "general" {
  code               = "general"
  attribute_as_label = akeneo_attribute.name.code
  attributes = [
    "sku",
    akeneo_attribute.description.code,
    akeneo_attribute.brand.code,
  ]
}