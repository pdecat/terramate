// Copyright 2022 Mineiros GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2etest

import (
	"fmt"
	"testing"

	"github.com/mineiros-io/terramate/stack/trigger"
	"github.com/mineiros-io/terramate/test"
	"github.com/mineiros-io/terramate/test/sandbox"
)

func TestListDetectAsChangedTriggeredStack(t *testing.T) {
	t.Parallel()

	s := sandbox.New(t)

	stack := s.CreateStack("stack")

	git := s.Git()
	git.CommitAll("all")
	git.Push("main")
	git.CheckoutNew("trigger-the-stack")

	cli := newCLI(t, s.RootDir())

	assertRunResult(t, cli.triggerStack("/stack"), runExpected{
		IgnoreStdout: true,
	})

	git.CommitAll("commit the trigger file")

	want := runExpected{
		Stdout: stack.RelPath() + "\n",
	}
	assertRunResult(t, cli.listChangedStacks(), want)
}

func TestRunChangedDetectionIgnoresDeletedTrigger(t *testing.T) {
	t.Parallel()

	const testfile = "testfile"

	s := sandbox.New(t)

	s.BuildTree([]string{
		"s:stack",
		fmt.Sprintf("f:stack/%s:stack\n", testfile),
	})

	cli := newCLI(t, s.RootDir())

	assertRunResult(t, cli.triggerStack("/stack"), runExpected{
		IgnoreStdout: true,
	})

	git := s.Git()
	git.CommitAll("all")
	git.Push("main")

	git.CheckoutNew("delete-stack-trigger")

	assertNoChanges := func() {
		t.Helper()

		assertRunResult(t, cli.run(
			"run",
			"--changed",
			testHelperBin,
			"cat",
			testfile,
		), runExpected{Stdout: ""})
	}

	assertNoChanges()

	triggerDir := trigger.Dir(s.RootDir())
	test.RemoveAll(t, triggerDir)

	git.CommitAll("removed trigger")

	assertNoChanges()
}

func TestRunChangedDetectsTriggeredStack(t *testing.T) {
	t.Parallel()

	const testfile = "testfile"

	s := sandbox.New(t)

	s.BuildTree([]string{
		"s:stack-1",
		"s:stack-2",
		fmt.Sprintf("f:stack-1/%s:stack-1\n", testfile),
		fmt.Sprintf("f:stack-2/%s:stack-2\n", testfile),
	})

	git := s.Git()
	git.CommitAll("all")
	git.Push("main")

	git.CheckoutNew("trigger-the-stack")

	cli := newCLI(t, s.RootDir())

	assertRunResult(t, cli.run(
		"run",
		"--changed",
		testHelperBin,
		"cat",
		testfile,
	), runExpected{Stdout: ""})

	assertRunResult(t, cli.triggerStack("/stack-1"), runExpected{
		IgnoreStdout: true,
	})
	git.CommitAll("commit the trigger file for stack-1")

	assertRunResult(t, cli.run(
		"run",
		"--changed",
		testHelperBin,
		"cat",
		testfile,
	), runExpected{Stdout: listStacks("stack-1")})

	assertRunResult(t, cli.triggerStack("/stack-2"), runExpected{
		IgnoreStdout: true,
	})
	git.CommitAll("commit the trigger file for stack-2")

	assertRunResult(t, cli.run(
		"run",
		"--changed",
		testHelperBin,
		"cat",
		testfile,
	), runExpected{Stdout: listStacks("stack-1", "stack-2")})
}

func TestRunWontDetectAsChangeDeletedTrigger(t *testing.T) {
	t.Parallel()

	const testfile = "testfile"

	s := sandbox.New(t)

	s.BuildTree([]string{
		"s:stack-1",
		"s:stack-2",
		fmt.Sprintf("f:stack-1/%s:stack-1\n", testfile),
		fmt.Sprintf("f:stack-2/%s:stack-2\n", testfile),
	})

	cli := newCLI(t, s.RootDir())
	assertRunResult(t, cli.triggerStack("/stacks/stack-1"), runExpected{
		IgnoreStdout: true,
	})

	git := s.Git()
	git.CommitAll("all")
	git.Push("main")

	git.CheckoutNew("delete-trigger")

	test.RemoveAll(t, trigger.Dir(s.RootDir()))
	git.CommitAll("removed trigger")

	assertRunResult(t, cli.run(
		"run",
		"--changed",
		testHelperBin,
		"cat",
		testfile,
	), runExpected{Stdout: ""})
}