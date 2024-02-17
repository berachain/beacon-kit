// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types

// DefaultLocalBuilderName holds the default local builder name.
const DefaultLocalBuilderName = "default-local-builder"

// BuilderEntry holds the BuilderServiceClient and a boolean indicating if it's local.
type BuilderEntry struct {
	Name string
	BuilderServiceClient
	IsLocal bool
}

// BuilderRegistry holds a map of builder names to their corresponding BuilderEntry.
type BuilderRegistry struct {
	localBuilders  map[string]*BuilderEntry
	remoteBuilders map[string]*BuilderEntry
}

// NewBuilderRegistry creates a new BuilderRegistry with an empty map of builders.
func NewBuilderRegistry() *BuilderRegistry {
	return &BuilderRegistry{
		localBuilders: make(map[string]*BuilderEntry),
	}
}

// RegisterBuilder adds a new BuilderEntry to the BuilderRegistry's map of builders.
func (br *BuilderRegistry) RegisterBuilder(
	name string, client BuilderServiceClient, isLocal bool,
) {
	var m map[string]*BuilderEntry
	if isLocal {
		m = br.localBuilders
	} else {
		m = br.remoteBuilders
	}
	m[name] = &BuilderEntry{
		Name:                 name,
		BuilderServiceClient: client,
	}
}

// GetBuilder returns the corresponding BuilderEntry for a given name.
func (br *BuilderRegistry) GetBuilder(name string) *BuilderEntry {
	if builder, ok := br.localBuilders[name]; ok {
		return builder
	}
	return br.remoteBuilders[name]
}

// LocalBuildersList returns a list of local builders.
func (br *BuilderRegistry) LocalBuilders() []*BuilderEntry {
	localBuildersList := make([]*BuilderEntry, 0, len(br.localBuilders))
	for _, builder := range br.localBuilders {
		localBuildersList = append(localBuildersList, builder)
	}
	return localBuildersList
}

// RemoteBuildersList returns a list of remote builders.
func (br *BuilderRegistry) RemoteBuilders() []*BuilderEntry {
	remoteBuildersList := make([]*BuilderEntry, 0, len(br.remoteBuilders))
	for _, builder := range br.remoteBuilders {
		remoteBuildersList = append(remoteBuildersList, builder)
	}
	return remoteBuildersList
}
