package entity

//type Data struct {
//	Id        int    `tsv:"n"`
//	Mqtt      string `tsv:"mqtt"`
//	Invid     string `tsv:"invid"`
//	UnitGuid  string `tsv:"unit_guid"`
//	MsgId     string `tsv:"msg_id"`
//	MsgText   string `tsv:"text"`
//	Context   string `tsv:"context"`
//	Class     string `tsv:"class"`
//	Level     string `tsv:"level"`
//	Area      string `tsv:"area"`
//	Addr      string `tsv:"addr"`
//	Block     string `tsv:"block"`
//	DataType  string `tsv:"type"`
//	Bit       string `tsv:"bit"`
//	InvertBit string `tsv:"invert_bit"`
//}

type Data struct {
	Id        int
	Mqtt      string
	Invid     string
	UnitGuid  string
	MsgId     string
	MsgText   string
	Context   string
	Class     string
	Level     string
	Area      string
	Addr      string
	Block     string
	DataType  string
	Bit       string
	InvertBit string
}
