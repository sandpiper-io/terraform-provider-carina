provider "carina" {
  // username and api key are in the environment
}

resource "carina_cluster" "test" {
  name = "test1"
}

provider "docker" {
  host = "${carina_cluster.test.docker_host}"
  cert_path = "${carina_cluster.test.docker_cert_path}"
}

# Create the network for our system
resource "docker_network" "mynetwork" {
  name = "${var.network_name}"
  depends_on = ["carina_cluster.test"]
}

# Create, start, and initialize the database
resource "docker_container" "mysql" {
  image = "${docker_image.mysql.latest}"
  name = "mysql"
  networks = [
    "${docker_network.mynetwork.name}"]
  env = [
    "MYSQL_ROOT_PASSWORD=${var.mysql_root_password}",
    "MYSQL_DATABASE=${var.mysql_database}",
    "MYSQL_USER=${var.mysql_user}",
    "MYSQL_PASSWORD=${var.mysql_password}",
  ]
}

# initialize the database
resource "docker_container" "database_init" {
  image = "${docker_image.guestbook.latest}"
  name = "database_init"
  networks = [ "${docker_network.mynetwork.name}"]
  env = [
    "MYSQL_ROOT_PASSWORD=${var.mysql_root_password}",
    "MYSQL_DATABASE=${var.mysql_database}",
    "MYSQL_USER=${var.mysql_user}",
    "MYSQL_PASSWORD=${var.mysql_password}",
    "MYSQL_HOST=${var.mysql_host}",
    "MYSQL_PORT=${var.mysql_port}",
  ]
  command = ["python", "app.py", "create_tables"]
}

# Create and start the app server
resource "docker_container" "guestbook" {
  image = "${docker_image.guestbook.latest}"
  name = "guestbook"
  networks = [
    "${docker_network.mynetwork.name}"]
  env = [
    "MYSQL_HOST=${var.mysql_host}",
    "MYSQL_PORT=${var.mysql_port}",
    "MYSQL_ROOT_PASSWORD=${var.mysql_root_password}",
    "MYSQL_DATABASE=${var.mysql_database}",
    "MYSQL_USER=${var.mysql_user}",
    "MYSQL_PASSWORD=${var.mysql_password}",
  ]
  ports {
    internal = 5000,
    external = 5000
  }
}

resource "docker_image" "mysql" {
  name = "mysql:5.6"
}

resource "docker_image" "guestbook" {
  name = "carinamarina/guestbook-mysql"
}

output "host_ip" {
  value = "${carina_cluster.test.docker_host}"
}