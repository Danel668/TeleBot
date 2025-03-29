package caches

import (
	"log"
	"fmt"
	"sync"
	"time"
	"io/ioutil"
	"os"
	"encoding/json"

	"TestBot/src/utils"

	"go.uber.org/zap"
)

type State int

const (
	Start State = iota
	CountOfCalories
    RegistrationTimeZone
    MyProgress
    SetGoal
    AddRation
    GetMyRation
    AddReminder
	AddRecommendation
	Admin
)

type StateInfo struct {
	State     State
	CreatedAt time.Time
}

type StateUserCache struct {
	userStates sync.Map
}

func NewStateUserCache() *StateUserCache {
	cache := &StateUserCache{}
	cache.loadFromDump()
	return cache
}

func (c *StateUserCache) SetState(userId int64, state State) {
	c.userStates.Store(userId, StateInfo{
		State: state,
		CreatedAt: time.Now(),
	})
}

func (c *StateUserCache) GetState(userId int64) State {
	value, ok := c.userStates.Load(userId)
	if !ok {
		return Start
	}
	return value.(StateInfo).State
}

func (c *StateUserCache) ClearExpiredState(duration time.Duration) {
	now := time.Now()
	c.userStates.Range(func(key, value interface{}) bool {
		stateInfo := value.(StateInfo)
		if now.Sub(stateInfo.CreatedAt) > duration {
			c.userStates.Delete(key)
		}
		return true
	})
}

func (c *StateUserCache) loadFromDump() {
	fileData, err := ioutil.ReadFile("caches/user_cache_dump.json")

	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		log.Fatalln(err)
		return
	}

	var data map[string]StateInfo
	if err := json.Unmarshal(fileData, &data); err != nil {
		log.Fatalln("failed to unmarshal caches/user_cache_dump.json:", err)
		return
	}

	for key, value := range data {
		userId := utils.StringToInt64(key)
		c.userStates.Store(userId, value)
	}
}

func (c *StateUserCache) DumpToFile(logger *zap.Logger) {
	data := make(map[string]StateInfo)

	c.userStates.Range(func(key, value interface{}) bool {
		userId := key.(int64)
		stateInfo := value.(StateInfo)
		data[utils.ToString(userId)] = stateInfo
		return true
	})

	fileData, err := json.Marshal(data)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to marshal data for user_cache_dump.json: %v", err))
	}

	if err := ioutil.WriteFile("caches/user_cache_dump.json", fileData, 0644); err != nil {
		logger.Error(fmt.Sprintf("failed to write data to user_cache_dump.json: %v", err))
	}
}
