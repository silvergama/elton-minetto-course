package beer

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// vamos usar agora o banco de dados SQLite
// para isso precisamos primeiro inicializar no nosso projeto o sistema de
// importação de modulos.
// execute o comando:
// go mod init github.com/silvergama/elton-minetto-course
// sendo o github.com/silvergama/elton-minetto-course o nome completo do
// projeto. É remomendado usar essa nomeclatura para evitar conflito de pacotes
// em projetos maiores. O próximo passo é importar o pacote SQLite  com o
// comando:
// go get github.com/mattn/go-sqlite3
// com isso vai ser feito o download do pacote e atualizar os arquivos go.mod
// go.sum com as dependencias.

//define a interface com as funções que serão usadas pelo restante do projeto
type UseCase interface {
	GetAll() ([]*Beer, error)
	Get(ID int64) (*Beer, error)
	Store(b *Beer) error
	Update(b *Beer) error
	Remove(ID int64) error
}

// Em Go qualquer coisa que implemente as funções de uma interface passa  a ser
// uma implementação válida. Não existe uma palavra "implements" como Java ou
// PHP desta forma uma struct, uma string, um inteiro, etc. qualquer coisa pode
// ser válido, desde que implemente todas as funcões.
type Service struct {
	DB *sql.DB
}

// esta função retorna um ponteiro em memória para uma estrutura
func NewService(db *sql.DB) *Service {
	return &Service{
		DB: db,
	}
}

// Vamos implementar as funcões na próxima etapa
func (s *Service) GetAll() ([]*Beer, error) {
	// result é um slice de ponteiros do tipo Beer
	var result []*Beer

	// vamos sempre usar a conexão que está dentro do Service
	rows, err := s.DB.Query("select id, name, type, style from beer")
	// se existe erro a função deve retorna-lo e ele vai ser tratado
	// por quem chamou o pacote. Esta é uma boa prática em Go
	if err != nil {
		return nil, err
	}

	// A funcão defer garante que o comando rows.Close vai ser executado na
	// saída da função, desta forma não precisamos nos preocupar em fechar a
	// conexão.
	defer rows.Close()

	for rows.Next() {
		var b Beer
		err = rows.Scan(&b.ID, &b.Name, &b.Type, &b.Style)
		if err != nil {
			return nil, err
		}
		// o comando append adiciona novos itens  a um slice, sempre no final.
		result = append(result, &b)
	}

	return result, nil
}

func (s *Service) Get(ID int64) (*Beer, error) {
	// b é um tipo Beer
	var b Beer

	// o comando Prepare verifica se a consulta está válida
	stmt, err := s.DB.Prepare("select id, name, type, style from beer where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(ID).Scan(&b.ID, &b.Name, &b.Type, &b.Style)
	if err != nil {
		return nil, err
	}

	// deve retornar a posição de memória de b
	return &b, nil
}

func (s *Service) Store(b *Beer) error {
	// iniciamos uma transação
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into beer(id, name, type, style) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	// O comando Exec retorn um Result, mas não temos interesse nele, por isso
	// podemos ignora-lo com o _
	_, err = stmt.Exec(b.ID, b.Name, b.Type, b.Style)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *Service) Update(b *Beer) error {
	if b.ID == 0 {
		// Podemos também um error de aplicação que criamos para definir uma
		// condição de erro com um possível update sem WHERE
		return fmt.Errorf("Invalid ID")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("update beer set name=?, type=?, style=? where id=?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	// O comando Exec retora um Result, mas não temos interesse nele, por isso
	// podemos ignora-lo com o _
	_, err = stmt.Exec(b.Name, b.Type, b.Style, b.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *Service) Remove(ID int64) error {
	if ID == 0 {
		// podemos também retornar um erro de aplicação que criamos para definir
		// uma condição de erro, como possível update sem WHERE
		return fmt.Errorf("Invalid ID")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	// o comando Exec retorna um Result, mas não temos interesse nele, por isso
	// podemos igonora-lo com o _
	_, err = tx.Exec("delete from beer where id=?", ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
