package main

import (
    "encoding/json"
    "io/ioutil"
    "net/http"

    "gopkg.in/yaml.v2"
)

type Config struct {
    Name    string `yaml:"name" json:"name"`
    Version string `yaml:"version" json:"version"`
    Settings struct {
        Enabled bool `yaml:"enabled" json:"enabled"`
        Timeout int  `yaml:"timeout" json:"timeout"`
    } `yaml:"settings" json:"settings"`
}

func readYAML(filePath string) (Config, error) {
    var config Config
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return config, err
    }
    err = yaml.Unmarshal(data, &config)
    return config, err
}

func writeYAML(filePath string, config Config) error {
    data, err := yaml.Marshal(&config)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(filePath, data, 0644)
}

func getConfigHandler(w http.ResponseWriter, r *http.Request) {
    config1, _ := readYAML("config1.yaml")
    config2, _ := readYAML("config2.yaml")
    response := map[string]Config{"config1": config1, "config2": config2}
    json.NewEncoder(w).Encode(response)
}

func saveConfigHandler(w http.ResponseWriter, r *http.Request) {
    var configData map[string]string
    if err := json.NewDecoder(r.Body).Decode(&configData); err != nil {
        http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Write the updated configuration to the specified YAML file
    fileName := configData["fileName"]
    content := configData["content"]

    // Write content to the file
    if err := ioutil.WriteFile(fileName, []byte(content), 0644); err != nil {
        http.Error(w, "Failed to write "+fileName+": "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte("File updated successfully: " + fileName))
}

func main() {
    // Serve static files
    fs := http.FileServer(http.Dir("./")) // Serve files from the current directory
    http.Handle("/", fs)

    // API routes
    http.HandleFunc("/api/config", getConfigHandler)
    http.HandleFunc("/api/config/save", saveConfigHandler)

    // Start the server
    http.ListenAndServe(":3000", nil)
}
