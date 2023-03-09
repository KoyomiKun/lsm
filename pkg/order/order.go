package order

type Order interface {
	Exec() Resp
}

type Resp struct {
	Status uint8
	Msg    string
	Data   [][]byte
}

type OrderAdd struct {
	Key   []byte
	Value []byte
}

type OrderDelete struct {
	Key []byte
}

type OrderUpdate struct {
	Key      []byte
	NewValue []byte
}

type OrderGet struct {
	Key []byte
}
