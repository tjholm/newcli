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
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/nitrictech/boxygen/pkg/backend/dockerfile"
)

type python struct {
	rte     RuntimeExt
	handler string
}

var _ Runtime = &python{}

func (t *python) DevImageName() string {
	return fmt.Sprintf("nitric-%s-dev", t.rte)
}

func (t *python) ContainerName() string {
	return strings.Replace(filepath.Base(t.handler), filepath.Ext(t.handler), "", 1)
}

func (t *python) BuildIgnore() []string {
	return append(commonIgnore, "__pycache__/", "*.py[cod]", "*$py.class")
}

func (t *python) FunctionDockerfileForCodeAsConfig(w io.Writer) error {
	css := dockerfile.NewStateStore()

	con, err := css.NewContainer(dockerfile.NewContainerOpts{
		From:   "python:3.9-slim",
		As:     layerFinal,
		Ignore: t.BuildIgnore(),
	})
	if err != nil {
		return err
	}

	con.Run(dockerfile.RunOptions{Command: []string{"pip", "install", "--upgrade", "pip"}})
	con.Config(dockerfile.ConfigOptions{
		WorkingDir: "/",
	})

	con.Run(dockerfile.RunOptions{Command: []string{"pip", "install", "jurigged"}})

	err = con.Copy(dockerfile.CopyOptions{Src: "requirements.txt", Dest: "requirements.txt"})
	if err != nil {
		return err
	}

	con.Run(dockerfile.RunOptions{Command: []string{"pip", "install", "--no-cache-dir", "-r", "requirements.txt"}})

	err = con.Copy(dockerfile.CopyOptions{Src: ".", Dest: "."})
	if err != nil {
		return err
	}

	con.Config(dockerfile.ConfigOptions{
		Env: map[string]string{
			"PYTHONPATH": "/app/:${PYTHONPATH}",
		},
		Ports: []int32{9001},
		Cmd:   []string{"python", t.handler},
	})

	_, err = w.Write([]byte(strings.Join(con.Lines(), "\n")))

	return err
}

func (t *python) LaunchOptsForFunctionCollect(runCtx string) (LaunchOpts, error) {
	return LaunchOpts{
		Image:      t.DevImageName(),
		Entrypoint: strslice.StrSlice{"python"},
		Cmd:        strslice.StrSlice{"/app/" + filepath.ToSlash(t.handler)},
		TargetWD:   "/app",
		Mounts: []mount.Mount{
			{
				Type:   "bind",
				Source: runCtx,
				Target: "/app",
			},
		},
	}, nil
}

func (t *python) LaunchOptsForFunction(runCtx string) (LaunchOpts, error) {
	return LaunchOpts{
		Image:      t.DevImageName(),
		Entrypoint: strslice.StrSlice{"python"},
		Cmd:        strslice.StrSlice{"-m", "jurigged", "-v", "/app/" + filepath.ToSlash(t.handler)},
		TargetWD:   "/app",
		Mounts: []mount.Mount{
			{
				Type:   "bind",
				Source: runCtx,
				Target: "/app",
			},
		},
	}, nil
}

func (t *python) FunctionDockerfile(funcCtxDir, version, provider string, w io.Writer) error {
	css := dockerfile.NewStateStore()

	con, err := css.NewContainer(dockerfile.NewContainerOpts{
		From:   "python:3.9-slim",
		As:     layerFinal,
		Ignore: t.BuildIgnore(),
	})
	if err != nil {
		return err
	}

	con.Run(dockerfile.RunOptions{Command: []string{"pip", "install", "--upgrade", "pip"}})
	con.Config(dockerfile.ConfigOptions{
		WorkingDir: "/",
	})

	err = con.Copy(dockerfile.CopyOptions{Src: "requirements.txt", Dest: "requirements.txt"})
	if err != nil {
		return err
	}

	con.Run(dockerfile.RunOptions{Command: []string{"pip", "install", "--no-cache-dir", "-r", "requirements.txt"}})

	err = con.Copy(dockerfile.CopyOptions{Src: ".", Dest: "."})
	if err != nil {
		return err
	}

	withMembrane(con, version, provider)

	con.Config(dockerfile.ConfigOptions{
		Env: map[string]string{
			"PYTHONPATH": "/app/:${PYTHONPATH}",
		},
		Ports: []int32{9001},
		Cmd:   []string{"python", t.handler},
	})

	_, err = w.Write([]byte(strings.Join(con.Lines(), "\n")))

	return err
}
