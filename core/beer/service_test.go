package beer

import (
	"database/sql"
	"testing"
)

func TestStore(t *testing.T) {
	b := &Beer{
		ID:    1,
		Name:  "Heineken",
		Type:  TypeLager,
		Style: StylePale,
	}

	db, err := sql.Open("sqlite3", "../../data/beer_test.db")
	if err != nil {
		t.Fatalf("Erro ao tentar conectar ao banco de dados: %s", err.Error())
	}

	defer db.Close()

	service := NewService(db)

	err = clearDB(db)
	if err != nil {
		t.Errorf("Limpando o banco de dados: %s", err.Error())
	}

	err = service.Store(b)
	if err != nil {
		t.Fatalf("Erro ao tentar salvar no banco de dados: %s", err.Error())
	}

	saved, err := service.Get(b.ID)
	if err != nil {
		t.Fatalf("Erro ao tentar buscar do banco de dados: %s", err.Error())
	}

	if saved.ID != b.ID {
		t.Fatalf("Dados inválidos. Esperando %d, recebido %d", 1, saved.ID)
	}

}

func TestUpdate(t *testing.T) {
	b := &Beer{
		ID:    1,
		Name:  "Heineken updated",
		Type:  TypeLager,
		Style: StylePale,
	}

	db := connection(t)
	defer db.Close()

	service := NewService(db)

	err := service.Update(b)
	if err != nil {
		t.Fatalf("Erro ao tentar atualizar dados no banco de dados. %s", err.Error())
	}

	updated, err := service.Get(b.ID)
	if err != nil {
		t.Fatalf("Erro ao tentar buscar do banco de dados. %s", err.Error())
	}

	if updated.ID != b.ID {
		t.Fatalf("Dados inválidos. Esperando %s, recebido %s", b.Name, updated.Name)
	}

}

func connection(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "../../data/beer_test.db")
	if err != nil {
		t.Fatalf("Erro ao conectar ao banco de dados. %s", err.Error())
	}
	return db
}

func clearDB(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("delete from beer")
	tx.Commit()

	return err
}
