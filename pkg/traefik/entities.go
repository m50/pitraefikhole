package traefik

import "strings"

type RouterList []Router

func (l RouterList) ToHosts() []string {
	hosts := make([]string, len(l))
	for i, router := range l {
		hosts[i] = router.Host()
	}
	return hosts
}

type Router struct {
	status string
	rule   string
}

func (r *Router) IsHost() bool {
	return strings.Contains(r.rule, "Host(") && r.status == "enabled"
}

func (r *Router) Host() string {
	h := strings.Replace(r.rule, "Host", "", -1)
	h = strings.Trim(h, "`()")
	return h
}
