package matcher

import "strings"

// MatchPath verifica se um padrao tipo "/accounts/:accountId/balances" casa
// com um path concreto tipo "/accounts/abc-123/balances".
// Retorna os parametros extraidos e true em caso de match.
//
// Suporta segmentos parametrizados no formato ":nome". Nao suporta wildcards.
func MatchPath(pattern, path string) (map[string]string, bool) {
	pSegs := splitPath(pattern)
	aSegs := splitPath(path)
	if len(pSegs) != len(aSegs) {
		return nil, false
	}
	params := make(map[string]string)
	for i, seg := range pSegs {
		if strings.HasPrefix(seg, ":") {
			params[seg[1:]] = aSegs[i]
			continue
		}
		if seg != aSegs[i] {
			return nil, false
		}
	}
	return params, true
}

func splitPath(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}
	return strings.Split(p, "/")
}
