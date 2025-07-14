package config

import (
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

var AppConfig Config

type Config struct {
	MongoDatabase      *mongo.Database
	RabbitMqConnection *amqp091.Connection
}

func (config *Config) InitialiseAppConfig() *Config {
	AppConfig = Config{}
	return &AppConfig
}

func (config *Config) SetMongodbDatabase(mongoDb *mongo.Database) {
	config.MongoDatabase = mongoDb
}

func (config *Config) SetRabbitMqConnection(rabbitConection *amqp091.Connection) {
	config.RabbitMqConnection = rabbitConection
}
