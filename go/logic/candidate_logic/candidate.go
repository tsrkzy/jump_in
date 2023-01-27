package candidate_logic

import (
	"context"
	"database/sql"
	"github.com/tsrkzy/jump_in/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

//func BulkInsertCandidate(ctx context.Context, tx *sql.Tx, eventId int64, openAtList []string) error {
//	q := `INSERT INTO candidate (event_id, open_at)
// VALUES %s
// ON CONFLICT DO NOTHING`
//	valueList := make([]string, 0)
//	for _, openAt := range openAtList {
//		value := fmt.Sprintf("( %d , %s )", eventId, escape.Literal(openAt))
//		valueList = append(valueList, value)
//	}
//	values := strings.Join(valueList, ", ")
//	query := fmt.Sprintf(q, values)
//
//	//lg.Debug(query)
//
//	_, err := queries.Raw(query).ExecContext(ctx, tx)
//
//	return err
//}

func FetchCandidateByDate(ctx context.Context, tx *sql.Tx, eId int64, openAt string) (*models.Candidate, error) {
	return models.Candidates(qm.Where("event_id = ? and open_at = ?", eId, openAt)).One(ctx, tx)
}

//func CandidateExistsByDate(ctx context.Context, tx *sql.Tx, eId int64, openAt string) (bool, error) {
//	return models.Candidates(qm.Where("event_id = ? and open_at = ?", eId, openAt)).Exists(ctx, tx)
//}

func FetchCandidateByID(ctx context.Context, tx *sql.Tx, cId int64, eId int64) (*models.Candidate, error) {
	return models.Candidates(qm.Where("event_id = ? and id = ?", eId, cId)).One(ctx, tx)
}
func CandidateExistsByID(ctx context.Context, tx *sql.Tx, cId int64, eId int64) (bool, error) {
	return models.Candidates(qm.Where("event_id = ? and id = ?", eId, cId)).Exists(ctx, tx)
}

func CreateCandidate(ctx context.Context, tx *sql.Tx, eventID int64, openAt string) (*models.Candidate, error) {
	c := models.Candidate{EventID: eventID, OpenAt: openAt}
	err := c.Insert(ctx, tx, boil.Infer())
	return &c, err
}

//func FetchVoteByUK(ctx context.Context, tx *sql.Tx, accountID int64, cId int64) (*models.Vote, error) {
//	v, err := models.Votes(qm.Where("account_id = ? and candidate_id = ?", accountID, cId)).One(ctx, tx)
//	return v, err
//}

func VoteExistsByUK(ctx context.Context, tx *sql.Tx, accountID int64, cId int64) (bool, error) {
	exists, err := models.Votes(qm.Where("account_id = ? and candidate_id = ?", accountID, cId)).Exists(ctx, tx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func CreateVote(ctx context.Context, tx *sql.Tx, accountID int64, cId int64) (*models.Vote, error) {
	v := models.Vote{
		AccountID:   accountID,
		CandidateID: cId,
	}
	err := v.Insert(ctx, tx, boil.Infer())

	return &v, err
}

func DeleteVote(ctx context.Context, tx *sql.Tx, accountID int64, candidateId int64) error {
	_, err := models.Votes(qm.Where("account_id = ? and candidate_id = ?", accountID, candidateId)).DeleteAll(ctx, tx)
	return err
}
