package repo

import (
	"biocad/internal/entity"
	"biocad/internal/usecase"
	"biocad/pkg/postgres"
	"context"
	"fmt"
	"log"
)

type UserRepo struct {
	*postgres.Postgres
}

var _ usecase.UserRp = (*UserRepo)(nil)

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (u *UserRepo) GetAllGuidList(ctx context.Context) ([]string, error) {
	query := `select distinct unit_guid from data`

	rows, err := u.Pool.Query(ctx, query)
	if err != nil {
		log.Println("Cannot execute query to get all guids")
		return nil, fmt.Errorf("cannot execute query to get all guids %w", err)
	}
	defer rows.Close()

	var allGuids []string
	for rows.Next() {
		var guid string
		err = rows.Scan(&guid)
		if err != nil {
			log.Printf("Cannot scan guid: %v\n", err)
			return nil, fmt.Errorf("cannot scan guid %w", err)
		}
		allGuids = append(allGuids, guid)
	}
	return allGuids, nil
}

func (u *UserRepo) CheckGuidExists(ctx context.Context, guid string) (bool, error) {
	query := `select exists(select * from data where unit_guid = $1)`

	rows, err := u.Pool.Query(ctx, query, guid)
	if err != nil {
		log.Println("Cannot execute query to check guid existence")
		return false, fmt.Errorf("cannot execute query to check guid existence %w", err)
	}
	defer rows.Close()

	var exists bool
	for rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			log.Printf("Cannot scan guid existence: %v\n", err)
			return false, fmt.Errorf("cannot scan guid existence %w", err)
		}
	}
	return exists, nil
}

func (u *UserRepo) GetDataByUnitGuid(ctx context.Context, guid string, limit uint64, offset uint64) ([]entity.Data, error) {
	query := `select n, mqtt, invid, unit_guid, msg_id, msg_text, context, class, level, area, addr, block, type, bit, invert_bit from data where unit_guid = $1 order by n LIMIT $2 OFFSET $3`

	rows, err := u.Pool.Query(ctx, query, guid, limit, offset)
	if err != nil {
		log.Println("Cannot execute query to get data by guid")
		return nil, fmt.Errorf("cannot execute query to get data by guid %w", err)
	}
	defer rows.Close()

	var allData []entity.Data
	for rows.Next() {
		var data entity.Data
		err = rows.Scan(&data.Id,
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
			log.Printf("Cannot scan data by guid: %v\n", err)
			return nil, fmt.Errorf("cannot scan data by guid %w", err)
		}
		allData = append(allData, data)
	}
	return allData, err
}
