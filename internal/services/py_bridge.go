package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// PythonBridge provides an interface to Python scripts for ML model inference
type PythonBridge struct {
	initialized      bool
	pythonExecutable string
	initMutex        sync.Mutex
	scriptDir        string
	virtualEnvPath   string
}

// PredictionResult represents the output from Python prediction model
type PredictionResult struct {
	PredictedPrice float64                 `json:"predicted_price"`
	Confidence     float64                 `json:"confidence"`
	Direction      string                  `json:"direction"`
	Factors        []string                `json:"factors"`
	Technical      map[string]interface{}  `json:"technical,omitempty"`
	Error          string                  `json:"error,omitempty"`
}

var defaultBridge *PythonBridge

// GetPythonBridge returns the shared PythonBridge instance
func GetPythonBridge() *PythonBridge {
	if defaultBridge == nil {
		defaultBridge = NewPythonBridge()
	}
	return defaultBridge
}

// NewPythonBridge creates a new Python bridge
func NewPythonBridge() *PythonBridge {
	return &PythonBridge{
		initialized:      false,
		pythonExecutable: detectPythonExecutable(),
		scriptDir:        detectScriptDirectory(),
		virtualEnvPath:   detectVirtualEnvPath(),
	}
}

// Initialize checks if Python is available and required packages are installed
func (pb *PythonBridge) Initialize() error {
	pb.initMutex.Lock()
	defer pb.initMutex.Unlock()
	
	if pb.initialized {
		return nil
	}
	
	if pb.pythonExecutable == "" {
		return fmt.Errorf("Python executable not found")
	}
	
	// Check if the virtual environment is active
	if pb.virtualEnvPath != "" {
		log.Printf("Using Python virtual environment: %s", pb.virtualEnvPath)
	}
	
	// Check if the simple analyzer script exists
	analyzerPath := filepath.Join(pb.scriptDir, "simple_analyzer.py")
	if _, err := os.Stat(analyzerPath); os.IsNotExist(err) {
		return fmt.Errorf("simple_analyzer.py script not found at %s", analyzerPath)
	}
	
	// Try a simple import test with python
	cmd := exec.Command(pb.pythonExecutable, "-c", "import numpy; import pandas; import yfinance; print('OK')")
	
	// Set the virtual environment's Python if available
	if pb.virtualEnvPath != "" {
		venvPython := filepath.Join(pb.virtualEnvPath, "bin", "python")
		if _, err := os.Stat(venvPython); err == nil {
			cmd = exec.Command(venvPython, "-c", "import numpy; import pandas; import yfinance; print('OK')")
		}
	}
	
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dependency check failed: %v: %s", err, stderr.String())
	}
	
	pb.initialized = true
	return nil
}

// PredictStockPrice predicts the stock price for a given ticker
func (pb *PythonBridge) PredictStockPrice(ticker string) (*PredictionResult, error) {
	if !pb.initialized {
		if err := pb.Initialize(); err != nil {
			return nil, fmt.Errorf("bridge not initialized: %v", err)
		}
	}
	
	// Validate ticker
	ticker = strings.TrimSpace(ticker)
	if ticker == "" {
		return nil, fmt.Errorf("empty ticker")
	}
	
	analyzerPath := filepath.Join(pb.scriptDir, "simple_analyzer.py")
	
	// Prepare the command
	var cmd *exec.Cmd
	
	// If we have a virtual environment, use its Python
	if pb.virtualEnvPath != "" {
		venvPython := filepath.Join(pb.virtualEnvPath, "bin", "python")
		if _, err := os.Stat(venvPython); err == nil {
			cmd = exec.Command(venvPython, analyzerPath, ticker)
		} else {
			cmd = exec.Command(pb.pythonExecutable, analyzerPath, ticker)
		}
	} else {
		cmd = exec.Command(pb.pythonExecutable, analyzerPath, ticker)
	}
	
	// Execute the command
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("prediction failed: %v: %s", err, stderr.String())
	}
	
	// Parse the result
	var result PredictionResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse prediction result: %v: %s", err, stdout.String())
	}
	
	if result.Error != "" {
		return nil, fmt.Errorf("prediction error: %s", result.Error)
	}
	
	return &result, nil
}

// PredictStockPriceWithSimpleAnalyzer predicts the stock price for a given ticker using simple_analyzer.py
func (pb *PythonBridge) PredictStockPriceWithSimpleAnalyzer(ticker string) (*PredictionResult, error) {
	if !pb.initialized {
		if err := pb.Initialize(); err != nil {
			return nil, fmt.Errorf("bridge not initialized: %v", err)
		}
	}
	
	// Validate ticker
	ticker = strings.TrimSpace(ticker)
	if ticker == "" {
		return nil, fmt.Errorf("empty ticker")
	}
	
	analyzerPath := filepath.Join(pb.scriptDir, "simple_analyzer.py")
	
	// Prepare the command
	var cmd *exec.Cmd
	
	// If we have a virtual environment, use its Python
	if pb.virtualEnvPath != "" {
		venvPython := filepath.Join(pb.virtualEnvPath, "bin", "python")
		if _, err := os.Stat(venvPython); err == nil {
			cmd = exec.Command(venvPython, analyzerPath, ticker)
		} else {
			cmd = exec.Command(pb.pythonExecutable, analyzerPath, ticker)
		}
	} else {
		cmd = exec.Command(pb.pythonExecutable, analyzerPath, ticker)
	}
	
	// Execute the command
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("prediction failed: %v: %s", err, stderr.String())
	}
	
	// Parse the result
	var result PredictionResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse prediction result: %v: %s", err, stdout.String())
	}
	
	if result.Error != "" {
		return nil, fmt.Errorf("prediction error: %s", result.Error)
	}
	
	return &result, nil
}

// Helper functions

// detectPythonExecutable tries to find the Python executable
func detectPythonExecutable() string {
	// Try several common Python executable names
	pythons := []string{"python3", "python"}
	
	for _, python := range pythons {
		path, err := exec.LookPath(python)
		if err == nil {
			return path
		}
	}
	
	return ""
}

// detectScriptDirectory finds where the Python scripts are located
func detectScriptDirectory() string {
	// Try several common relative paths
	candidates := []string{
		"models",                // Run from project root
		"../models",             // Run from the internal dir
		"../../models",          // Run from internal/services
	}
	
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	
	for _, candidate := range candidates {
		path := filepath.Join(cwd, candidate)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	// Fall back to the executable's directory
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	exeDir := filepath.Dir(exePath)
	
	for _, candidate := range candidates {
		path := filepath.Join(exeDir, candidate)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	return ""
}

// detectVirtualEnvPath tries to find the Python virtual environment
func detectVirtualEnvPath() string {
	// Try to detect a virtual environment
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	
	// Check for virtual environment in the project root
	candidates := []string{
		filepath.Join(cwd, ".venv"),                  // Run from project root (.venv)
		filepath.Join(cwd, "venv"),                   // Run from project root (venv)
		filepath.Join(cwd, "..", ".venv"),            // Run from subdirectory (.venv)
		filepath.Join(cwd, "..", "venv"),             // Run from subdirectory (venv)
		filepath.Join(cwd, "../..", ".venv"),         // Run from subsubdirectory (.venv)
		filepath.Join(cwd, "../..", "venv"),          // Run from subsubdirectory (venv)
	}
	
	for _, candidate := range candidates {
		// Look for bin/python to confirm it's a valid venv
		binPython := filepath.Join(candidate, "bin", "python")
		if _, err := os.Stat(binPython); err == nil {
			return candidate
		}
	}
	
	// Check if we're already in a virtual environment
	if os.Getenv("VIRTUAL_ENV") != "" {
		return os.Getenv("VIRTUAL_ENV")
	}
	
	return ""
}
