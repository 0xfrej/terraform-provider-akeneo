# Terraform Akeneo Provider

This provider aims to implement Akeneo catalog structure resources for easy management and keeping all of your environments in sync.
At the moment Akeneo does not support deletion of most of these types, so does this provider (upon such try, the provider will throw an error).

## Features

- Importing resources from Akeneo (possible, not tested)

### Currently supported resources

- Measurement Family
- Channel
- Family
- Family Variant
- Attribute
- Attribute Option
- Attribute Group

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make acctest
```
