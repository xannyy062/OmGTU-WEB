package connection

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateConnectionPool(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig("postgres://postgres:1234@localhost:5432/postgres")
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %v", err)
	}

	// Настройки пула соединений
	config.MaxConns = 20                      // максимальное количество соединений
	config.MinConns = 5                       // минимальное количество соединений
	config.MaxConnLifetime = time.Hour        // максимальное время жизни соединения
	config.MaxConnIdleTime = time.Minute * 30 // максимальное время простоя соединения
	config.HealthCheckPeriod = time.Minute    // периодичность проверки здоровья

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пула: %v", err)
	}

	// Проверка подключения
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("ошибка ping: %v", err)
	}

	fmt.Printf("✅ Пул соединений создан. Активных соединений: %d\n", pool.Stat().TotalConns())
	return pool, nil
}

// Старая функция для обратной совместимости
func CreateConnection(ctx context.Context) (*pgxpool.Pool, error) {
	return CreateConnectionPool(ctx)
}
