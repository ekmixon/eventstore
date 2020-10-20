/*
Copyright (c) 2020 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sharedmain

import (
	"go.uber.org/zap"

	"knative.dev/pkg/logging"
)

// EnvConfigConstructor returns an Env Config Accesor object
type EnvConfigConstructor func() EnvConfigAccessor

// EnvConfig is the minimal set of configuration parameters
// events storage server should support.
type EnvConfig struct {
	// Component is the kind of this events storage server.
	Component string `envconfig:"K_COMPONENT"`

	// Environment variable containing the namespace of the events storage server.
	Namespace string `envconfig:"NAMESPACE" required:"true"`

	// Environment variable containing the name of the events storage server.
	Name string `envconfig:"NAME" default:"server"`

	// LoggingConfigJson is a json string of logging.Config.
	// This is used to configure the logging config, the config is stored in
	// a config map inside the controllers namespace and copied here.
	LoggingConfigJson string `envconfig:"K_LOGGING_CONFIG" required:"true"`
}

// EnvConfigAccessor defines accessors for the minimal
// set of events storage server configuration parameters.
type EnvConfigAccessor interface {
	// Set the component name.
	SetComponent(string)

	// Get the namespace of the server.
	GetNamespace() string

	// Get the name of the server.
	GetName() string

	// Get the parsed logger.
	GetLogger() *zap.SugaredLogger
}

var _ EnvConfigAccessor = (*EnvConfig)(nil)

// SetComponent for target Kind
func (e *EnvConfig) SetComponent(component string) {
	e.Component = component
}

// GetNamespace for target server
func (e *EnvConfig) GetNamespace() string {
	return e.Namespace
}

// GetName for target server
func (e *EnvConfig) GetName() string {
	return e.Name
}

// GetLogger retrieves the configured logger
func (e *EnvConfig) GetLogger() *zap.SugaredLogger {
	loggingConfig, err := logging.JsonToLoggingConfig(e.LoggingConfigJson)
	if err != nil {
		// Use default logging config.
		if loggingConfig, err = logging.NewConfigFromMap(map[string]string{}); err != nil {
			// If this fails, there is no recovering.
			panic(err)
		}
	}

	logger, _ := logging.NewLoggerFromConfig(loggingConfig, e.Component)

	return logger
}
