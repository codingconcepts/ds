package repo

import (
	"context"
	"ds/internal/pkg/model"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestInsertTable(t *testing.T) {
	if !integrationTests {
		t.Skipf("not running integration tests")
	}

	sourceTable := model.Table{
		Name: "person",
		Columns: []model.Column{
			{Name: "id"},
			{Name: "full_name"},
			{Name: "created_at"},
		},
		ReadLimit: 3,
		Filter:    "WHERE created_at > '2023-01-01T01:01:01Z'",
	}

	targetTable := model.Table{
		Name: "person",
		Columns: []model.Column{
			{Name: "id"},
			{Name: "full_name"},
			{Name: "created_at"},
		},
	}

	assert.Nil(t, InsertTable(source, target, sourceTable, targetTable))

	act := fetchTargetPeople(t)
	act = lo.Map(act, func(p person, i int) person {
		p.createdAt = p.createdAt.In(time.UTC)
		return p
	})

	assert.Equal(t, person{id: "bc229cee-5387-4c36-b83e-7a46613071de", fullName: "b b", createdAt: time.Date(2023, 1, 1, 1, 1, 2, 0, time.UTC)}, act[0])
	assert.Equal(t, person{id: "ccb45142-cba0-4f97-9179-33c7e3d51e92", fullName: "c c", createdAt: time.Date(2023, 1, 1, 1, 1, 3, 0, time.UTC)}, act[1])
	assert.Equal(t, person{id: "dfbad599-b5b6-4c59-8859-b19db1dd2ac0", fullName: "d d", createdAt: time.Date(2023, 1, 1, 1, 1, 4, 0, time.UTC)}, act[2])
	assert.Equal(t, person{id: "eba7ea84-e57b-4816-8806-2faae31c2830", fullName: "e e", createdAt: time.Date(2023, 1, 1, 1, 1, 5, 0, time.UTC)}, act[3])
}

func TestUpdateTable(t *testing.T) {
	if !integrationTests {
		t.Skipf("not running integration tests")
	}

	sourceTable := model.Table{
		Name: "person",
		Columns: []model.Column{
			{Name: "id"},
			{Name: "full_name"},
			{Name: "created_at"},
		},
		ReadLimit: 3,
		Filter:    "WHERE created_at > '2023-01-01T01:01:01Z'",
	}

	targetTable := model.Table{
		Name: "person",
		Columns: []model.Column{
			{Name: "id"},
			{Name: "full_name"},
			{Name: "created_at"},
		},
	}

	assert.Nil(t, InsertTable(source, target, sourceTable, targetTable))

	makeUpdate(t)

	act := fetchTargetPeople(t)
	act = lo.Map(act, func(p person, i int) person {
		p.createdAt = p.createdAt.In(time.UTC)
		return p
	})

	assert.Equal(t, person{id: "bc229cee-5387-4c36-b83e-7a46613071de", fullName: "B B", createdAt: time.Date(2023, 1, 1, 1, 1, 2, 0, time.UTC)}, act[0])
	assert.Equal(t, person{id: "ccb45142-cba0-4f97-9179-33c7e3d51e92", fullName: "C C", createdAt: time.Date(2023, 1, 1, 1, 1, 3, 0, time.UTC)}, act[1])
	assert.Equal(t, person{id: "dfbad599-b5b6-4c59-8859-b19db1dd2ac0", fullName: "D D", createdAt: time.Date(2023, 1, 1, 1, 1, 4, 0, time.UTC)}, act[2])
	assert.Equal(t, person{id: "eba7ea84-e57b-4816-8806-2faae31c2830", fullName: "E E", createdAt: time.Date(2023, 1, 1, 1, 1, 5, 0, time.UTC)}, act[3])
	assert.Equal(t, person{id: "ee807359-2a2c-4f6b-a753-0b3cddc3729a", fullName: "F F", createdAt: time.Date(2023, 1, 1, 1, 1, 5, 0, time.UTC)}, act[4])
}

func makeUpdate(t *testing.T) {
	insertStmt := `INSERT INTO person (id, full_name, created_at) VALUES
		('ee807359-2a2c-4f6b-a753-0b3cddc3729a', 'f f', '2023-01-01T01:01:05Z')`

	if _, err := target.Exec(context.Background(), insertStmt); err != nil {
		t.Fatalf("error inserting row: %v", err)
	}

	updateStmt := `UPDATE person SET full_name = upper(full_name)`

	if _, err := target.Exec(context.Background(), updateStmt); err != nil {
		t.Fatalf("error updating rows: %v", err)
	}
}

type person struct {
	id        string
	fullName  string
	createdAt time.Time
}

func fetchTargetPeople(t *testing.T) []person {
	stmt := `SELECT id, full_name, created_at FROM person`
	rows, err := target.Query(context.Background(), stmt)
	assert.Nil(t, err)

	var people []person
	var p person
	for rows.Next() {
		if err = rows.Scan(&p.id, &p.fullName, &p.createdAt); err != nil {
			t.Fatalf("error scanning target person: %v", err)
		}

		people = append(people, p)
	}

	return people
}
