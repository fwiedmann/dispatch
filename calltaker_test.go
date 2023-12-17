package dispatch_test

import (
	"fmt"
	"github.com/fwiedmann/dispatch"
	mocks "github.com/fwiedmann/dispatch/mock"
	"github.com/fwiedmann/dispatch/pkg/code"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup(t *testing.T, owners []code.Owner) (dispatch.CallTaker, *gomock.Controller, *mocks.Mockanalyzer, *mocks.Mockdispatcher) {
	c := gomock.NewController(t)
	analyzer := mocks.NewMockanalyzer(c)
	dispatcher := mocks.NewMockdispatcher(c)
	callTaker := dispatch.NewCallTaker(analyzer, dispatcher, owners)

	return callTaker, c, analyzer, dispatcher
}

func TestCallTaker_Dispose_return_error_when_analyzer_fails(t *testing.T) {
	// given
	calltaker, mockcontroller, analyzerMock, _ := setup(t, make([]code.Owner, 0))
	defer mockcontroller.Finish()

	analyzerMock.EXPECT().Analyze().Return(nil, fmt.Errorf("analyzer error"))

	// when // then
	assert.Error(t, calltaker.Dispose())
}

func TestCallTaker_Dispose_successful(t *testing.T) {
	// given

	owner := code.Owner{
		Id:      "1",
		Name:    "number uno",
		Members: []string{"felix", "alex", "ferenc"},
	}

	calltaker, mockcontroller, analyzerMock, dispatcherMock := setup(t, []code.Owner{owner})
	defer mockcontroller.Finish()

	analyzerMock.EXPECT().Analyze().Return([]code.Info{
		{
			OwnerId:      owner.Id,
			LocationName: "tour de offenburg",
		},
	}, nil)

	dispatcherMock.EXPECT().Dispatch([]dispatch.Note{
		{
			Info: code.Info{
				OwnerId:      owner.Id,
				LocationName: "tour de offenburg",
			},
			Owner: owner,
		},
	}).Return(nil)

	// when // then
	assert.NoError(t, calltaker.Dispose())
}

func TestCallTaker_Dispose_return_error_on_dispatch(t *testing.T) {
	// given

	owner := code.Owner{
		Id:      "1",
		Name:    "number uno",
		Members: []string{"felix", "alex", "ferenc"},
	}

	calltaker, mockcontroller, analyzerMock, dispatcherMock := setup(t, []code.Owner{owner})
	defer mockcontroller.Finish()

	analyzerMock.EXPECT().Analyze().Return([]code.Info{
		{
			OwnerId:      owner.Id,
			LocationName: "tour de offenburg",
		},
	}, nil)

	dispatcherMock.EXPECT().Dispatch([]dispatch.Note{
		{
			Info: code.Info{
				OwnerId:      owner.Id,
				LocationName: "tour de offenburg",
			},
			Owner: owner,
		},
	}).Return(fmt.Errorf("error on dispatching"))

	// when // then
	assert.Error(t, calltaker.Dispose())
}
