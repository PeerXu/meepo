package tracker_interface

type GetCandidatesRequest struct {
	Protocol string
	Target   string
	Requests int
	Excludes []string
}

type GetCandidatesResponse struct {
	Candidates []string
}
