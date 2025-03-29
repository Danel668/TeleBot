package sources

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"errors"
	"strconv"
	"fmt"

	"go.uber.org/zap"
)

func NewConfig(logger *zap.Logger) *Config {
	return LoadConfig(logger)
}

func LoadConfig(logger *zap.Logger) *Config {
	data, err := ioutil.ReadFile("config/production.json")

	if err != nil {
		logger.Warn(fmt.Sprintf("failed to read file config/productions.json: %v", err))
		return nil
	}

	var config Config
	err = json.Unmarshal(data, &config)

	if err != nil {
		logger.Warn(fmt.Sprintf("failed to unmarshal config/productions.json: %v", err))
		return nil
	}

	return &config
}

func GetConfigToString() (string, error) {
	data, err := ioutil.ReadFile("config/production.json")

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func updateValueInConfigByKey(s interface{}, key string, newValue string, logger *zap.Logger) bool {
	val := reflect.ValueOf(s)

	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		logger.Warn("getting value isn`t ptr or *ptr isn`t struct")
		return false
	}

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if fieldType.Name == key && field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(newValue)
				return true
			case reflect.Bool:
				newValueParsed, err := strconv.ParseBool(newValue)
				if err != nil {
					logger.Warn("error parsing string to bool")
					return false
				}
				field.SetBool(newValueParsed)
				return true
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				newValueParsed, err := strconv.ParseInt(newValue, 10, 64)
				if err != nil {
					logger.Warn("error parsing string to int64")
					return false
				}
				field.SetInt(newValueParsed)
				return true
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				newValueParsed, err := strconv.ParseUint(newValue, 10, 64)
				if err != nil {
					logger.Warn("error parsing string to uint64")
					return false
				}
				field.SetUint(newValueParsed)
				return true
			case reflect.Float32, reflect.Float64:
				newValueParsed, err := strconv.ParseFloat(newValue, 64)
				if err != nil {
					logger.Warn("error parsing string to float64")
					return false
				}
				field.SetFloat(newValueParsed)
				return true
			}
		}

		if field.Kind() == reflect.Struct && field.CanAddr() {
			if updateValueInConfigByKey(field.Addr().Interface(), key, newValue, logger) {
				return true
			}
		}
	}
	return false
}

func ChangeConfigValueByKey(key string, value string, logger *zap.Logger) error {
	data, err := ioutil.ReadFile("config/production.json")

	if err != nil {
		return err
	}

	var config Config
	err = json.Unmarshal(data, &config)

	if err != nil {
		return err
	}

	isUpdate := updateValueInConfigByKey(&config, key, value, logger)

	if isUpdate {
		updatedConfig, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			return errors.New("config wasn`t updated: failed to marshal json")
		}

		if err := ioutil.WriteFile("config/production.json", updatedConfig, 0644); err != nil {
			return errors.New("config wasn`t updated: failed to write json")
		}

		return nil
	}

	return errors.New("config wasn`t updated: key wasn`t found")
}
