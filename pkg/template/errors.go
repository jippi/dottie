// Copyright 2020 The Compose Specification Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package template

import "fmt"

// InvalidTemplateError is returned when a variable template is not in a valid
// format
type InvalidTemplateError struct {
	Template string
	Wrapped  error
}

func (e InvalidTemplateError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("Invalid template: %#v (%s)", e.Template, e.Wrapped)
	}

	return fmt.Sprintf("Invalid template: %#v", e.Template)
}

// MissingRequiredError is returned when a variable template is missing
type MissingRequiredError struct {
	Variable string
	Reason   string
}

func (e MissingRequiredError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("required variable %s is missing a value: %s", e.Variable, e.Reason)
	}

	return fmt.Sprintf("required variable %s is missing a value", e.Variable)
}
