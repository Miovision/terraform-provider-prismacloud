package prismacloud_client

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token         string     `json:"token"`
	Message       string     `json:"message"`
	CustomerNames []Customer `json:"customerNames"`
}

type Customer struct {
	CustomerName string `json:"customerName"`
	TosAccepted  bool   `json:"tosAccepted"`
}

type AccessKey struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	CreatedBy    string `json:"CreatedBy"`
	CreatedTs    string `json:"createdTs"`
	LastUsedTime string `json:"lastUsedTime"`
	Status       string `json:"status"`
	ExpiresOn    string `json:"expiresOn"`
	Role         Role   `json:"role"`
	RoleType     string `json:"roleType"`
}

type Role struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AccountRequest struct {
	AccountId  string   `json:"accountId"`
	Enabled    bool     `json:"enabled"`
	ExternalId string   `json:"externalId"`
	GroupIds   []string `json:"groupIds"`
	Name       string   `json:"name"`
	RoleArn    string   `json:"roleArn"`
}

type Account struct {
	Name           string   `json:"name"`
	Enabled        bool     `json:"enabled"`
	CloudType      string   `json:"cloudType"`
	LastModifiedTs int      `json:"lastModifiedTs"`
	LastModifiedBy string   `json:"lastModifiedBy"`
	IngestionMode  int      `json:"ingestionMode"`
	AccountType    string   `json:"accountType"`
	GroupIds       []string `json:"groupIds"`
	AccountId      string   `json:"accountId"`
	AddedOn        int      `json:"addedOn"`
	ExternalId     string   `json:"externalId"`
	RoleArn        string   `json:"roleArn"`
}

type AccountGroupRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	AccountIds  []string `json:"accountIds"`
}

type AccountGroup struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	LastModifiedTs int      `json:"lastModifiedTs"`
	LastModifiedBy string   `json:"lastModifiedBy"`
	AccountIds     []string `json:"accountIds"`
}

//Implements the sort.Interface for an []AccountGroup based on name
type AccountGroupByName []AccountGroup

func (a AccountGroupByName) Len() int           { return len(a) }
func (a AccountGroupByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AccountGroupByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
