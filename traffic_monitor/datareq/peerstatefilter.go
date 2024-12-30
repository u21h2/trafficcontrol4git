/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package datareq

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// PeerStateFilter fulfills the cache.Filter interface, for filtering stats. See the `NewPeerStateFilter` documentation for details on which query parameters are used to filter.
type PeerStateFilter struct {
	historyCount int
	cachesToUse  map[tc.CacheName]struct{}
	peersToUse   map[tc.TrafficMonitorName]struct{}
	wildcard     bool
	cacheType    tc.CacheType
	cacheTypes   map[tc.CacheName]tc.CacheType
}

// UsePeer returns whether the given Traffic Monitor peer is in this filter.
func (f *PeerStateFilter) UsePeer(name tc.TrafficMonitorName) bool {
	if _, inPeers := f.peersToUse[name]; len(f.peersToUse) != 0 && !inPeers {
		return false
	}
	return true
}

// UseCache returns whether the given cache is in this filter.
func (f *PeerStateFilter) UseCache(name tc.CacheName) bool {
	if f.cacheType != tc.CacheTypeInvalid && f.cacheTypes[name] != f.cacheType {
		return false
	}

	if len(f.cachesToUse) == 0 {
		return true
	}

	if !f.wildcard {
		_, ok := f.cachesToUse[name]
		return ok
	}
	for cacheToUse := range f.cachesToUse {
		if strings.Contains(string(name), string(cacheToUse)) {
			return true
		}
	}
	return false
}

// WithinStatHistoryMax returns whether the given history index is less than the max history of this filter.
func (f *PeerStateFilter) WithinStatHistoryMax(n int) bool {
	if f.historyCount == 0 {
		return true
	}
	if n <= f.historyCount {
		return true
	}
	return false
}

// NewPeerStateFilter takes the HTTP query parameters and creates a cache.Filter, filtering according to the query parameters passed.
// Query parameters used are `hc`, `stats`, `wildcard`, `typep`, and `hosts`. The `stats` param filters caches. The `hosts` param filters peer Traffic Monitors. The `type` param filters cache types (edge, mid).
// If `hc` is 0, all history is returned. If `hc` is empty, 1 history is returned.
// If `stats` is empty, all stats are returned.
// If `wildcard` is empty, `stats` is considered exact.
// If `type` is empty, all cache types are returned.
func NewPeerStateFilter(path string, params url.Values, cacheTypes map[tc.CacheName]tc.CacheType) (*PeerStateFilter, error) {
	// TODO change legacy `stats` and `hosts` to `caches` and `monitors` (or `peers`).
	validParams := map[string]struct{}{"hc": struct{}{}, "stats": struct{}{}, "wildcard": struct{}{}, "type": struct{}{}, "peers": struct{}{}}
	if len(params) > len(validParams) {
		return nil, fmt.Errorf("invalid query parameters")
	}
	for param := range params {
		if _, ok := validParams[param]; !ok {
			return nil, fmt.Errorf("invalid query parameter '%v'", param)
		}
	}

	historyCount := 1
	if paramHc, exists := params["hc"]; exists && len(paramHc) > 0 {
		v, err := strconv.Atoi(paramHc[0])
		if err == nil {
			historyCount = v
		}
	}

	cachesToUse := map[tc.CacheName]struct{}{}
	// TODO rename 'stats' to 'caches'
	if paramStats, exists := params["stats"]; exists && len(paramStats) > 0 {
		commaStats := strings.Split(paramStats[0], ",")
		for _, stat := range commaStats {
			cachesToUse[tc.CacheName(stat)] = struct{}{}
		}
	}

	wildcard := false
	if paramWildcard, exists := params["wildcard"]; exists && len(paramWildcard) > 0 {
		wildcard, _ = strconv.ParseBool(paramWildcard[0]) // ignore errors, error => false
	}

	cacheType := tc.CacheTypeInvalid
	if paramType, exists := params["type"]; exists && len(paramType) > 0 {
		cacheType = tc.CacheTypeFromString(paramType[0])
		if cacheType == tc.CacheTypeInvalid {
			return nil, fmt.Errorf("invalid query parameter type '%v' - valid types are: {edge, mid}", paramType[0])
		}
	}

	peersToUse := map[tc.TrafficMonitorName]struct{}{}
	if paramNames, exists := params["peers"]; exists && len(paramNames) > 0 {
		commaNames := strings.Split(paramNames[0], ",")
		for _, name := range commaNames {
			peersToUse[tc.TrafficMonitorName(name)] = struct{}{}
		}
	}

	pathArgument := getPathArgument(path)
	if pathArgument != "" {
		peersToUse[tc.TrafficMonitorName(pathArgument)] = struct{}{}
	}

	// parameters without values are considered names, e.g. `?my-cache-0` or `?my-delivery-service`
	for maybeName, val := range params {
		if len(val) == 0 || (len(val) == 1 && val[0] == "") {
			peersToUse[tc.TrafficMonitorName(maybeName)] = struct{}{}
		}
	}

	return &PeerStateFilter{
		historyCount: historyCount,
		cachesToUse:  cachesToUse,
		wildcard:     wildcard,
		cacheType:    cacheType,
		peersToUse:   peersToUse,
		cacheTypes:   cacheTypes,
	}, nil
}
