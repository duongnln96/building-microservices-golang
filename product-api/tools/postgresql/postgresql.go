package tools

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type RowFuncHandler func(*sql.Row) error
type RowsFuncHandler func(*sql.Rows) bool

type PsqlConnectorI interface {
	Start()
	Close() error
	DBExec(string, ...interface{}) error
	DBQueryRow(RowFuncHandler, string, ...interface{}) error
	DBQueryRows(RowsFuncHandler, string, ...interface{}) error
}

type PsqlConnectorDeps struct {
	Log *zap.SugaredLogger
	Ctx context.Context
	Cfg PsqlConnectorConfig
}

type psqlConnector struct {
	log      *zap.SugaredLogger
	ctx      context.Context
	cfg      PsqlConnectorConfig
	db       *sql.DB
	isOK     bool
	stopChan chan bool
}

func NewPsqlConnector(deps PsqlConnectorDeps) PsqlConnectorI {
	dbConfig := deps.Cfg.GetPsqlConfig()
	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		deps.Log.Fatalf("Cannot establish the connection to database %s. Error %+v", dbConfig, err)
	}

	ctx, cancel := context.WithTimeout(deps.Ctx, deps.Cfg.Timeout)
	if err := db.PingContext(ctx); err != nil {
		deps.Log.Fatalf("Cannot ping to database %+v", err)
	}
	defer cancel()

	return &psqlConnector{
		log:      deps.Log,
		ctx:      deps.Ctx,
		cfg:      deps.Cfg,
		db:       db,
		isOK:     true,
		stopChan: make(chan bool),
	}
}

func (psql *psqlConnector) healthCheck() {
	go func() {
		ticker := time.NewTicker(psql.cfg.HealthCheckInterval)
		for {
			select {
			case <-psql.stopChan:
				return
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(psql.ctx, psql.cfg.Timeout)
				if err := psql.db.PingContext(ctx); err != nil {
					psql.isOK = false
				} else {
					psql.isOK = true
				}
				cancel()
			}
		}
	}()
}

func (psql *psqlConnector) connIsOK() bool {
	return psql.isOK
}

func (psql *psqlConnector) Start() {
	psql.healthCheck()
}

func (psql *psqlConnector) Close() error {
	psql.stopChan <- true
	if err := psql.db.Close(); err != nil {
		return err
	}
	return nil
}

func (psql *psqlConnector) DBExec(cmd string, args ...interface{}) error {
	if !psql.connIsOK() {
		return fmt.Errorf("Connection to database not ok")
	}

	ctx, cancel := context.WithTimeout(psql.ctx, psql.cfg.Timeout)
	defer cancel()
	_, err := psql.db.ExecContext(ctx, cmd, args...)

	return err
}

func (psql *psqlConnector) DBQueryRow(handleFunc RowFuncHandler, cmd string, args ...interface{}) error {
	if !psql.connIsOK() {
		return fmt.Errorf("Connection to database not ok")
	}

	ctx, cancel := context.WithTimeout(psql.ctx, psql.cfg.Timeout)
	defer cancel()

	return handleFunc(psql.db.QueryRowContext(ctx, cmd, args...))
}

func (psql *psqlConnector) DBQueryRows(handleFunc RowsFuncHandler, cmd string, args ...interface{}) error {
	if !psql.connIsOK() {
		return fmt.Errorf("Connection to database not ok")
	}

	ctx, cancel := context.WithTimeout(psql.ctx, psql.cfg.Timeout)
	defer cancel()

	rows, err := psql.db.QueryContext(ctx, cmd, args...)
	if err != nil {
		return err
	}

	for rows.Next() {
		if !handleFunc(rows) {
			return fmt.Errorf("Error while processing multiple rows from database")
		}
	}

	return nil
}
