package main

import (
	"backend/config"
	"backend/internal"
	"backend/internal/database"
	"context"
	"flag"
	"log"
)

func main() {
    // Parse command line flags
    migrate := flag.Bool("migrate", false, "Run database migrations")
    migrateDown := flag.Bool("migrate-down", false, "Revert database migrations")
    flag.Parse()

    // If migration flags are set, run migrations and exit
    if *migrate || *migrateDown {
        cfg, err := config.LoadConfig()
        if err != nil {
            log.Fatal("Failed to load configuration:", err)
        }

        db, err := config.InitDB(cfg)
        if err != nil {
            log.Fatal("Failed to initialize database:", err)
        }
        defer db.Close()

        migrator := database.NewMigrator(db)
        ctx := context.Background()

        if *migrateDown {
            if err := migrator.MigrateDown(ctx); err != nil {
                log.Fatal("Failed to revert migrations:", err)
            }
            log.Println("Successfully reverted all migrations")
        } else {
            if err := migrator.MigrateUp(ctx); err != nil {
                log.Fatal("Failed to apply migrations:", err)
            }
            log.Println("Successfully applied all migrations")
        }
        return
    }

    // Start the server normally if no migration flags
    internal.StartServer()
}