package models

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type Hero struct {
	ID         int    `dynamo:"ID"`
	Name       string `dynamo:"name"`
	Superpower string `dynamo:"superpower"`
}

type DB struct {
	table dynamo.Table
}

func Connect() *DB {
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	return &DB{table: db.Table("heroes")}
}

func (db *DB) Insert(hero *Hero) error {
	return db.table.Put(hero).Run()
}

func (db *DB) GetAllHeroes() ([]*Hero, error) {
	var heroes []*Hero
	err := db.table.Scan().All(&heroes)
	return heroes, err
}

func (db *DB) GetAHero(ID int64) (*Hero, error) {
	hero := new(Hero)
	err := db.table.Get(strconv.FormatInt(ID, 10), nil).One(hero)
	return hero, err
}
