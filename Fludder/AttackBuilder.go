package Fludder

type Attack struct {
	Url           string
	Host          string
	AttackMethod  string
	PostData      interface{}
	RequestsPerIP int
	Cookie        interface{}
	Ja3           string
}

type System struct {
	Banner       string
	HTTP2Timeout int
	Attack       *Attack
}
