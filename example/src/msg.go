package src

type HandleMsg struct {
	Register *RegisterDomain `json:"register,omitempty"`
	Sell     *SellDomain     `json:"sell,omitempty"`
}

type SellDomain struct {
	Buyer  string `json:"buyer"`
	Domain string `json:"domain"`
}

type RegisterDomain struct {
	Domain string `json:"domain"`
}

type QueryMsg struct {
	Get *GetOwner `json:"get,omitempty"`
}

type GetOwner struct {
	Domain string `json:"domain"`
}
