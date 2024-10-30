package manager

//import "github.com/harness/runner/utils"
//
//type ManagerClient struct {
//	utils.HTTPClient
//	AccountID  string
//	Token      string
//	TokenCache *delegate.TokenCache
//}

// this is not needed because manager client from runner would be set to this
// or you can use the runner code for all this and only focus on the logger for now
//func NewManagerClient(endpoint, id, secret string, skipverify bool, additionalCertsDir string) *ManagerClient {
//	return &ManagerClient{
//		HTTPClient: *utils.New(endpoint, skipverify, additionalCertsDir),
//		AccountID:  id,
//		TokenCache: delegate.NewTokenCache(id, secret),
//	}
//}
