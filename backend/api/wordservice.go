package api

import (
	"context"
	"database/sql"
	"log"

	db "github.com/tayl0r89/bgdict/db/generated"
)

type WordRepository interface {
	FindWords(query string) ([]*DerivativeForm, error)
	GetWordById(id int) (*Word, error)
}

type dbWordLoader struct {
	queries *db.Queries
}

func NewWordRepository(dbConn *sql.DB) WordRepository {
	queries := db.New(dbConn)
	return &dbWordLoader{queries: queries}
}

func getOrEmpty(s sql.NullString) string {
	if !s.Valid {
		return ""
	}

	return s.String
}

func getOrZero(s sql.NullInt32) int {
	if !s.Valid {
		return 0
	}

	return int(s.Int32)
}

func wordRowToDto(r *db.Word) *Word {
	return &Word{
		Id:           int(r.ID),
		Name:         getOrEmpty(r.Name),
		NameStressed: getOrEmpty(r.NameStressed),
		NameBroken:   getOrEmpty(r.NameBroken),
		TypeId:       getOrZero(r.TypeID),
	}
}

func findWordRowToDto(r *db.FindWordsRow) *DerivativeForm {
	return &DerivativeForm{
		Id:           getOrZero(r.DerivativeForm.ID),
		Name:         getOrEmpty(r.DerivativeForm.Name),
		NameBroken:   getOrEmpty(r.DerivativeForm.NameBroken),
		NameStressed: getOrEmpty(r.DerivativeForm.NameStressed),
		IsInfinitive: getOrZero(r.DerivativeForm.IsInfinitive),
		BaseWordId:   getOrZero(r.DerivativeForm.BaseWordID),
		BaseWord: Word{
			Id:           int(r.Word.ID),
			Name:         getOrEmpty(r.Word.Name),
			NameStressed: getOrEmpty(r.Word.NameStressed),
			NameBroken:   getOrEmpty(r.Word.NameBroken),
			TypeId:       getOrZero(r.Word.TypeID),
		},
	}
}

func mapResultsToDto(results []db.FindWordsRow) (r []*DerivativeForm) {
	for _, result := range results {
		r = append(r, findWordRowToDto(&result))
	}
	return r
}

func (l *dbWordLoader) FindWords(query string) ([]*DerivativeForm, error) {
	ctx := context.Background()
	res, err := l.queries.FindWords(ctx, sql.NullString{String: query, Valid: true})

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return mapResultsToDto(res), nil
}

func (l *dbWordLoader) GetWordById(id int) (*Word, error) {
	ctx := context.Background()
	res, err := l.queries.GetWord(ctx, int32(id))

	if err != err {
		log.Println(err.Error())
		return nil, err
	}
	return wordRowToDto(&res), nil
}
