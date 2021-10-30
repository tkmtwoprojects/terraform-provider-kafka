

terraform {
  required_providers {
    connect = {
      source = "github.com/tkmtwoprojects/kafka"
      version = "0.1.0"
    }
  }
}


provider "connect" {
  url = "http://localhost:8083/"
}



resource "connect_connector" "datagenAlpha" {
  name = "DatagenUsersAlpha"
  config = {
    "name": "DatagenUsersAlpha"
		"connector.class": "io.confluent.kafka.connect.datagen.DatagenConnector",
		"key.converter":   "org.apache.kafka.connect.storage.StringConverter",
		"kafka.topic":     "datagen.users.Alpha",
		"max.interval":    "5000",
		"quickstart":      "users"
    }
}

resource "connect_connector" "datagenBravo" {
  name = "DatagenUsersBravo"
  config = {
    "name": "DatagenUsersBravo"
		"connector.class": "io.confluent.kafka.connect.datagen.DatagenConnector",
		"key.converter":   "org.apache.kafka.connect.storage.StringConverter",
		"kafka.topic":     "datagen.users.Bravo",
		"max.interval":    "5000",
		"quickstart":      "users"
    }
}

