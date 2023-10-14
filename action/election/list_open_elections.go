package election

type ListOpenElections struct{}

type ListOpenElectionsResponse struct {
	// TODO: Add support for slices
}

type listOpenElectionsHandler struct{}

func NewListOpenElectionsHandler() *listOpenElectionsHandler {
	return &listOpenElectionsHandler{}
}

func (h *listOpenElectionsHandler) On(query ListOpenElections) (ListOpenElectionsResponse, error) {
	return ListOpenElectionsResponse{}, nil
}
