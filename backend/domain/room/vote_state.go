package room

type VoteStatus int

const (
	Upvoted   VoteStatus = iota
	Downvoted VoteStatus = iota
	NotVoted  VoteStatus = iota
)

var voteStatusToString = map[VoteStatus]string{
	Upvoted:   "upvoted",
	Downvoted: "downvoted",
	NotVoted:  "not_voted",
}
var stringToVoteStatus = map[string]VoteStatus{
	"upvoted":   Upvoted,
	"downvoted": Downvoted,
	"not_voted": NotVoted,
}

func (v VoteStatus) String() string {
	return voteStatusToString[v]
}
func ParseVoteStatus(s string) *VoteStatus {
	voteStatus, ok := stringToVoteStatus[s]
	if !ok {
		return nil
	}
	return &voteStatus
}
