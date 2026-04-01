package dto

type GenerateInput struct {
	Password   string `json:"password"`
	MasterPass string `json:"master_pass"`
}

type GenerateOutput struct {
	SecurePassword string `json:"secure_password"`
	Length         int    `json:"length"`
}

type SaveInput struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	MasterPass string `json:"master_pass"`
}

type SaveOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ShowInput struct {
	Name       string `json:"name"`
	MasterPass string `json:"master_pass"`
}

type ShowOutput struct {
	Password string `json:"password"`
	Copied   bool   `json:"copied"`
}
