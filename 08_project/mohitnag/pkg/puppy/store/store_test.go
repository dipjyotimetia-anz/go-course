package store

import (
	"fmt"
	"testing"

	puppy "github.com/anz-bank/go-course/08_project/mohitnag/pkg/puppy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type storerSuite struct {
	suite.Suite
	store      puppy.Storer
	makeStorer func() puppy.Storer
	puppy      puppy.Puppy
}

func (s *storerSuite) SetupTest() {
	s.store = s.makeStorer()
	s.puppy = puppy.Puppy{ID: 1, Breed: "dog", Colour: "white", Value: "$2"}
	err := s.store.CreatePuppy(s.puppy)
	if err != nil {
		s.FailNow("Error in Test Setup")
	}
}

func (s *storerSuite) TestCreatePuppy() {
	assert := assert.New(s.T())
	s.puppy = puppy.Puppy{ID: 2, Breed: "dog", Colour: "black", Value: "2"}
	freePuppy := puppy.Puppy{ID: 3, Breed: "dog", Colour: "white", Value: "-1"}
	testCases := []struct {
		Scenario      string
		Input         puppy.Puppy
		ExpectedError string
	}{
		{"Create Puppy", s.puppy, "<nil>"},
		{"Puppy with value less than 0 results error", freePuppy,
			"puppy with value less than 0 not allowed :(code: invalid input)"},
		{"Creating already existing Puppy should fail", s.puppy,
			"puppy with Id 2 already exists :(code: puppy already exists)"},
	}
	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.Scenario, func(t *testing.T) {
			err := s.store.CreatePuppy(tc.Input)
			assert.Equal(tc.ExpectedError, fmt.Sprint(err))
		})
	}
}

func (s *storerSuite) TestReadPuppy() {
	assert := assert.New(s.T())
	testCases := []struct {
		Scenario      string
		ID            uint32
		Expected      puppy.Puppy
		ExpectedError string
	}{
		{"Read already existing Puppy", 1, s.puppy, "<nil>"},
		{"Reading a non-existing Puppy should fail", 32, puppy.Puppy{},
			"puppy with Id 32 does not exists :(code: puppy not found)"},
	}
	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.Scenario, func(t *testing.T) {
			puppy, err := s.store.ReadPuppy(tc.ID)
			assert.Equal(tc.Expected, puppy)
			assert.Equal(tc.ExpectedError, fmt.Sprint(err))
		})
	}
}

func (s *storerSuite) TestUpdatePuppy() {
	assert := assert.New(s.T())
	puppyUpdate := puppy.Puppy{ID: 1, Breed: "dog", Colour: "black", Value: "2"}
	puppyNonExist := puppy.Puppy{ID: 32, Breed: "dog", Colour: "black", Value: "2"}

	testCases := []struct {
		Scenario      string
		Puppy         puppy.Puppy
		ExpectedError string
	}{
		{"Update already existing Puppy", puppyUpdate, "<nil>"},
		{"Update a non-existing Puppy should fail", puppyNonExist,
			"puppy with Id 32 does not exists :(code: puppy not found)"},
	}
	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.Scenario, func(t *testing.T) {
			err := s.store.UpdatePuppy(tc.Puppy)
			if err == nil {
				puppy, _ := s.store.ReadPuppy(tc.Puppy.ID)
				assert.Equal("black", puppy.Colour)
			}
			assert.Equal(tc.ExpectedError, fmt.Sprint(err))
		})
	}

}

func (s *storerSuite) TestDeletePuppy() {
	assert := assert.New(s.T())
	testCases := []struct {
		Scenario      string
		ID            uint32
		ExpectedError string
	}{
		{"Delete already existing Puppy", 1, "<nil>"},
		{"Delete a non-existing Puppy should fail", 32, "puppy with Id 32 does not exists :(code: puppy not found)"},
	}
	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.Scenario, func(t *testing.T) {
			err := s.store.DeletePuppy(tc.ID)
			assert.Equal(tc.ExpectedError, fmt.Sprint(err))
		})
	}
}

func TestStorers(t *testing.T) {
	suite.Run(t, &storerSuite{
		makeStorer: func() puppy.Storer { return &MapStore{} },
	})
	suite.Run(t, &storerSuite{
		makeStorer: func() puppy.Storer { return &SyncStore{} },
	})
}
