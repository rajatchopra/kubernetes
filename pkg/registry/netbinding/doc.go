/*
Copyright 2014 Google Inc. All rights reserved.

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

// Package netbinding contains the middle layer logic for netbindings.
// NetBindings are objects containing instructions for how a pod ought to
// be networked within the cluster. This allows a registry object which supports this
// action (ApplyNetBinding) to be served through an apiserver.
package netbinding
