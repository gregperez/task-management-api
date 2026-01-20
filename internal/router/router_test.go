package router

import "testing"

func TestRouter_New(t *testing.T) {
	// Configuración mínima para las pruebas
	cfg := &Config{}

	r := New(cfg)
	if r == nil {
		t.Fatal("Router instance is nil")
	}
}
