package postgre

import (
	"context"
	"encoding/json"
)

type IntegrationRow struct {
	ID     string            `json:"id"`
	Type   string            `json:"type"`
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
	Status string            `json:"status"`
}

func UpsertIntegration(ctx context.Context, i *IntegrationRow) error {
	configJSON, err := json.Marshal(i.Config)
	if err != nil {
		return err
	}
	_, err = Pool.Exec(ctx,
		`INSERT INTO supportflow.integrations (id, type, name, config, status)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (id) DO UPDATE SET type=$2, name=$3, config=$4, status=$5, updated_at=now()`,
		i.ID, i.Type, i.Name, configJSON, i.Status)
	return err
}

func UpdateIntegrationStatus(ctx context.Context, id, status string) error {
	_, err := Pool.Exec(ctx,
		`UPDATE supportflow.integrations SET status=$1, updated_at=now() WHERE id=$2`,
		status, id)
	return err
}

func ListIntegrations(ctx context.Context) ([]IntegrationRow, error) {
	rows, err := Pool.Query(ctx,
		`SELECT id, type, name, config, status FROM supportflow.integrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []IntegrationRow
	for rows.Next() {
		var r IntegrationRow
		var configJSON []byte
		if err := rows.Scan(&r.ID, &r.Type, &r.Name, &configJSON, &r.Status); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(configJSON, &r.Config); err != nil {
			r.Config = map[string]string{}
		}
		result = append(result, r)
	}
	return result, nil
}

func DeleteIntegration(ctx context.Context, id string) error {
	_, err := Pool.Exec(ctx,
		`DELETE FROM supportflow.integrations WHERE id=$1`, id)
	return err
}

func UpsertChannelMapping(ctx context.Context, customerID, channel, externalID string) error {
	_, err := Pool.Exec(ctx,
		`INSERT INTO supportflow.channel_mappings (customer_id, channel, external_id)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (customer_id, channel) DO UPDATE SET external_id=$3, updated_at=now()`,
		customerID, channel, externalID)
	return err
}

func GetChannelMapping(ctx context.Context, customerID, channel string) (string, error) {
	var externalID string
	err := Pool.QueryRow(ctx,
		`SELECT external_id FROM supportflow.channel_mappings WHERE customer_id=$1 AND channel=$2`,
		customerID, channel).Scan(&externalID)
	if err != nil {
		return "", err
	}
	return externalID, nil
}
