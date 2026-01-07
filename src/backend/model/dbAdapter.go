package model

// this is technical service for database connection
import (
	"context"
	"time"

	config "project/Config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var globalAdapterInstance *gorm.DB

type Adapter struct {
	db *gorm.DB
}

func (adapter *Adapter) newAdapter() {}

func (a *Adapter) WithTransaction(
	ctx context.Context,
	fn func(txAdapter *Adapter) error,
) error {

	return globalAdapterInstance.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&Adapter{db: tx})
	})
}

func (adapter *Adapter) GetAdapterIntance() *Adapter {
	if globalAdapterInstance == nil {
		cfg := config.LoadConfig()

		db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}

		globalAdapterInstance = db
		adapter = new(Adapter)
		adapter.newAdapter()

		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}

		// ✅ connection pool
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour)

	}
	return adapter
}

func (adater *Adapter) GetGormIntance() *gorm.DB {
	adater.GetAdapterIntance()
	return globalAdapterInstance
}
