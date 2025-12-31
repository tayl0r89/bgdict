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
	GetDerivedForms(id int) ([]*DerivativeForm, error)
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

func translationToDto(r *db.WordTranslation) *WordTranslation {
	if r == nil {
		return nil
	}
	return &WordTranslation{
		Id:      int(r.ID),
		Lang:    getOrEmpty(r.Lang),
		Content: getOrEmpty(r.Content),
		WordId:  getOrZero(r.WordID),
	}
}

func wordTypeToDto(r *db.WordType) *WordType {
	if r == nil {
		return nil
	}
	return &WordType{
		Id:         int(r.ID),
		Name:       getOrEmpty(r.Name),
		SpeechPart: getOrEmpty(r.SpeechPart),
	}
}

func wordRowToDto(r *db.Word, wt *db.WordType, wtr *db.WordTranslation) *Word {
	return &Word{
		Id:           int(r.ID),
		Name:         getOrEmpty(r.Name),
		NameStressed: getOrEmpty(r.NameStressed),
		NameBroken:   getOrEmpty(r.NameBroken),
		TypeId:       getOrZero(r.TypeID),
		Type:         wordTypeToDto(wt),
		Translation:  translationToDto(wtr),
	}
}

func findWordRowToDto(r *db.FindWordsRow) *DerivativeForm {
	return &DerivativeForm{
		Id:           int(r.DerivativeForm.ID),
		Name:         getOrEmpty(r.DerivativeForm.Name),
		NameBroken:   getOrEmpty(r.DerivativeForm.NameBroken),
		NameStressed: getOrEmpty(r.DerivativeForm.NameStressed),
		Description:  getOrEmpty(r.DerivativeForm.Description),
		IsInfinitive: getOrZero(r.DerivativeForm.IsInfinitive),
		BaseWordId:   getOrZero(r.DerivativeForm.BaseWordID),
		BaseWord:     wordRowToDto(&r.Word, &r.WordType, &r.WordTranslation),
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
	return wordRowToDto(&res.Word, &res.WordType, &res.WordTranslation), nil
}

func wordByNameRowToResult(res *db.GetWordByNameRow) WordResult {
	return WordResult{
		BaseWord: Word{
			Id:           int(res.Word.ID),
			Name:         getOrEmpty(res.Word.Name),
			NameStressed: getOrEmpty(res.Word.NameStressed),
			NameBroken:   getOrEmpty(res.Word.NameBroken),
			TypeId:       int(res.Word.TypeID.Int32),
			Type:         wordTypeToDto(&res.WordType),
			Translation:  translationToDto(&res.WordTranslation),
		},
		Derivative: nil,
	}
}

func findWordRowToResult(res *db.FindWordsRow) WordResult {
	return WordResult{
		BaseWord: Word{
			Id:           int(res.Word.ID),
			Name:         getOrEmpty(res.Word.Name),
			NameStressed: getOrEmpty(res.Word.NameStressed),
			NameBroken:   getOrEmpty(res.Word.NameBroken),
			TypeId:       int(res.Word.TypeID.Int32),
			Type:         wordTypeToDto(&res.WordType),
			Translation:  translationToDto(&res.WordTranslation),
		},
		Derivative: &DerivativeForm{
			Id:           int(res.DerivativeForm.ID),
			Name:         getOrEmpty(res.DerivativeForm.Name),
			NameBroken:   getOrEmpty(res.DerivativeForm.NameBroken),
			NameStressed: getOrEmpty(res.DerivativeForm.NameStressed),
			IsInfinitive: getOrZero(res.DerivativeForm.IsInfinitive),
			BaseWordId:   getOrZero(res.DerivativeForm.BaseWordID),
			BaseWord:     wordRowToDto(&res.Word, &res.WordType, &res.WordTranslation),
			Description:  getOrEmpty(res.DerivativeForm.Description),
		},
	}
}

func (l *dbWordLoader) GetDerivedForms(id int) ([]*DerivativeForm, error) {
	ctx := context.Background()

	word, word_err := l.GetWordById(id)
	if word_err != nil {
		log.Println(word_err.Error())
		return nil, word_err
	}

	var result = make([]*DerivativeForm, 0)

	derived, err := l.queries.GetDerived(ctx, sql.NullInt32{Valid: true, Int32: int32(id)})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	for _, item := range derived {
		result = append(result, &DerivativeForm{
			Id:           int(item.DerivativeForm.ID),
			Name:         item.DerivativeForm.Name.String,
			NameBroken:   item.DerivativeForm.NameBroken.String,
			NameStressed: item.DerivativeForm.NameStressed.String,
			Description:  item.DerivativeForm.Description.String,
			IsInfinitive: int(item.DerivativeForm.IsInfinitive.Int32),
			BaseWord:     word,
			BaseWordId:   id,
		})
	}

	return result, nil
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
