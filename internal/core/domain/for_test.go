package domain

type Test struct{}

func NewTest() Test {
	return Test{}
}

func NewTestUninitialized(fullName string, phoneNumber *string) Test {
	return NewTest()
}

func (u *Test) Validate() error {
	return nil
}

type TestPatch struct{}

func NewTestPatch() TestPatch {
	return TestPatch{}
}

func (p *TestPatch) Validate() error {
	return nil
}

func (u *Test) ApplyPatch(patch TestPatch) error {
	return nil
}
