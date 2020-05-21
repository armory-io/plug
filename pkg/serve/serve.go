/*
Copyright © Armory, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package serve

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/armory-io/plug/pkg/plugin"
)

type Server struct {
	port      int
	pluginDir string
	loader    *plugin.Loader
	Server    *http.Server
}

func New(host string, port int, pluginDir string) (*Server, error) {
	l := &plugin.Loader{
		PluginDir:     pluginDir,
		BinaryAddress: fmt.Sprintf("%v:%v/pluginBinary.zip", host, port),
	}

	s := &Server{
		port:      port,
		pluginDir: pluginDir,
		loader:    l,
		Server: &http.Server{
			Addr: fmt.Sprintf(":%v", port),
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/plugins.json", s.metadataHandler)
	mux.HandleFunc("/pluginBinary.zip", s.binaryHandler)

	s.Server.Handler = mux

	return s, nil
}

func (s *Server) metadataHandler(w http.ResponseWriter, r *http.Request) {
	m, err := s.loader.LoadMetadata()
	if err != nil {
		// TODO: return HTTP error code like a real server...
		log.Fatalf("Could not load metadata: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode([]*plugin.Metadata{m}); err != nil {
		log.Fatalf("Could not encode metadata in response: %v", err)
	}
}

func (s *Server) binaryHandler(w http.ResponseWriter, r *http.Request) {
	p, err := s.loader.LocateBinary()
	if err != nil {
		// TODO: return HTTP error code like a real server...
		log.Fatalf("Could not load binary: %v", err)
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename='%s'", path.Base(p)))
	http.ServeFile(w, r, p)
}
