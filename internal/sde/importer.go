package sde

import (
	"archive/zip"
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const LatestSDEURL = "https://developers.eveonline.com/static-data/eve-online-static-data-latest-jsonl.zip"

type Importer struct {
	db *sql.DB
}

func NewImporter(db *sql.DB) *Importer {
	return &Importer{db: db}
}

type solarSystem struct {
	Name            string  `json:"name"`
	ConstellationID int64   `json:"constellationID"`
	RegionID        int64   `json:"regionID"`
	SecurityStatus  float64 `json:"securityStatus"`
	WormholeClassID *int    `json:"wormholeClassID"`
}

type sdeRecord struct {
	Key   json.RawMessage `json:"_key"`
	Value solarSystem     `json:"_value"`
}

type tempFile struct {
	*os.File
}

func (i *Importer) Import(ctx context.Context) error {
	resp, err := http.Get(LatestSDEURL)
	if err != nil {
		return fmt.Errorf("download SDE: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"download SDE: unexpected HTTP status %s",
			resp.Status,
		)
	}

	tmpFile, err := createTempFile(resp.Body)
	if err != nil {
		return fmt.Errorf("save SDE: %w", err)
	}
	defer tmpFile.Close()

	fileInfo, err := tmpFile.Stat()
	if err != nil {
		return fmt.Errorf("stat SDE zip: %w", err)
	}

	reader, err := zip.NewReader(
		tmpFile,
		fileInfo.Size(),
	)
	if err != nil {
		return fmt.Errorf("open SDE zip: %w", err)
	}

	var systemsFile *zip.File

	for _, file := range reader.File {
		if file.Name == "mapSolarSystems.jsonl" {
			systemsFile = file
			break
		}
	}

	if systemsFile == nil {
		return fmt.Errorf("mapSolarSystems.jsonl not found in SDE")
	}

	return i.importSystems(ctx, systemsFile)
}

func (i *Importer) importSystems(
	ctx context.Context,
	file *zip.File,
) error {
	reader, err := file.Open()
	if err != nil {
		return fmt.Errorf("open systems file: %w", err)
	}
	defer reader.Close()

	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.PrepareContext(
		ctx,
		`
		INSERT INTO systems (
			id,
			name,
			constellation_id,
			region_id,
			security_status,
			wormhole_class_id,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (id)
		DO UPDATE SET
			name = EXCLUDED.name,
			constellation_id = EXCLUDED.constellation_id,
			region_id = EXCLUDED.region_id,
			security_status = EXCLUDED.security_status,
			wormhole_class_id = EXCLUDED.wormhole_class_id,
			updated_at = NOW()
		`,
	)
	if err != nil {
		return fmt.Errorf("prepare system upsert: %w", err)
	}
	defer stmt.Close()

	scanner := bufio.NewScanner(reader)

	// SDE JSONL records can be large.
	scanner.Buffer(
		make([]byte, 64*1024),
		4*1024*1024,
	)

	count := 0

	for scanner.Scan() {
		var record sdeRecord

		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return fmt.Errorf(
				"decode system record: %w",
				err,
			)
		}

		systemID, err := parseID(record.Key)
		if err != nil {
			return fmt.Errorf(
				"parse system ID: %w",
				err,
			)
		}

		_, err = stmt.ExecContext(
			ctx,
			systemID,
			record.Value.Name,
			record.Value.ConstellationID,
			record.Value.RegionID,
			record.Value.SecurityStatus,
			record.Value.WormholeClassID,
		)
		if err != nil {
			return fmt.Errorf(
				"insert system %d: %w",
				systemID,
				err,
			)
		}

		count++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read systems: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit systems: %w", err)
	}

	fmt.Printf("imported %d systems\n", count)

	return nil
}

func parseID(raw json.RawMessage) (int64, error) {
	var number int64

	if err := json.Unmarshal(raw, &number); err == nil {
		return number, nil
	}

	var stringValue string

	if err := json.Unmarshal(raw, &stringValue); err != nil {
		return 0, err
	}

	return strconv.ParseInt(stringValue, 10, 64)
}

func createTempFile(src io.Reader) (*tempFile, error) {
	return newTempFile(src)
}

func newTempFile(src io.Reader) (*tempFile, error) {
	file, err := os.CreateTemp("", "eve-sde-*.zip")
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(file, src); err != nil {
		file.Close()
		os.Remove(file.Name())
		return nil, err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		file.Close()
		os.Remove(file.Name())
		return nil, err
	}

	return &tempFile{File: file}, nil
}

func (f *tempFile) Close() error {
	name := f.Name()

	if err := f.File.Close(); err != nil {
		os.Remove(name)
		return err
	}

	return os.Remove(name)
}
