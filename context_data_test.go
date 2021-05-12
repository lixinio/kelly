package kelly

import (
	"testing"
)

func TestContextData(t *testing.T) {
	contextData := newContextMapData()
	v := contextData.GetDefault("a", "b")
	if v != "b" {
		t.Errorf("ContextData GetDefault fail %v", v)
		return
	}

	vi := contextData.Get("a")
	if vi != nil {
		t.Errorf("ContextData Get fail %v", vi)
		return
	}

	contextData.Set("a", "c")

	v = contextData.GetDefault("a", "b")
	if v != "c" {
		t.Errorf("ContextData GetDefault fail %v", v)
		return
	}

	vi = contextData.Get("a")
	if vi != "c" {
		t.Errorf("ContextData Get fail %v", vi)
		return
	}

	vi = contextData.MustGet("a")
	if v != "c" {
		t.Errorf("ContextData MustGet fail %v", v)
		return
	}

	// 抛异常
	defer checkError(t, ErrNoContextData)
	contextData.MustGet("b")
}
