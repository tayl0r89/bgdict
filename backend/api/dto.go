package api

type WordResult struct {
	BaseWord   Word            `json:"baseWord"`
	Derivative *DerivativeForm `json:"derivative"`
}

type Word struct {
	Id           int      `json:"id"`
	Name         string   `json:"name"`
	NameStressed string   `json:"nameStressed"`
	NameBroken   string   `json:"nameBroken"`
	TypeId       int      `json:"typeId"`
	Type         WordType `json:"wordType"`
}

type WordType struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	SpeechPart string `json:"SpeechPart"`
}

type DerivativeForm struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	NameBroken   string `json:"nameBroken"`
	NameStressed string `json:"nameStressed"`
	IsInfinitive int    `json:"isInfinitive"`
	BaseWordId   int    `json:"baseWordId"`
	BaseWord     Word   `json:"word"`
}

type IncorrectForm struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	CorrectWordId int    `json:"correctWordId"`
}
