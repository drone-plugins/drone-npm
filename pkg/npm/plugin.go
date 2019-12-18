// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package npm

import (
	"github.com/drone-plugins/drone-plugin-lib/pkg/plugin"
	"github.com/drone-plugins/drone-plugin-lib/pkg/urfave"
)

type pluginImpl struct {
	settings Settings
	pipeline plugin.Pipeline
	network  urfave.Network
}

// New Plugin from the given Settings, Pipeline, and Network.
func New(settings Settings, pipeline plugin.Pipeline, network urfave.Network) plugin.Plugin {
	return &pluginImpl{
		settings: settings,
		pipeline: pipeline,
		network:  network,
	}
}
