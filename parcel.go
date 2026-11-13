package main

import (
	"database/sql"
	"fmt"
	"log"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(`INSERT INTO parcel (client, status, address, created_at)
		VALUES (:client, :status, :address, :created_at);`,
		sql.Named("client", p.Client), sql.Named("status", p.Status),
		sql.Named("address", p.Address), sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, fmt.Errorf("Не удалось добавить посылку. %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Не удалось добавить посылку. %w", err)
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (*Parcel, error) {
	p := &Parcel{}

	row := s.db.QueryRow(`SELECT number, client, status, address, created_at
		FROM parcel WHERE number = :number;`, sql.Named("number", number))

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("Не удалось найти посылку №%d : %w", number, err)
	}

	if p.Number == 0 {
		return nil, nil
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client;",
		sql.Named("client", client))

	if err != nil {
		return res, fmt.Errorf("Не удалось получить посылки клиента %d : %w", client, err)
	}
	defer rows.Close()

	for rows.Next() {
		p := Parcel{}
		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			log.Println(err)
			continue
		}
		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return res, fmt.Errorf("Не удалось получить посылки клиента %d : %w", client, err)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec(`UPDATE parcel
	SET status = :status
	WHERE number = :number;`,
		sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		return fmt.Errorf("Не удалось обновить статус посылки №%d : %w", number, err)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec(`UPDATE parcel
	SET address = :address
	WHERE number = :number;`,
		sql.Named("address", address), sql.Named("number", number))
	if err != nil {
		return fmt.Errorf("Не удалось обновить адрес посылки №%d : %w", number, err)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	res, err := s.db.Exec(`DELETE FROM parcel 
		WHERE number = :number
		AND status = 'registered';`,
		sql.Named("number", number))
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("Посылка %d не найдена или ее статус не registered", number)
	}

	return nil
}
