# terraform-provider-carina
Carina provider for Terraform.  Can be used with the Terraform Docker provider to build Carina cluseters and then provision Docker swarms within the Carina cluster (with a small caveat see below).

## Installation

### Requirements

1. [Terraform](https://www.terraform.io/downloads.html). Make sure you have it installed and it's accessible from your `$PATH`.
2. [Carina Login](https://getcarina.com/) and API key.  

### Build From Source

* Install Go and [configure](https://golang.org/doc/code.html) your workspace.

* Download this repo:

```shell
$ go get github.com/sandpiper-io/terraform-provider-carina
```

* Compile and install it:

```shell
$ go install github.com/sandpiper-io/terraform-provider-carina
```

* Now let Terraform know where your plugin is by adding the following to your `~/.terraformrc` file:

```shell
providers {
    "carina" = "YOUR_GOPATH/bin/terraform-provider-carina"
}
```

## Usage

Here's a simple Terraform file to get you started:

```ruby
provider "carina" {
  // username and api key are in the environment
  // CARINA_USERNAME & CARINA_APIKEY
  // or can be added here as 
  // username & api_key
}

resource "carina_cluster" "test1" {
  name = "test1"
}
```

Save this to a `.tf` file and run:

```shell
$ terraform plan
$ terraform apply
$ terraform show
```

[Here](example.tf) is an example terraform file using that builds a carina cluster then builds a docker swarm inside of it based on the carina getting started example.

## Reference

### Provider

#### Example

```ruby
provider "carina" {
  username = "jim@example.com"
  api_key = "000011112222333344445555"
}
```

#### Parameters

* `username`: Optional. Directly provide your username.  Can also come from the environment variable CARINA_USERNAME.
* `api_key`: Optional. Directly provide your api_key.  Can also come from the environment variable CARINA_APIKEY.

### carina_cluster

This resource can be used to create a Carina cluster.  If you have a Docker provider that will control the Swarm inside of the cluster you must run terraform with the -target command first to create the carina cluster before the swarm is created.  using the example from above run the following:

```shell
$ terraform plan -target test1
$ terraform apply -target test1
```

Once the cluster is created you can run terraform without the -target parameter.

#### Example

```ruby
resource "carina_cluster" "test1" {
  name = "test1"
}
```

#### Parameters

* `name`: Required. The name of the cluster.
* `nodes`: Optional. The number of nodes in the cluster.


#### Exported Parameters

* `docker_host`: Docker host endpoint of the Carina node.  Can be handed to the Docker provider to control the Docker Swarm inside of the Carina node.  Note you will need to use -target on your first run of terraform if you plan to instantiate a Docker provider.
* `docker_certpath`: The directory containing the certs to connect to the Docker Swarm of the Carina node.  Can be handed to the Docker provider to control the Docker Swarm inside of the Carina node.

