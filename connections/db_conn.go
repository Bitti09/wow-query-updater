package connections

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
    log "github.com/sirupsen/logrus"
	"io/ioutil"
    "crypto/tls"
    "crypto/x509"
)

var dbConn *pg.DB
var ReportingMode bool

type dbLogger struct { }

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}

// Connect is used to create the Postgres connection pool
func Connect(username string, password string, db string, poolSize int, schema string, hostname string, port int) {
    CAFile := "root.crt"
    CACert, err := ioutil.ReadFile(CAFile)
    if err != nil {
        log.Errorf("failed to load server certificate: %v", err)
    }

    CACertPool := x509.NewCertPool()
	CACertPool.AppendCertsFromPEM(CACert)
	    tlsConfig := &tls.Config{
        RootCAs:            CACertPool,
		InsecureSkipVerify: true,
		// ServerName:         "localhost",
    }
	createSchemaStatement := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS \"%s\";", schema)
	useSchemaStatement := fmt.Sprintf("SET SEARCH_PATH = \"%s\";", schema)
	dbConn = pg.Connect(&pg.Options{
		Addr: fmt.Sprintf("%s:%d", hostname, port),
		User:     username,
		Password: password,
		Database: db,
		PoolSize: poolSize,
        TLSConfig: tlsConfig,
		OnConnect: func(ctx context.Context, conn *pg.Conn) error {
			_, err := conn.Exec(createSchemaStatement)
			if err != nil {
				log.Fatal(err)
			}
			_, err = conn.Exec(useSchemaStatement)
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	})
}

// Disconnect is used to disconnect all Postgres connections
func Disconnect() {
	err := dbConn.Close()
	if err != nil {
		log.Panic(err)
	}
}

// GetDBConn is a safer way to access the connection pool
func GetDBConn() *pg.DB {
	if dbConn != nil {
		return dbConn
	}
	log.Println("database connection not active yet")
	return nil
}
