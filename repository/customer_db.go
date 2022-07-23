package repository

import "github.com/jmoiron/sqlx"

type cutomerRepositoryDb struct {
	db *sqlx.DB
}

func NewCustomerRepositoryDB(db *sqlx.DB) cutomerRepositoryDb {
	return cutomerRepositoryDb{db: db}
}

func (r cutomerRepositoryDb) GetAll() ([]Customer, error) {
	customers := []Customer{}
	query := "select * from customers"
	err := r.db.Select(&customers, query)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (r cutomerRepositoryDb) GetById(id int) (*Customer, error) {
	customer := Customer{}
	query := "select * from customers where customer_id = ?"
	err := r.db.Get(&customer, query, id)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}
