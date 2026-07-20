package test_service

type TestService struct {
	testRepository TestRepository
}

type TestRepository interface {
	// Methods from repo
}

func NewTestService(testRepository TestRepository) *TestService {
	return &TestService{
		testRepository: testRepository,
	}
}
