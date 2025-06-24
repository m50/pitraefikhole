package pihole

type Session struct {
	Valid    bool   `json:"valid"`
	Validity int    `json:"validity"`
	SID      string `json:"sid"`
	Message  string `json:"message"`
}

type AuthResponse struct {
	Session Session `json:"session"`
}

type DNS struct {
	CNameRecords []string `json:"cnameRecords"`
}

type Config struct {
	DNS DNS `json:"dns"`
}

type ConfigResponse struct {
	Config Config `json:"config"`
}
