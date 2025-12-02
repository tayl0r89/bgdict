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
	SearchWord(query string) ([]*WordResult, error)
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

func wordRowToDto(r *db.Word, wt *db.WordType) *Word {
	return &Word{
		Id:           int(r.ID),
		Name:         getOrEmpty(r.Name),
		NameStressed: getOrEmpty(r.NameStressed),
		NameBroken:   getOrEmpty(r.NameBroken),
		TypeId:       getOrZero(r.TypeID),
		Type: WordType{
			Id:         int(wt.ID),
			Name:       getOrEmpty(wt.Name),
			SpeechPart: getOrEmpty(wt.SpeechPart),
		},
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
			Type: WordType{
				Id:         int(r.WordType.ID),
				Name:       getOrEmpty(r.WordType.Name),
				SpeechPart: getOrEmpty(r.WordType.SpeechPart),
			},
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
	return wordRowToDto(&res.Word, &res.WordType), nil
}

func wordByNameRowToResult(res *db.GetWordByNameRow) WordResult {
	return WordResult{
		BaseWord: Word{
			Id:           int(res.Word.ID),
			Name:         getOrEmpty(res.Word.Name),
			NameStressed: getOrEmpty(res.Word.NameStressed),
			NameBroken:   getOrEmpty(res.Word.NameBroken),
			TypeId:       int(res.Word.TypeID.Int32),
			Type: WordType{
				Id:         int(res.WordType.ID),
				Name:       getOrEmpty(res.WordType.Name),
				SpeechPart: getOrEmpty(res.WordType.SpeechPart),
			},
		},
		Derivative: nil,
	}
}

func findWordRowToResult(res *db.FindWordsRow) WordResult {
	var derivativeForm *DerivativeForm = nil
	if res.DerivativeForm.ID.Valid {
		derivativeForm = &DerivativeForm{
			Id:           getOrZero(res.DerivativeForm.ID),
			Name:         getOrEmpty(res.DerivativeForm.Name),
			NameBroken:   getOrEmpty(res.DerivativeForm.NameBroken),
			NameStressed: getOrEmpty(res.DerivativeForm.NameStressed),
			IsInfinitive: getOrZero(res.DerivativeForm.IsInfinitive),
			BaseWordId:   getOrZero(res.DerivativeForm.BaseWordID),
			BaseWord: Word{
				Id:           int(res.Word.ID),
				Name:         getOrEmpty(res.Word.Name),
				NameStressed: getOrEmpty(res.Word.NameStressed),
				NameBroken:   getOrEmpty(res.Word.NameBroken),
				TypeId:       int(res.Word.TypeID.Int32),
				Type: WordType{
					Id:         int(res.WordType.ID),
					Name:       getOrEmpty(res.WordType.Name),
					SpeechPart: getOrEmpty(res.WordType.SpeechPart),
				},
			},
		}
	}
	return WordResult{
		BaseWord: Word{
			Id:           int(res.Word.ID),
			Name:         getOrEmpty(res.Word.Name),
			NameStressed: getOrEmpty(res.Word.NameStressed),
			NameBroken:   getOrEmpty(res.Word.NameBroken),
			TypeId:       int(res.Word.TypeID.Int32),
			Type: WordType{
				Id:         int(res.WordType.ID),
				Name:       getOrEmpty(res.WordType.Name),
				SpeechPart: getOrEmpty(res.WordType.SpeechPart),
			},
		},
		Derivative: derivativeForm,
	}
}

func (l *dbWordLoader) SearchWord(query string) ([]*WordResult, error) {
	ctx := context.Background()
	var result = make([]*WordResult, 0)
	// Search in derived words & get derived + base word if available

	log.Println("Looking at derived words")
	derived, derived_err := l.queries.FindWords(ctx, sql.NullString{String: query, Valid: true})
	if derived_err != nil {
		log.Println(derived_err.Error())
	}

	if len(derived) > 0 {
		converted := findWordRowToResult(&derived[0])
		result = append(result, &converted)
		return result, nil
	}

	res, err := l.queries.GetWordByName(ctx, sql.NullString{String: query, Valid: true})

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if len(res) > 0 {
		results := make([]*WordResult, 0)
		for _, item := range res {
			word := wordByNameRowToResult(&item)
			results = append(results, &word)
		}
		log.Println("Word result found")
		return results, nil
	}

	return make([]*WordResult, 0), nil
}
