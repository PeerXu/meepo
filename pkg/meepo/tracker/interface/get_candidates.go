package tracker_interface

type GetCandidatesRequest struct {
	Target   string
	Requests int
	Excludes []string
}

type GetCandidatesResponse struct {
	Candidates []string
}
