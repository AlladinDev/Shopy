package main

import (
	config "UserService/Pkg/Config"
	mongodb "UserService/Pkg/Mongodb"
	rabbitmq "UserService/Pkg/Rabbitmq"
	shoplogsmodule "UserService/ShopLogsModule"
	shopmodule "UserService/ShopModule"
	suppliermodule "UserService/Suppliers"
	usermodule "UserService/UserModule"
	_ "net/http/pprof"
	"time"

	"github.com/gofiber/fiber/v2/middleware/pprof"

	"log"

	"github.com/go-playground/validator/v10"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	//load the envs
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a new Fiber app
	//set go-json as json library for fast performance
	app := fiber.New(fiber.Config{
		JSONEncoder:  gojson.Marshal,
		JSONDecoder:  gojson.Unmarshal,
		WriteTimeout: 5 * time.Second,  //this ensures max req handlers should take max of 6 seconds to process request after this req  will be cancelled
		IdleTimeout:  10 * time.Second, //this means each connection can wait for max 5 seconds for another request
	})

	// Initialize default config
	app.Use(pprof.New())

	//create an empty app config which can be shared in the app
	appConfig := config.AppConfig.InitialiseAppConfig()

	//connect to mongodb
	mongoDatabaseHandler, err := mongodb.ConnectToMongodb()

	if err != nil {
		log.Fatal(err)
	}

	//connect to rabbitmq
	rabbitMqChannel, rabbitConnection := rabbitmq.ConnectToRabbitMq()

	//##############  set configurations like mongodb rabbitmq which ever will be used by modules
	appConfig.SetMongodbDatabase(mongoDatabaseHandler)

	////amqps://ygiiabba:j0lIhtyQIfRTR6PDfipk4uUDWTLVCt3i@rabbit.lmq.cloudamqp.com/ygiiabba
	appConfig.SetRabbitMqChannel(rabbitMqChannel)

	appConfig.SetRabbitMqConnection(rabbitConnection)

	//########### set validator instance to appconfig
	validator := validator.New()
	appConfig.Validator = validator

	//assign this validator instance to this config level variables because in future we will remove di for config related things they can be direclty imported
	config.Validator = validator

	//close rabbitmq connection at last
	defer func() {
		if err := rabbitConnection.Close(); err != nil {
			log.Println("failed to close rabbit mq connection", err)
		}
	}()

	v1Api := app.Group("/v1")

	//#########  initialize the modules
	usermodule.InitialiseUserModule(appConfig, v1Api)
	shopmodule.InitialiseShopModule(appConfig, v1Api)
	shoplogsmodule.InitialiseShopLogsModule(appConfig, v1Api)
	suppliermodule.InitialiseSupplierModule(v1Api)

	// fi, err := os.Stat("./dukaandar.exe")
	// if err != nil {
	// 	fmt.Println("Error:", err)

	// }
	// fmt.Printf("Binary size: %d bytes (%.2f MB)\n", fi.Size(), float64(fi.Size())/(1024*1024))

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
