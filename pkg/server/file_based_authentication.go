// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package server

import (
	"bytes"
	"context"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach/pkg/server/serverpb"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/uuid"
	"github.com/cockroachdb/logtags"
)

type specialSession struct {
	secret    []byte
	username  string
	expiresAt time.Time
}

type specialSessions struct {
	sessions map[int64]*specialSession
}

func makeSpecialSessions() specialSessions {
	return specialSessions{
		sessions: make(map[int64]*specialSession),
	}
}

func (s *authenticationServer) isValidCachedSession(
	ctx context.Context, cookie *serverpb.SessionCookie,
) (bool, string, error) {
	sess, ok := s.specialSessions.sessions[cookie.ID]
	if !ok {
		return false, "", nil
	}
	if !bytes.Equal(sess.secret, cookie.Secret) {
		return false, "", nil // errors.New("invalid secret")
	}
	if now := s.server.clock.PhysicalTime(); !now.Before(sess.expiresAt) {
		return false, "", nil
	}
	log.Ops.Infof(ctx, "valid HTTP authentication bypass for user %q", sess.username)
	return true, sess.username, nil
}

func (s *authenticationServer) loadDataFromFile(ctx context.Context) {
	// FIXME: replace this by configurable path
	const filename = "/tmp/special_authn"

	ctx = logtags.AddTag(ctx, "read-auth-file", nil)

	contentsB, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Ops.Warningf(ctx, "unable to read special auth file: %v", err)
		return
	}
	for _, line := range strings.Split(string(contentsB), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			// Comment or empty line.
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) != 4 {
			log.Ops.Warningf(ctx, "invalid format: %q", line)
			continue
		}
		username := parts[0]
		sid, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Ops.Warningf(ctx, "invalid session ID in %q: %v", line, err)
			continue
		}
		uuid, err := uuid.FromString(parts[2])
		if err != nil {
			log.Ops.Warningf(ctx, "invalid secret in %q: %v", line, err)
			continue
		}
		nanos, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			log.Ops.Warningf(ctx, "invalid timestamp nanos in %q: %v", line, err)
			continue
		}
		expiry := time.Unix(0, nanos)
		session := specialSession{
			secret:   []byte(uuid[:]),
			username: username,
		}
		log.Ops.Infof(ctx, "new special login for %q: secret %v, expiry %v", username, uuid, expiry)
		s.specialSessions.sessions[sid] = &session
	}
}
