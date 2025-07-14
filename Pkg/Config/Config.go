// Package config provides functions related to configuration like connecting to mongodb rabbitmq etc
package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

var AppConfig Config

var Validator *validator.Validate

type Config struct {
	MongoDatabase      *mongo.Database
	RabbitMqChannel    *amqp091.Channel
	Validator          *validator.Validate
	RabbitMqConnection *amqp091.Connection
}

var MongoDbDatabase *mongo.Database

var RabbitConnection *amqp091.Connection

func (config *Config) InitialiseAppConfig() *Config {
	AppConfig = Config{}
	return &AppConfig
}

func (config *Config) SetMongodbDatabase(mongoDB *mongo.Database) {
	config.MongoDatabase = mongoDB
	MongoDbDatabase = mongoDB
}

func (config *Config) SetRabbitMqChannel(rabbitChannel *amqp091.Channel) {
	config.RabbitMqChannel = rabbitChannel

}

func (config *Config) SetRabbitMqConnection(rabbitConn *amqp091.Connection) {
	config.RabbitMqConnection = rabbitConn
	RabbitConnection = rabbitConn
}
