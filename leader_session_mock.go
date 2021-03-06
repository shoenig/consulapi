package consulapi

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// LeaderSessionMock implements LeaderSession
type LeaderSessionMock struct {
	t minimock.Tester

	funcAbdicate          func(c1 Ctx) (err error)
	inspectFuncAbdicate   func(c1 Ctx)
	afterAbdicateCounter  uint64
	beforeAbdicateCounter uint64
	AbdicateMock          mLeaderSessionMockAbdicate

	funcCurrent          func(c1 Ctx) (s1 string, err error)
	inspectFuncCurrent   func(c1 Ctx)
	afterCurrentCounter  uint64
	beforeCurrentCounter uint64
	CurrentMock          mLeaderSessionMockCurrent

	funcSessionID          func(c1 Ctx) (s1 string)
	inspectFuncSessionID   func(c1 Ctx)
	afterSessionIDCounter  uint64
	beforeSessionIDCounter uint64
	SessionIDMock          mLeaderSessionMockSessionID
}

// NewLeaderSessionMock returns a mock for LeaderSession
func NewLeaderSessionMock(t minimock.Tester) *LeaderSessionMock {
	m := &LeaderSessionMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AbdicateMock = mLeaderSessionMockAbdicate{mock: m}
	m.AbdicateMock.callArgs = []*LeaderSessionMockAbdicateParams{}

	m.CurrentMock = mLeaderSessionMockCurrent{mock: m}
	m.CurrentMock.callArgs = []*LeaderSessionMockCurrentParams{}

	m.SessionIDMock = mLeaderSessionMockSessionID{mock: m}
	m.SessionIDMock.callArgs = []*LeaderSessionMockSessionIDParams{}

	return m
}

type mLeaderSessionMockAbdicate struct {
	mock               *LeaderSessionMock
	defaultExpectation *LeaderSessionMockAbdicateExpectation
	expectations       []*LeaderSessionMockAbdicateExpectation

	callArgs []*LeaderSessionMockAbdicateParams
	mutex    sync.RWMutex
}

// LeaderSessionMockAbdicateExpectation specifies expectation struct of the LeaderSession.Abdicate
type LeaderSessionMockAbdicateExpectation struct {
	mock    *LeaderSessionMock
	params  *LeaderSessionMockAbdicateParams
	results *LeaderSessionMockAbdicateResults
	Counter uint64
}

// LeaderSessionMockAbdicateParams contains parameters of the LeaderSession.Abdicate
type LeaderSessionMockAbdicateParams struct {
	c1 Ctx
}

// LeaderSessionMockAbdicateResults contains results of the LeaderSession.Abdicate
type LeaderSessionMockAbdicateResults struct {
	err error
}

// Expect sets up expected params for LeaderSession.Abdicate
func (mmAbdicate *mLeaderSessionMockAbdicate) Expect(c1 Ctx) *mLeaderSessionMockAbdicate {
	if mmAbdicate.mock.funcAbdicate != nil {
		mmAbdicate.mock.t.Fatalf("LeaderSessionMock.Abdicate mock is already set by Set")
	}

	if mmAbdicate.defaultExpectation == nil {
		mmAbdicate.defaultExpectation = &LeaderSessionMockAbdicateExpectation{}
	}

	mmAbdicate.defaultExpectation.params = &LeaderSessionMockAbdicateParams{c1}
	for _, e := range mmAbdicate.expectations {
		if minimock.Equal(e.params, mmAbdicate.defaultExpectation.params) {
			mmAbdicate.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmAbdicate.defaultExpectation.params)
		}
	}

	return mmAbdicate
}

// Inspect accepts an inspector function that has same arguments as the LeaderSession.Abdicate
func (mmAbdicate *mLeaderSessionMockAbdicate) Inspect(f func(c1 Ctx)) *mLeaderSessionMockAbdicate {
	if mmAbdicate.mock.inspectFuncAbdicate != nil {
		mmAbdicate.mock.t.Fatalf("Inspect function is already set for LeaderSessionMock.Abdicate")
	}

	mmAbdicate.mock.inspectFuncAbdicate = f

	return mmAbdicate
}

// Return sets up results that will be returned by LeaderSession.Abdicate
func (mmAbdicate *mLeaderSessionMockAbdicate) Return(err error) *LeaderSessionMock {
	if mmAbdicate.mock.funcAbdicate != nil {
		mmAbdicate.mock.t.Fatalf("LeaderSessionMock.Abdicate mock is already set by Set")
	}

	if mmAbdicate.defaultExpectation == nil {
		mmAbdicate.defaultExpectation = &LeaderSessionMockAbdicateExpectation{mock: mmAbdicate.mock}
	}
	mmAbdicate.defaultExpectation.results = &LeaderSessionMockAbdicateResults{err}
	return mmAbdicate.mock
}

//Set uses given function f to mock the LeaderSession.Abdicate method
func (mmAbdicate *mLeaderSessionMockAbdicate) Set(f func(c1 Ctx) (err error)) *LeaderSessionMock {
	if mmAbdicate.defaultExpectation != nil {
		mmAbdicate.mock.t.Fatalf("Default expectation is already set for the LeaderSession.Abdicate method")
	}

	if len(mmAbdicate.expectations) > 0 {
		mmAbdicate.mock.t.Fatalf("Some expectations are already set for the LeaderSession.Abdicate method")
	}

	mmAbdicate.mock.funcAbdicate = f
	return mmAbdicate.mock
}

// When sets expectation for the LeaderSession.Abdicate which will trigger the result defined by the following
// Then helper
func (mmAbdicate *mLeaderSessionMockAbdicate) When(c1 Ctx) *LeaderSessionMockAbdicateExpectation {
	if mmAbdicate.mock.funcAbdicate != nil {
		mmAbdicate.mock.t.Fatalf("LeaderSessionMock.Abdicate mock is already set by Set")
	}

	expectation := &LeaderSessionMockAbdicateExpectation{
		mock:   mmAbdicate.mock,
		params: &LeaderSessionMockAbdicateParams{c1},
	}
	mmAbdicate.expectations = append(mmAbdicate.expectations, expectation)
	return expectation
}

// Then sets up LeaderSession.Abdicate return parameters for the expectation previously defined by the When method
func (e *LeaderSessionMockAbdicateExpectation) Then(err error) *LeaderSessionMock {
	e.results = &LeaderSessionMockAbdicateResults{err}
	return e.mock
}

// Abdicate implements LeaderSession
func (mmAbdicate *LeaderSessionMock) Abdicate(c1 Ctx) (err error) {
	mm_atomic.AddUint64(&mmAbdicate.beforeAbdicateCounter, 1)
	defer mm_atomic.AddUint64(&mmAbdicate.afterAbdicateCounter, 1)

	if mmAbdicate.inspectFuncAbdicate != nil {
		mmAbdicate.inspectFuncAbdicate(c1)
	}

	mm_params := &LeaderSessionMockAbdicateParams{c1}

	// Record call args
	mmAbdicate.AbdicateMock.mutex.Lock()
	mmAbdicate.AbdicateMock.callArgs = append(mmAbdicate.AbdicateMock.callArgs, mm_params)
	mmAbdicate.AbdicateMock.mutex.Unlock()

	for _, e := range mmAbdicate.AbdicateMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmAbdicate.AbdicateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmAbdicate.AbdicateMock.defaultExpectation.Counter, 1)
		mm_want := mmAbdicate.AbdicateMock.defaultExpectation.params
		mm_got := LeaderSessionMockAbdicateParams{c1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmAbdicate.t.Errorf("LeaderSessionMock.Abdicate got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmAbdicate.AbdicateMock.defaultExpectation.results
		if mm_results == nil {
			mmAbdicate.t.Fatal("No results are set for the LeaderSessionMock.Abdicate")
		}
		return (*mm_results).err
	}
	if mmAbdicate.funcAbdicate != nil {
		return mmAbdicate.funcAbdicate(c1)
	}
	mmAbdicate.t.Fatalf("Unexpected call to LeaderSessionMock.Abdicate. %v", c1)
	return
}

// AbdicateAfterCounter returns a count of finished LeaderSessionMock.Abdicate invocations
func (mmAbdicate *LeaderSessionMock) AbdicateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAbdicate.afterAbdicateCounter)
}

// AbdicateBeforeCounter returns a count of LeaderSessionMock.Abdicate invocations
func (mmAbdicate *LeaderSessionMock) AbdicateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAbdicate.beforeAbdicateCounter)
}

// Calls returns a list of arguments used in each call to LeaderSessionMock.Abdicate.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmAbdicate *mLeaderSessionMockAbdicate) Calls() []*LeaderSessionMockAbdicateParams {
	mmAbdicate.mutex.RLock()

	argCopy := make([]*LeaderSessionMockAbdicateParams, len(mmAbdicate.callArgs))
	copy(argCopy, mmAbdicate.callArgs)

	mmAbdicate.mutex.RUnlock()

	return argCopy
}

// MinimockAbdicateDone returns true if the count of the Abdicate invocations corresponds
// the number of defined expectations
func (m *LeaderSessionMock) MinimockAbdicateDone() bool {
	for _, e := range m.AbdicateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AbdicateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterAbdicateCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAbdicate != nil && mm_atomic.LoadUint64(&m.afterAbdicateCounter) < 1 {
		return false
	}
	return true
}

// MinimockAbdicateInspect logs each unmet expectation
func (m *LeaderSessionMock) MinimockAbdicateInspect() {
	for _, e := range m.AbdicateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LeaderSessionMock.Abdicate with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AbdicateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterAbdicateCounter) < 1 {
		if m.AbdicateMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to LeaderSessionMock.Abdicate")
		} else {
			m.t.Errorf("Expected call to LeaderSessionMock.Abdicate with params: %#v", *m.AbdicateMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAbdicate != nil && mm_atomic.LoadUint64(&m.afterAbdicateCounter) < 1 {
		m.t.Error("Expected call to LeaderSessionMock.Abdicate")
	}
}

type mLeaderSessionMockCurrent struct {
	mock               *LeaderSessionMock
	defaultExpectation *LeaderSessionMockCurrentExpectation
	expectations       []*LeaderSessionMockCurrentExpectation

	callArgs []*LeaderSessionMockCurrentParams
	mutex    sync.RWMutex
}

// LeaderSessionMockCurrentExpectation specifies expectation struct of the LeaderSession.Current
type LeaderSessionMockCurrentExpectation struct {
	mock    *LeaderSessionMock
	params  *LeaderSessionMockCurrentParams
	results *LeaderSessionMockCurrentResults
	Counter uint64
}

// LeaderSessionMockCurrentParams contains parameters of the LeaderSession.Current
type LeaderSessionMockCurrentParams struct {
	c1 Ctx
}

// LeaderSessionMockCurrentResults contains results of the LeaderSession.Current
type LeaderSessionMockCurrentResults struct {
	s1  string
	err error
}

// Expect sets up expected params for LeaderSession.Current
func (mmCurrent *mLeaderSessionMockCurrent) Expect(c1 Ctx) *mLeaderSessionMockCurrent {
	if mmCurrent.mock.funcCurrent != nil {
		mmCurrent.mock.t.Fatalf("LeaderSessionMock.Current mock is already set by Set")
	}

	if mmCurrent.defaultExpectation == nil {
		mmCurrent.defaultExpectation = &LeaderSessionMockCurrentExpectation{}
	}

	mmCurrent.defaultExpectation.params = &LeaderSessionMockCurrentParams{c1}
	for _, e := range mmCurrent.expectations {
		if minimock.Equal(e.params, mmCurrent.defaultExpectation.params) {
			mmCurrent.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmCurrent.defaultExpectation.params)
		}
	}

	return mmCurrent
}

// Inspect accepts an inspector function that has same arguments as the LeaderSession.Current
func (mmCurrent *mLeaderSessionMockCurrent) Inspect(f func(c1 Ctx)) *mLeaderSessionMockCurrent {
	if mmCurrent.mock.inspectFuncCurrent != nil {
		mmCurrent.mock.t.Fatalf("Inspect function is already set for LeaderSessionMock.Current")
	}

	mmCurrent.mock.inspectFuncCurrent = f

	return mmCurrent
}

// Return sets up results that will be returned by LeaderSession.Current
func (mmCurrent *mLeaderSessionMockCurrent) Return(s1 string, err error) *LeaderSessionMock {
	if mmCurrent.mock.funcCurrent != nil {
		mmCurrent.mock.t.Fatalf("LeaderSessionMock.Current mock is already set by Set")
	}

	if mmCurrent.defaultExpectation == nil {
		mmCurrent.defaultExpectation = &LeaderSessionMockCurrentExpectation{mock: mmCurrent.mock}
	}
	mmCurrent.defaultExpectation.results = &LeaderSessionMockCurrentResults{s1, err}
	return mmCurrent.mock
}

//Set uses given function f to mock the LeaderSession.Current method
func (mmCurrent *mLeaderSessionMockCurrent) Set(f func(c1 Ctx) (s1 string, err error)) *LeaderSessionMock {
	if mmCurrent.defaultExpectation != nil {
		mmCurrent.mock.t.Fatalf("Default expectation is already set for the LeaderSession.Current method")
	}

	if len(mmCurrent.expectations) > 0 {
		mmCurrent.mock.t.Fatalf("Some expectations are already set for the LeaderSession.Current method")
	}

	mmCurrent.mock.funcCurrent = f
	return mmCurrent.mock
}

// When sets expectation for the LeaderSession.Current which will trigger the result defined by the following
// Then helper
func (mmCurrent *mLeaderSessionMockCurrent) When(c1 Ctx) *LeaderSessionMockCurrentExpectation {
	if mmCurrent.mock.funcCurrent != nil {
		mmCurrent.mock.t.Fatalf("LeaderSessionMock.Current mock is already set by Set")
	}

	expectation := &LeaderSessionMockCurrentExpectation{
		mock:   mmCurrent.mock,
		params: &LeaderSessionMockCurrentParams{c1},
	}
	mmCurrent.expectations = append(mmCurrent.expectations, expectation)
	return expectation
}

// Then sets up LeaderSession.Current return parameters for the expectation previously defined by the When method
func (e *LeaderSessionMockCurrentExpectation) Then(s1 string, err error) *LeaderSessionMock {
	e.results = &LeaderSessionMockCurrentResults{s1, err}
	return e.mock
}

// Current implements LeaderSession
func (mmCurrent *LeaderSessionMock) Current(c1 Ctx) (s1 string, err error) {
	mm_atomic.AddUint64(&mmCurrent.beforeCurrentCounter, 1)
	defer mm_atomic.AddUint64(&mmCurrent.afterCurrentCounter, 1)

	if mmCurrent.inspectFuncCurrent != nil {
		mmCurrent.inspectFuncCurrent(c1)
	}

	mm_params := &LeaderSessionMockCurrentParams{c1}

	// Record call args
	mmCurrent.CurrentMock.mutex.Lock()
	mmCurrent.CurrentMock.callArgs = append(mmCurrent.CurrentMock.callArgs, mm_params)
	mmCurrent.CurrentMock.mutex.Unlock()

	for _, e := range mmCurrent.CurrentMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.s1, e.results.err
		}
	}

	if mmCurrent.CurrentMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmCurrent.CurrentMock.defaultExpectation.Counter, 1)
		mm_want := mmCurrent.CurrentMock.defaultExpectation.params
		mm_got := LeaderSessionMockCurrentParams{c1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmCurrent.t.Errorf("LeaderSessionMock.Current got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmCurrent.CurrentMock.defaultExpectation.results
		if mm_results == nil {
			mmCurrent.t.Fatal("No results are set for the LeaderSessionMock.Current")
		}
		return (*mm_results).s1, (*mm_results).err
	}
	if mmCurrent.funcCurrent != nil {
		return mmCurrent.funcCurrent(c1)
	}
	mmCurrent.t.Fatalf("Unexpected call to LeaderSessionMock.Current. %v", c1)
	return
}

// CurrentAfterCounter returns a count of finished LeaderSessionMock.Current invocations
func (mmCurrent *LeaderSessionMock) CurrentAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCurrent.afterCurrentCounter)
}

// CurrentBeforeCounter returns a count of LeaderSessionMock.Current invocations
func (mmCurrent *LeaderSessionMock) CurrentBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCurrent.beforeCurrentCounter)
}

// Calls returns a list of arguments used in each call to LeaderSessionMock.Current.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmCurrent *mLeaderSessionMockCurrent) Calls() []*LeaderSessionMockCurrentParams {
	mmCurrent.mutex.RLock()

	argCopy := make([]*LeaderSessionMockCurrentParams, len(mmCurrent.callArgs))
	copy(argCopy, mmCurrent.callArgs)

	mmCurrent.mutex.RUnlock()

	return argCopy
}

// MinimockCurrentDone returns true if the count of the Current invocations corresponds
// the number of defined expectations
func (m *LeaderSessionMock) MinimockCurrentDone() bool {
	for _, e := range m.CurrentMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CurrentMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterCurrentCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCurrent != nil && mm_atomic.LoadUint64(&m.afterCurrentCounter) < 1 {
		return false
	}
	return true
}

// MinimockCurrentInspect logs each unmet expectation
func (m *LeaderSessionMock) MinimockCurrentInspect() {
	for _, e := range m.CurrentMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LeaderSessionMock.Current with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CurrentMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterCurrentCounter) < 1 {
		if m.CurrentMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to LeaderSessionMock.Current")
		} else {
			m.t.Errorf("Expected call to LeaderSessionMock.Current with params: %#v", *m.CurrentMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCurrent != nil && mm_atomic.LoadUint64(&m.afterCurrentCounter) < 1 {
		m.t.Error("Expected call to LeaderSessionMock.Current")
	}
}

type mLeaderSessionMockSessionID struct {
	mock               *LeaderSessionMock
	defaultExpectation *LeaderSessionMockSessionIDExpectation
	expectations       []*LeaderSessionMockSessionIDExpectation

	callArgs []*LeaderSessionMockSessionIDParams
	mutex    sync.RWMutex
}

// LeaderSessionMockSessionIDExpectation specifies expectation struct of the LeaderSession.SessionID
type LeaderSessionMockSessionIDExpectation struct {
	mock    *LeaderSessionMock
	params  *LeaderSessionMockSessionIDParams
	results *LeaderSessionMockSessionIDResults
	Counter uint64
}

// LeaderSessionMockSessionIDParams contains parameters of the LeaderSession.SessionID
type LeaderSessionMockSessionIDParams struct {
	c1 Ctx
}

// LeaderSessionMockSessionIDResults contains results of the LeaderSession.SessionID
type LeaderSessionMockSessionIDResults struct {
	s1 string
}

// Expect sets up expected params for LeaderSession.SessionID
func (mmSessionID *mLeaderSessionMockSessionID) Expect(c1 Ctx) *mLeaderSessionMockSessionID {
	if mmSessionID.mock.funcSessionID != nil {
		mmSessionID.mock.t.Fatalf("LeaderSessionMock.SessionID mock is already set by Set")
	}

	if mmSessionID.defaultExpectation == nil {
		mmSessionID.defaultExpectation = &LeaderSessionMockSessionIDExpectation{}
	}

	mmSessionID.defaultExpectation.params = &LeaderSessionMockSessionIDParams{c1}
	for _, e := range mmSessionID.expectations {
		if minimock.Equal(e.params, mmSessionID.defaultExpectation.params) {
			mmSessionID.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmSessionID.defaultExpectation.params)
		}
	}

	return mmSessionID
}

// Inspect accepts an inspector function that has same arguments as the LeaderSession.SessionID
func (mmSessionID *mLeaderSessionMockSessionID) Inspect(f func(c1 Ctx)) *mLeaderSessionMockSessionID {
	if mmSessionID.mock.inspectFuncSessionID != nil {
		mmSessionID.mock.t.Fatalf("Inspect function is already set for LeaderSessionMock.SessionID")
	}

	mmSessionID.mock.inspectFuncSessionID = f

	return mmSessionID
}

// Return sets up results that will be returned by LeaderSession.SessionID
func (mmSessionID *mLeaderSessionMockSessionID) Return(s1 string) *LeaderSessionMock {
	if mmSessionID.mock.funcSessionID != nil {
		mmSessionID.mock.t.Fatalf("LeaderSessionMock.SessionID mock is already set by Set")
	}

	if mmSessionID.defaultExpectation == nil {
		mmSessionID.defaultExpectation = &LeaderSessionMockSessionIDExpectation{mock: mmSessionID.mock}
	}
	mmSessionID.defaultExpectation.results = &LeaderSessionMockSessionIDResults{s1}
	return mmSessionID.mock
}

//Set uses given function f to mock the LeaderSession.SessionID method
func (mmSessionID *mLeaderSessionMockSessionID) Set(f func(c1 Ctx) (s1 string)) *LeaderSessionMock {
	if mmSessionID.defaultExpectation != nil {
		mmSessionID.mock.t.Fatalf("Default expectation is already set for the LeaderSession.SessionID method")
	}

	if len(mmSessionID.expectations) > 0 {
		mmSessionID.mock.t.Fatalf("Some expectations are already set for the LeaderSession.SessionID method")
	}

	mmSessionID.mock.funcSessionID = f
	return mmSessionID.mock
}

// When sets expectation for the LeaderSession.SessionID which will trigger the result defined by the following
// Then helper
func (mmSessionID *mLeaderSessionMockSessionID) When(c1 Ctx) *LeaderSessionMockSessionIDExpectation {
	if mmSessionID.mock.funcSessionID != nil {
		mmSessionID.mock.t.Fatalf("LeaderSessionMock.SessionID mock is already set by Set")
	}

	expectation := &LeaderSessionMockSessionIDExpectation{
		mock:   mmSessionID.mock,
		params: &LeaderSessionMockSessionIDParams{c1},
	}
	mmSessionID.expectations = append(mmSessionID.expectations, expectation)
	return expectation
}

// Then sets up LeaderSession.SessionID return parameters for the expectation previously defined by the When method
func (e *LeaderSessionMockSessionIDExpectation) Then(s1 string) *LeaderSessionMock {
	e.results = &LeaderSessionMockSessionIDResults{s1}
	return e.mock
}

// SessionID implements LeaderSession
func (mmSessionID *LeaderSessionMock) SessionID(c1 Ctx) (s1 string) {
	mm_atomic.AddUint64(&mmSessionID.beforeSessionIDCounter, 1)
	defer mm_atomic.AddUint64(&mmSessionID.afterSessionIDCounter, 1)

	if mmSessionID.inspectFuncSessionID != nil {
		mmSessionID.inspectFuncSessionID(c1)
	}

	mm_params := &LeaderSessionMockSessionIDParams{c1}

	// Record call args
	mmSessionID.SessionIDMock.mutex.Lock()
	mmSessionID.SessionIDMock.callArgs = append(mmSessionID.SessionIDMock.callArgs, mm_params)
	mmSessionID.SessionIDMock.mutex.Unlock()

	for _, e := range mmSessionID.SessionIDMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.s1
		}
	}

	if mmSessionID.SessionIDMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmSessionID.SessionIDMock.defaultExpectation.Counter, 1)
		mm_want := mmSessionID.SessionIDMock.defaultExpectation.params
		mm_got := LeaderSessionMockSessionIDParams{c1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmSessionID.t.Errorf("LeaderSessionMock.SessionID got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmSessionID.SessionIDMock.defaultExpectation.results
		if mm_results == nil {
			mmSessionID.t.Fatal("No results are set for the LeaderSessionMock.SessionID")
		}
		return (*mm_results).s1
	}
	if mmSessionID.funcSessionID != nil {
		return mmSessionID.funcSessionID(c1)
	}
	mmSessionID.t.Fatalf("Unexpected call to LeaderSessionMock.SessionID. %v", c1)
	return
}

// SessionIDAfterCounter returns a count of finished LeaderSessionMock.SessionID invocations
func (mmSessionID *LeaderSessionMock) SessionIDAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSessionID.afterSessionIDCounter)
}

// SessionIDBeforeCounter returns a count of LeaderSessionMock.SessionID invocations
func (mmSessionID *LeaderSessionMock) SessionIDBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSessionID.beforeSessionIDCounter)
}

// Calls returns a list of arguments used in each call to LeaderSessionMock.SessionID.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmSessionID *mLeaderSessionMockSessionID) Calls() []*LeaderSessionMockSessionIDParams {
	mmSessionID.mutex.RLock()

	argCopy := make([]*LeaderSessionMockSessionIDParams, len(mmSessionID.callArgs))
	copy(argCopy, mmSessionID.callArgs)

	mmSessionID.mutex.RUnlock()

	return argCopy
}

// MinimockSessionIDDone returns true if the count of the SessionID invocations corresponds
// the number of defined expectations
func (m *LeaderSessionMock) MinimockSessionIDDone() bool {
	for _, e := range m.SessionIDMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SessionIDMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterSessionIDCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSessionID != nil && mm_atomic.LoadUint64(&m.afterSessionIDCounter) < 1 {
		return false
	}
	return true
}

// MinimockSessionIDInspect logs each unmet expectation
func (m *LeaderSessionMock) MinimockSessionIDInspect() {
	for _, e := range m.SessionIDMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LeaderSessionMock.SessionID with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SessionIDMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterSessionIDCounter) < 1 {
		if m.SessionIDMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to LeaderSessionMock.SessionID")
		} else {
			m.t.Errorf("Expected call to LeaderSessionMock.SessionID with params: %#v", *m.SessionIDMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSessionID != nil && mm_atomic.LoadUint64(&m.afterSessionIDCounter) < 1 {
		m.t.Error("Expected call to LeaderSessionMock.SessionID")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *LeaderSessionMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockAbdicateInspect()

		m.MinimockCurrentInspect()

		m.MinimockSessionIDInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *LeaderSessionMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *LeaderSessionMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockAbdicateDone() &&
		m.MinimockCurrentDone() &&
		m.MinimockSessionIDDone()
}
