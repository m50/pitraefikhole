package traefik

import "strings"

type RouterList []Router

func (l RouterList) ToHosts() []string {
	hosts := make([]string, 0, len(l))
	for _, router := range l {
		if !router.IsHost() {
			continue
		}
		hosts = append(hosts, router.Host())
	}
	return hosts
}

type Router struct {
	Status string `json:"status"`
	Rule   string `json:"rule"`
}

func (r *Router) IsHost() bool {
	return strings.Contains(r.Rule, "Host(") && r.Status == "enabled"
}

func (r *Router) Host() string {
	h := strings.Replace(r.Rule, "Host", "", -1)
	h = strings.Trim(h, "`()")
	return h
}
