package main

import (
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app"
	"mahir-trade-be/internal/app/controller"
	"mahir-trade-be/internal/app/infra"
	"mahir-trade-be/internal/app/repo/discord"
	"mahir-trade-be/internal/app/repo/google"
	"mahir-trade-be/internal/app/repo/midtrans"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/internal/app/repo/sendgrid"
	"mahir-trade-be/internal/app/service"
	"mahir-trade-be/internal/app/service/utils"
	"mahir-trade-be/pkg/di"
	"mahir-trade-be/pkg/middleware"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	utils.InitValidator()

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

	err = di.Provide(infra.LoadJwtCfg)
	if err != nil {
		return fmt.Errorf("LoadJwtCfg: %s", err.Error())
	}

	err = di.Provide(infra.LoadGoogleCfg)
	if err != nil {
		return fmt.Errorf("LoadGoogleCfg: %s", err.Error())
	}

	err = di.Provide(infra.NewScheduler)
	if err != nil {
		fmt.Println("NewScheduler: ", err.Error())
		return fmt.Errorf("NewScheduler: %s", err.Error())
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

	err = di.Provide(infra.NewOauth)
	if err != nil {
		fmt.Println("NewOauth: ", err.Error())
		return fmt.Errorf("NewOauth: %s", err.Error())
	}

	return nil
}

func LoadApplicationRepository() error {
	err := di.Provide(postgres.NewUserRepo)
	if err != nil {
		return fmt.Errorf("NewMoneyTransferRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewGroupRepo)
	if err != nil {
		return fmt.Errorf("NewGroupRepo: %s", err.Error())
	}

	err = di.Provide(discord.NewDiscordRepo)
	if err != nil {
		return fmt.Errorf("NewDiscordRepo: %s", err.Error())
	}

	err = di.Provide(sendgrid.NewSendgridRepo)
	if err != nil {
		return fmt.Errorf("NewSendgridRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewModuleRepo)
	if err != nil {
		return fmt.Errorf("NewModuleRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewAdminRepo)
	if err != nil {
		return fmt.Errorf("NewAdminRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewDiscordAccountRepo)
	if err != nil {
		return fmt.Errorf("NewDiscordAccountRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewPackageRepo)
	if err != nil {
		return fmt.Errorf("NewPackageRepo: %s", err.Error())
	}

	err = di.Provide(google.NewGoogleRepo)
	if err != nil {
		return fmt.Errorf("NewGoogleRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewSubModuleRepo)
	if err != nil {
		return fmt.Errorf("NewSubModuleRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewReportRepo)
	if err != nil {
		return fmt.Errorf("NewReportRepo: %s", err.Error())
	}

	err = di.Provide(midtrans.NewMidtransRepo)
	if err != nil {
		return fmt.Errorf("NewMidtransRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewOrderRepo)
	if err != nil {
		return fmt.Errorf("NewOrderRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewTransactionRepo)
	if err != nil {
		return fmt.Errorf("NewTransactionRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewUserMembershipRepo)
	if err != nil {
		return fmt.Errorf("NewUserMembershipRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewGeneralLogRepo)
	if err != nil {
		return fmt.Errorf("NewGeneralLogRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewUserSubModuleRepo)
	if err != nil {
		return fmt.Errorf("NewUserSubModuleRepo: %s", err.Error())
	}

	err = di.Provide(google.NewBucketRepo)
	if err != nil {
		return fmt.Errorf("NewBucketRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewEmailTemplateRepo)
	if err != nil {
		return fmt.Errorf("NewEmailTemplateRepo: %s", err.Error())
	}

	err = di.Provide(postgres.NewConfigRepo)
	if err != nil {
		return fmt.Errorf("NewConfigRepo: %s", err.Error())
	}

	return nil
}

func LoadApplicationService() error {
	err := di.Provide(service.NewUserSvc)
	if err != nil {
		return fmt.Errorf("NewUserSvc: %s", err.Error())
	}

	err = di.Provide(service.NewGroupSvc)
	if err != nil {
		return fmt.Errorf("NewGroupSvc: %s", err.Error())
	}

	err = di.Provide(service.NewModuleSvc)
	if err != nil {
		return fmt.Errorf("NewModuleSvc: %s", err.Error())
	}

	err = di.Provide(service.NewAdminSvc)
	if err != nil {
		return fmt.Errorf("NewAdminSvc: %s", err.Error())
	}

	err = di.Provide(service.NewPackageSvc)
	if err != nil {
		return fmt.Errorf("NewPackageSvc: %s", err.Error())
	}

	err = di.Provide(service.NewSubModuleSvc)
	if err != nil {
		return fmt.Errorf("NewSubModuleSvc: %s", err.Error())
	}

	err = di.Provide(service.NewReportSvc)
	if err != nil {
		return fmt.Errorf("NewReportSvc: %s", err.Error())
	}

	err = di.Provide(service.NewPaymentSvc)
	if err != nil {
		return fmt.Errorf("NewPaymentSvc: %s", err.Error())
	}

	return nil
}

func LoadApplicationController() error {
	err := di.Provide(controller.NewAuthCtrl)
	if err != nil {
		return fmt.Errorf("NewAuthCtrl: %s", err.Error())
	}

	err = di.Provide(controller.NewGroupCtrl)
	if err != nil {
		return fmt.Errorf("NewGroupCtrl: %s", err.Error())
	}

	err = di.Provide(controller.NewModuleCtrl)
	if err != nil {
		return fmt.Errorf("NewModuleCtrl: %s", err.Error())
	}

	err = di.Provide(controller.NewAdminCtrl)
	if err != nil {
		return fmt.Errorf("NewAdminCtrl: %s", err.Error())
	}

	err = di.Provide(middleware.NewMiddleWare)
	if err != nil {
		return fmt.Errorf("NewMiddleWare: %s", err.Error())
	}
	err = di.Provide(controller.NewPackageCtrl)
	if err != nil {
		return fmt.Errorf("NewPackageCtrl: %s", err.Error())
	}

	err = di.Provide(controller.NewSubModuleCtrl)
	if err != nil {
		return fmt.Errorf("NewSubModuleCtrl: %s", err.Error())
	}

	err = di.Provide(controller.NewReportCtrl)
	if err != nil {
		return fmt.Errorf("NewReportCtrl: %s", err.Error())
	}

	err = di.Provide(controller.NewPaymentCtrl)
	if err != nil {
		return fmt.Errorf("NewPaymentCtrl: %s", err.Error())
	}

	err = di.Provide(controller.NewSchedulerCtrl)
	if err != nil {
		return fmt.Errorf("NewSchedulerCtrl: %s", err.Error())
	}

	return nil
}
