// Copyright 2017 Bo-Yi Wu. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build !(sonic && avx && (linux || windows || darwin) && amd64)

package json

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
)

type RawMessage = json.RawMessage

var (
	jsonlib = jsoniter.ConfigCompatibleWithStandardLibrary
	// Marshal is exported by gin/json package.
	Marshal = jsonlib.Marshal
	// Unmarshal is exported by gin/json package.
	Unmarshal = jsonlib.Unmarshal
	// MarshalIndent is exported by gin/json package.
	MarshalIndent = jsonlib.MarshalIndent
	// NewDecoder is exported by gin/json package.
	NewDecoder = jsonlib.NewDecoder
	// NewEncoder is exported by gin/json package.
	NewEncoder = jsonlib.NewEncoder
)
