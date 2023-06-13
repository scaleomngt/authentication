package model

type CardInfo struct {
	Gender    string `json:"gender"`    //性别
	Name      string `json:"name"`      //姓名
	Nation    string `json:"nation"`    //民族
	Addr      string `json:"addr"`      //地址
	Birthdate string `json:"birthdate"` //出生年月
}

type Transaction struct {
	Id      string `json:"id"`
	Address string `json:"address"`
	Result  string `json:"result"`
}

type Account struct {
	Index      uint32 `json:"index"`      // 序号
	PrivateKey string `json:"privateKey"` // Private Key
	ViewKey    string `json:"viewKey"`    // View Key
	Address    string `json:"address"`    // Address
}
