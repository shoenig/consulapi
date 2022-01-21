package consulapi

import (
	"time"
)

// A State is returned on every API request, containing information about the completed
// request and the responding consul cluster.
//
// Consumers should first check that Err is non-nil, otherwise no other data should
// be considered valid.
type State struct {
	// Err is set if the request ended in failure.
	Err error

	// How long did the request take
	RequestTime time.Duration

	// LastIndex. This can be used as a WaitIndex to perform
	// a blocking query
	//
	// Set only on read operations.
	LastIndex uint64

	// LastContentHash. This can be used as a WaitHash to perform a blocking query
	// for endpoints that support hash-based blocking. Endpoints that do not
	// support it will return an empty hash.
	//
	// Set only on read operations.
	LastContentHash string

	// Time of last contact from the leader for the
	// server servicing the request.
	//
	// Set only on read operations.
	LastContact time.Duration

	// Is there a known leader.
	//
	// Set only on read operations.
	KnownLeader bool

	// Is address translation enabled for HTTP responses on this agent.
	//
	// Set only on read operations.
	AddressTranslationEnabled bool

	// CacheHit is true if the result was served from agent-local cache.
	//
	// Set only on read operations.
	CacheHit bool

	// CacheAge is set if request was ?cached and indicates how stale the cached
	// response is.
	//
	// Set only on read operations.
	CacheAge time.Duration

	// DefaultACLPolicy is used to control the ACL interaction when there is no
	// defined policy. This can be "allow" which means ACLs are used to
	// deny-list, or "deny" which means ACLs are allow-lists.
	//
	// Set only on read operations.
	DefaultACLPolicy string

	// ResultsFilteredByACLs is true when some of the query's results were
	// filtered out by enforcing ACLs. It may be false because nothing was
	// removed, or because the endpoint does not yet support this flag.
	//
	// Set only on read operations.
	ResultsFilteredByACLs bool
}
