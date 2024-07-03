package Global

type CodeName struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

type KeyData struct {
	Cons []CodeName `json:"cons"`
	Key  CodeName   `json:"key"`
}

type Data struct {
	K1 KeyData `json:"K1"`
	K2 KeyData `json:"K2"`
	K3 KeyData `json:"K3"`
}
