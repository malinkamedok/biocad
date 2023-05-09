package repo

import (
	"biocad/internal/entity"
	"biocad/internal/usecase"
	"biocad/pkg/postgres"
	"context"
	"fmt"
	"log"
	"time"
)

type ParserRepo struct {
	*postgres.Postgres
}

var _ usecase.ParserRp = (*ParserRepo)(nil)

func NewParserRepo(pg *postgres.Postgres) *ParserRepo {
	return &ParserRepo{pg}
}

func (p *ParserRepo) InsertData(data entity.Data) error {
	ctx := context.TODO()

	query := `insert into data (n, mqtt, invid, unit_guid, msg_id, msg_text, context, class, level, area, addr, block, type, bit, invert_bit) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	_, err := p.Pool.Exec(ctx, query, data.Id, data.Mqtt, data.Invid, data.UnitGuid, data.MsgId, data.MsgText, data.Context, data.Class, data.Level, data.Area, data.Addr, data.Block, data.DataType, data.Bit, data.InvertBit)
	if err != nil {
		log.Println("Cannot execute query to insert data")
		return fmt.Errorf("cannot execute query to insert data %w", err)
	}
	return nil
}

func (p *ParserRepo) GetAllGUIDList() ([]string, error) {
	ctx := context.TODO()

	query := `select distinct unit_guid from data`

	rows, err := p.Pool.Query(ctx, query)
	if err != nil {
		log.Println("Cannot execute query to get all GUIDs")
		return nil, fmt.Errorf("cannot execute query to get GUIDs %w", err)
	}
	defer rows.Close()

	var listGUIDs []string
	for rows.Next() {
		var guid string
		err = rows.Scan(&guid)
		if err != nil {
			log.Println("Cannot scan GUID")
			return nil, fmt.Errorf("cannot scan GUID %w", err)
		}
		listGUIDs = append(listGUIDs, guid)
	}
	return listGUIDs, nil
}

func (p *ParserRepo) GetAllDataByGUID(guid string) ([]entity.Data, error) {
	ctx := context.TODO()

	query := `select n, mqtt, invid, unit_guid, msg_id, msg_text, context, class, level, area, addr, block, type, bit, invert_bit from data where unit_guid = $1`

	rows, err := p.Pool.Query(ctx, query, guid)
	if err != nil {
		log.Println("Cannot execute query to get data by GUID")
		return nil, fmt.Errorf("cannot execute query to get data by GUID %w", err)
	}
	defer rows.Close()

	var allData []entity.Data
	for rows.Next() {
		var data entity.Data
		err = rows.Scan(
			&data.Id,
			&data.Mqtt,
			&data.Invid,
			&data.UnitGuid,
			&data.MsgId,
			&data.MsgText,
			&data.Context,
			&data.Class,
			&data.Level,
			&data.Area,
			&data.Addr,
			&data.Block,
			&data.DataType,
			&data.Bit,
			&data.InvertBit)
		if err != nil {
			log.Printf("Cannot scan data by GUID: %v\n", err)
			return nil, fmt.Errorf("cannot scan data by GUID %w", err)
		}
		allData = append(allData, data)
	}
	return allData, nil
}

func (p *ParserRepo) GetAllFileNames() ([]string, error) {
	ctx := context.TODO()

	query := `select file_name from processed_data`

	rows, err := p.Pool.Query(ctx, query)
	if err != nil {
		log.Println("Cannot execute query to get all file names")
		return nil, fmt.Errorf("cannot execute query to get all file names %w", err)
	}
	defer rows.Close()

	var fileNames []string
	for rows.Next() {
		var fn string
		err = rows.Scan(&fn)
		if err != nil {
			log.Println("Cannot scan file name")
			return nil, fmt.Errorf("cannot scan file name %w", err)
		}
		fileNames = append(fileNames, fn)
	}
	return fileNames, nil
}

func (p *ParserRepo) InsertFileName(fileName string) error {
	ctx := context.TODO()

	query := `insert into processed_data (file_name) values ($1)`

	_, err := p.Pool.Exec(ctx, query, fileName)
	if err != nil {
		log.Println("Cannot execute query to save history")
		return fmt.Errorf("cannot execute query to save history %w", err)
	}
	return nil
}

func (p *ParserRepo) ChangeProcessedStatus(fileName string) error {
	ctx := context.TODO()

	query := `update processed_data set is_processed = $1, date_processed = $2 where file_name = $3`

	_, err := p.Pool.Exec(ctx, query, true, time.Now(), fileName)
	if err != nil {
		log.Println("Cannot execute query to change status")
		return fmt.Errorf("cannot execute query to change status %w", err)
	}
	return nil
}

func (p *ParserRepo) FailureReport(fileName string) error {
	ctx := context.TODO()

	query := `update processed_data set failure = $1, date_processed = $2 where file_name = $3`

	_, err := p.Pool.Exec(ctx, query, true, time.Now(), fileName)
	if err != nil {
		log.Println("Cannot execute query to set failure report")
		return fmt.Errorf("cannot execute query to set failure report %w", err)
	}
	return nil
}
