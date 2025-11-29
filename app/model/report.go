package model

type StatisticByType struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type StatisticByLevel struct {
	Level string `json:"level"`
	Count int    `json:"count"`
}

type StatisticByPeriod struct {
	Period string `json:"period"`
	Count  int    `json:"count"`
}

type TopStudent struct {
	StudentID string `json:"studentId"`
	Count     int    `json:"count"`
	Points    int    `json:"points"`
}

type StatisticResponse struct {
	TotalAchievements int                 `json:"totalAchievements"`
	Verified          int                 `json:"verified"`
	Draft             int                 `json:"draft"`
	Submitted         int                 `json:"submitted"`
	Rejected          int                 `json:"rejected"`
	ByType            []StatisticByType   `json:"byType"`
	ByPeriod          []StatisticByPeriod `json:"byPeriod"`
	TopStudents       []TopStudent        `json:"topStudents"`
	ByLevel           []StatisticByLevel  `json:"byLevel"`
}

type StudentReportResponse struct {
	StudentID         string `json:"studentId"`
	TotalAchievements int    `json:"totalAchievements"`
	Verified          int    `json:"verified"`
	Draft             int    `json:"draft"`
	Submitted         int    `json:"submitted"`
	Rejected          int    `json:"rejected"`
}
