package twitch

import "testing"

type TestStruct struct {
	String      string   `query:"string"`
	StringSlice []string `query:"strings"`
	Int         int      `query:"number"`
}

func TestStructToQuery(t *testing.T) {
	s := TestStruct{
		String:      "test string",
		StringSlice: []string{"one", "two"},
		Int:         4,
	}
	values, err := StructToQuery(s)
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	if values.Encode() != "number=4&string=test+string&strings=one&strings=two" {
		t.Fail()
	}

	t.Log(values.Encode())
}
