package cmdb

import (
	"encoding/csv"
	"errors"
	"io"
)

type ImportedAsset struct {
	UniqueKey  string         `json:"unique_key"`
	Attributes map[string]any `json:"attributes"`
}

func ParseCSVAssets(reader io.Reader) ([]ImportedAsset, error) {
	csvReader := csv.NewReader(reader)
	headers, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	uniqueKeyIndex := -1
	for index, header := range headers {
		if header == "unique_key" {
			uniqueKeyIndex = index
			break
		}
	}
	if uniqueKeyIndex == -1 {
		return nil, errors.New("unique_key column is required")
	}

	rows := make([]ImportedAsset, 0)
	for {
		record, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		row := ImportedAsset{
			UniqueKey:  record[uniqueKeyIndex],
			Attributes: make(map[string]any),
		}
		for index, header := range headers {
			if index == uniqueKeyIndex {
				continue
			}
			row.Attributes[header] = record[index]
		}
		rows = append(rows, row)
	}
	return rows, nil
}
