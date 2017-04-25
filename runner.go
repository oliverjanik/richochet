package ricochet

import (
	"fmt"
	"sync"
)

var (
	suits []*Suite
)

// Run all the Suites
func Run(suites ...*Suite) {
	var wg sync.WaitGroup
	wg.Add(len(suites))

	for _, suite := range suites {
		runSuite(suite, &wg)
	}

	wg.Wait()
}

func runSuite(s *Suite, wg *sync.WaitGroup) {
	defer wg.Done()

	if s.setUp != nil {
		s.setUp(s)
	}

	var groupWg sync.WaitGroup
	groupWg.Add(len(s.groups) + 1)

	// run self
	runGroup(&s.TestGroup, s, &groupWg)

	for _, g := range s.groups {
		go runGroup(g, s, &groupWg)
	}

	groupWg.Wait()

	if s.tearDown != nil {
		s.tearDown(s)
	}
}

func runGroup(g *TestGroup, s *Suite, wg *sync.WaitGroup) {
	fmt.Println(g.indent+"Running", g.name)

	defer wg.Done()

	defer func() {
		if msg := recover(); msg != nil {
			fmt.Printf(g.indent+"\t\tError: %v", msg)
			g.failed = true
		}
	}()

	for _, t := range g.tests {
		fmt.Println(g.indent+"\t", "...", t.name)
		t.f(&R{
			baseURL: s.baseURL,
			Header:  s.header,
		})
	}
}
