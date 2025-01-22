// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package json // import "go.opentelemetry.io/collector/pdata/internal/json"

import (
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var opts = &protojson.MarshalOptions{}

func Marshal(out io.Writer, pb proto.Message) error {
	jsonBytes, err := opts.Marshal(pb)
	if err != nil {
		return err
	}
	if _, err := out.Write(jsonBytes); err != nil {
		return err
	}
	return nil
}
