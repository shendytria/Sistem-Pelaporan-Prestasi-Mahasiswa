package model

type StatisticByType struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type StatisticByLevel struct {
	Level string `json:"level"`
	Count int    `json:"count"`
}

type StatisticResponse struct {
	TotalAchievements int               `json:"totalAchievements"`
	Verified          int               `json:"verified"`
	Draft             int               `json:"draft"`
	Submitted         int               `json:"submitted"`
	Rejected          int               `json:"rejected"`
}

type StudentReportResponse struct {
	StudentID         string `json:"studentId"`
	TotalAchievements int    `json:"totalAchievements"`
	Verified          int    `json:"verified"`
	Draft             int    `json:"draft"`
	Submitted         int    `json:"submitted"`
	Rejected          int    `json:"rejected"`
}
