package provider

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"strconv"
	"time"
)

type Config struct {
	Database string
	Host     string
	Port     int
	Username string
	Password string
	Timeout  int
}

func Connect(config *Config) (driver.Conn, error) {
	addr := fmt.Sprintf("%s:%s", config.Host, strconv.Itoa(config.Port))
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.Username,
			Password: config.Password,
		},
		Debug:           true,
		DialTimeout:     time.Second * time.Duration(config.Timeout),
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})

	if err != nil {
		return nil, err
	}

	ctx := clickhouse.Context(
		context.Background(),
		clickhouse.WithSettings(
			clickhouse.Settings{
				"max_block_size": 10,
			},
		),
		clickhouse.WithProgress(
			func(p *clickhouse.Progress) {
				fmt.Println("progress: ", p)
			},
		),
	)

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	return conn, nil
}
