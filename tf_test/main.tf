terraform {
  required_providers {
    clickhouse = {
      source  = "terraform.example.com/local/clickhouse"
      version = "~> 1.0.0"
    }
  }
}

provider "clickhouse" {
  database = "default"
  host = "localhost"
  port = 9000
  username = "default"
  password = ""
  timeout = 10
}

resource "clickhouse_user" "example" {
  name = "foo2"
}


resource "clickhouse_user" "example2" {
  name = "lol"
}
