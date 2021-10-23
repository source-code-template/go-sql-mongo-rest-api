package app

import (
	"context"
	"github.com/core-go/health"
	"github.com/core-go/log"
	mgo "github.com/core-go/mongo"
	mq "github.com/core-go/mongo/query"
	"github.com/core-go/search"
	sv "github.com/core-go/service"
	v "github.com/core-go/service/v10"
	q "github.com/core-go/sql"
	"github.com/core-go/sql/query"
	_ "github.com/go-sql-driver/mysql"
	"reflect"

	"go-service/internal/usecase/user"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	UserHandler   user.UserHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	db, err := q.OpenByConfig(root.Sql)
	if err != nil {
		return nil, err
	}
	mongoDb, err := mgo.Setup(ctx, root.Mongo)
	if err != nil {
		return nil, err
	}
	logError := log.ErrorMsg
	status := sv.InitializeStatus(root.Status)

	var userRepository sv.Repository
	var searchUser func(context.Context, interface{}, interface{}, int64,...int64) (int64, string, error)
	userType := reflect.TypeOf(user.User{})
	if root.Provider != "mongo" {
		userQueryBuilder := query.NewBuilder(db, "users", userType)
		userSearchBuilder, err := q.NewSearchBuilder(db, userType, userQueryBuilder.BuildQuery)
		if err != nil {
			return nil, err
		}
		searchUser = userSearchBuilder.Search
		userRepository, err = q.NewRepository(db, "users", userType)
		if err != nil {
			return nil, err
		}
	} else {
		userQueryBuilder := mq.NewBuilder(userType)
		userSearchBuilder := mgo.NewSearchBuilder(mongoDb, "users", userQueryBuilder.BuildQuery, search.GetSort)
		searchUser = userSearchBuilder.Search
		userRepository = mgo.NewRepository(mongoDb, "users", userType)
	}
	userService := user.NewUserService(userRepository)
	validator := v.NewValidator()
	userHandler := user.NewUserHandler(searchUser, userService, status, validator.Validate, logError)

	sqlChecker := q.NewHealthChecker(db)
	mongoChecker := mgo.NewHealthChecker(mongoDb)
	healthHandler := health.NewHandler(sqlChecker, mongoChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}, nil
}
