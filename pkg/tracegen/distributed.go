// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package tracegen

import (
	"context"
	"fmt"

	"go.elastic.co/apm/v2"
)

// SendDistributedTrace sends events generated by both APM Go Agent and OTEL library and
// link them with the same traceID so that they are linked and can be shown in the same trace view
func SendDistributedTrace(ctx context.Context, cfg Config) (apm.TraceID, error) {
	if err := cfg.validate(); err != nil {
		return apm.TraceID{}, err
	}
	txCtx, err := SendIntakeV2Trace(ctx, cfg)
	if err != nil {
		return txCtx.Trace, err
	}

	traceparent := formatTraceparentHeader(txCtx)
	tracestate := txCtx.State.String()
	ctx = SetOTLPTracePropagator(ctx, traceparent, tracestate)
	err = SendOTLPTrace(ctx, cfg)
	return txCtx.Trace, err
}

func formatTraceparentHeader(c apm.TraceContext) string {
	return fmt.Sprintf("%02x-%032x-%016x-%02x", 0, c.Trace[:], c.Span[:], c.Options)
}
