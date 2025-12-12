package model

type LoginReq struct {
    Username string `json:"username" example:"admin"`
    Password string `json:"password" example:"123"`
}

type LoginResponse struct {
    Token        string   `json:"token"`
    RefreshToken string   `json:"refreshToken"`
    ID           string   `json:"id"`
    Username     string   `json:"username"`
    FullName     string   `json:"fullName"`
    Role         string   `json:"role"`
    Permissions  []string `json:"permissions"`
}

type RefreshReq struct {
    RefreshToken string `json:"refreshToken"`
}
