// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtime

import (
	"fmt"
	"io"
	"path/filepath"
	osruntime "runtime"
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"

	"github.com/nitrictech/cli/pkg/utils"
)

type golang struct {
	rte     RuntimeExt
	handler string
}

var _ Runtime = &golang{}

func (t *golang) DevImageName() string {
	return fmt.Sprintf("nitric-%s-dev", t.rte)
}

func (t *golang) BuildIgnore() []string {
	return []string{}
}

func (t *golang) ContainerName() string {
	// get the abs dir in case user provides "."
	absH, err := filepath.Abs(t.handler)
	if err != nil {
		return ""
	}

	return filepath.Base(filepath.Dir(absH))
}

// dockerfile for running requirements collection and development runs
const devDockerfile = `# syntax = docker/dockerfile:1.3
FROM golang:alpine

RUN apk add --no-cache git
RUN go install github.com/asalkeld/CompileDaemon@d4b10de
`

// final production image for running in the cloud
const prodDockerfile = `# syntax = docker/dockerfile:1.3
FROM golang:alpine as build
RUN apk update
RUN apk upgrade
RUN apk add --no-cache git gcc g++ make

WORKDIR /app/

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build go build -o /bin/main ./%s/...

FROM alpine

RUN apk add --no-cache tzdata
ADD %s /bin/membrane
RUN chmod +x /bin/membrane

COPY --from=build /bin/main /bin/main

RUN chmod +x-rw /bin/main

EXPOSE 9001

ENTRYPOINT /bin/membrane
CMD /bin/main
`

func (t *golang) FunctionDockerfile(funcCtxDir, version, provider string, w io.Writer) error {
	dockerfile := fmt.Sprintf(prodDockerfile, filepath.Dir(t.handler), membraneUrl(version, provider))

	_, err := w.Write([]byte(dockerfile))
	return err
}

func (t *golang) FunctionDockerfileForCodeAsConfig(w io.Writer) error {
	_, err := w.Write([]byte(devDockerfile))
	return err
}

func (t *golang) LaunchOptsForFunctionCollect(runCtx string) (LaunchOpts, error) {
	module, err := utils.GoModule(runCtx)
	if err != nil {
		return LaunchOpts{}, err
	}

	goPath, err := utils.GoPath()
	if err != nil {
		return LaunchOpts{}, err
	}

	return LaunchOpts{
		Image:    t.DevImageName(),
		TargetWD: filepath.ToSlash(filepath.Join("/go/src", module)),
		Cmd:      strslice.StrSlice{"go", "run", "./" + filepath.ToSlash(filepath.Dir(t.handler)) + "/..."},
		Mounts: []mount.Mount{
			{
				Type:   "bind",
				Source: filepath.Join(goPath, "pkg"),
				Target: "/go/pkg",
			},
			{
				Type:   "bind",
				Source: runCtx,
				Target: filepath.ToSlash(filepath.Join("/go/src", module)),
			},
		},
	}, nil
}

func (t *golang) LaunchOptsForFunction(runCtx string) (LaunchOpts, error) {
	module, err := utils.GoModule(runCtx)
	if err != nil {
		return LaunchOpts{}, err
	}
	containerRunCtx := filepath.ToSlash(filepath.Join("/go/src", module))
	relHandler := t.handler
	if strings.HasPrefix(t.handler, runCtx) {
		relHandler, err = filepath.Rel(runCtx, t.handler)
		if err != nil {
			return LaunchOpts{}, err
		}
	}

	goPath, err := utils.GoPath()
	if err != nil {
		return LaunchOpts{}, err
	}

	opts := LaunchOpts{
		TargetWD: containerRunCtx,
		Cmd: strslice.StrSlice{
			"/go/bin/CompileDaemon",
			"-verbose",
			"-exclude-dir=.git",
			"-exclude-dir=.nitric",
			"-directory=.",
			fmt.Sprintf("-polling=%t", osruntime.GOOS == "windows"),
			fmt.Sprintf("-build=go build -buildvcs=false -o %s ./%s/...", t.ContainerName(), filepath.ToSlash(filepath.Dir(relHandler))),
			"-command=./" + t.ContainerName(),
		},
		Mounts: []mount.Mount{
			{
				Type:   "bind",
				Source: filepath.Join(goPath, "pkg"),
				Target: "/go/pkg",
			},
			{
				Type:   "bind",
				Source: runCtx,
				Target: containerRunCtx,
			},
		},
	}

	return opts, nil
}
