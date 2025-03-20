package bootstrap

import (
	"fmt"
	"log"
	"my-graphql-project/core/logs"
	"my-graphql-project/domain/entities"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//Postgres connection
func NewSqlConnection(env *Env) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.Sql.Username,
		env.Sql.Password,
		env.Sql.DBHost,
		env.Sql.DBPort,
		env.Sql.DBName,
	)
	// logs.Info(dsn)
	// Set up a custom pool configuration
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("can not connect to database")
	}
	logs.Info("database connection success")

	Migrate(db)

	return db
}
//Postgres connection
func NewPostgresConnection(env *Env) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Vientiane",
		env.Postgres.DBHost,
		env.Postgres.Username,
		env.Postgres.Password,
		env.Postgres.DBName,
		env.Postgres.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("can not connect to database")
	}
	logs.Info("database connection success")

	// Migrate(db)

	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
	//Migrate entities
		&entities.User{},
	)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		logs.Error(err.Error())
	}
	logs.Info("Migrate successfully")
}
