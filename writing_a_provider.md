# Writing a Terraform Provider

Having just finished my first Terraform provider (and first Go program) I thought I would write down a few lessons learned from the experiance in hopes that it would help other people get the job done a little quicker. 
 
The goal of this project was to write a provider for the RackSpace Carina service.  I knew that once a Carina provider was created I could use the Docker provider to control the Swarm inside a Carina cluster.  Of course there was a catch... providers with interpolated values from a different provider don't respect the terraform dependency graph for some reason.  You can read more about that in the 

## Starting

I started with the [Hashi documentation](https://www.terraform.io/docs/plugins/index.html) which is a good place to start but there is no full example included showing the process of creating the provider.  

Next I looked into the [Terraform source](https://github.com/hashicorp/terraform) for an existing provider that would help.  The providers are hidden down in [builtin](https://github.com/hashicorp/terraform/tree/master/builtin) and I picked the [postgresql provider](https://github.com/hashicorp/terraform/tree/master/builtin/providers/postgresql) since it's simple.  This helped with the code side of things but not so much with the development environment setup.

Finally, I found the [LXC provider](https://github.com/jtopjian/terraform-provider-lxc) which talked about how to setup your environment for plugin development and gave a little more sophisticated example of a provider.  This example was super helpful!

## Development Environment Setup

The Hashi doc on [basics](https://www.terraform.io/docs/plugins/basics.html) plus the Go [How to Write Go Code](https://golang.org/doc/code.html) were very helpful.  My project was setup like this

```shell
$ export GOPATH=~/work/poc

# terraform binary is here the plugin will be too
$ export PATH=$PATH:$GOPATH/bin 

$ mkdir -p $GOPATH/src/github.com/sandpiper-io/terraform-provider-carina
$ cd $GOPATH/src/github.com/sandpiper-io/terraform-provider-carina

# main.go is here... rest of the code is in the carina package
$ ls
LICENSE			carina			main.go
README.md		example.tf		writing_a_provider.md

```

I setup the ~/.terraformrc file as shown in the basic example:

```shell
$ cat ~/.terraformrc 
providers {
    carina = "MYGOPATH/bin/terraform-provider-carina"
}
```

My main.go file is simple and looks like [this](main.go).  It imports the terraform plugin api and my package and passes the main entry point to my package back to terraform.  (The connection between a plugin and the terraform commandline is through http so the plugin library is doing some magic under the hood here.)

## Coding

The meat of the provider is in the [carina package](carina).  The [Provider](carina/provider.go#L10) (class?) passes back a structure when instantiated that defines its [parameters](carina/provider.go#L12), the [resources](carina/provider.go#L27) that it provides, and its [configuration function](carina/provider.go#L35).  The configuration function defined in [config.go](carina/config.go) is where the Carina API client is created.
 
Creating, refreshing, updating, and deleting a cluster is all handled in the [resource_carina_cluster.go](carina/resource_carina_cluster.go).  These entry are provided by the resource constructor and are well documented in the [Terraform docs for providers](https://www.terraform.io/docs/plugins/provider.html).

## Debugging

Make sure you setup a [test structure](carina/provider_test.go) for your provider as recommended in the documentation.  Once this is setup you can run the tests from the commandline with `go test`.  

```shell
$ cd $GOPATH/src/github.com/sandpiper-io/terraform-provider-carina
$ go test ./carina
ok  	github.com/sandpiper-io/terraform-provider-carina/carina	0.017s
```

You may also find that your plugin crashes the terraform commandline.  Turning on logging was helpful here so that you can add print statements to help isolate the problem.  The terraform docs have more on [debugging](https://www.terraform.io/docs/internals/debugging.html).

```shell
$ export TF_LOG=TRACE
```

I looked around for a good go debugger but wasn't able to set anything up quick enough for it to be worth it. 