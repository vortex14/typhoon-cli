package data

import "typhoon-cli/src/interfaces"

type StructData struct {
	Fields []string
}

func (s *StructData) GetFields()  {

}

func TestFunc () interfaces.TestData {
	test := &StructData{
		Fields: []string{"1", "2"},
	}

	return test

}
