package dormitoryElectricity

type buildCode struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	Roomlist []struct {
		ID          string      `json:"id"`
		Name        string      `json:"name"`
		Factorycode interface{} `json:"factorycode"`
	} `json:"roomlist"`
	Custparam string `json:"custparam"`
}

type roomStats struct {
	Returncode   string      `json:"returncode"`
	Returnmsg    string      `json:"returnmsg"`
	Quantity     string      `json:"quantity"`
	Quantityunit string      `json:"quantityunit"`
	Canbuy       string      `json:"canbuy"`
	Description  string      `json:"description"`
	Custparams   interface{} `json:"custparams"`
}
