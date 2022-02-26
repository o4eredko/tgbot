package sheets

import (
	"context"
	"fmt"

	"google.golang.org/api/sheets/v4"
)

type Service struct {
	spreadsheetID string
	sheetsClient  *sheets.Service
}

func NewService(spreadsheetID string, sheetsClient *sheets.Service) *Service {
	return &Service{spreadsheetID: spreadsheetID, sheetsClient: sheetsClient}
}

func (s *Service) GetHeaders(_ context.Context, tableName string) ([]string, error) {
	resp, err := s.sheetsClient.Spreadsheets.Values.Get(s.spreadsheetID, tableName).Do()
	if err != nil {
		return nil, fmt.Errorf("get spreadsheets: %w", err)
	}

	if len(resp.Values) == 0 {
		return nil, fmt.Errorf("empty google sheets response")
	}

	headers := make([]string, len(resp.Values[0]))
	for i, item := range resp.Values[0] {
		headers[i] = item.(string)
	}

	return headers, nil
}

func (s *Service) GetTableRows(_ context.Context, tableName string) ([]map[string]string, error) {
	resp, err := s.sheetsClient.Spreadsheets.Values.Get(s.spreadsheetID, tableName).Do()
	if err != nil {
		return nil, fmt.Errorf("get spreadsheets: %w", err)
	}

	if len(resp.Values) == 0 {
		return nil, fmt.Errorf("empty google sheets response")
	}

	var rows []map[string]string

	headers := NewHeadersFromRow(resp.Values[0])

	for _, value := range resp.Values[1:] {
		row := make(map[string]string)

		for i, item := range value {
			row[headers.GetByIndex(i)] = item.(string)
		}

		if len(row) == 0 {
			continue
		}

		rows = append(rows, row)
	}

	return rows, nil
}

func (s *Service) GetColumn(_ context.Context, tableName string, colName string) ([]string, error) {
	resp, err := s.sheetsClient.Spreadsheets.Values.Get(s.spreadsheetID, tableName).Do()
	if err != nil {
		return nil, fmt.Errorf("get spreadsheets: %w", err)
	}

	if len(resp.Values) == 0 {
		return nil, fmt.Errorf("empty google sheets response")
	}

	var column []string

	headers := NewHeadersFromRow(resp.Values[0])

	for _, value := range resp.Values[1:] {
		for i, item := range value {
			if headers.GetByIndex(i) != colName {
				continue
			}

			column = append(column, item.(string))
		}
	}

	return column, nil
}
