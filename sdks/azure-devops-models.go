package sdks

type WorkItemDetailsResponse struct {
	Id        int                               `json:"id"`
	Rev       int                               `json:"rev"`
	Url       string                            `json:"url"`
	Fields    WorkItemDetailsFieldsResponse     `json:"fields"`
	Relations []WorkItemDetailsRelationResponse `json:"relations"`
}

type WorkItemDetailsFieldsResponse struct {
	AreaPath      string `json:"System.AreaPath"`
	TeamProject   string `json:"System.TeamProject"`
	IterationPath string `json:"System.IterationPath"`
	State         string `json:"System.State"`
	WorkItemType  string `json:"System.WorkItemType"`
	Title         string `json:"System.Title"`
	Description   string `json:"System.Description"`
}

type WorkItemDetailsRelationResponse struct {
	Rel string `json:"rel"`
	Url string `json:"url"`
}

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
