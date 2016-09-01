provider "carina" {
    username = "YOUR_USERNAME_HERE"
    api_key = "YOUR_API_KEY_HERE"
}

resource "carina_cluster" "test" {
  name = "test"
}

provider "docker" {
  host = "${carina_cluster.test.docker_host}"
  cert_path = "${carina_cluster.test.docker_cert_path}"
}

module "mysql_example" {
    source = "github.com/sandpiper-io/terraform/carina//mysql"
}
