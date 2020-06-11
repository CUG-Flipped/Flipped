// Copyright 2018 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nats-io/jwt"
	"github.com/nats-io/nkeys"
)

var (
	// This matches ./configs/nkeys_jwts/test.seed
	oSeed = []byte("SOAFYNORQLQFJYBYNUGC5D7SH2MXMUX5BFEWWGHN3EK4VGG5TPT5DZP7QU")
)

func opTrustBasicSetup() *Server {
	kp, _ := nkeys.FromSeed(oSeed)
	pub, _ := kp.PublicKey()
	opts := defaultServerOptions
	opts.TrustedKeys = []string{pub}
	s, c, _, _ := rawSetup(opts)
	c.close()
	return s
}

func buildMemAccResolver(s *Server) {
	mr := &MemAccResolver{}
	s.SetAccountResolver(mr)
}

func addAccountToMemResolver(s *Server, pub, jwtclaim string) {
	s.AccountResolver().Store(pub, jwtclaim)
}

func createClient(t *testing.T, s *Server, akp nkeys.KeyPair) (*testAsyncClient, *bufio.Reader, string) {
	return createClientWithIssuer(t, s, akp, "")
}

func createClientWithIssuer(t *testing.T, s *Server, akp nkeys.KeyPair, optIssuerAccount string) (*testAsyncClient, *bufio.Reader, string) {
	t.Helper()
	nkp, _ := nkeys.CreateUser()
	pub, _ := nkp.PublicKey()
	nuc := jwt.NewUserClaims(pub)
	if optIssuerAccount != "" {
		nuc.IssuerAccount = optIssuerAccount
	}
	ujwt, err := nuc.Encode(akp)
	if err != nil {
		t.Fatalf("Error generating user JWT: %v", err)
	}
	c, cr, l := newClientForServer(s)

	// Sign Nonce
	var info nonceInfo
	json.Unmarshal([]byte(l[5:]), &info)
	sigraw, _ := nkp.Sign([]byte(info.Nonce))
	sig := base64.RawURLEncoding.EncodeToString(sigraw)

	cs := fmt.Sprintf("CONNECT {\"jwt\":%q,\"sig\":\"%s\"}\r\nPING\r\n", ujwt, sig)
	return c, cr, cs
}

func setupJWTTestWithClaims(t *testing.T, nac *jwt.AccountClaims, nuc *jwt.UserClaims, expected string) (*Server, nkeys.KeyPair, *testAsyncClient, *bufio.Reader) {
	t.Helper()

	okp, _ := nkeys.FromSeed(oSeed)

	akp, _ := nkeys.CreateAccount()
	apub, _ := akp.PublicKey()
	if nac == nil {
		nac = jwt.NewAccountClaims(apub)
	} else {
		nac.Subject = apub
	}
	ajwt, err := nac.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	nkp, _ := nkeys.CreateUser()
	pub, _ := nkp.PublicKey()
	if nuc == nil {
		nuc = jwt.NewUserClaims(pub)
	} else {
		nuc.Subject = pub
	}
	jwt, err := nuc.Encode(akp)
	if err != nil {
		t.Fatalf("Error generating user JWT: %v", err)
	}

	s := opTrustBasicSetup()
	buildMemAccResolver(s)
	addAccountToMemResolver(s, apub, ajwt)

	c, cr, l := newClientForServer(s)

	// Sign Nonce
	var info nonceInfo
	json.Unmarshal([]byte(l[5:]), &info)
	sigraw, _ := nkp.Sign([]byte(info.Nonce))
	sig := base64.RawURLEncoding.EncodeToString(sigraw)

	// PING needed to flush the +OK/-ERR to us.
	cs := fmt.Sprintf("CONNECT {\"jwt\":%q,\"sig\":\"%s\",\"verbose\":true,\"pedantic\":true}\r\nPING\r\n", jwt, sig)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		c.parse([]byte(cs))
		wg.Done()
	}()
	l, _ = cr.ReadString('\n')
	if !strings.HasPrefix(l, expected) {
		t.Fatalf("Expected %q, got %q", expected, l)
	}
	wg.Wait()

	return s, akp, c, cr
}

func setupJWTTestWitAccountClaims(t *testing.T, nac *jwt.AccountClaims, expected string) (*Server, nkeys.KeyPair, *testAsyncClient, *bufio.Reader) {
	t.Helper()
	return setupJWTTestWithClaims(t, nac, nil, expected)
}

// This is used in test to create account claims and pass it
// to setupJWTTestWitAccountClaims.
func newJWTTestAccountClaims() *jwt.AccountClaims {
	// We call NewAccountClaims() because it sets some defaults.
	// However, this call needs a subject, but the real subject will
	// be set in setupJWTTestWitAccountClaims(). Use some temporary one
	// here.
	return jwt.NewAccountClaims("temp")
}

func setupJWTTestWithUserClaims(t *testing.T, nuc *jwt.UserClaims, expected string) (*Server, *testAsyncClient, *bufio.Reader) {
	t.Helper()
	s, _, c, cr := setupJWTTestWithClaims(t, nil, nuc, expected)
	return s, c, cr
}

// This is used in test to create user claims and pass it
// to setupJWTTestWithUserClaims.
func newJWTTestUserClaims() *jwt.UserClaims {
	// As of now, tests could simply do &jwt.UserClaims{}, but in
	// case some defaults are later added, we call NewUserClaims().
	// However, this call needs a subject, but the real subject will
	// be set in setupJWTTestWithUserClaims(). Use some temporary one
	// here.
	return jwt.NewUserClaims("temp")
}

func TestJWTUser(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()

	// Check to make sure we would have an authTimer
	if !s.info.AuthRequired {
		t.Fatalf("Expect the server to require auth")
	}

	c, cr, _ := newClientForServer(s)
	defer c.close()

	// Don't send jwt field, should fail.
	c.parseAsync("CONNECT {\"verbose\":true,\"pedantic\":true}\r\nPING\r\n")
	l, _ := cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}

	okp, _ := nkeys.FromSeed(oSeed)

	// Create an account that will be expired.
	akp, _ := nkeys.CreateAccount()
	apub, _ := akp.PublicKey()
	nac := jwt.NewAccountClaims(apub)
	ajwt, err := nac.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	c, cr, cs := createClient(t, s, akp)
	defer c.close()

	// PING needed to flush the +OK/-ERR to us.
	// This should fail too since no account resolver is defined.
	c.parseAsync(cs)
	l, _ = cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}

	// Ok now let's walk through and make sure all is good.
	// We will set the account resolver by hand to a memory resolver.
	buildMemAccResolver(s)
	addAccountToMemResolver(s, apub, ajwt)

	c, cr, cs = createClient(t, s, akp)
	defer c.close()

	c.parseAsync(cs)
	l, _ = cr.ReadString('\n')
	if !strings.HasPrefix(l, "PONG") {
		t.Fatalf("Expected a PONG, got %q", l)
	}
}

func TestJWTUserBadTrusted(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()

	// Check to make sure we would have an authTimer
	if !s.info.AuthRequired {
		t.Fatalf("Expect the server to require auth")
	}
	// Now place bad trusted key
	s.mu.Lock()
	s.trustedKeys = []string{"bad"}
	s.mu.Unlock()

	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Create an account that will be expired.
	akp, _ := nkeys.CreateAccount()
	apub, _ := akp.PublicKey()
	nac := jwt.NewAccountClaims(apub)
	ajwt, err := nac.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, apub, ajwt)

	c, cr, cs := createClient(t, s, akp)
	defer c.close()
	c.parseAsync(cs)
	l, _ := cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}
}

// Test that if a user tries to connect with an expired user JWT we do the right thing.
func TestJWTUserExpired(t *testing.T) {
	nuc := newJWTTestUserClaims()
	nuc.IssuedAt = time.Now().Add(-10 * time.Second).Unix()
	nuc.Expires = time.Now().Add(-2 * time.Second).Unix()
	s, c, _ := setupJWTTestWithUserClaims(t, nuc, "-ERR ")
	c.close()
	s.Shutdown()
}

func TestJWTUserExpiresAfterConnect(t *testing.T) {
	nuc := newJWTTestUserClaims()
	nuc.IssuedAt = time.Now().Unix()
	nuc.Expires = time.Now().Add(time.Second).Unix()
	s, c, cr := setupJWTTestWithUserClaims(t, nuc, "+OK")
	defer s.Shutdown()
	defer c.close()
	l, _ := cr.ReadString('\n')
	if !strings.HasPrefix(l, "PONG") {
		t.Fatalf("Expected a PONG")
	}

	// Now we should expire after 1 second or so.
	time.Sleep(1250 * time.Millisecond)

	l, _ = cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}
	if !strings.Contains(l, "Expired") {
		t.Fatalf("Expected 'Expired' to be in the error")
	}
}

func TestJWTUserPermissionClaims(t *testing.T) {
	nuc := newJWTTestUserClaims()
	nuc.Permissions.Pub.Allow.Add("foo")
	nuc.Permissions.Pub.Allow.Add("bar")
	nuc.Permissions.Pub.Deny.Add("baz")
	nuc.Permissions.Sub.Allow.Add("foo")
	nuc.Permissions.Sub.Allow.Add("bar")
	nuc.Permissions.Sub.Deny.Add("baz")

	s, c, _ := setupJWTTestWithUserClaims(t, nuc, "+OK")
	defer s.Shutdown()
	defer c.close()

	// Now check client to make sure permissions transferred.
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.perms == nil {
		t.Fatalf("Expected client permissions to be set")
	}

	if lpa := c.perms.pub.allow.Count(); lpa != 2 {
		t.Fatalf("Expected 2 publish allow subjects, got %d", lpa)
	}
	if lpd := c.perms.pub.deny.Count(); lpd != 1 {
		t.Fatalf("Expected 1 publish deny subjects, got %d", lpd)
	}
	if lsa := c.perms.sub.allow.Count(); lsa != 2 {
		t.Fatalf("Expected 2 subscribe allow subjects, got %d", lsa)
	}
	if lsd := c.perms.sub.deny.Count(); lsd != 1 {
		t.Fatalf("Expected 1 subscribe deny subjects, got %d", lsd)
	}
}

func TestJWTUserResponsePermissionClaims(t *testing.T) {
	nuc := newJWTTestUserClaims()
	nuc.Permissions.Resp = &jwt.ResponsePermission{
		MaxMsgs: 22,
		Expires: 100 * time.Millisecond,
	}
	s, c, _ := setupJWTTestWithUserClaims(t, nuc, "+OK")
	defer s.Shutdown()
	defer c.close()

	// Now check client to make sure permissions transferred.
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.perms == nil {
		t.Fatalf("Expected client permissions to be set")
	}
	if c.perms.pub.allow == nil {
		t.Fatalf("Expected client perms for pub allow to be non-nil")
	}
	if lpa := c.perms.pub.allow.Count(); lpa != 0 {
		t.Fatalf("Expected 0 publish allow subjects, got %d", lpa)
	}
	if c.perms.resp == nil {
		t.Fatalf("Expected client perms for response permissions to be non-nil")
	}
	if c.perms.resp.MaxMsgs != nuc.Permissions.Resp.MaxMsgs {
		t.Fatalf("Expected client perms for response permissions MaxMsgs to be same as jwt: %d vs %d",
			c.perms.resp.MaxMsgs, nuc.Permissions.Resp.MaxMsgs)
	}
	if c.perms.resp.Expires != nuc.Permissions.Resp.Expires {
		t.Fatalf("Expected client perms for response permissions Expires to be same as jwt: %v vs %v",
			c.perms.resp.Expires, nuc.Permissions.Resp.Expires)
	}
}

func TestJWTUserResponsePermissionClaimsDefaultValues(t *testing.T) {
	nuc := newJWTTestUserClaims()
	nuc.Permissions.Resp = &jwt.ResponsePermission{}
	s, c, _ := setupJWTTestWithUserClaims(t, nuc, "+OK")
	defer s.Shutdown()
	defer c.close()

	// Now check client to make sure permissions transferred
	// and defaults are set.
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.perms == nil {
		t.Fatalf("Expected client permissions to be set")
	}
	if c.perms.pub.allow == nil {
		t.Fatalf("Expected client perms for pub allow to be non-nil")
	}
	if lpa := c.perms.pub.allow.Count(); lpa != 0 {
		t.Fatalf("Expected 0 publish allow subjects, got %d", lpa)
	}
	if c.perms.resp == nil {
		t.Fatalf("Expected client perms for response permissions to be non-nil")
	}
	if c.perms.resp.MaxMsgs != DEFAULT_ALLOW_RESPONSE_MAX_MSGS {
		t.Fatalf("Expected client perms for response permissions MaxMsgs to be default %v, got %v",
			DEFAULT_ALLOW_RESPONSE_MAX_MSGS, c.perms.resp.MaxMsgs)
	}
	if c.perms.resp.Expires != DEFAULT_ALLOW_RESPONSE_EXPIRATION {
		t.Fatalf("Expected client perms for response permissions Expires to be default %v, got %v",
			DEFAULT_ALLOW_RESPONSE_EXPIRATION, c.perms.resp.Expires)
	}
}

func TestJWTUserResponsePermissionClaimsNegativeValues(t *testing.T) {
	nuc := newJWTTestUserClaims()
	nuc.Permissions.Resp = &jwt.ResponsePermission{
		MaxMsgs: -1,
		Expires: -1 * time.Second,
	}
	s, c, _ := setupJWTTestWithUserClaims(t, nuc, "+OK")
	defer s.Shutdown()
	defer c.close()

	// Now check client to make sure permissions transferred
	// and negative values are transferred.
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.perms == nil {
		t.Fatalf("Expected client permissions to be set")
	}
	if c.perms.pub.allow == nil {
		t.Fatalf("Expected client perms for pub allow to be non-nil")
	}
	if lpa := c.perms.pub.allow.Count(); lpa != 0 {
		t.Fatalf("Expected 0 publish allow subjects, got %d", lpa)
	}
	if c.perms.resp == nil {
		t.Fatalf("Expected client perms for response permissions to be non-nil")
	}
	if c.perms.resp.MaxMsgs != -1 {
		t.Fatalf("Expected client perms for response permissions MaxMsgs to be %v, got %v",
			-1, c.perms.resp.MaxMsgs)
	}
	if c.perms.resp.Expires != -1*time.Second {
		t.Fatalf("Expected client perms for response permissions Expires to be %v, got %v",
			-1*time.Second, c.perms.resp.Expires)
	}
}

func TestJWTAccountExpired(t *testing.T) {
	nac := newJWTTestAccountClaims()
	nac.IssuedAt = time.Now().Add(-10 * time.Second).Unix()
	nac.Expires = time.Now().Add(-2 * time.Second).Unix()
	s, _, c, _ := setupJWTTestWitAccountClaims(t, nac, "-ERR ")
	defer s.Shutdown()
	defer c.close()
}

func TestJWTAccountExpiresAfterConnect(t *testing.T) {
	nac := newJWTTestAccountClaims()
	now := time.Now()
	nac.IssuedAt = now.Add(-10 * time.Second).Unix()
	nac.Expires = now.Round(time.Second).Add(time.Second).Unix()
	s, akp, c, cr := setupJWTTestWitAccountClaims(t, nac, "+OK")
	defer s.Shutdown()
	defer c.close()

	apub, _ := akp.PublicKey()
	acc, err := s.LookupAccount(apub)
	if acc == nil || err != nil {
		t.Fatalf("Expected to retrieve the account")
	}

	if l, _ := cr.ReadString('\n'); !strings.HasPrefix(l, "PONG") {
		t.Fatalf("Expected PONG, got %q", l)
	}

	// Wait for the account to be expired.
	checkFor(t, 3*time.Second, 100*time.Millisecond, func() error {
		if acc.IsExpired() {
			return nil
		}
		return fmt.Errorf("Account not expired yet")
	})

	l, _ := cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error, got %q", l)
	}
	if !strings.Contains(l, "Expired") {
		t.Fatalf("Expected 'Expired' to be in the error")
	}

	// Now make sure that accounts that have expired return an error.
	c, cr, cs := createClient(t, s, akp)
	defer c.close()
	c.parseAsync(cs)
	l, _ = cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}
}

func TestJWTAccountRenew(t *testing.T) {
	nac := newJWTTestAccountClaims()
	// Create an account that has expired.
	nac.IssuedAt = time.Now().Add(-10 * time.Second).Unix()
	nac.Expires = time.Now().Add(-2 * time.Second).Unix()
	// Expect an error
	s, akp, c, _ := setupJWTTestWitAccountClaims(t, nac, "-ERR ")
	defer s.Shutdown()
	defer c.close()

	okp, _ := nkeys.FromSeed(oSeed)
	apub, _ := akp.PublicKey()

	// Now update with new expiration
	nac.IssuedAt = time.Now().Unix()
	nac.Expires = time.Now().Add(5 * time.Second).Unix()
	ajwt, err := nac.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	// Update the account
	addAccountToMemResolver(s, apub, ajwt)
	acc, _ := s.LookupAccount(apub)
	if acc == nil {
		t.Fatalf("Expected to retrieve the account")
	}
	s.updateAccountClaims(acc, nac)

	// Now make sure we can connect.
	c, cr, cs := createClient(t, s, akp)
	defer c.close()
	c.parseAsync(cs)
	if l, _ := cr.ReadString('\n'); !strings.HasPrefix(l, "PONG") {
		t.Fatalf("Expected a PONG, got: %q", l)
	}
}

func TestJWTAccountRenewFromResolver(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	akp, _ := nkeys.CreateAccount()
	apub, _ := akp.PublicKey()
	nac := jwt.NewAccountClaims(apub)
	nac.IssuedAt = time.Now().Add(-10 * time.Second).Unix()
	nac.Expires = time.Now().Add(time.Second).Unix()
	ajwt, err := nac.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	addAccountToMemResolver(s, apub, ajwt)
	// Force it to be loaded by the server and start the expiration timer.
	acc, _ := s.LookupAccount(apub)
	if acc == nil {
		t.Fatalf("Could not retrieve account for %q", apub)
	}

	// Create a new user
	c, cr, cs := createClient(t, s, akp)
	defer c.close()
	// Wait for expiration.
	time.Sleep(1250 * time.Millisecond)

	c.parseAsync(cs)
	l, _ := cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}

	// Now update with new expiration
	nac.IssuedAt = time.Now().Unix()
	nac.Expires = time.Now().Add(5 * time.Second).Unix()
	ajwt, err = nac.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	// Update the account
	addAccountToMemResolver(s, apub, ajwt)
	// Make sure the too quick update suppression does not bite us.
	acc.mu.Lock()
	acc.updated = time.Now().Add(-1 * time.Hour)
	acc.mu.Unlock()

	// Do not update the account directly. The resolver should
	// happen automatically.

	// Now make sure we can connect.
	c, cr, cs = createClient(t, s, akp)
	defer c.close()
	c.parseAsync(cs)
	l, _ = cr.ReadString('\n')
	if !strings.HasPrefix(l, "PONG") {
		t.Fatalf("Expected a PONG, got: %q", l)
	}
}

func TestJWTAccountBasicImportExport(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Create accounts and imports/exports.
	fooKP, _ := nkeys.CreateAccount()
	fooPub, _ := fooKP.PublicKey()
	fooAC := jwt.NewAccountClaims(fooPub)

	// Now create Exports.
	streamExport := &jwt.Export{Subject: "foo", Type: jwt.Stream}
	streamExport2 := &jwt.Export{Subject: "private", Type: jwt.Stream, TokenReq: true}
	serviceExport := &jwt.Export{Subject: "req.echo", Type: jwt.Service, TokenReq: true}
	serviceExport2 := &jwt.Export{Subject: "req.add", Type: jwt.Service, TokenReq: true}

	fooAC.Exports.Add(streamExport, streamExport2, serviceExport, serviceExport2)
	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	addAccountToMemResolver(s, fooPub, fooJWT)

	acc, _ := s.LookupAccount(fooPub)
	if acc == nil {
		t.Fatalf("Expected to retrieve the account")
	}

	// Check to make sure exports transferred over.
	if les := len(acc.exports.streams); les != 2 {
		t.Fatalf("Expected exports streams len of 2, got %d", les)
	}
	if les := len(acc.exports.services); les != 2 {
		t.Fatalf("Expected exports services len of 2, got %d", les)
	}
	_, ok := acc.exports.streams["foo"]
	if !ok {
		t.Fatalf("Expected to map a stream export")
	}
	se, ok := acc.exports.services["req.echo"]
	if !ok || se == nil {
		t.Fatalf("Expected to map a service export")
	}
	if !se.tokenReq {
		t.Fatalf("Expected the service export to require tokens")
	}

	barKP, _ := nkeys.CreateAccount()
	barPub, _ := barKP.PublicKey()
	barAC := jwt.NewAccountClaims(barPub)

	streamImport := &jwt.Import{Account: fooPub, Subject: "foo", To: "import.foo", Type: jwt.Stream}
	serviceImport := &jwt.Import{Account: fooPub, Subject: "req.echo", Type: jwt.Service}
	barAC.Imports.Add(streamImport, serviceImport)
	barJWT, err := barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)

	acc, _ = s.LookupAccount(barPub)
	if acc == nil {
		t.Fatalf("Expected to retrieve the account")
	}
	if les := len(acc.imports.streams); les != 1 {
		t.Fatalf("Expected imports streams len of 1, got %d", les)
	}
	// Our service import should have failed without a token.
	if les := len(acc.imports.services); les != 0 {
		t.Fatalf("Expected imports services len of 0, got %d", les)
	}

	// Now add in a bad activation token.
	barAC = jwt.NewAccountClaims(barPub)
	serviceImport = &jwt.Import{Account: fooPub, Subject: "req.echo", Token: "not a token", Type: jwt.Service}
	barAC.Imports.Add(serviceImport)
	barJWT, err = barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)

	s.updateAccountClaims(acc, barAC)

	// Our service import should have failed with a bad token.
	if les := len(acc.imports.services); les != 0 {
		t.Fatalf("Expected imports services len of 0, got %d", les)
	}

	// Now make a correct one.
	barAC = jwt.NewAccountClaims(barPub)
	serviceImport = &jwt.Import{Account: fooPub, Subject: "req.echo", Type: jwt.Service}

	activation := jwt.NewActivationClaims(barPub)
	activation.ImportSubject = "req.echo"
	activation.ImportType = jwt.Service
	actJWT, err := activation.Encode(fooKP)
	if err != nil {
		t.Fatalf("Error generating activation token: %v", err)
	}
	serviceImport.Token = actJWT
	barAC.Imports.Add(serviceImport)
	barJWT, err = barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)
	s.updateAccountClaims(acc, barAC)
	// Our service import should have succeeded.
	if les := len(acc.imports.services); les != 1 {
		t.Fatalf("Expected imports services len of 1, got %d", les)
	}

	// Now test url
	barAC = jwt.NewAccountClaims(barPub)
	serviceImport = &jwt.Import{Account: fooPub, Subject: "req.add", Type: jwt.Service}

	activation = jwt.NewActivationClaims(barPub)
	activation.ImportSubject = "req.add"
	activation.ImportType = jwt.Service
	actJWT, err = activation.Encode(fooKP)
	if err != nil {
		t.Fatalf("Error generating activation token: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(actJWT))
	}))
	defer ts.Close()

	serviceImport.Token = ts.URL
	barAC.Imports.Add(serviceImport)
	barJWT, err = barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)
	s.updateAccountClaims(acc, barAC)
	// Our service import should have succeeded. Should be the only one since we reset.
	if les := len(acc.imports.services); les != 1 {
		t.Fatalf("Expected imports services len of 1, got %d", les)
	}

	// Now streams
	barAC = jwt.NewAccountClaims(barPub)
	streamImport = &jwt.Import{Account: fooPub, Subject: "private", To: "import.private", Type: jwt.Stream}

	barAC.Imports.Add(streamImport)
	barJWT, err = barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)
	s.updateAccountClaims(acc, barAC)
	// Our stream import should have not succeeded.
	if les := len(acc.imports.streams); les != 0 {
		t.Fatalf("Expected imports services len of 0, got %d", les)
	}

	// Now add in activation.
	barAC = jwt.NewAccountClaims(barPub)
	streamImport = &jwt.Import{Account: fooPub, Subject: "private", To: "import.private", Type: jwt.Stream}

	activation = jwt.NewActivationClaims(barPub)
	activation.ImportSubject = "private"
	activation.ImportType = jwt.Stream
	actJWT, err = activation.Encode(fooKP)
	if err != nil {
		t.Fatalf("Error generating activation token: %v", err)
	}
	streamImport.Token = actJWT
	barAC.Imports.Add(streamImport)
	barJWT, err = barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)
	s.updateAccountClaims(acc, barAC)
	// Our stream import should have not succeeded.
	if les := len(acc.imports.streams); les != 1 {
		t.Fatalf("Expected imports services len of 1, got %d", les)
	}
}

func TestJWTAccountExportWithResponseType(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Create accounts and imports/exports.
	fooKP, _ := nkeys.CreateAccount()
	fooPub, _ := fooKP.PublicKey()
	fooAC := jwt.NewAccountClaims(fooPub)

	// Now create Exports.
	serviceStreamExport := &jwt.Export{Subject: "test.stream", Type: jwt.Service, ResponseType: jwt.ResponseTypeStream, TokenReq: false}
	serviceChunkExport := &jwt.Export{Subject: "test.chunk", Type: jwt.Service, ResponseType: jwt.ResponseTypeChunked, TokenReq: false}
	serviceSingletonExport := &jwt.Export{Subject: "test.single", Type: jwt.Service, ResponseType: jwt.ResponseTypeSingleton, TokenReq: true}
	serviceDefExport := &jwt.Export{Subject: "test.def", Type: jwt.Service, TokenReq: true}
	serviceOldExport := &jwt.Export{Subject: "test.old", Type: jwt.Service, TokenReq: false}

	fooAC.Exports.Add(serviceStreamExport, serviceSingletonExport, serviceChunkExport, serviceDefExport, serviceOldExport)
	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	addAccountToMemResolver(s, fooPub, fooJWT)

	fooAcc, _ := s.LookupAccount(fooPub)
	if fooAcc == nil {
		t.Fatalf("Expected to retrieve the account")
	}

	services := fooAcc.exports.services

	if len(services) != 5 {
		t.Fatalf("Expected 4 services")
	}

	se, ok := services["test.stream"]
	if !ok || se == nil {
		t.Fatalf("Expected to map a service export")
	}
	if se.tokenReq {
		t.Fatalf("Expected the service export to not require tokens")
	}
	if se.respType != Stream {
		t.Fatalf("Expected the service export to respond with a stream")
	}

	se, ok = services["test.chunk"]
	if !ok || se == nil {
		t.Fatalf("Expected to map a service export")
	}
	if se.tokenReq {
		t.Fatalf("Expected the service export to not require tokens")
	}
	if se.respType != Chunked {
		t.Fatalf("Expected the service export to respond with a stream")
	}

	se, ok = services["test.def"]
	if !ok || se == nil {
		t.Fatalf("Expected to map a service export")
	}
	if !se.tokenReq {
		t.Fatalf("Expected the service export to not require tokens")
	}
	if se.respType != Singleton {
		t.Fatalf("Expected the service export to respond with a stream")
	}

	se, ok = services["test.single"]
	if !ok || se == nil {
		t.Fatalf("Expected to map a service export")
	}
	if !se.tokenReq {
		t.Fatalf("Expected the service export to not require tokens")
	}
	if se.respType != Singleton {
		t.Fatalf("Expected the service export to respond with a stream")
	}

	se, ok = services["test.old"]
	if !ok || se != nil {
		t.Fatalf("Service with a singleton response and no tokens should be nil in the map")
	}
}

func expectPong(t *testing.T, cr *bufio.Reader) {
	t.Helper()
	l, _ := cr.ReadString('\n')
	if !strings.HasPrefix(l, "PONG") {
		t.Fatalf("Expected a PONG, got %q", l)
	}
}

func expectMsg(t *testing.T, cr *bufio.Reader, sub, payload string) {
	t.Helper()
	l, _ := cr.ReadString('\n')
	expected := "MSG " + sub
	if !strings.HasPrefix(l, expected) {
		t.Fatalf("Expected %q, got %q", expected, l)
	}
	l, _ = cr.ReadString('\n')
	if l != payload+"\r\n" {
		t.Fatalf("Expected %q, got %q", payload, l)
	}
	expectPong(t, cr)
}

func TestJWTAccountImportExportUpdates(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Create accounts and imports/exports.
	fooKP, _ := nkeys.CreateAccount()
	fooPub, _ := fooKP.PublicKey()
	fooAC := jwt.NewAccountClaims(fooPub)
	streamExport := &jwt.Export{Subject: "foo", Type: jwt.Stream}

	fooAC.Exports.Add(streamExport)
	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, fooPub, fooJWT)

	barKP, _ := nkeys.CreateAccount()
	barPub, _ := barKP.PublicKey()
	barAC := jwt.NewAccountClaims(barPub)
	streamImport := &jwt.Import{Account: fooPub, Subject: "foo", To: "import", Type: jwt.Stream}

	barAC.Imports.Add(streamImport)
	barJWT, err := barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)

	// Create a client.
	c, cr, cs := createClient(t, s, barKP)
	defer c.close()

	c.parseAsync(cs)
	expectPong(t, cr)

	c.parseAsync("SUB import.foo 1\r\nPING\r\n")
	expectPong(t, cr)

	checkShadow := func(expected int) {
		t.Helper()
		c.mu.Lock()
		defer c.mu.Unlock()
		sub := c.subs["1"]
		if ls := len(sub.shadow); ls != expected {
			t.Fatalf("Expected shadows to be %d, got %d", expected, ls)
		}
	}

	// We created a SUB on foo which should create a shadow subscription.
	checkShadow(1)

	// Now update bar and remove the import which should make the shadow go away.
	barAC = jwt.NewAccountClaims(barPub)
	barJWT, _ = barAC.Encode(okp)
	addAccountToMemResolver(s, barPub, barJWT)
	acc, _ := s.LookupAccount(barPub)
	s.updateAccountClaims(acc, barAC)

	checkShadow(0)

	// Now add it back and make sure the shadow comes back.
	streamImport = &jwt.Import{Account: string(fooPub), Subject: "foo", To: "import", Type: jwt.Stream}
	barAC.Imports.Add(streamImport)
	barJWT, _ = barAC.Encode(okp)
	addAccountToMemResolver(s, barPub, barJWT)
	s.updateAccountClaims(acc, barAC)

	checkShadow(1)

	// Now change export and make sure it goes away as well. So no exports anymore.
	fooAC = jwt.NewAccountClaims(fooPub)
	fooJWT, _ = fooAC.Encode(okp)
	addAccountToMemResolver(s, fooPub, fooJWT)
	acc, _ = s.LookupAccount(fooPub)
	s.updateAccountClaims(acc, fooAC)
	checkShadow(0)

	// Now add it in but with permission required.
	streamExport = &jwt.Export{Subject: "foo", Type: jwt.Stream, TokenReq: true}
	fooAC.Exports.Add(streamExport)
	fooJWT, _ = fooAC.Encode(okp)
	addAccountToMemResolver(s, fooPub, fooJWT)
	s.updateAccountClaims(acc, fooAC)

	checkShadow(0)

	// Now put it back as normal.
	fooAC = jwt.NewAccountClaims(fooPub)
	streamExport = &jwt.Export{Subject: "foo", Type: jwt.Stream}
	fooAC.Exports.Add(streamExport)
	fooJWT, _ = fooAC.Encode(okp)
	addAccountToMemResolver(s, fooPub, fooJWT)
	s.updateAccountClaims(acc, fooAC)

	checkShadow(1)
}

func TestJWTAccountImportActivationExpires(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Create accounts and imports/exports.
	fooKP, _ := nkeys.CreateAccount()
	fooPub, _ := fooKP.PublicKey()
	fooAC := jwt.NewAccountClaims(fooPub)
	streamExport := &jwt.Export{Subject: "foo", Type: jwt.Stream, TokenReq: true}
	fooAC.Exports.Add(streamExport)

	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	addAccountToMemResolver(s, fooPub, fooJWT)
	acc, _ := s.LookupAccount(fooPub)
	if acc == nil {
		t.Fatalf("Expected to retrieve the account")
	}

	barKP, _ := nkeys.CreateAccount()
	barPub, _ := barKP.PublicKey()
	barAC := jwt.NewAccountClaims(barPub)
	streamImport := &jwt.Import{Account: fooPub, Subject: "foo", To: "import.", Type: jwt.Stream}

	activation := jwt.NewActivationClaims(barPub)
	activation.ImportSubject = "foo"
	activation.ImportType = jwt.Stream
	now := time.Now()
	activation.IssuedAt = now.Add(-10 * time.Second).Unix()
	// These are second resolution. So round up before adding a second.
	activation.Expires = now.Round(time.Second).Add(time.Second).Unix()
	actJWT, err := activation.Encode(fooKP)
	if err != nil {
		t.Fatalf("Error generating activation token: %v", err)
	}
	streamImport.Token = actJWT
	barAC.Imports.Add(streamImport)
	barJWT, err := barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)
	if acc, _ := s.LookupAccount(barPub); acc == nil {
		t.Fatalf("Expected to retrieve the account")
	}

	// Create a client.
	c, cr, cs := createClient(t, s, barKP)
	defer c.close()

	c.parseAsync(cs)
	expectPong(t, cr)

	c.parseAsync("SUB import.foo 1\r\nPING\r\n")
	expectPong(t, cr)

	checkShadow := func(t *testing.T, expected int) {
		t.Helper()
		checkFor(t, 3*time.Second, 15*time.Millisecond, func() error {
			c.mu.Lock()
			defer c.mu.Unlock()
			sub := c.subs["1"]
			if ls := len(sub.shadow); ls != expected {
				return fmt.Errorf("Expected shadows to be %d, got %d", expected, ls)
			}
			return nil
		})
	}

	// We created a SUB on foo which should create a shadow subscription.
	checkShadow(t, 1)

	time.Sleep(1250 * time.Millisecond)

	// Should have expired and been removed.
	checkShadow(t, 0)
}

func TestJWTAccountLimitsSubs(t *testing.T) {
	fooAC := newJWTTestAccountClaims()
	fooAC.Limits.Subs = 10
	s, fooKP, c, _ := setupJWTTestWitAccountClaims(t, fooAC, "+OK")
	defer s.Shutdown()
	defer c.close()

	okp, _ := nkeys.FromSeed(oSeed)
	fooPub, _ := fooKP.PublicKey()

	// Create a client.
	c, cr, cs := createClient(t, s, fooKP)
	defer c.close()

	c.parseAsync(cs)
	expectPong(t, cr)

	// Check to make sure we have the limit set.
	// Account first
	fooAcc, _ := s.LookupAccount(fooPub)
	fooAcc.mu.RLock()
	if fooAcc.msubs != 10 {
		fooAcc.mu.RUnlock()
		t.Fatalf("Expected account to have msubs of 10, got %d", fooAcc.msubs)
	}
	fooAcc.mu.RUnlock()
	// Now test that the client has limits too.
	c.mu.Lock()
	if c.msubs != 10 {
		c.mu.Unlock()
		t.Fatalf("Expected client msubs to be 10, got %d", c.msubs)
	}
	c.mu.Unlock()

	// Now make sure its enforced.
	/// These should all work ok.
	for i := 0; i < 10; i++ {
		c.parseAsync(fmt.Sprintf("SUB foo %d\r\nPING\r\n", i))
		expectPong(t, cr)
	}

	// This one should fail.
	c.parseAsync("SUB foo 22\r\n")
	l, _ := cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR") {
		t.Fatalf("Expected an ERR, got: %v", l)
	}
	if !strings.Contains(l, "maximum subscriptions exceeded") {
		t.Fatalf("Expected an ERR for max subscriptions exceeded, got: %v", l)
	}

	// Now update the claims and expect if max is lower to be disconnected.
	fooAC.Limits.Subs = 5
	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, fooPub, fooJWT)
	s.updateAccountClaims(fooAcc, fooAC)
	l, _ = cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR") {
		t.Fatalf("Expected an ERR, got: %v", l)
	}
	if !strings.Contains(l, "maximum subscriptions exceeded") {
		t.Fatalf("Expected an ERR for max subscriptions exceeded, got: %v", l)
	}
}

func TestJWTAccountLimitsSubsButServerOverrides(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	// override with server setting of 2.
	opts := s.getOpts()
	opts.MaxSubs = 2

	okp, _ := nkeys.FromSeed(oSeed)

	// Create accounts and imports/exports.
	fooKP, _ := nkeys.CreateAccount()
	fooPub, _ := fooKP.PublicKey()
	fooAC := jwt.NewAccountClaims(fooPub)
	fooAC.Limits.Subs = 10
	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, fooPub, fooJWT)
	fooAcc, _ := s.LookupAccount(fooPub)
	fooAcc.mu.RLock()
	if fooAcc.msubs != 10 {
		fooAcc.mu.RUnlock()
		t.Fatalf("Expected account to have msubs of 10, got %d", fooAcc.msubs)
	}
	fooAcc.mu.RUnlock()

	// Create a client.
	c, cr, cs := createClient(t, s, fooKP)
	defer c.close()

	c.parseAsync(cs)
	expectPong(t, cr)

	c.parseAsync("SUB foo 1\r\nSUB bar 2\r\nSUB baz 3\r\nPING\r\n")
	l, _ := cr.ReadString('\n')

	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}
	if !strings.Contains(l, "maximum subscriptions exceeded") {
		t.Fatalf("Expected an ERR for max subscriptions exceeded, got: %v", l)
	}
	// Read last PONG so does not hold up test.
	cr.ReadString('\n')
}

func TestJWTAccountLimitsMaxPayload(t *testing.T) {
	fooAC := newJWTTestAccountClaims()
	fooAC.Limits.Payload = 8
	s, fooKP, c, _ := setupJWTTestWitAccountClaims(t, fooAC, "+OK")
	defer s.Shutdown()
	defer c.close()

	fooPub, _ := fooKP.PublicKey()

	// Create a client.
	c, cr, cs := createClient(t, s, fooKP)
	defer c.close()

	c.parseAsync(cs)
	expectPong(t, cr)

	// Check to make sure we have the limit set.
	// Account first
	fooAcc, _ := s.LookupAccount(fooPub)
	fooAcc.mu.RLock()
	if fooAcc.mpay != 8 {
		fooAcc.mu.RUnlock()
		t.Fatalf("Expected account to have mpay of 8, got %d", fooAcc.mpay)
	}
	fooAcc.mu.RUnlock()
	// Now test that the client has limits too.
	c.mu.Lock()
	if c.mpay != 8 {
		c.mu.Unlock()
		t.Fatalf("Expected client to have mpay of 10, got %d", c.mpay)
	}
	c.mu.Unlock()

	c.parseAsync("PUB foo 4\r\nXXXX\r\nPING\r\n")
	expectPong(t, cr)

	c.parseAsync("PUB foo 10\r\nXXXXXXXXXX\r\nPING\r\n")
	l, _ := cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}
	if !strings.Contains(l, "Maximum Payload") {
		t.Fatalf("Expected an ERR for max payload violation, got: %v", l)
	}
}

func TestJWTAccountLimitsMaxPayloadButServerOverrides(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	// override with server setting of 4.
	opts := s.getOpts()
	opts.MaxPayload = 4

	okp, _ := nkeys.FromSeed(oSeed)

	// Create accounts and imports/exports.
	fooKP, _ := nkeys.CreateAccount()
	fooPub, _ := fooKP.PublicKey()
	fooAC := jwt.NewAccountClaims(fooPub)
	fooAC.Limits.Payload = 8
	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, fooPub, fooJWT)

	// Create a client.
	c, cr, cs := createClient(t, s, fooKP)
	defer c.close()

	c.parseAsync(cs)
	expectPong(t, cr)

	c.parseAsync("PUB foo 6\r\nXXXXXX\r\nPING\r\n")
	l, _ := cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}
	if !strings.Contains(l, "Maximum Payload") {
		t.Fatalf("Expected an ERR for max payload violation, got: %v", l)
	}
}

func TestJWTAccountLimitsMaxConns(t *testing.T) {
	fooAC := newJWTTestAccountClaims()
	fooAC.Limits.Conn = 8
	s, fooKP, c, _ := setupJWTTestWitAccountClaims(t, fooAC, "+OK")
	defer s.Shutdown()
	defer c.close()

	newClient := func(expPre string) *testAsyncClient {
		t.Helper()
		// Create a client.
		c, cr, cs := createClient(t, s, fooKP)
		c.parseAsync(cs)
		l, _ := cr.ReadString('\n')
		if !strings.HasPrefix(l, expPre) {
			t.Fatalf("Expected a response starting with %q, got %q", expPre, l)
		}
		return c
	}

	// A connection is created in setupJWTTestWitAccountClaims(), so limit
	// to 7 here (8 total).
	for i := 0; i < 7; i++ {
		c := newClient("PONG")
		defer c.close()
	}
	// Now this one should fail.
	c = newClient("-ERR ")
	c.close()
}

// This will test that we can switch from a public export to a private
// one and back with export claims to make sure the claim update mechanism
// is working properly.
func TestJWTAccountServiceImportAuthSwitch(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Create accounts and imports/exports.
	fooKP, _ := nkeys.CreateAccount()
	fooPub, _ := fooKP.PublicKey()
	fooAC := jwt.NewAccountClaims(fooPub)
	serviceExport := &jwt.Export{Subject: "ngs.usage.*", Type: jwt.Service}
	fooAC.Exports.Add(serviceExport)
	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, fooPub, fooJWT)

	barKP, _ := nkeys.CreateAccount()
	barPub, _ := barKP.PublicKey()
	barAC := jwt.NewAccountClaims(barPub)
	serviceImport := &jwt.Import{Account: fooPub, Subject: "ngs.usage", To: "ngs.usage.DEREK", Type: jwt.Service}
	barAC.Imports.Add(serviceImport)
	barJWT, err := barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)

	// Create a client that will send the request
	ca, cra, csa := createClient(t, s, barKP)
	defer ca.close()
	ca.parseAsync(csa)
	expectPong(t, cra)

	// Create the client that will respond to the requests.
	cb, crb, csb := createClient(t, s, fooKP)
	defer cb.close()
	cb.parseAsync(csb)
	expectPong(t, crb)

	// Create Subscriber.
	cb.parseAsync("SUB ngs.usage.* 1\r\nPING\r\n")
	expectPong(t, crb)

	// Send Request
	ca.parseAsync("PUB ngs.usage 2\r\nhi\r\nPING\r\n")
	expectPong(t, cra)

	// We should receive the request mapped into our account. PING needed to flush.
	cb.parseAsync("PING\r\n")
	expectMsg(t, crb, "ngs.usage.DEREK", "hi")

	// Now update to make the export private.
	fooACPrivate := jwt.NewAccountClaims(fooPub)
	serviceExport = &jwt.Export{Subject: "ngs.usage.*", Type: jwt.Service, TokenReq: true}
	fooACPrivate.Exports.Add(serviceExport)
	fooJWTPrivate, err := fooACPrivate.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, fooPub, fooJWTPrivate)
	acc, _ := s.LookupAccount(fooPub)
	s.updateAccountClaims(acc, fooACPrivate)

	// Send Another Request
	ca.parseAsync("PUB ngs.usage 2\r\nhi\r\nPING\r\n")
	expectPong(t, cra)

	// We should not receive the request this time.
	cb.parseAsync("PING\r\n")
	expectPong(t, crb)

	// Now put it back again to public and make sure it works again.
	addAccountToMemResolver(s, fooPub, fooJWT)
	s.updateAccountClaims(acc, fooAC)

	// Send Request
	ca.parseAsync("PUB ngs.usage 2\r\nhi\r\nPING\r\n")
	expectPong(t, cra)

	// We should receive the request mapped into our account. PING needed to flush.
	cb.parseAsync("PING\r\n")
	expectMsg(t, crb, "ngs.usage.DEREK", "hi")
}

func TestJWTAccountServiceImportExpires(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Create accounts and imports/exports.
	fooKP, _ := nkeys.CreateAccount()
	fooPub, _ := fooKP.PublicKey()
	fooAC := jwt.NewAccountClaims(fooPub)
	serviceExport := &jwt.Export{Subject: "foo", Type: jwt.Service}

	fooAC.Exports.Add(serviceExport)
	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, fooPub, fooJWT)

	barKP, _ := nkeys.CreateAccount()
	barPub, _ := barKP.PublicKey()
	barAC := jwt.NewAccountClaims(barPub)
	serviceImport := &jwt.Import{Account: fooPub, Subject: "foo", Type: jwt.Service}

	barAC.Imports.Add(serviceImport)
	barJWT, err := barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)

	// Create a client that will send the request
	ca, cra, csa := createClient(t, s, barKP)
	defer ca.close()
	ca.parseAsync(csa)
	expectPong(t, cra)

	// Create the client that will respond to the requests.
	cb, crb, csb := createClient(t, s, fooKP)
	defer cb.close()
	cb.parseAsync(csb)
	expectPong(t, crb)

	// Create Subscriber.
	cb.parseAsync("SUB foo 1\r\nPING\r\n")
	expectPong(t, crb)

	// Send Request
	ca.parseAsync("PUB foo 2\r\nhi\r\nPING\r\n")
	expectPong(t, cra)

	// We should receive the request. PING needed to flush.
	cb.parseAsync("PING\r\n")
	expectMsg(t, crb, "foo", "hi")

	// Now update the exported service to require auth.
	fooAC = jwt.NewAccountClaims(fooPub)
	serviceExport = &jwt.Export{Subject: "foo", Type: jwt.Service, TokenReq: true}

	fooAC.Exports.Add(serviceExport)
	fooJWT, err = fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, fooPub, fooJWT)
	acc, _ := s.LookupAccount(fooPub)
	s.updateAccountClaims(acc, fooAC)

	// Send Another Request
	ca.parseAsync("PUB foo 2\r\nhi\r\nPING\r\n")
	expectPong(t, cra)

	// We should not receive the request this time.
	cb.parseAsync("PING\r\n")
	expectPong(t, crb)

	// Now get an activation token such that it will work, but will expire.
	barAC = jwt.NewAccountClaims(barPub)
	serviceImport = &jwt.Import{Account: fooPub, Subject: "foo", Type: jwt.Service}

	now := time.Now()
	activation := jwt.NewActivationClaims(barPub)
	activation.ImportSubject = "foo"
	activation.ImportType = jwt.Service
	activation.IssuedAt = now.Add(-10 * time.Second).Unix()
	activation.Expires = now.Add(time.Second).Round(time.Second).Unix()
	actJWT, err := activation.Encode(fooKP)
	if err != nil {
		t.Fatalf("Error generating activation token: %v", err)
	}
	serviceImport.Token = actJWT

	barAC.Imports.Add(serviceImport)
	barJWT, err = barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)
	acc, _ = s.LookupAccount(barPub)
	s.updateAccountClaims(acc, barAC)

	// Now it should work again.
	// Send Another Request
	ca.parseAsync("PUB foo 3\r\nhi2\r\nPING\r\n")
	expectPong(t, cra)

	// We should receive the request. PING needed to flush.
	cb.parseAsync("PING\r\n")
	expectMsg(t, crb, "foo", "hi2")

	// Now wait for it to expire, then retry.
	waitTime := time.Duration(activation.Expires-time.Now().Unix()) * time.Second
	time.Sleep(waitTime + 250*time.Millisecond)

	// Send Another Request
	ca.parseAsync("PUB foo 3\r\nhi3\r\nPING\r\n")
	expectPong(t, cra)

	// We should NOT receive the request. PING needed to flush.
	cb.parseAsync("PING\r\n")
	expectPong(t, crb)
}

func TestAccountURLResolver(t *testing.T) {
	for _, test := range []struct {
		name   string
		useTLS bool
	}{
		{"plain", false},
		{"tls", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			kp, _ := nkeys.FromSeed(oSeed)
			akp, _ := nkeys.CreateAccount()
			apub, _ := akp.PublicKey()
			nac := jwt.NewAccountClaims(apub)
			ajwt, err := nac.Encode(kp)
			if err != nil {
				t.Fatalf("Error generating account JWT: %v", err)
			}

			hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(ajwt))
			})
			var ts *httptest.Server
			if test.useTLS {
				ts = httptest.NewTLSServer(hf)
			} else {
				ts = httptest.NewServer(hf)
			}
			defer ts.Close()

			confTemplate := `
				listen: -1
				resolver: URL("%s/ngs/v1/accounts/jwt/")
				resolver_tls {
					insecure: true
				}
			`
			conf := createConfFile(t, []byte(fmt.Sprintf(confTemplate, ts.URL)))
			defer os.Remove(conf)

			s, opts := RunServerWithConfig(conf)
			pub, _ := kp.PublicKey()
			opts.TrustedKeys = []string{pub}
			defer s.Shutdown()

			acc, _ := s.LookupAccount(apub)
			if acc == nil {
				t.Fatalf("Expected to receive an account")
			}
			if acc.Name != apub {
				t.Fatalf("Account name did not match claim key")
			}
		})
	}
}

func TestAccountURLResolverTimeout(t *testing.T) {
	kp, _ := nkeys.FromSeed(oSeed)
	akp, _ := nkeys.CreateAccount()
	apub, _ := akp.PublicKey()
	nac := jwt.NewAccountClaims(apub)
	ajwt, err := nac.Encode(kp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	basePath := "/ngs/v1/accounts/jwt/"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == basePath {
			w.Write([]byte("ok"))
			return
		}
		// Purposely be slow on account lookup.
		time.Sleep(200 * time.Millisecond)
		w.Write([]byte(ajwt))
	}))
	defer ts.Close()

	confTemplate := `
		listen: -1
		resolver: URL("%s%s")
    `
	conf := createConfFile(t, []byte(fmt.Sprintf(confTemplate, ts.URL, basePath)))
	defer os.Remove(conf)

	s, opts := RunServerWithConfig(conf)
	pub, _ := kp.PublicKey()
	opts.TrustedKeys = []string{pub}
	defer s.Shutdown()

	// Lower default timeout to speed-up test
	s.AccountResolver().(*URLAccResolver).c.Timeout = 50 * time.Millisecond

	acc, _ := s.LookupAccount(apub)
	if acc != nil {
		t.Fatalf("Expected to not receive an account due to timeout")
	}
}

func TestAccountURLResolverNoFetchOnReload(t *testing.T) {
	kp, _ := nkeys.FromSeed(oSeed)
	akp, _ := nkeys.CreateAccount()
	apub, _ := akp.PublicKey()
	nac := jwt.NewAccountClaims(apub)
	ajwt, err := nac.Encode(kp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(ajwt))
	}))
	defer ts.Close()

	confTemplate := `
		listen: -1
		resolver: URL("%s/ngs/v1/accounts/jwt/")
    `
	conf := createConfFile(t, []byte(fmt.Sprintf(confTemplate, ts.URL)))
	defer os.Remove(conf)

	s, _ := RunServerWithConfig(conf)
	defer s.Shutdown()

	acc, _ := s.LookupAccount(apub)
	if acc == nil {
		t.Fatalf("Expected to receive an account")
	}

	// Reload would produce a DATA race during the DeepEqual check for the account resolver,
	// so close the current one and we will create a new one that keeps track of fetch calls.
	ts.Close()

	fetch := int32(0)
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&fetch, 1)
		w.Write([]byte(ajwt))
	}))
	defer ts.Close()

	changeCurrentConfigContentWithNewContent(t, conf, []byte(fmt.Sprintf(confTemplate, ts.URL)))

	if err := s.Reload(); err != nil {
		t.Fatalf("Error on reload: %v", err)
	}
	if atomic.LoadInt32(&fetch) != 0 {
		t.Fatalf("Fetch invoked during reload")
	}

	// Now stop the resolver and make sure that on startup, we report URL resolver failure
	s.Shutdown()
	s = nil
	ts.Close()

	opts := LoadConfig(conf)
	if s, err := NewServer(opts); err == nil || !strings.Contains(err.Error(), "could not fetch") {
		if s != nil {
			s.Shutdown()
		}
		t.Fatalf("Expected error regarding account resolver, got %v", err)
	}
}

func TestJWTUserSigningKey(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()

	// Check to make sure we would have an authTimer
	if !s.info.AuthRequired {
		t.Fatalf("Expect the server to require auth")
	}

	c, cr, _ := newClientForServer(s)
	defer c.close()
	// Don't send jwt field, should fail.
	c.parseAsync("CONNECT {\"verbose\":true,\"pedantic\":true}\r\nPING\r\n")
	l, _ := cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}

	okp, _ := nkeys.FromSeed(oSeed)

	// Create an account
	akp, _ := nkeys.CreateAccount()
	apub, _ := akp.PublicKey()

	// Create a signing key for the account
	askp, _ := nkeys.CreateAccount()
	aspub, _ := askp.PublicKey()

	nac := jwt.NewAccountClaims(apub)
	ajwt, err := nac.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	// Create a client with the account signing key
	c, cr, cs := createClientWithIssuer(t, s, askp, apub)
	defer c.close()

	// PING needed to flush the +OK/-ERR to us.
	// This should fail too since no account resolver is defined.
	c.parseAsync(cs)
	l, _ = cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}

	// Ok now let's walk through and make sure all is good.
	// We will set the account resolver by hand to a memory resolver.
	buildMemAccResolver(s)
	addAccountToMemResolver(s, apub, ajwt)

	// Create a client with a signing key
	c, cr, cs = createClientWithIssuer(t, s, askp, apub)
	defer c.close()
	// should fail because the signing key is not known
	c.parseAsync(cs)
	l, _ = cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error: %v", l)
	}

	// add a signing key
	nac.SigningKeys.Add(aspub)
	// update the memory resolver
	acc, _ := s.LookupAccount(apub)
	s.updateAccountClaims(acc, nac)

	// Create a client with a signing key
	c, cr, cs = createClientWithIssuer(t, s, askp, apub)
	defer c.close()

	// expect this to work
	c.parseAsync(cs)
	l, _ = cr.ReadString('\n')
	if !strings.HasPrefix(l, "PONG") {
		t.Fatalf("Expected a PONG, got %q", l)
	}

	isClosed := func() bool {
		c.mu.Lock()
		defer c.mu.Unlock()
		return c.isClosed()
	}

	if isClosed() {
		t.Fatal("expected client to be alive")
	}
	// remove the signing key should bounce client
	nac.SigningKeys = nil
	acc, _ = s.LookupAccount(apub)
	s.updateAccountClaims(acc, nac)

	if !isClosed() {
		t.Fatal("expected client to be gone")
	}
}

func TestJWTAccountImportSignerRemoved(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Exporter keys
	srvKP, _ := nkeys.CreateAccount()
	srvPK, _ := srvKP.PublicKey()
	srvSignerKP, _ := nkeys.CreateAccount()
	srvSignerPK, _ := srvSignerKP.PublicKey()

	// Importer keys
	clientKP, _ := nkeys.CreateAccount()
	clientPK, _ := clientKP.PublicKey()

	createSrvJwt := func(signingKeys ...string) (string, *jwt.AccountClaims) {
		ac := jwt.NewAccountClaims(srvPK)
		ac.SigningKeys.Add(signingKeys...)
		ac.Exports.Add(&jwt.Export{Subject: "foo", Type: jwt.Service, TokenReq: true})
		ac.Exports.Add(&jwt.Export{Subject: "bar", Type: jwt.Stream, TokenReq: true})
		token, err := ac.Encode(okp)
		if err != nil {
			t.Fatalf("Error generating exporter JWT: %v", err)
		}
		return token, ac
	}

	createImportToken := func(sub string, kind jwt.ExportType) string {
		actC := jwt.NewActivationClaims(clientPK)
		actC.IssuerAccount = srvPK
		actC.ImportType = kind
		actC.ImportSubject = jwt.Subject(sub)
		token, err := actC.Encode(srvSignerKP)
		if err != nil {
			t.Fatal(err)
		}
		return token
	}

	createClientJwt := func() string {
		ac := jwt.NewAccountClaims(clientPK)
		ac.Imports.Add(&jwt.Import{Account: srvPK, Subject: "foo", Type: jwt.Service, Token: createImportToken("foo", jwt.Service)})
		ac.Imports.Add(&jwt.Import{Account: srvPK, Subject: "bar", Type: jwt.Stream, Token: createImportToken("bar", jwt.Stream)})
		token, err := ac.Encode(okp)
		if err != nil {
			t.Fatalf("Error generating importer JWT: %v", err)
		}
		return token
	}

	srvJWT, _ := createSrvJwt(srvSignerPK)
	addAccountToMemResolver(s, srvPK, srvJWT)

	clientJWT := createClientJwt()
	addAccountToMemResolver(s, clientPK, clientJWT)

	// Create a client that will send the request
	client, clientReader, clientCS := createClient(t, s, clientKP)
	defer client.close()
	client.parseAsync(clientCS)
	expectPong(t, clientReader)

	checkShadow := func(expected int) {
		t.Helper()
		client.mu.Lock()
		defer client.mu.Unlock()
		sub := client.subs["1"]
		count := 0
		if sub != nil {
			count = len(sub.shadow)
		}
		if count != expected {
			t.Fatalf("Expected shadows to be %d, got %d", expected, count)
		}
	}

	checkShadow(0)
	// Create the client that will respond to the requests.
	srv, srvReader, srvCS := createClient(t, s, srvKP)
	defer srv.close()
	srv.parseAsync(srvCS)
	expectPong(t, srvReader)

	// Create Subscriber.
	srv.parseAsync("SUB foo 1\r\nPING\r\n")
	expectPong(t, srvReader)

	// Send Request
	client.parseAsync("PUB foo 2\r\nhi\r\nPING\r\n")
	expectPong(t, clientReader)

	// We should receive the request. PING needed to flush.
	srv.parseAsync("PING\r\n")
	expectMsg(t, srvReader, "foo", "hi")

	client.parseAsync("SUB bar 1\r\nPING\r\n")
	expectPong(t, clientReader)
	checkShadow(1)

	srv.parseAsync("PUB bar 2\r\nhi\r\nPING\r\n")
	expectPong(t, srvReader)

	// We should receive from stream. PING needed to flush.
	client.parseAsync("PING\r\n")
	expectMsg(t, clientReader, "bar", "hi")

	// Now update the exported service no signer
	srvJWT, srvAC := createSrvJwt()
	addAccountToMemResolver(s, srvPK, srvJWT)
	acc, _ := s.LookupAccount(srvPK)
	s.updateAccountClaims(acc, srvAC)

	// Send Another Request
	client.parseAsync("PUB foo 2\r\nhi\r\nPING\r\n")
	expectPong(t, clientReader)

	// We should not receive the request this time.
	srv.parseAsync("PING\r\n")
	expectPong(t, srvReader)

	// Publish on the stream
	srv.parseAsync("PUB bar 2\r\nhi\r\nPING\r\n")
	expectPong(t, srvReader)

	// We should not receive from the stream this time
	client.parseAsync("PING\r\n")
	expectPong(t, clientReader)
	checkShadow(0)
}

func TestJWTAccountImportSignerDeadlock(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Exporter keys
	srvKP, _ := nkeys.CreateAccount()
	srvPK, _ := srvKP.PublicKey()
	srvSignerKP, _ := nkeys.CreateAccount()
	srvSignerPK, _ := srvSignerKP.PublicKey()

	// Importer keys
	clientKP, _ := nkeys.CreateAccount()
	clientPK, _ := clientKP.PublicKey()

	createSrvJwt := func(signingKeys ...string) (string, *jwt.AccountClaims) {
		ac := jwt.NewAccountClaims(srvPK)
		ac.SigningKeys.Add(signingKeys...)
		ac.Exports.Add(&jwt.Export{Subject: "foo", Type: jwt.Service, TokenReq: true})
		ac.Exports.Add(&jwt.Export{Subject: "bar", Type: jwt.Stream, TokenReq: true})
		token, err := ac.Encode(okp)
		if err != nil {
			t.Fatalf("Error generating exporter JWT: %v", err)
		}
		return token, ac
	}

	createImportToken := func(sub string, kind jwt.ExportType) string {
		actC := jwt.NewActivationClaims(clientPK)
		actC.IssuerAccount = srvPK
		actC.ImportType = kind
		actC.ImportSubject = jwt.Subject(sub)
		token, err := actC.Encode(srvSignerKP)
		if err != nil {
			t.Fatal(err)
		}
		return token
	}

	createClientJwt := func() string {
		ac := jwt.NewAccountClaims(clientPK)
		ac.Imports.Add(&jwt.Import{Account: srvPK, Subject: "foo", Type: jwt.Service, Token: createImportToken("foo", jwt.Service)})
		ac.Imports.Add(&jwt.Import{Account: srvPK, Subject: "bar", Type: jwt.Stream, Token: createImportToken("bar", jwt.Stream)})
		token, err := ac.Encode(okp)
		if err != nil {
			t.Fatalf("Error generating importer JWT: %v", err)
		}
		return token
	}

	srvJWT, _ := createSrvJwt(srvSignerPK)
	addAccountToMemResolver(s, srvPK, srvJWT)

	clientJWT := createClientJwt()
	addAccountToMemResolver(s, clientPK, clientJWT)

	acc, _ := s.LookupAccount(srvPK)
	// Have a go routine that constantly gets/releases the acc's write lock.
	// There was a bug that could cause AddServiceImportWithClaim to deadlock.
	ch := make(chan bool, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ch:
				return
			default:
				acc.mu.Lock()
				acc.mu.Unlock()
				time.Sleep(time.Millisecond)
			}
		}
	}()

	// Create a client that will send the request
	client, clientReader, clientCS := createClient(t, s, clientKP)
	defer client.close()
	client.parseAsync(clientCS)
	expectPong(t, clientReader)

	close(ch)
	wg.Wait()
}

func TestJWTAccountImportWrongIssuerAccount(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	l := &captureErrorLogger{errCh: make(chan string, 2)}
	s.SetLogger(l, false, false)

	okp, _ := nkeys.FromSeed(oSeed)

	// Exporter keys
	srvKP, _ := nkeys.CreateAccount()
	srvPK, _ := srvKP.PublicKey()
	srvSignerKP, _ := nkeys.CreateAccount()
	srvSignerPK, _ := srvSignerKP.PublicKey()

	// Importer keys
	clientKP, _ := nkeys.CreateAccount()
	clientPK, _ := clientKP.PublicKey()

	createSrvJwt := func(signingKeys ...string) (string, *jwt.AccountClaims) {
		ac := jwt.NewAccountClaims(srvPK)
		ac.SigningKeys.Add(signingKeys...)
		ac.Exports.Add(&jwt.Export{Subject: "foo", Type: jwt.Service, TokenReq: true})
		ac.Exports.Add(&jwt.Export{Subject: "bar", Type: jwt.Stream, TokenReq: true})
		token, err := ac.Encode(okp)
		if err != nil {
			t.Fatalf("Error generating exporter JWT: %v", err)
		}
		return token, ac
	}

	createImportToken := func(sub string, kind jwt.ExportType) string {
		actC := jwt.NewActivationClaims(clientPK)
		// Reference ourselves, which is wrong.
		actC.IssuerAccount = clientPK
		actC.ImportType = kind
		actC.ImportSubject = jwt.Subject(sub)
		token, err := actC.Encode(srvSignerKP)
		if err != nil {
			t.Fatal(err)
		}
		return token
	}

	createClientJwt := func() string {
		ac := jwt.NewAccountClaims(clientPK)
		ac.Imports.Add(&jwt.Import{Account: srvPK, Subject: "foo", Type: jwt.Service, Token: createImportToken("foo", jwt.Service)})
		ac.Imports.Add(&jwt.Import{Account: srvPK, Subject: "bar", Type: jwt.Stream, Token: createImportToken("bar", jwt.Stream)})
		token, err := ac.Encode(okp)
		if err != nil {
			t.Fatalf("Error generating importer JWT: %v", err)
		}
		return token
	}

	srvJWT, _ := createSrvJwt(srvSignerPK)
	addAccountToMemResolver(s, srvPK, srvJWT)

	clientJWT := createClientJwt()
	addAccountToMemResolver(s, clientPK, clientJWT)

	// Create a client that will send the request
	client, clientReader, clientCS := createClient(t, s, clientKP)
	defer client.close()
	client.parseAsync(clientCS)
	expectPong(t, clientReader)

	for i := 0; i < 2; i++ {
		select {
		case e := <-l.errCh:
			if !strings.HasPrefix(e, fmt.Sprintf("Invalid issuer account %q in activation claim", clientPK)) {
				t.Fatalf("Unexpected error: %v", e)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("Did not get error regarding issuer account")
		}
	}
}

func TestJWTUserRevokedOnAccountUpdate(t *testing.T) {
	nac := newJWTTestAccountClaims()
	s, akp, c, cr := setupJWTTestWitAccountClaims(t, nac, "+OK")
	defer s.Shutdown()
	defer c.close()

	expectPong(t, cr)

	okp, _ := nkeys.FromSeed(oSeed)
	apub, _ := akp.PublicKey()

	c.mu.Lock()
	pub := c.user.Nkey
	c.mu.Unlock()

	// Now revoke the user.
	nac.Revoke(pub)

	ajwt, err := nac.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	// Update the account on the server.
	addAccountToMemResolver(s, apub, ajwt)
	acc, err := s.LookupAccount(apub)
	if err != nil {
		t.Fatalf("Error looking up the account: %v", err)
	}

	// This is simulating a system update for the account claims.
	go s.updateAccountWithClaimJWT(acc, ajwt)

	l, _ := cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}
	if !strings.Contains(l, "Revoked") {
		t.Fatalf("Expected 'Revoked' to be in the error")
	}
}

func TestJWTUserRevoked(t *testing.T) {
	okp, _ := nkeys.FromSeed(oSeed)

	// Create a new user that we will make sure has been revoked.
	nkp, _ := nkeys.CreateUser()
	pub, _ := nkp.PublicKey()
	nuc := jwt.NewUserClaims(pub)

	akp, _ := nkeys.CreateAccount()
	apub, _ := akp.PublicKey()
	nac := jwt.NewAccountClaims(apub)
	// Revoke the user right away.
	nac.Revoke(pub)
	ajwt, err := nac.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	// Sign for the user.
	jwt, err := nuc.Encode(akp)
	if err != nil {
		t.Fatalf("Error generating user JWT: %v", err)
	}

	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)
	addAccountToMemResolver(s, apub, ajwt)

	c, cr, l := newClientForServer(s)
	defer c.close()

	// Sign Nonce
	var info nonceInfo
	json.Unmarshal([]byte(l[5:]), &info)
	sigraw, _ := nkp.Sign([]byte(info.Nonce))
	sig := base64.RawURLEncoding.EncodeToString(sigraw)

	// PING needed to flush the +OK/-ERR to us.
	cs := fmt.Sprintf("CONNECT {\"jwt\":%q,\"sig\":\"%s\"}\r\nPING\r\n", jwt, sig)

	c.parseAsync(cs)

	l, _ = cr.ReadString('\n')
	if !strings.HasPrefix(l, "-ERR ") {
		t.Fatalf("Expected an error")
	}
	if !strings.Contains(l, "Authorization") {
		t.Fatalf("Expected 'Revoked' to be in the error")
	}
}

// Test that an account update that revokes an import authorization cancels the import.
func TestJWTImportTokenRevokedAfter(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Create accounts and imports/exports.
	fooKP, _ := nkeys.CreateAccount()
	fooPub, _ := fooKP.PublicKey()
	fooAC := jwt.NewAccountClaims(fooPub)

	// Now create Exports.
	export := &jwt.Export{Subject: "foo.private", Type: jwt.Stream, TokenReq: true}

	fooAC.Exports.Add(export)
	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	addAccountToMemResolver(s, fooPub, fooJWT)

	barKP, _ := nkeys.CreateAccount()
	barPub, _ := barKP.PublicKey()
	barAC := jwt.NewAccountClaims(barPub)
	simport := &jwt.Import{Account: fooPub, Subject: "foo.private", Type: jwt.Stream}

	activation := jwt.NewActivationClaims(barPub)
	activation.ImportSubject = "foo.private"
	activation.ImportType = jwt.Stream
	actJWT, err := activation.Encode(fooKP)
	if err != nil {
		t.Fatalf("Error generating activation token: %v", err)
	}

	simport.Token = actJWT
	barAC.Imports.Add(simport)
	barJWT, err := barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)

	// Now revoke the export.
	decoded, _ := jwt.DecodeActivationClaims(actJWT)
	export.Revoke(decoded.Subject)

	fooJWT, err = fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	addAccountToMemResolver(s, fooPub, fooJWT)

	fooAcc, _ := s.LookupAccount(fooPub)
	if fooAcc == nil {
		t.Fatalf("Expected to retrieve the account")
	}

	// Now lookup bar account and make sure it was revoked.
	acc, _ := s.LookupAccount(barPub)
	if acc == nil {
		t.Fatalf("Expected to retrieve the account")
	}
	if les := len(acc.imports.streams); les != 0 {
		t.Fatalf("Expected imports streams len of 0, got %d", les)
	}
}

// Test that an account update that revokes an import authorization cancels the import.
func TestJWTImportTokenRevokedBefore(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Create accounts and imports/exports.
	fooKP, _ := nkeys.CreateAccount()
	fooPub, _ := fooKP.PublicKey()
	fooAC := jwt.NewAccountClaims(fooPub)

	// Now create Exports.
	export := &jwt.Export{Subject: "foo.private", Type: jwt.Stream, TokenReq: true}

	fooAC.Exports.Add(export)

	// Import account
	barKP, _ := nkeys.CreateAccount()
	barPub, _ := barKP.PublicKey()
	barAC := jwt.NewAccountClaims(barPub)
	simport := &jwt.Import{Account: fooPub, Subject: "foo.private", Type: jwt.Stream}

	activation := jwt.NewActivationClaims(barPub)
	activation.ImportSubject = "foo.private"
	activation.ImportType = jwt.Stream
	actJWT, err := activation.Encode(fooKP)
	if err != nil {
		t.Fatalf("Error generating activation token: %v", err)
	}

	simport.Token = actJWT
	barAC.Imports.Add(simport)

	// Now revoke the export.
	decoded, _ := jwt.DecodeActivationClaims(actJWT)
	export.Revoke(decoded.Subject)

	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	addAccountToMemResolver(s, fooPub, fooJWT)

	barJWT, err := barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, barPub, barJWT)

	fooAcc, _ := s.LookupAccount(fooPub)
	if fooAcc == nil {
		t.Fatalf("Expected to retrieve the account")
	}

	// Now lookup bar account and make sure it was revoked.
	acc, _ := s.LookupAccount(barPub)
	if acc == nil {
		t.Fatalf("Expected to retrieve the account")
	}
	if les := len(acc.imports.streams); les != 0 {
		t.Fatalf("Expected imports streams len of 0, got %d", les)
	}
}

func TestJWTCircularAccountServiceImport(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Create accounts
	fooKP, _ := nkeys.CreateAccount()
	fooPub, _ := fooKP.PublicKey()
	fooAC := jwt.NewAccountClaims(fooPub)

	barKP, _ := nkeys.CreateAccount()
	barPub, _ := barKP.PublicKey()
	barAC := jwt.NewAccountClaims(barPub)

	// Create service export/import for account foo
	serviceExport := &jwt.Export{Subject: "foo", Type: jwt.Service, TokenReq: true}
	serviceImport := &jwt.Import{Account: barPub, Subject: "bar", Type: jwt.Service}

	fooAC.Exports.Add(serviceExport)
	fooAC.Imports.Add(serviceImport)
	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	addAccountToMemResolver(s, fooPub, fooJWT)

	// Create service export/import for account bar
	serviceExport = &jwt.Export{Subject: "bar", Type: jwt.Service, TokenReq: true}
	serviceImport = &jwt.Import{Account: fooPub, Subject: "foo", Type: jwt.Service}

	barAC.Exports.Add(serviceExport)
	barAC.Imports.Add(serviceImport)
	barJWT, err := barAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}

	addAccountToMemResolver(s, barPub, barJWT)

	c, cr, cs := createClient(t, s, fooKP)
	defer c.close()

	c.parseAsync(cs)
	expectPong(t, cr)

	c.parseAsync("SUB foo 1\r\nPING\r\n")
	expectPong(t, cr)
}

// This test ensures that connected clients are properly evicted
// (no deadlock) if the max conns of an account has been lowered
// and the account is being updated (following expiration during
// a lookup).
func TestJWTAccountLimitsMaxConnsAfterExpired(t *testing.T) {
	s := opTrustBasicSetup()
	defer s.Shutdown()
	buildMemAccResolver(s)

	okp, _ := nkeys.FromSeed(oSeed)

	// Create accounts and imports/exports.
	fooKP, _ := nkeys.CreateAccount()
	fooPub, _ := fooKP.PublicKey()
	fooAC := jwt.NewAccountClaims(fooPub)
	fooAC.Limits.Conn = 10
	fooJWT, err := fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, fooPub, fooJWT)

	newClient := func(expPre string) *testAsyncClient {
		t.Helper()
		// Create a client.
		c, cr, cs := createClient(t, s, fooKP)
		c.parseAsync(cs)
		l, _ := cr.ReadString('\n')
		if !strings.HasPrefix(l, expPre) {
			t.Fatalf("Expected a response starting with %q, got %q", expPre, l)
		}
		go func() {
			for {
				if _, _, err := cr.ReadLine(); err != nil {
					return
				}
			}
		}()
		return c
	}

	for i := 0; i < 4; i++ {
		c := newClient("PONG")
		defer c.close()
	}

	// We will simulate that the account has expired. When
	// a new client will connect, the server will do a lookup
	// and find the account expired, which then will cause
	// a fetch and a rebuild of the account. Since max conns
	// is now lower, some clients should have been removed.
	acc, _ := s.LookupAccount(fooPub)
	acc.mu.Lock()
	acc.expired = true
	acc.mu.Unlock()

	// Now update with new expiration and max connections lowered to 2
	fooAC.Limits.Conn = 2
	fooJWT, err = fooAC.Encode(okp)
	if err != nil {
		t.Fatalf("Error generating account JWT: %v", err)
	}
	addAccountToMemResolver(s, fooPub, fooJWT)

	// Cause the lookup that will detect that account was expired
	// and rebuild it, and kick clients out.
	c := newClient("-ERR ")
	defer c.close()

	acc, _ = s.LookupAccount(fooPub)
	checkFor(t, 2*time.Second, 15*time.Millisecond, func() error {
		acc.mu.RLock()
		numClients := len(acc.clients)
		acc.mu.RUnlock()
		if numClients != 2 {
			return fmt.Errorf("Should have 2 clients, got %v", numClients)
		}
		return nil
	})
}
