// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"testing"
)

func TestSystemInfo(t *testing.T) {
	s, err := getSchema(&Resolver{
		si: SystemInfo{
			HostName:        "hostname",
			CPUArchitecture: "cpuArch",
			OSPlatform:      "osPlatform",
			Memory:          1234,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
	}{
		{"OK"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `
			query OpName {
			 systemInfo {
			    hostname
			    cpuArchitecture
			    osPlatform
			    memory
			    capabilities {
			  	  key
			  	  value
			    }
			  }
			}`

			res := s.Exec(context.Background(), q, "", nil)

			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}
