// Contains utility functions around the database

package database

import (
	"archives/pkg/config"
	"archives/pkg/models"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// DBCon is the connection handle
// for the database
var (
	DBCon *pg.DB
)

// CreateSchema creates the tables in the database
// in case they don't alreay exist
func CreateSchema() error {
	if !tableExists("messages") {

		for _, model := range []interface{}{(*models.Message)(nil),
			(*models.MessageToReferences)(nil)} {

			err := DBCon.CreateTable(model, &orm.CreateTableOptions{
				IfNotExists: true,
			})
			if err != nil {
				return err
			}
		}

		// Add tsvector column for subjects
		DBCon.Exec("ALTER TABLE messages ADD COLUMN tsv_subject tsvector;")
		DBCon.Exec("CREATE INDEX subject_idx ON messages USING gin(tsv_subject);")

		// Add tsvector column for bodies
		DBCon.Exec("ALTER TABLE messages ADD COLUMN tsv_body tsvector;")
		DBCon.Exec("CREATE INDEX body_idx ON messages USING gin(tsv_body);")

		return nil
	}
	return nil
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

// AfterQuery is used to log SQL queries
func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	// logger.Debug.Println(q.FormattedQuery())
	return nil
}

// Connect is used to connect to the database
// and turn on logging if desired
func Connect() {
	DBCon = pg.Connect(&pg.Options{
		User:     config.PostgresUser(),
		Password: config.PostgresPass(),
		Database: config.PostgresDb(),
		Addr:     config.PostgresHost() + ":" + config.PostgresPort(),
	})

	DBCon.AddQueryHook(dbLogger{})

	orm.RegisterTable((*models.MessageToReferences)(nil))

	err := CreateSchema()
	if err != nil {
		// logger.Error.Println("ERROR: Could not create database schema")
		// logger.Error.Println(err)
	}

}

// utility methods

func tableExists(tableName string) bool {
	_, err := DBCon.Exec("select * from " + tableName + " limit 1;")
	return err == nil
}
