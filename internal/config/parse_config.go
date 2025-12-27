package config

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "strings"
    
    "gopkg.in/yaml.v3"
    "github.com/kroyser123/ParserGW/internal/models"
)

type JSONPathConfig struct {
    Name        string `yaml:"name"`
    Description string `yaml:"description"`
    Version     string `yaml:"version"`
    Author      string `yaml:"author"`
    Tags        string `yaml:"tags"`
}

func LoadJSONPaths(configPath string) ([]JSONPathConfig, error) {
    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("read JSON paths config error: %v", err)
    }
    
    var config struct {
        Schemas []JSONPathConfig `yaml:"json_schemas"`
    }
    
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("parse JSON paths config error: %v", err)
    }
    
    return config.Schemas, nil
}

func ParseYAMLFile(filePath string) (models.Config, error) {
    var config models.Config
    
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return config, fmt.Errorf("read YAML file error: %v", err)
    }
    
    if err := yaml.Unmarshal(data, &config); err != nil {
        return config, fmt.Errorf("parse YAML error: %v", err)
    }
    
    return config, nil
}

func ParseJSONFile(filePath string, jsonPaths []JSONPathConfig) (models.Config, error) {
    var config models.Config
    var data map[string]interface{}
    
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        return config, fmt.Errorf("read JSON file error: %v", err)
    }
    
    if err := json.Unmarshal(content, &data); err != nil {
        return config, fmt.Errorf("parse JSON error: %v", err)
    }
    
    // Пробуем каждую схему
    for _, schema := range jsonPaths {
        if name := getValueByPath(data, schema.Name); name != nil {
            config.Name = fmt.Sprintf("%v", name)
            config.Description = fmt.Sprintf("%v", getValueByPath(data, schema.Description))
            
            if version := getValueByPath(data, schema.Version); version != nil {
                if v, ok := version.(float64); ok {
                    config.Version = int(v)
                }
            }
            
            config.Author = fmt.Sprintf("%v", getValueByPath(data, schema.Author))
            
            if tags := getValueByPath(data, schema.Tags); tags != nil {
                if tagSlice, ok := tags.([]interface{}); ok {
                    for _, tag := range tagSlice {
                        config.Tags = append(config.Tags, fmt.Sprintf("%v", tag))
                    }
                }
            }
            
            return config, nil
        }
    }
    
    return config, fmt.Errorf("no matching JSON schema found for file: %s", filePath)
}

func getValueByPath(data map[string]interface{}, path string) interface{} {
    parts := strings.Split(path, ".")
    var current interface{} = data
    
    for _, part := range parts {
        if m, ok := current.(map[string]interface{}); ok {
            current = m[part]
        } else {
            return nil
        }
    }
    
    return current
}
