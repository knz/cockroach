// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package sqltelemetry

import "github.com/cockroachdb/cockroach/pkg/server/telemetry"

// DistSQLExecCounter is to be incremented whenever a query is distributed
// across multiple nodes.
var DistSQLExecCounter = telemetry.GetCounterOnce("sql.exec.query.is-distributed")

// VecExecCounter is to be incremented whenever a query runs with the vectorized
// execution engine.
var VecExecCounter = telemetry.GetCounterOnce("sql.exec.query.is-vectorized")

// CascadesLimitReached is to be incremented whenever the limit of foreign key
// cascade for a single query is exceeded.
var CascadesLimitReached = telemetry.GetCounterOnce("sql.exec.cascade-limit-reached")

// AutoStatsClusterSettingName is the name of the automatic stats collection
// cluster setting.
const AutoStatsClusterSettingName = "sql.stats.automatic_collection.enabled"

// VectorizeClusterSettingName is the name for the cluster setting that controls
// the VectorizeClusterMode.
const VectorizeClusterSettingName = "sql.defaults.vectorize"

// ReorderJoinsLimitClusterSettingName is the name of the cluster setting for
// the maximum number of joins to reorder.
const ReorderJoinsLimitClusterSettingName = "sql.defaults.reorder_joins_limit"
