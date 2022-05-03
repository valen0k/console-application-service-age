package test

type test1 struct {
	Key string `json:"key"`
	Val int64  `json:"val"`
}

type test2 struct {
	SLine string `json:"s"`
	Key   string `json:"key"`
}

type test3 struct {
	ALine string `json:"a"`
	BLine string `json:"b"`
	Key   string `json:"key"`
}
