package config

import (
	"fmt"
	"strings"

	"github.com/Jrp0h/backuper/action"
	"github.com/Jrp0h/backuper/utils"
)

type actionsConfig map[string]action.Action


func loadActions(config *configFile) (actionsConfig, error) {
	actions := make(map[string]action.Action) 

	for k, v := range config.Actions {
		k = strings.ToLower(k)
		action, err := action.LoadFromConfig(&v)
		if err != nil {
			return nil, err
		}

		_, exists := actions[k]
		if exists {
			return nil, fmt.Errorf("config/actions: action %s already exists", k)
		}

		actions[k] = action
	}

	return actions, nil
}

func (actions *actionsConfig) Only(names []string) (actionsConfig, error) {
	only := make(map[string]action.Action) 

	for _, name := range names {
		name = strings.ToLower(name)

		value, exists := (*actions)[name]
		if !exists {
			alternatives := utils.Levenshtein(name, actions.Names(), true).AsQuestion()
			return nil, fmt.Errorf("config/actions: Name %s does not exist in actions\nDid you mean %s?", name, alternatives)
		}

		_, exists = only[name]
		if exists {
			return nil, fmt.Errorf("config/actions: Name %s has already been seen", name)
		}

		only[name] = value
	}

	return only, nil
}

func (actions *actionsConfig) Except(names []string) (actionsConfig, error) {
	except := *actions

	for _, name := range names {
		name = strings.ToLower(name)

		_, exists := (*actions)[name]
		if !exists {
			alternatives := utils.Levenshtein(name, actions.Names(), true).AsQuestion()
			return nil, fmt.Errorf("config/actions: Name %s does not exist in actions\nDid you mean %s?", name, alternatives)
		}

		// Already removed
		_, exists = except[name]
		if !exists {
			return nil, fmt.Errorf("config/actions: Name %s has already been seen", name)
		}

		delete(except, name)
	}

	return except, nil
}

func (actions *actionsConfig) Names() []string {
	names := make([]string, len(*actions))

	i := 0
	for k := range *actions {
		names[i] = k
		i += 1
	}

	return names
}

func (actions *actionsConfig) OnlyOrExcept(only, except []string) actionsConfig {
	if len(only) > 0 && len(except) > 0 {
		utils.Log.Fatal("config/actions: only and except can't be used together")
	}

	if len(only) > 0 {
		o, err := actions.Only(only)
		if err != nil {
			utils.Log.Fatal(err.Error())
		}
		return o
	}
	if len(except) > 0 {
		o, err := actions.Except(except)
		if err != nil {
			utils.Log.Fatal(err.Error())
		}

		return o
	}

	return *actions
}