package services

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type ClickHouse struct {
	db DBProvider
}

type StatRecord struct {
	BytesWrittenGB         float64
	InfoHash, OriginalPath string
}

func NewClickHouse(c *cli.Context, db DBProvider) *ClickHouse {
	return &ClickHouse{
		db: db,
	}
}

func (s *ClickHouse) GetTopContent() ([]StatRecord, error) {
	db, err := s.db.Get()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get ClickHouse DB")
	}

	recs := []StatRecord{}

	rows, err := db.Query(`
		select * from (
			select * from (
				select infohash, original_path, sum(bytes_written) / 1024 / 1024 / 1024 as downloaded_gb from proxy_stat_all
				where edge = 'transcode-web-cache' and timestamp > now() - interval 3 hour
				group by infohash, original_path 
			) where downloaded_gb > 3
			union all
			select * from (
				select infohash, original_path as full_path, sum(bytes_written) / 1024 / 1024 / 1024 as downloaded_gb from proxy_stat_all
				where edge = 'nginx-vod' and timestamp > now() - interval 3 hour
				group by infohash, original_path 
			) where downloaded_gb > 3
		) order by downloaded_gb asc
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			bytesWrittenGB         float64
			infoHash, originalPath string
		)
		if err := rows.Scan(&infoHash, &originalPath, &bytesWrittenGB); err != nil {
			return nil, err
		}
		recs = append(recs, StatRecord{
			InfoHash:       infoHash,
			OriginalPath:   originalPath,
			BytesWrittenGB: bytesWrittenGB,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recs, nil
}
