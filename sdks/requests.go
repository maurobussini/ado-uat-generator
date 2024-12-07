package sdks

type WorkItemAddFieldRequest struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

type WorkItemAddPlainFieldRequest struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

type WorkItemAddComplexFieldRequest struct {
	Op    string                          `json:"op"`
	Path  string                          `json:"path"`
	Value WorkItemAddRelationFieldRequest `json:"value"`
}

type WorkItemAddRelationFieldRequest struct {
	Rel string `json:"rel"`
	Url string `json:"url"`
}

type IWorkItemFieldRequest interface {
	GetOp() string
	GetPath() string
}

func (in WorkItemAddPlainFieldRequest) GetOp() string {
	return in.Op
}

func (in WorkItemAddComplexFieldRequest) GetOp() string {
	return in.Op
}

func (in WorkItemAddPlainFieldRequest) GetPath() string {
	return in.Path
}

func (in WorkItemAddComplexFieldRequest) GetPath() string {
	return in.Path
}
