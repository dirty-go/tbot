package tbot

import (
	"encoding/json"
	"testing"
)

func TestChatID_String(t *testing.T) {
	cases := []struct {
		name string
		id   ChatID
		want string
	}{
		{"int", Int64(42), "42"},
		{"negative supergroup", Int64(-1001234567890), "-1001234567890"},
		{"username", Username("@channel"), "@channel"},
		{"username takes precedence", ChatID{id: 5, username: "@c"}, "@c"},
		{"zero", ChatID{}, "0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.id.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestChatID_MarshalJSON(t *testing.T) {
	cases := []struct {
		name string
		id   ChatID
		want string
	}{
		{"int emits number", Int64(42), `42`},
		{"username emits string", Username("@channel"), `"@channel"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.id)
			if err != nil {
				t.Fatalf("Marshal() error: %v", err)
			}
			if string(b) != tc.want {
				t.Errorf("Marshal() = %s, want %s", b, tc.want)
			}
		})
	}
}

func TestChatID_IsZero(t *testing.T) {
	if !(ChatID{}).IsZero() {
		t.Error("zero value should be zero")
	}
	if Int64(1).IsZero() {
		t.Error("Int64(1) should not be zero")
	}
	if Username("@c").IsZero() {
		t.Error("Username should not be zero")
	}
}
