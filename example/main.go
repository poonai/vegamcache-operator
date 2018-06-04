/*
Copyright 2018 The vegamcache Authors.
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

package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sch00lb0y/vegamcache"
)

func main() {
	vg, _ := vegamcache.NewVegam(&vegamcache.VegamConfig{
		Port:   8087,
		Logger: log.New(ioutil.Discard, "", 0),
	})
	vg.Start()
	go vegamcache.ListenServer(vg, ":8000")
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		val, exist := vg.Get(key)
		if !exist {
			w.Write([]byte(`illeh`))
			return
		}
		w.Write([]byte(val.(string)))
	})
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		// TODO: convert this into POST.
		// lazy to write marshal
		key := r.URL.Query().Get("key")
		val := r.URL.Query().Get("val")
		vg.Put(key, val, 0)
		w.Write([]byte(`cached`))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`I'm home`))
	})
	log.Print("Server Stared")
	err := http.ListenAndServe(":90", nil)
	if err != nil {
		panic(err)
	}
}
