package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Define what is a "Spatial Piece" for our API
type AerospacePart struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Geometry   string  `json:"geometry"`
	GapSpacing float64 `json:"gap_spacing,omitempty"`
}

// Helpers
// Send format text in jSON with the right HTTP code
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Typer", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error for the JSON
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// The Middleware Cors
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Endpoints
// GET 1 : Send pieces
func getPartsHandler(w http.ResponseWriter, r *http.Request) {
	parts := []AerospacePart{
		{ID: "A6-BST", Name: "Ariane6_Booster", Geometry: "half-cone-rectangular"},
		{ID: "JWST-MIR", Name: "JWST_Mirror_Assembly", GapSpacing: 0.519},
	}
	writeJSON(w, http.StatusOK, parts)
}

// GET 2 : Send telemetry
func getTelemetryHandler(w http.ResponseWriter, r *http.Request) {
	telemetry := map[string]interface{}{
		"status":    "NOMINAL",
		"alititude": 408.5,
		"velocity":  27600,
	}
	writeJSON(w, http.StatusOK, telemetry)
}

// GET 3
func getWelcomeHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		user = "Engineer"
	}
	message := map[string]string{
		"message": fmt.Sprintf("Welcome to Aeroforge Control, %s!", user),
	}
	writeJSON(w, http.StatusOK, message)
}

// POST
func validatePartHandler(w http.ResponseWriter, r *http.Request) {
	var part AerospacePart

	if err := json.NewDecoder(r.Body).Decode(&part); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	if part.Name == "Ariane6_Booster" && part.Geometry != "half-cone-rectangular" {
		writeError(w, http.StatusNotAcceptable, "CRITICAL ERROR: Booster geometry must be a half-cone with rectangular ends.")
		return
	}

	if part.Name == "JWST_Mirror_Assembly" && part.GapSpacing != 0.519 {
		writeError(w, http.StatusNotAcceptable, "CRITICAL ERROR: Mirror gap spacing must be exactly 0.519mm for thermal expansion.")
		return
	}

	// If coorect
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "APPROVED",
		"message": fmt.Sprintf("Design %s passed all tolerances and is ready for manufacturing.", part.Name),
	})
}

// Main func
func main() {
	// Create the router
	mux := http.NewServeMux()

	// Link URLs to the functions
	mux.HandleFunc("GET /parts", getPartsHandler)
	mux.HandleFunc("GET /telemetry", getTelemetryHandler)
	mux.HandleFunc("GET /welcome", getWelcomeHandler)
	mux.HandleFunc("POST /validate", validatePartHandler)

	// ROutes for the documentation
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Starting the server
	fmt.Println("AeroForge Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", corsMiddleware(mux))

	if err != nil {
		fmt.Println("Critical error of the server :", err)
	}

}
