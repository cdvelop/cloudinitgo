package cloudinitgo

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestCloudInit(t *testing.T) {
	// Configuración de prueba
	testDir := "testdata"
	config := Config{
		UserName:     "testuser",
		Password:     "testpass",
		SshPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
		DataDir:      testDir,
	}

	// Limpiar directorio de prueba
	defer os.RemoveAll(testDir)

	// Crear nueva instancia
	ci := New(config)

	// Test: Crear archivos de configuración
	err := ci.CreateConfigFiles()
	if err != nil {
		t.Fatalf("CreateConfigFiles failed: %v", err)
	}

	// Verificar contenido del archivo YML
	userDataPath := filepath.Join(testDir, "user-data")
	if _, err := os.Stat(userDataPath); os.IsNotExist(err) {
		t.Errorf("user-data file not created")
	}
	// Read and verify content
	content, err := os.ReadFile(userDataPath)
	if err != nil {
		t.Fatalf("Error reading user-data: %v", err)
	}

	newExpectedConfig := NewCloudConfig(config)
	expectedContent, err := ci.MarshalYAML(newExpectedConfig)
	if err != nil {
		t.Fatalf("Error marshaling expected YAML: %v", err)
	}

	normalizedExpected := string(expectedContent)
	normalizedActual := string(content)

	if normalizedActual != normalizedExpected {
		t.Errorf("YAML content mismatch:\nGot:\n%s\nWant:\n%s", normalizedActual, normalizedExpected)
	}
	// Test: Iniciar servidor
	err = ci.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer ci.Stop()

	// Test: Verificar que el servidor responde
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ci.createMux().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != normalizedExpected {
		t.Errorf("Unexpected response body:\nGot:\n%s\nWant:\n%s", w.Body.String(), expectedContent)
	}
}
