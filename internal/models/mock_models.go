package models

type Definition struct {
	Examples     []string `json:"examples" firestore:"examples"`
	Meaning      string   `json:"meaning" firestore:"meaning"`
	PartOfSpeech string   `json:"partOfSpeech" firestore:"partOfSpeech"`
}

type Query struct {
}

type Tenses struct {
	Continuous *string `json:"continuous,omitempty" firestore:"continuous,omitempty"`
	Future     *string `json:"future,omitempty" firestore:"future,omitempty"`
	Past       *string `json:"past,omitempty" firestore:"past,omitempty"`
	Present    *string `json:"present,omitempty" firestore:"present,omitempty"`
}

type Word struct {
	ID            string        `json:"id" firestore:"id"`
	Word          string        `json:"word" firestore:"word"`
	Synonyms      []string      `json:"synonyms,omitempty" firestore:"synonyms,omitempty"`
	Antonyms      []string      `json:"antonyms,omitempty" firestore:"antonyms,omitempty"`
	RelatedTerms  []string      `json:"relatedTerms,omitempty" firestore:"relatedTerms,omitempty"`
	Pronunciation string        `json:"pronunciation" firestore:"pronunciation"`
	Tenses        *Tenses       `json:"tenses,omitempty" firestore:"tenses,omitempty"`
	Definitions   []*Definition `json:"definitions" firestore:"definitions"`
}
