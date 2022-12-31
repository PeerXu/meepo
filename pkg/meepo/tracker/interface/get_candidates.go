package tracker_interface

type GetCandidatesRequest struct {
	Target   string
	Count    int
	Excludes []string
}

type GetCandidatesResponse struct {
	Candidates []string
}
