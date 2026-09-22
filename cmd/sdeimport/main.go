package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/istvzsig/eve-gate-watch/external/sde"
	"github.com/istvzsig/eve-gate-watch/internal/data"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting SDE import")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("database connection: %v", err)
	}

	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("create data directory: %v", err)
	}

	sdePath := filepath.Join("data", "eve-sde-latest.zip")

	client := sde.NewClient(http.DefaultClient)

	if _, err := os.Stat(sdePath); os.IsNotExist(err) {
		log.Println("Downloading latest EVE SDE...")

		file, err := os.Create(sdePath)
		if err != nil {
			log.Fatalf("create SDE file: %v", err)
		}

		if err := client.Download(sde.LatestURL, file); err != nil {
			file.Close()
			log.Fatalf("download SDE: %v", err)
		}

		if err := file.Close(); err != nil {
			log.Fatalf("close SDE file: %v", err)
		}
	} else if err != nil {
		log.Fatalf("check SDE file: %v", err)
	} else {
		log.Println("SDE archive already exists")
	}

	repository := data.NewPostgresRepository(db)
	ctx := context.Background()

	log.Println("Importing regions...")

	regions := 0

	err = client.Regions(sdePath, func(region sde.Region) error {
		if err := repository.UpsertRegion(ctx, region); err != nil {
			return err
		}

		regions++
		return nil
	})

	if err != nil {
		log.Fatalf("import regions: %v", err)
	}

	log.Printf("Imported %d regions", regions)

	log.Println("Importing constellations...")

	constellations := 0

	err = client.Constellations(sdePath, func(constellation sde.Constellation) error {
		if err := repository.UpsertConstellation(ctx, constellation); err != nil {
			return err
		}

		constellations++
		return nil
	})

	if err != nil {
		log.Fatalf("import constellations: %v", err)
	}

	log.Printf("Imported %d constellations", constellations)

	log.Println("Importing solar systems...")

	systems := 0

	err = client.SolarSystems(sdePath, func(system sde.SolarSystem) error {
		if err := repository.UpsertSystem(ctx, system); err != nil {
			return err
		}

		systems++
		return nil
	})

	if err != nil {
		log.Fatalf("import solar systems: %v", err)
	}

	log.Printf("Imported %d solar systems", systems)

	log.Println("SDE import complete")
}
