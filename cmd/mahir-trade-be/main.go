package main

import (
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app"
	"mahir-trade-be/internal/app/infra"
	"mahir-trade-be/pkg/di"

	"github.com/joho/godotenv"
)

func main()  {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	err = infra.InitTimezone()
	if err != nil {
		slog.Error("Error loading timezone")
	}

	err = LoadApplicationConfig()
	if err != nil {
		fmt.Println("LoadApplicationConfig: ", err.Error())
		slog.Error(err.Error())
	}

	err = LoadApplicationPackage()
	if err != nil {
		fmt.Println("LoadApplicationPackage: ", err.Error())
		slog.Error(err.Error())
	}

	err = LoadApplicationRepository()
	if err != nil {
		fmt.Println("LoadApplicationRepository: ", err.Error())
		slog.Error(err.Error())
	}

	err = LoadApplicationService()
	if err != nil {
		fmt.Println("LoadApplicationService: ", err.Error())
		slog.Error(err.Error())
	}

	err = LoadApplicationController()
	if err != nil {
		fmt.Println("LoadApplicationController: ", err.Error())
		slog.Error(err.Error())
	}

	app.Start()
}

func LoadApplicationConfig() error {
	err := di.Provide(infra.LoadPgDatabaseCfg)
	if err != nil {
		return fmt.Errorf("LoadPgDatabaseCfg: %s", err.Error())
	}

	err = di.Provide(infra.LoadAppCfg)
	if err != nil {
		return fmt.Errorf("LoadAppCfg: %s", err.Error())
	}

	return nil
}

func LoadApplicationPackage() error {
	err := di.Provide(infra.NewEcho)
	if err != nil {
		return fmt.Errorf("NewEcho: %s", err.Error())
	}

	err = di.Provide(infra.NewDatabases)
	if err != nil {
		fmt.Println("NewDatabases: ", err.Error())
		return fmt.Errorf("NewDatabases: %s", err.Error())
	}

	return nil
}

func LoadApplicationRepository() error {
	// fill after create repo
	// Example
	//err := di.Provide(postgres.NewMoneyTransferRepo)
	//if err != nil {
	//	return fmt.Errorf("NewMoneyTransferRepo: %s", err.Error())
	//}

	return nil
}

func LoadApplicationService() error {
	// fill after create repo
	// Example
	//err := di.Provide(service.NewMoneyTransferSvc)
	//if err != nil {
	//	return fmt.Errorf("NewMoneyTransferService: %s", err.Error())
	//}

	return nil
}

func LoadApplicationController() error {
	// fill after create repo
	// Example
	//err := di.Provide(controller.NewMoneyTransferCtrl)
	//if err != nil {
	//	return fmt.Errorf("NewMoneyTransferController: %s", err.Error())
	//}

	return nil
}