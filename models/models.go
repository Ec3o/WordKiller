package models

type Question struct {
	PaperDetailId string  `json:"paperDetailId"`
	Title         string  `json:"title"`
	AnswerA       string  `json:"answerA"`
	AnswerB       string  `json:"answerB"`
	AnswerC       string  `json:"answerC"`
	AnswerD       string  `json:"answerD"`
	QuestionId    *string `json:"questionId"`
	QuestionNum   int     `json:"questionNum"`
	Answer        *string `json:"answer"`
	Input         *string `json:"input"`
	Level         int     `json:"level"`
	Cet           int     `json:"cet"`
	Right         *bool   `json:"right"`
}
type Paper struct {
	PaperId string     `json:"paperId"`
	Type    string     `json:"type"`
	List    []Question `json:"list"`
}

// ExamScore 用于解析考试得分的响应
type ExamScore struct {
	Mark int `json:"mark"`
}

type AnswerBank struct {
	Word string `json:"Word"`
	Mean string `json:"Mean"`
}

type Answer struct {
	PaperId string         `json:"paperId"`
	Type    string         `json:"type"`
	List    []AnswerDetail `json:"list"`
}

type AnswerDetail struct {
	Input         string `json:"input"`
	PaperDetailId string `json:"paperDetailId"`
}
