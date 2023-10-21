package task

import "github.com/xh3b4sd/rescue/matcher"

type Host map[string]string

func (h *Host) All(key ...string) *Host {
	hos := Host(matcher.All(*h, key...))
	return &hos
}

func (h *Host) Any(key ...string) *Host {
	hos := Host(matcher.Any(*h, key...))
	return &hos
}

func (h *Host) Emp() bool {
	return h.Len() == 0
}

func (h *Host) Eql(x *Host) bool {
	return h != nil && x != nil && h.Len() == x.Len() && h.Has(*x)
}

func (h *Host) Exi(key string) bool {
	if h == nil {
		return false
	}

	hos := *h
	return key != "" && hos[key] != ""
}

func (h *Host) Get(key string) string {
	if h == nil {
		return ""
	}

	hos := *h
	return hos[key]
}

func (h *Host) Has(lab map[string]string) bool {
	return matcher.Has(*h, lab)
}

func (h *Host) Key() []string {
	if h == nil {
		return nil
	}

	var key []string

	for k := range *h {
		key = append(key, k)
	}

	return key
}

func (h *Host) Len() int {
	if h == nil {
		return 0
	}

	hos := *h
	return len(hos)
}

func (h *Host) Set(key string, val string) {
	hos := *h
	hos[key] = val
}
