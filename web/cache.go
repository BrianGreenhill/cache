package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	// driver for mysql
	"github.com/bradfitz/gomemcache/memcache"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"
)

var (
	memcachedHost = "10.0.0.23"
	memcachedPort = "11211"
	sqlHost       = "10.0.0.24"
	sqlPort       = "3306"
	sqlUser       = "root"
	sqlDB         = "fallback"
)

func main() {
	writer := os.Stdout
	logger := slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{}))
	logger.Info("Starting cache service")
	logger.Info("Connecting to memcached server")
	mc := memcache.New(memcachedHost + ":" + memcachedPort)
	if mc == nil {
		logger.Error("Error connecting to memcached server", "error", fmt.Errorf("memcached server not found"))
		os.Exit(1)
	}

	logger.Info("Connecting to database")
	db, err := newDB()
	if err != nil {
		logger.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}

	for {
		if err := run(context.Background(), writer, logger, mc, db); err != nil {
			logger.Error("Error running cache service", "error", err)
			os.Exit(1)
		}
		time.Sleep(5 * time.Second)
	}
}

func run(ctx context.Context, writer io.Writer, logger *slog.Logger, mc *memcache.Client, db *sql.DB) error {
	var cacheKey = "foo"
	var cacheValue = "my value"
	g := errgroup.Group{}

	g.Go(func() error {
		// try to get the item from memcached
		item, err := mc.Get(cacheKey)
		if err != nil {
			if err == memcache.ErrCacheMiss {
				logger.Info("Item not found in memcached, falling back to SQL")
			} else {
				return err
			}
		}

		// if item is found in memcached, delete it
		if item != nil {
			logger.Info("Item found", "key", item.Key, "value", string(item.Value))
			if err := mc.Delete(cacheKey); err != nil {
				return err
			}

			logger.Info("Item deleted")
			return nil
		}
		// if item is not found in memcached, fallback to database
		type dbi struct {
			Key   string
			Value string
		}

		var dbItem dbi

		row := db.QueryRowContext(ctx, "SELECT `key`, `value` FROM `keys` WHERE `key`= '"+cacheKey+"'")
		if err := row.Scan(&dbItem.Key, &dbItem.Value); err != nil {
			if err == sql.ErrNoRows {
				logger.Info("Item not found in database", "key", cacheKey, "value", cacheValue)
				logger.Info("Setting value in database")
				res, err := db.ExecContext(ctx, "INSERT INTO `keys` (`key`, `value`) VALUES ('"+cacheKey+"', '"+cacheValue+"')")
				if err != nil {
					return err
				}
				rows, err := res.RowsAffected()
				if err != nil {
					return err
				}
				if rows != 1 {
					logger.Error("Error inserting row", "rows", rows)
				}
			} else {
				return fmt.Errorf("error scanning row: %w", err)
			}
		}

		if row.Err() != nil {
			return fmt.Errorf("error scanning row: %w", row.Err())
		}

		if dbItem.Key != "" {
			logger.Info("Item found in database", "key", dbItem.Key, "value", dbItem.Value)
		}

		logger.Info("Caching item in memcached")
		if err := setItem(mc, cacheKey, cacheValue); err != nil {
			return err
		}

		return nil
	})

	return g.Wait()
}

func setItem(mc *memcache.Client, key, val string) error {
	expiration := time.Now().Add(24 * time.Hour).Unix()
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(val),
		Expiration: int32(expiration),
	}
	if err := mc.Set(item); err != nil {
		return err
	}
	return nil
}

// newDB creates a new database connection
// and provisions the database with a table and some data
func newDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", sqlUser+"@tcp("+sqlHost+":"+sqlPort+")/"+sqlDB)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(context.Background(), "USE "+sqlDB)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(context.Background(), "DROP TABLE IF EXISTS `keys`")
	if err != nil {
		return nil, err
	}

	create := "CREATE TABLE IF NOT EXISTS `keys` (id INT AUTO_INCREMENT PRIMARY KEY, `key` TEXT, `value` TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"
	_, err = db.ExecContext(context.Background(), create)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(context.Background(), "INSERT INTO `keys` (`key`, `value`) VALUES ('foo', 'value')")
	if err != nil {
		return nil, err
	}

	return db, nil
}
