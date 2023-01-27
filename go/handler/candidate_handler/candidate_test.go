package candidate_handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tsrkzy/jump_in/helper/testhelper"
	"github.com/tsrkzy/jump_in/types/candidate_types"
	"github.com/tsrkzy/jump_in/types/event_types"
	"gopkg.in/resty.v1"
	"net/http"
	"testing"
)

//const testDebug = false
const DEFAULT_EMAIL = "tsrmix+echo@gmail.com"

func TestCreate(t *testing.T) {
	respMl1, w, err := testhelper.Login(t, DEFAULT_EMAIL)
	accountId := fmt.Sprintf("%s", w.ID)
	assert.NoError(t, err)

	// イベント作成
	cEc := testhelper.MakeClient(respMl1)
	reqEc := event_types.CreateRequest{
		Name:        "テスト用イベント名",
		Description: "/Users/tsrkzy/dev/go/github.com/tsrkzy/jump_in/go/handler/candidate_handler/candidate_test.go - TestCreate",
		AccountID:   accountId,
	}
	resEc := event_types.CreateResponse{}
	respEc, err := cEc.R().
		SetBody(reqEc).
		SetResult(&resEc).
		Post("http://localhost:80/api/event/create")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respEc.StatusCode())

	// 候補日作成
	eventId := resEc.Event.ID
	clientCandidateCreate := testhelper.MakeClient(respMl1)

	//fmt.Printf("accountId: %s", accountId)
	//fmt.Printf("eventId: %s", eventId)
	reqCandidateCreate := candidate_types.CreateRequest{
		EventID:   eventId,
		AccountID: accountId,
		OpenAt:    "202301241200",
	}
	resCandidateCreate := candidate_types.CreateResponse{}
	respCandidateCreate, err := clientCandidateCreate.R().
		SetBody(&reqCandidateCreate).
		SetResult(&resCandidateCreate).
		Post("http://localhost:80/api/candidate/create")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respCandidateCreate.StatusCode())

	candidateId := resCandidateCreate.Candidates[0].ID

	// 投票
	createVote(t, respMl1, eventId, accountId, candidateId)

	// 二重投票
	createVote(t, respMl1, eventId, accountId, candidateId)

	// 投票の削除
	deleteVote(t, respMl1, eventId, accountId, candidateId)

	// 候補日の削除
	clientCandidateDelete := testhelper.MakeClient(respMl1)

	reqCandidateDelete := candidate_types.DeleteRequest{
		EventID:     eventId,
		CandidateID: candidateId,
		AccountID:   accountId,
	}
	resCandidateDelete := candidate_types.DeleteResponse{}
	respCandidateDelete, err := clientCandidateDelete.R().
		SetBody(&reqCandidateDelete).
		SetResult(&resCandidateDelete).
		Post("http://localhost:80/api/candidate/delete")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respCandidateDelete.StatusCode())

}

func deleteVote(t *testing.T, respMl1 *resty.Response, eventId string, accountId string, candidateId string) {
	clientVoteDelete := testhelper.MakeClient(respMl1)

	reqBodyVoteDelete := candidate_types.DownvoteRequest{
		EventID:     eventId,
		AccountID:   accountId,
		CandidateID: candidateId,
	}
	resBodyVoteDelete := candidate_types.UpvoteResponse{}
	respVoteDelete, err := clientVoteDelete.R().
		SetBody(&reqBodyVoteDelete).
		SetResult(&resBodyVoteDelete).
		Post("http://localhost:80/api/vote/delete")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respVoteDelete.StatusCode())
}

func createVote(t *testing.T, respMl1 *resty.Response, eventId string, accountId string, candidateId string) *candidate_types.UpvoteResponse {
	clientVoteCreate := testhelper.MakeClient(respMl1)

	reqBodyVoteCreate := candidate_types.UpvoteRequest{
		EventID:     eventId,
		AccountID:   accountId,
		CandidateID: candidateId,
	}
	resBodyVoteCreate := candidate_types.UpvoteResponse{}
	respVoteCreate, err := clientVoteCreate.R().
		SetBody(&reqBodyVoteCreate).
		SetResult(&resBodyVoteCreate).
		Post("http://localhost:80/api/vote/create")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respVoteCreate.StatusCode())

	return &resBodyVoteCreate
}
