package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type DBRowHandler func(*sql.Row) error
type DBRowsHandler func(*sql.Rows) bool

type DBConnectorI interface {
	Start()
	Stop() error
	DBExec(string, ...interface{}) error
	DBQueryRow(DBRowHandler, string, ...interface{}) error
	DBQueryRows(DBRowsHandler, string, ...interface{}) error
}

type PsqlConnectorDeps struct {
	Log zap.SugaredLogger
	Ctx context.Context
	Cfg PsqlConfig
}

type psqlConnector struct {
	log      zap.SugaredLogger
	ctx      context.Context
	cfg      PsqlConfig
	conn     *sql.DB
	isOK     bool
	stopChan chan bool
}

func NewPsqlConnector(deps PsqlConnectorDeps) DBConnectorI {
	dbInfo := deps.Cfg.GetPsqlInfo()
	conn, err := sql.Open("postgres", dbInfo)
	if err != nil {
		deps.Log.Fatalf("Cannot establish the connection to database %s. Error %+v", dbInfo, err)
	}

	ctx, cancel := context.WithTimeout(deps.Ctx, deps.Cfg.TimeOut)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		deps.Log.Fatalf("Cannot ping to database %+v", err)
	}

	return &psqlConnector{
		log:      deps.Log,
		ctx:      deps.Ctx,
		cfg:      deps.Cfg,
		conn:     conn,
		isOK:     true,
		stopChan: make(chan bool),
	}
}

func (p *psqlConnector) connIsOk() bool {
	return p.isOK
}

func (p *psqlConnector) healthCheck(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(p.ctx, p.cfg.TimeOut)
			if err := p.conn.PingContext(ctx); err != nil {
				p.isOK = false
			} else {
				p.isOK = true
			}
			cancel()
		case <-p.stopChan:
			close(p.stopChan)
			return
		}
	}
}

func (p *psqlConnector) Start() {
	go p.healthCheck(p.cfg.HealthCheckInterval)
}

func (p *psqlConnector) Stop() error {
	p.stopChan <- true
	if err := p.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (p *psqlConnector) DBExec(query string, args ...interface{}) error {
	if p.connIsOk() {
		return fmt.Errorf("The connection is not ok")
	}
	ctx, cancel := context.WithTimeout(p.ctx, p.cfg.TimeOut)
	defer cancel()
	_, err := p.conn.ExecContext(ctx, query, args...)

	return err
}

func (p *psqlConnector) DBQueryRow(handleFunc DBRowHandler, query string, args ...interface{}) error {
	if p.connIsOk() {
		return fmt.Errorf("The connection is not ok")
	}

	ctx, cancel := context.WithTimeout(p.ctx, p.cfg.TimeOut)
	defer cancel()

	return handleFunc(p.conn.QueryRowContext(ctx, query, args...))
}

func (p *psqlConnector) DBQueryRows(handleFunc DBRowsHandler, query string, args ...interface{}) error {
	if p.connIsOk() {
		return fmt.Errorf("The connection is not ok")
	}

	ctx, cancel := context.WithTimeout(p.ctx, p.cfg.TimeOut)
	defer cancel()

	rows, err := p.conn.QueryContext(ctx, query, args...)
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
