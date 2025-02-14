/*
Copyright 2017 Google LLC

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

package spanner

import (
	"container/heap"
	"container/list"
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/smyte/google-cloud-go/internal/trace"
	vkit "github.com/smyte/google-cloud-go/spanner/apiv1"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// sessionHandle is an interface for transactions to access Cloud Spanner
// sessions safely. It is generated by sessionPool.take().
type sessionHandle struct {
	// mu guarantees that the inner session object is returned / destroyed only
	// once.
	mu sync.Mutex
	// session is a pointer to a session object. Transactions never need to
	// access it directly.
	session *session
}

// recycle gives the inner session object back to its home session pool. It is
// safe to call recycle multiple times but only the first one would take effect.
func (sh *sessionHandle) recycle() {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.session == nil {
		// sessionHandle has already been recycled.
		return
	}
	sh.session.recycle()
	sh.session = nil
}

// getID gets the Cloud Spanner session ID from the internal session object.
// getID returns empty string if the sessionHandle is nil or the inner session
// object has been released by recycle / destroy.
func (sh *sessionHandle) getID() string {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.session == nil {
		// sessionHandle has already been recycled/destroyed.
		return ""
	}
	return sh.session.getID()
}

// getClient gets the Cloud Spanner RPC client associated with the session ID
// in sessionHandle.
func (sh *sessionHandle) getClient() *vkit.Client {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.session == nil {
		return nil
	}
	return sh.session.client
}

// getMetadata returns the metadata associated with the session in sessionHandle.
func (sh *sessionHandle) getMetadata() metadata.MD {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.session == nil {
		return nil
	}
	return sh.session.md
}

// getTransactionID returns the transaction id in the session if available.
func (sh *sessionHandle) getTransactionID() transactionID {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.session == nil {
		return nil
	}
	return sh.session.tx
}

// destroy destroys the inner session object. It is safe to call destroy
// multiple times and only the first call would attempt to
// destroy the inner session object.
func (sh *sessionHandle) destroy() {
	sh.mu.Lock()
	s := sh.session
	sh.session = nil
	sh.mu.Unlock()
	if s == nil {
		// sessionHandle has already been destroyed..
		return
	}
	s.destroy(false)
}

// session wraps a Cloud Spanner session ID through which transactions are
// created and executed.
type session struct {
	// client is the RPC channel to Cloud Spanner. It is set only once during
	// session's creation.
	client *vkit.Client
	// id is the unique id of the session in Cloud Spanner. It is set only once
	// during session's creation.
	id string
	// pool is the session's home session pool where it was created. It is set
	// only once during session's creation.
	pool *sessionPool
	// createTime is the timestamp of the session's creation. It is set only
	// once during session's creation.
	createTime time.Time

	// mu protects the following fields from concurrent access: both
	// healthcheck workers and transactions can modify them.
	mu sync.Mutex
	// valid marks the validity of a session.
	valid bool
	// hcIndex is the index of the session inside the global healthcheck queue.
	// If hcIndex < 0, session has been unregistered from the queue.
	hcIndex int
	// idleList is the linkedlist node which links the session to its home
	// session pool's idle list. If idleList == nil, the
	// session is not in idle list.
	idleList *list.Element
	// nextCheck is the timestamp of next scheduled healthcheck of the session.
	// It is maintained by the global health checker.
	nextCheck time.Time
	// checkingHelath is true if currently this session is being processed by
	// health checker. Must be modified under health checker lock.
	checkingHealth bool
	// md is the Metadata to be sent with each request.
	md metadata.MD
	// tx contains the transaction id if the session has been prepared for
	// write.
	tx transactionID
}

// isValid returns true if the session is still valid for use.
func (s *session) isValid() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.valid
}

// isWritePrepared returns true if the session is prepared for write.
func (s *session) isWritePrepared() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tx != nil
}

// String implements fmt.Stringer for session.
func (s *session) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("<id=%v, hcIdx=%v, idleList=%p, valid=%v, create=%v, nextcheck=%v>",
		s.id, s.hcIndex, s.idleList, s.valid, s.createTime, s.nextCheck)
}

// ping verifies if the session is still alive in Cloud Spanner.
func (s *session) ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// s.getID is safe even when s is invalid.
	_, err := s.client.GetSession(contextWithOutgoingMetadata(ctx, s.pool.md), &sppb.GetSessionRequest{Name: s.getID()})
	return err
}

// setHcIndex atomically sets the session's index in the healthcheck queue and
// returns the old index.
func (s *session) setHcIndex(i int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	oi := s.hcIndex
	s.hcIndex = i
	return oi
}

// setIdleList atomically sets the session's idle list link and returns the old
// link.
func (s *session) setIdleList(le *list.Element) *list.Element {
	s.mu.Lock()
	defer s.mu.Unlock()
	old := s.idleList
	s.idleList = le
	return old
}

// invalidate marks a session as invalid and returns the old validity.
func (s *session) invalidate() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	ov := s.valid
	s.valid = false
	return ov
}

// setNextCheck sets the timestamp for next healthcheck on the session.
func (s *session) setNextCheck(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextCheck = t
}

// setTransactionID sets the transaction id in the session
func (s *session) setTransactionID(tx transactionID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tx = tx
}

// getID returns the session ID which uniquely identifies the session in Cloud
// Spanner.
func (s *session) getID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.id
}

// getHcIndex returns the session's index into the global healthcheck priority
// queue.
func (s *session) getHcIndex() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hcIndex
}

// getIdleList returns the session's link in its home session pool's idle list.
func (s *session) getIdleList() *list.Element {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.idleList
}

// getNextCheck returns the timestamp for next healthcheck on the session.
func (s *session) getNextCheck() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.nextCheck
}

// recycle turns the session back to its home session pool.
func (s *session) recycle() {
	s.setTransactionID(nil)
	if !s.pool.recycle(s) {
		// s is rejected by its home session pool because it expired and the
		// session pool currently has enough open sessions.
		s.destroy(false)
	}
}

// destroy removes the session from its home session pool, healthcheck queue
// and Cloud Spanner service.
func (s *session) destroy(isExpire bool) bool {
	// Remove s from session pool.
	if !s.pool.remove(s, isExpire) {
		return false
	}
	// Unregister s from healthcheck queue.
	s.pool.hc.unregister(s)
	// Remove s from Cloud Spanner service.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	s.delete(ctx)
	return true
}

func (s *session) delete(ctx context.Context) {
	// Ignore the error because even if we fail to explicitly destroy the
	// session, it will be eventually garbage collected by Cloud Spanner.
	err := s.client.DeleteSession(ctx, &sppb.DeleteSessionRequest{Name: s.getID()})
	if err != nil {
		log.Printf("Failed to delete session %v. Error: %v", s.getID(), err)
	}
}

// prepareForWrite prepares the session for write if it is not already in that
// state.
func (s *session) prepareForWrite(ctx context.Context) error {
	if s.isWritePrepared() {
		return nil
	}
	tx, err := beginTransaction(ctx, s.getID(), s.client)
	if err != nil {
		return err
	}
	s.setTransactionID(tx)
	return nil
}

// SessionPoolConfig stores configurations of a session pool.
type SessionPoolConfig struct {
	// getRPCClient is the caller supplied method for getting a gRPC client to
	// Cloud Spanner, this makes session pool able to use client pooling.
	getRPCClient func() (*vkit.Client, error)

	// MaxOpened is the maximum number of opened sessions allowed by the session
	// pool. If the client tries to open a session and there are already
	// MaxOpened sessions, it will block until one becomes available or the
	// context passed to the client method is canceled or times out.
	//
	// Defaults to NumChannels * 100.
	MaxOpened uint64

	// MinOpened is the minimum number of opened sessions that the session pool
	// tries to maintain. Session pool won't continue to expire sessions if
	// number of opened connections drops below MinOpened. However, if a session
	// is found to be broken, it will still be evicted from the session pool,
	// therefore it is posssible that the number of opened sessions drops below
	// MinOpened.
	//
	// Defaults to 0.
	MinOpened uint64

	// MaxIdle is the maximum number of idle sessions, pool is allowed to keep.
	//
	// Defaults to 0.
	MaxIdle uint64

	// MaxBurst is the maximum number of concurrent session creation requests.
	//
	// Defaults to 10.
	MaxBurst uint64

	// WriteSessions is the fraction of sessions we try to keep prepared for
	// write.
	//
	// Defaults to 0.
	WriteSessions float64

	// HealthCheckWorkers is number of workers used by health checker for this
	// pool.
	//
	// Defaults to 10.
	HealthCheckWorkers int

	// HealthCheckInterval is how often the health checker pings a session.
	//
	// Defaults to 5m.
	HealthCheckInterval time.Duration

	// healthCheckSampleInterval is how often the health checker samples live
	// session (for use in maintaining session pool size).
	//
	// Defaults to 1m.
	healthCheckSampleInterval time.Duration

	// sessionLabels for the sessions created in the session pool.
	sessionLabels map[string]string
}

// errNoRPCGetter returns error for SessionPoolConfig missing getRPCClient method.
func errNoRPCGetter() error {
	return spannerErrorf(codes.InvalidArgument, "require SessionPoolConfig.getRPCClient != nil, got nil")
}

// errMinOpenedGTMapOpened returns error for SessionPoolConfig.MaxOpened < SessionPoolConfig.MinOpened when SessionPoolConfig.MaxOpened is set.
func errMinOpenedGTMaxOpened(maxOpened, minOpened uint64) error {
	return spannerErrorf(codes.InvalidArgument,
		"require SessionPoolConfig.MaxOpened >= SessionPoolConfig.MinOpened, got %v and %v", maxOpened, minOpened)
}

// validate verifies that the SessionPoolConfig is good for use.
func (spc *SessionPoolConfig) validate() error {
	if spc.getRPCClient == nil {
		return errNoRPCGetter()
	}
	if spc.MinOpened > spc.MaxOpened && spc.MaxOpened > 0 {
		return errMinOpenedGTMaxOpened(spc.MaxOpened, spc.MinOpened)
	}
	return nil
}

// sessionPool creates and caches Cloud Spanner sessions.
type sessionPool struct {
	// mu protects sessionPool from concurrent access.
	mu sync.Mutex
	// valid marks the validity of the session pool.
	valid bool
	// db is the database name that all sessions in the pool are associated with.
	db string
	// idleList caches idle session IDs. Session IDs in this list can be
	// allocated for use.
	idleList list.List
	// idleWriteList caches idle sessions which have been prepared for write.
	idleWriteList list.List
	// mayGetSession is for broadcasting that session retrival/creation may
	// proceed.
	mayGetSession chan struct{}
	// numOpened is the total number of open sessions from the session pool.
	numOpened uint64
	// createReqs is the number of ongoing session creation requests.
	createReqs uint64
	// prepareReqs is the number of ongoing session preparation request.
	prepareReqs uint64
	// configuration of the session pool.
	SessionPoolConfig
	// Metadata to be sent with each request
	md metadata.MD
	// hc is the health checker
	hc *healthChecker
}

// newSessionPool creates a new session pool.
func newSessionPool(db string, config SessionPoolConfig, md metadata.MD) (*sessionPool, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}
	pool := &sessionPool{
		db:                db,
		valid:             true,
		mayGetSession:     make(chan struct{}),
		SessionPoolConfig: config,
		md:                md,
	}
	if config.HealthCheckWorkers == 0 {
		// With 10 workers and assuming average latency of 5ms for
		// BeginTransaction, we will be able to prepare 2000 tx/sec in advance.
		// If the rate of takeWriteSession is more than that, it will degrade to
		// doing BeginTransaction inline.
		//
		// TODO: consider resizing the worker pool dynamically according to the load.
		config.HealthCheckWorkers = 10
	}
	if config.HealthCheckInterval == 0 {
		config.HealthCheckInterval = 5 * time.Minute
	}
	if config.healthCheckSampleInterval == 0 {
		config.healthCheckSampleInterval = time.Minute
	}
	// On GCE VM, within the same region an healthcheck ping takes on average
	// 10ms to finish, given a 5 minutes interval and 10 healthcheck workers, a
	// healthChecker can effectively mantain
	// 100 checks_per_worker/sec * 10 workers * 300 seconds = 300K sessions.
	pool.hc = newHealthChecker(config.HealthCheckInterval, config.HealthCheckWorkers, config.healthCheckSampleInterval, pool)
	close(pool.hc.ready)
	return pool, nil
}

// isValid checks if the session pool is still valid.
func (p *sessionPool) isValid() bool {
	if p == nil {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.valid
}

// close marks the session pool as closed.
func (p *sessionPool) close() {
	if p == nil {
		return
	}
	p.mu.Lock()
	if !p.valid {
		p.mu.Unlock()
		return
	}
	p.valid = false
	p.mu.Unlock()
	p.hc.close()
	// destroy all the sessions
	p.hc.mu.Lock()
	allSessions := make([]*session, len(p.hc.queue.sessions))
	copy(allSessions, p.hc.queue.sessions)
	p.hc.mu.Unlock()
	for _, s := range allSessions {
		s.destroy(false)
	}
}

// errInvalidSessionPool returns error for using an invalid session pool.
func errInvalidSessionPool() error {
	return spannerErrorf(codes.InvalidArgument, "invalid session pool")
}

// errGetSessionTimeout returns error for context timeout during
// sessionPool.take().
func errGetSessionTimeout() error {
	return spannerErrorf(codes.Canceled, "timeout / context canceled during getting session")
}

// shouldPrepareWrite returns true if we should prepare more sessions for write.
func (p *sessionPool) shouldPrepareWrite() bool {
	return float64(p.numOpened)*p.WriteSessions > float64(p.idleWriteList.Len()+int(p.prepareReqs))
}

func (p *sessionPool) createSession(ctx context.Context) (*session, error) {
	trace.TracePrintf(ctx, nil, "Creating a new session")
	doneCreate := func(done bool) {
		p.mu.Lock()
		if !done {
			// Session creation failed, give budget back.
			p.numOpened--
			recordStat(ctx, OpenSessionCount, int64(p.numOpened))
		}
		p.createReqs--
		// Notify other waiters blocking on session creation.
		close(p.mayGetSession)
		p.mayGetSession = make(chan struct{})
		p.mu.Unlock()
	}
	sc, err := p.getRPCClient()
	if err != nil {
		doneCreate(false)
		return nil, err
	}
	s, err := createSession(ctx, sc, p.db, p.sessionLabels, p.md)
	if err != nil {
		doneCreate(false)
		// Should return error directly because of the previous retries on
		// CreateSession RPC.
		// If the error is a timeout, there is a chance that the session was
		// created on the server but is not known to the session pool. This
		// session will then be garbage collected by the server after 1 hour.
		return nil, err
	}
	s.pool = p
	p.hc.register(s)
	doneCreate(true)
	return s, nil
}

func createSession(ctx context.Context, sc *vkit.Client, db string, labels map[string]string, md metadata.MD) (*session, error) {
	var s *session
	sid, e := sc.CreateSession(ctx, &sppb.CreateSessionRequest{
		Database: db,
		Session:  &sppb.Session{Labels: labels},
	})
	if e != nil {
		return nil, toSpannerError(e)
	}
	// If no error, construct the new session.
	s = &session{valid: true, client: sc, id: sid.Name, createTime: time.Now(), md: md}
	return s, nil
}

func (p *sessionPool) isHealthy(s *session) bool {
	if s.getNextCheck().Add(2 * p.hc.getInterval()).Before(time.Now()) {
		// TODO: figure out if we need to schedule a new healthcheck worker here.
		if err := s.ping(); shouldDropSession(err) {
			// The session is already bad, continue to fetch/create a new one.
			s.destroy(false)
			return false
		}
		p.hc.scheduledHC(s)
	}
	return true
}

// take returns a cached session if there are available ones; if there isn't
// any, it tries to allocate a new one. Session returned by take should be used
// for read operations.
func (p *sessionPool) take(ctx context.Context) (*sessionHandle, error) {
	trace.TracePrintf(ctx, nil, "Acquiring a read-only session")
	ctx = contextWithOutgoingMetadata(ctx, p.md)
	for {
		var (
			s   *session
			err error
		)

		p.mu.Lock()
		if !p.valid {
			p.mu.Unlock()
			return nil, errInvalidSessionPool()
		}
		if p.idleList.Len() > 0 {
			// Idle sessions are available, get one from the top of the idle
			// list.
			s = p.idleList.Remove(p.idleList.Front()).(*session)
			trace.TracePrintf(ctx, map[string]interface{}{"sessionID": s.getID()},
				"Acquired read-only session")
		} else if p.idleWriteList.Len() > 0 {
			s = p.idleWriteList.Remove(p.idleWriteList.Front()).(*session)
			trace.TracePrintf(ctx, map[string]interface{}{"sessionID": s.getID()},
				"Acquired read-write session")
		}
		if s != nil {
			s.setIdleList(nil)
			p.mu.Unlock()
			// From here, session is no longer in idle list, so healthcheck
			// workers won't destroy it. If healthcheck workers failed to
			// schedule healthcheck for the session timely, do the check here.
			// Because session check is still much cheaper than session
			// creation, they should be reused as much as possible.
			if !p.isHealthy(s) {
				continue
			}
			return &sessionHandle{session: s}, nil
		}

		// Idle list is empty, block if session pool has reached max session
		// creation concurrency or max number of open sessions.
		if (p.MaxOpened > 0 && p.numOpened >= p.MaxOpened) || (p.MaxBurst > 0 && p.createReqs >= p.MaxBurst) {
			mayGetSession := p.mayGetSession
			p.mu.Unlock()
			trace.TracePrintf(ctx, nil, "Waiting for read-only session to become available")
			select {
			case <-ctx.Done():
				trace.TracePrintf(ctx, nil, "Context done waiting for session")
				return nil, errGetSessionTimeout()
			case <-mayGetSession:
			}
			continue
		}

		// Take budget before the actual session creation.
		p.numOpened++
		recordStat(ctx, OpenSessionCount, int64(p.numOpened))
		p.createReqs++
		p.mu.Unlock()
		if s, err = p.createSession(ctx); err != nil {
			trace.TracePrintf(ctx, nil, "Error creating session: %v", err)
			return nil, toSpannerError(err)
		}
		trace.TracePrintf(ctx, map[string]interface{}{"sessionID": s.getID()},
			"Created session")
		return &sessionHandle{session: s}, nil
	}
}

// takeWriteSession returns a write prepared cached session if there are
// available ones; if there isn't any, it tries to allocate a new one. Session
// returned should be used for read write transactions.
func (p *sessionPool) takeWriteSession(ctx context.Context) (*sessionHandle, error) {
	trace.TracePrintf(ctx, nil, "Acquiring a read-write session")
	ctx = contextWithOutgoingMetadata(ctx, p.md)
	for {
		var (
			s   *session
			err error
		)

		p.mu.Lock()
		if !p.valid {
			p.mu.Unlock()
			return nil, errInvalidSessionPool()
		}
		if p.idleWriteList.Len() > 0 {
			// Idle sessions are available, get one from the top of the idle
			// list.
			s = p.idleWriteList.Remove(p.idleWriteList.Front()).(*session)
			trace.TracePrintf(ctx, map[string]interface{}{"sessionID": s.getID()}, "Acquired read-write session")
		} else if p.idleList.Len() > 0 {
			s = p.idleList.Remove(p.idleList.Front()).(*session)
			trace.TracePrintf(ctx, map[string]interface{}{"sessionID": s.getID()}, "Acquired read-only session")
		}
		if s != nil {
			s.setIdleList(nil)
			p.mu.Unlock()
			// From here, session is no longer in idle list, so healthcheck
			// workers won't destroy it. If healthcheck workers failed to
			// schedule healthcheck for the session timely, do the check here.
			// Because session check is still much cheaper than session
			// creation, they should be reused as much as possible.
			if !p.isHealthy(s) {
				continue
			}
		} else {
			// Idle list is empty, block if session pool has reached max session
			// creation concurrency or max number of open sessions.
			if (p.MaxOpened > 0 && p.numOpened >= p.MaxOpened) || (p.MaxBurst > 0 && p.createReqs >= p.MaxBurst) {
				mayGetSession := p.mayGetSession
				p.mu.Unlock()
				trace.TracePrintf(ctx, nil, "Waiting for read-write session to become available")
				select {
				case <-ctx.Done():
					trace.TracePrintf(ctx, nil, "Context done waiting for session")
					return nil, errGetSessionTimeout()
				case <-mayGetSession:
				}
				continue
			}

			// Take budget before the actual session creation.
			p.numOpened++
			recordStat(ctx, OpenSessionCount, int64(p.numOpened))
			p.createReqs++
			p.mu.Unlock()
			if s, err = p.createSession(ctx); err != nil {
				trace.TracePrintf(ctx, nil, "Error creating session: %v", err)
				return nil, toSpannerError(err)
			}
			trace.TracePrintf(ctx, map[string]interface{}{"sessionID": s.getID()},
				"Created session")
		}
		if !s.isWritePrepared() {
			if err = s.prepareForWrite(ctx); err != nil {
				s.recycle()
				trace.TracePrintf(ctx, map[string]interface{}{"sessionID": s.getID()},
					"Error preparing session for write")
				return nil, toSpannerError(err)
			}
		}
		return &sessionHandle{session: s}, nil
	}
}

// recycle puts session s back to the session pool's idle list, it returns true
// if the session pool successfully recycles session s.
func (p *sessionPool) recycle(s *session) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !s.isValid() || !p.valid {
		// Reject the session if session is invalid or pool itself is invalid.
		return false
	}
	// Put session at the back of the list to round robin for load balancing
	// across channels.
	if s.isWritePrepared() {
		s.setIdleList(p.idleWriteList.PushBack(s))
	} else {
		s.setIdleList(p.idleList.PushBack(s))
	}
	// Broadcast that a session has been returned to idle list.
	close(p.mayGetSession)
	p.mayGetSession = make(chan struct{})
	return true
}

// remove atomically removes session s from the session pool and invalidates s.
// If isExpire == true, the removal is triggered by session expiration and in
// such cases, only idle sessions can be removed.
func (p *sessionPool) remove(s *session, isExpire bool) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if isExpire && (p.numOpened <= p.MinOpened || s.getIdleList() == nil) {
		// Don't expire session if the session is not in idle list (in use), or
		// if number of open sessions is going below p.MinOpened.
		return false
	}
	ol := s.setIdleList(nil)
	// If the session is in the idlelist, remove it.
	if ol != nil {
		// Remove from whichever list it is in.
		p.idleList.Remove(ol)
		p.idleWriteList.Remove(ol)
	}
	if s.invalidate() {
		// Decrease the number of opened sessions.
		p.numOpened--
		recordStat(context.Background(), OpenSessionCount, int64(p.numOpened))
		// Broadcast that a session has been destroyed.
		close(p.mayGetSession)
		p.mayGetSession = make(chan struct{})
		return true
	}
	return false
}

// hcHeap implements heap.Interface. It is used to create the priority queue for
// session healthchecks.
type hcHeap struct {
	sessions []*session
}

// Len impelemnts heap.Interface.Len.
func (h hcHeap) Len() int {
	return len(h.sessions)
}

// Less implements heap.Interface.Less.
func (h hcHeap) Less(i, j int) bool {
	return h.sessions[i].getNextCheck().Before(h.sessions[j].getNextCheck())
}

// Swap implements heap.Interface.Swap.
func (h hcHeap) Swap(i, j int) {
	h.sessions[i], h.sessions[j] = h.sessions[j], h.sessions[i]
	h.sessions[i].setHcIndex(i)
	h.sessions[j].setHcIndex(j)
}

// Push implements heap.Interface.Push.
func (h *hcHeap) Push(s interface{}) {
	ns := s.(*session)
	ns.setHcIndex(len(h.sessions))
	h.sessions = append(h.sessions, ns)
}

// Pop implements heap.Interface.Pop.
func (h *hcHeap) Pop() interface{} {
	old := h.sessions
	n := len(old)
	s := old[n-1]
	h.sessions = old[:n-1]
	s.setHcIndex(-1)
	return s
}

// healthChecker performs periodical healthchecks on registered sessions.
type healthChecker struct {
	// mu protects concurrent access to healthChecker.
	mu sync.Mutex
	// queue is the priority queue for session healthchecks. Sessions with lower
	// nextCheck rank higher in the queue.
	queue hcHeap
	// interval is the average interval between two healthchecks on a session.
	interval time.Duration
	// workers is the number of concurrent healthcheck workers.
	workers int
	// waitWorkers waits for all healthcheck workers to exit
	waitWorkers sync.WaitGroup
	// pool is the underlying session pool.
	pool *sessionPool
	// sampleInterval is the interval of sampling by the maintainer.
	sampleInterval time.Duration
	// ready is used to signal that maintainer can start running.
	ready chan struct{}
	// done is used to signal that health checker should be closed.
	done chan struct{}
	// once is used for closing channel done only once.
	once             sync.Once
	maintainerCancel func()
}

// newHealthChecker initializes new instance of healthChecker.
func newHealthChecker(interval time.Duration, workers int, sampleInterval time.Duration, pool *sessionPool) *healthChecker {
	if workers <= 0 {
		workers = 1
	}
	hc := &healthChecker{
		interval:         interval,
		workers:          workers,
		pool:             pool,
		sampleInterval:   sampleInterval,
		ready:            make(chan struct{}),
		done:             make(chan struct{}),
		maintainerCancel: func() {},
	}
	hc.waitWorkers.Add(1)
	go hc.maintainer()
	for i := 1; i <= hc.workers; i++ {
		hc.waitWorkers.Add(1)
		go hc.worker(i)
	}
	return hc
}

// close closes the healthChecker and waits for all healthcheck workers to exit.
func (hc *healthChecker) close() {
	hc.mu.Lock()
	hc.maintainerCancel()
	hc.mu.Unlock()
	hc.once.Do(func() { close(hc.done) })
	hc.waitWorkers.Wait()
}

// isClosing checks if a healthChecker is already closing.
func (hc *healthChecker) isClosing() bool {
	select {
	case <-hc.done:
		return true
	default:
		return false
	}
}

// getInterval gets the healthcheck interval.
func (hc *healthChecker) getInterval() time.Duration {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	return hc.interval
}

// scheduledHCLocked schedules next healthcheck on session s with the assumption
// that hc.mu is being held.
func (hc *healthChecker) scheduledHCLocked(s *session) {
	// The next healthcheck will be scheduled after
	// [interval*0.5, interval*1.5) ns.
	nsFromNow := rand.Int63n(int64(hc.interval)) + int64(hc.interval)/2
	s.setNextCheck(time.Now().Add(time.Duration(nsFromNow)))
	if hi := s.getHcIndex(); hi != -1 {
		// Session is still being tracked by healthcheck workers.
		heap.Fix(&hc.queue, hi)
	}
}

// scheduledHC schedules next healthcheck on session s. It is safe to be called
// concurrently.
func (hc *healthChecker) scheduledHC(s *session) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.scheduledHCLocked(s)
}

// register registers a session with healthChecker for periodical healthcheck.
func (hc *healthChecker) register(s *session) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.scheduledHCLocked(s)
	heap.Push(&hc.queue, s)
}

// unregister unregisters a session from healthcheck queue.
func (hc *healthChecker) unregister(s *session) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	oi := s.setHcIndex(-1)
	if oi >= 0 {
		heap.Remove(&hc.queue, oi)
	}
}

// markDone marks that health check for session has been performed.
func (hc *healthChecker) markDone(s *session) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	s.checkingHealth = false
}

// healthCheck checks the health of the session and pings it if needed.
func (hc *healthChecker) healthCheck(s *session) {
	defer hc.markDone(s)
	if !s.pool.isValid() {
		// Session pool is closed, perform a garbage collection.
		s.destroy(false)
		return
	}
	if err := s.ping(); shouldDropSession(err) {
		// Ping failed, destroy the session.
		s.destroy(false)
	}
}

// worker performs the healthcheck on sessions in healthChecker's priority
// queue.
func (hc *healthChecker) worker(i int) {
	// Returns a session which we should ping to keep it alive.
	getNextForPing := func() *session {
		hc.pool.mu.Lock()
		defer hc.pool.mu.Unlock()
		hc.mu.Lock()
		defer hc.mu.Unlock()
		if hc.queue.Len() <= 0 {
			// Queue is empty.
			return nil
		}
		s := hc.queue.sessions[0]
		if s.getNextCheck().After(time.Now()) && hc.pool.valid {
			// All sessions have been checked recently.
			return nil
		}
		hc.scheduledHCLocked(s)
		if !s.checkingHealth {
			s.checkingHealth = true
			return s
		}
		return nil
	}

	// Returns a session which we should prepare for write.
	getNextForTx := func() *session {
		hc.pool.mu.Lock()
		defer hc.pool.mu.Unlock()
		if hc.pool.shouldPrepareWrite() {
			if hc.pool.idleList.Len() > 0 && hc.pool.valid {
				hc.mu.Lock()
				defer hc.mu.Unlock()
				if hc.pool.idleList.Front().Value.(*session).checkingHealth {
					return nil
				}
				session := hc.pool.idleList.Remove(hc.pool.idleList.Front()).(*session)
				session.checkingHealth = true
				hc.pool.prepareReqs++
				return session
			}
		}
		return nil
	}

	for {
		if hc.isClosing() {
			// Exit when the pool has been closed and all sessions have been
			// destroyed or when health checker has been closed.
			hc.waitWorkers.Done()
			return
		}
		ws := getNextForTx()
		if ws != nil {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := ws.prepareForWrite(contextWithOutgoingMetadata(ctx, hc.pool.md))
			cancel()
			if err != nil {
				// Skip handling prepare error, session can be prepared in next
				// cycle.
				log.Printf("Failed to prepare session, error: %v", toSpannerError(err))
			}
			hc.pool.recycle(ws)
			hc.pool.mu.Lock()
			hc.pool.prepareReqs--
			hc.pool.mu.Unlock()
			hc.markDone(ws)
		}
		rs := getNextForPing()
		if rs == nil {
			if ws == nil {
				// No work to be done so sleep to avoid burning CPU.
				pause := int64(100 * time.Millisecond)
				if pause > int64(hc.interval) {
					pause = int64(hc.interval)
				}
				select {
				case <-time.After(time.Duration(rand.Int63n(pause) + pause/2)):
				case <-hc.done:
				}

			}
			continue
		}
		hc.healthCheck(rs)
	}
}

// maintainer maintains the maxSessionsInUse by a window of
// kWindowSize * sampleInterval. Based on this information, health checker will
// try to maintain the number of sessions by hc.
func (hc *healthChecker) maintainer() {
	// Wait so that pool is ready.
	<-hc.ready

	windowSize := uint64(10)

	for iteration := uint64(0); ; iteration++ {
		if hc.isClosing() {
			hc.waitWorkers.Done()
			return
		}

		// maxSessionsInUse is the maximum number of sessions in use
		// concurrently over a period of time.
		var maxSessionsInUse uint64

		// Updates metrics.
		hc.pool.mu.Lock()
		currSessionsInUse := hc.pool.numOpened - uint64(hc.pool.idleList.Len()) - uint64(hc.pool.idleWriteList.Len())
		currSessionsOpened := hc.pool.numOpened
		hc.pool.mu.Unlock()

		hc.mu.Lock()
		if iteration%windowSize == 0 || maxSessionsInUse < currSessionsInUse {
			maxSessionsInUse = currSessionsInUse
		}
		sessionsToKeep := maxUint64(hc.pool.MinOpened,
			minUint64(currSessionsOpened, hc.pool.MaxIdle+maxSessionsInUse))
		ctx, cancel := context.WithTimeout(context.Background(), hc.sampleInterval)
		hc.maintainerCancel = cancel
		hc.mu.Unlock()

		// Replenish or Shrink pool if needed.
		//
		// Note: we don't need to worry about pending create session requests,
		// we only need to sample the current sessions in use. The routines will
		// not try to create extra / delete creating sessions.
		if sessionsToKeep > currSessionsOpened {
			hc.replenishPool(ctx, sessionsToKeep)
		} else {
			hc.shrinkPool(ctx, sessionsToKeep)
		}

		select {
		case <-ctx.Done():
		case <-hc.done:
			cancel()
		}
	}
}

// replenishPool is run if numOpened is less than sessionsToKeep, timeouts on
// sampleInterval.
func (hc *healthChecker) replenishPool(ctx context.Context, sessionsToKeep uint64) {
	for {
		if ctx.Err() != nil {
			return
		}

		p := hc.pool
		p.mu.Lock()
		// Take budget before the actual session creation.
		if sessionsToKeep <= p.numOpened {
			p.mu.Unlock()
			break
		}
		p.numOpened++
		recordStat(ctx, OpenSessionCount, int64(p.numOpened))
		p.createReqs++
		shouldPrepareWrite := p.shouldPrepareWrite()
		p.mu.Unlock()
		var (
			s   *session
			err error
		)
		createContext, cancel := context.WithTimeout(context.Background(), time.Minute)
		if s, err = p.createSession(createContext); err != nil {
			cancel()
			log.Printf("Failed to create session, error: %v", toSpannerError(err))
			continue
		}
		cancel()
		if shouldPrepareWrite {
			prepareContext, cancel := context.WithTimeout(context.Background(), time.Minute)
			if err = s.prepareForWrite(prepareContext); err != nil {
				cancel()
				p.recycle(s)
				log.Printf("Failed to prepare session, error: %v", toSpannerError(err))
				continue
			}
			cancel()
		}
		p.recycle(s)
	}
}

// shrinkPool, scales down the session pool.
func (hc *healthChecker) shrinkPool(ctx context.Context, sessionsToKeep uint64) {
	for {
		if ctx.Err() != nil {
			return
		}

		p := hc.pool
		p.mu.Lock()

		if sessionsToKeep >= p.numOpened {
			p.mu.Unlock()
			break
		}

		var s *session
		if p.idleList.Len() > 0 {
			s = p.idleList.Front().Value.(*session)
		} else if p.idleWriteList.Len() > 0 {
			s = p.idleWriteList.Front().Value.(*session)
		}
		p.mu.Unlock()
		if s != nil {
			// destroy session as expire.
			s.destroy(true)
		} else {
			break
		}
	}
}

// shouldDropSession returns true if a particular error leads to the removal of
// a session
func shouldDropSession(err error) bool {
	if err == nil {
		return false
	}
	// If a Cloud Spanner can no longer locate the session (for example, if
	// session is garbage collected), then caller should not try to return the
	// session back into the session pool.
	//
	// TODO: once gRPC can return auxiliary error information, stop parsing the error message.
	if ErrCode(err) == codes.NotFound && strings.Contains(ErrDesc(err), "Session not found") {
		return true
	}
	return false
}

// maxUint64 returns the maximum of two uint64.
func maxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

// minUint64 returns the minimum of two uint64.
func minUint64(a, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}
