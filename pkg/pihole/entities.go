package pihole

type Session struct {
	Valid    bool
	Validity int
	SID      string
	Message  string
}

type AuthResponse struct {
	Session Session
}

type DNS struct {
	CNameRecords []string `json:"cnameRecords"`
}

type Config struct {
	DNS DNS `json:"dns"`
}

type ConfigResponse struct {
	Config Config
}
