package sde

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

const LatestURL = "https://developers.eveonline.com/static-data/eve-online-static-data-latest-jsonl.zip"

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) Download(url string, destination io.Writer) error {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("download SDE: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download SDE: HTTP %s", resp.Status)
	}

	if _, err := io.Copy(destination, resp.Body); err != nil {
		return fmt.Errorf("copy SDE: %w", err)
	}

	return nil
}

func OpenJSONL(zipPath, filename string) (io.ReadCloser, error) {
	z, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("open SDE archive: %w", err)
	}

	for _, file := range z.File {
		if filepath.Base(file.Name) != filename {
			continue
		}

		reader, err := file.Open()
		if err != nil {
			z.Close()
			return nil, fmt.Errorf("open %s: %w", filename, err)
		}

		return &zipFileReader{
			ReadCloser: reader,
			zipReader:  z,
		}, nil
	}

	z.Close()

	return nil, fmt.Errorf("%s not found in SDE archive", filename)
}

type zipFileReader struct {
	io.ReadCloser
	zipReader *zip.ReadCloser
}

func (r *zipFileReader) Close() error {
	readerErr := r.ReadCloser.Close()
	zipErr := r.zipReader.Close()

	if readerErr != nil {
		return readerErr
	}

	return zipErr
}

type regionRecord struct {
	Key  int64 `json:"_key"`
	Name struct {
		EN string `json:"en"`
	} `json:"name"`
}

type constellationRecord struct {
	Key      int64 `json:"_key"`
	RegionID int64 `json:"regionID"`
	Name     struct {
		EN string `json:"en"`
	} `json:"name"`
}

type solarSystemRecord struct {
	Key             int64   `json:"_key"`
	ConstellationID int64   `json:"constellationID"`
	RegionID        int64   `json:"regionID"`
	SecurityStatus  float64 `json:"securityStatus"`
	WormholeClassID *int    `json:"wormholeClassID"`
	Name            struct {
		EN string `json:"en"`
	} `json:"name"`
}

func ScanJSONL[T any](reader io.Reader, fn func(T) error) error {
	scanner := bufio.NewScanner(reader)

	scanner.Buffer(
		make([]byte, 64*1024),
		16*1024*1024,
	)

	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}

		var record T

		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return fmt.Errorf("decode SDE JSONL: %w", err)
		}

		if err := fn(record); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan SDE JSONL: %w", err)
	}

	return nil
}

func (c *Client) Regions(zipPath string, fn func(Region) error) error {
	reader, err := OpenJSONL(zipPath, "mapRegions.jsonl")
	if err != nil {
		return err
	}
	defer reader.Close()

	return ScanJSONL(reader, func(record regionRecord) error {
		return fn(Region{
			ID:   record.Key,
			Name: record.Name.EN,
		})
	})
}

func (c *Client) Constellations(zipPath string, fn func(Constellation) error) error {
	reader, err := OpenJSONL(zipPath, "mapConstellations.jsonl")
	if err != nil {
		return err
	}
	defer reader.Close()

	return ScanJSONL(reader, func(record constellationRecord) error {
		return fn(Constellation{
			ID:       record.Key,
			Name:     record.Name.EN,
			RegionID: record.RegionID,
		})
	})
}

func (c *Client) SolarSystems(zipPath string, fn func(SolarSystem) error) error {
	reader, err := OpenJSONL(zipPath, "mapSolarSystems.jsonl")
	if err != nil {
		return err
	}
	defer reader.Close()

	return ScanJSONL(reader, func(record solarSystemRecord) error {
		return fn(SolarSystem{
			ID:              record.Key,
			Name:            record.Name.EN,
			ConstellationID: record.ConstellationID,
			RegionID:        record.RegionID,
			SecurityStatus:  record.SecurityStatus,
			WormholeClassID: record.WormholeClassID,
		})
	})
}
