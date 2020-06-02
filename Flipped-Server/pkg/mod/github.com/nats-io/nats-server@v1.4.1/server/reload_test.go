// Copyright 2017-2018 The NATS Authors
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
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/go-nats"
)

func newServerWithConfig(t *testing.T, configFile string) (*Server, *Options, string) {
	t.Helper()
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Error loading file: %v", err)
	}
	return newServerWithContent(t, content)
}

func newServerWithContent(t *testing.T, content []byte) (*Server, *Options, string) {
	t.Helper()
	opts, tmpFile := newOptionsFromContent(t, content)
	return New(opts), opts, tmpFile
}

func newOptionsFromContent(t *testing.T, content []byte) (*Options, string) {
	t.Helper()
	tmpFile := createConfFile(t, content)
	opts, err := ProcessConfigFile(tmpFile)
	if err != nil {
		t.Fatalf("Error processing config file: %v", err)
	}
	opts.NoSigs = true
	return opts, tmpFile
}

func createConfFile(t *testing.T, content []byte) string {
	t.Helper()
	conf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Error creating conf file: %v", err)
	}
	fName := conf.Name()
	conf.Close()
	if err := ioutil.WriteFile(fName, content, 0666); err != nil {
		os.Remove(fName)
		t.Fatalf("Error writing conf file: %v", err)
	}
	return fName
}

func runReloadServerWithConfig(t *testing.T, configFile string) (*Server, *Options, string) {
	t.Helper()
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Error loading file: %v", err)
	}
	return runReloadServerWithContent(t, content)
}

func runReloadServerWithContent(t *testing.T, content []byte) (*Server, *Options, string) {
	t.Helper()
	opts, tmpFile := newOptionsFromContent(t, content)
	opts.NoLog = true
	opts.NoSigs = true
	s := RunServer(opts)
	return s, opts, tmpFile
}

func changeCurrentConfigContent(t *testing.T, curConfig, newConfig string) {
	t.Helper()
	content, err := ioutil.ReadFile(newConfig)
	if err != nil {
		t.Fatalf("Error loading file: %v", err)
	}
	changeCurrentConfigContentWithNewContent(t, curConfig, content)
}

func changeCurrentConfigContentWithNewContent(t *testing.T, curConfig string, content []byte) {
	t.Helper()
	if err := ioutil.WriteFile(curConfig, content, 0666); err != nil {
		t.Fatalf("Error writing config: %v", err)
	}
}

// Ensure Reload returns an error when attempting to reload a server that did
// not start with a config file.
func TestConfigReloadNoConfigFile(t *testing.T) {
	server := New(&Options{NoSigs: true})
	loaded := server.ConfigTime()
	if server.Reload() == nil {
		t.Fatal("Expected Reload to return an error")
	}
	if reloaded := server.ConfigTime(); reloaded != loaded {
		t.Fatalf("ConfigTime is incorrect.\nexpected: %s\ngot: %s", loaded, reloaded)
	}
}

// Ensure Reload returns an error when attempting to change an option which
// does not support reloading.
func TestConfigReloadUnsupported(t *testing.T) {
	server, _, config := newServerWithConfig(t, "./configs/reload/test.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	loaded := server.ConfigTime()

	golden := &Options{
		ConfigFile:     config,
		Host:           "0.0.0.0",
		Port:           2233,
		AuthTimeout:    1.0,
		Debug:          false,
		Trace:          false,
		Logtime:        false,
		MaxControlLine: 4096,
		MaxPayload:     1048576,
		MaxConn:        65536,
		PingInterval:   2 * time.Minute,
		MaxPingsOut:    2,
		WriteDeadline:  2 * time.Second,
		Cluster: ClusterOpts{
			Host: "127.0.0.1",
			Port: -1,
		},
		NoSigs: true,
	}
	processOptions(golden)

	checkOptionsEqual(t, golden, server.getOpts())

	// Change config file to bad config.
	changeCurrentConfigContent(t, config, "./configs/reload/reload_unsupported.conf")

	// This should fail because `cluster` host cannot be changed.
	if err := server.Reload(); err == nil {
		t.Fatal("Expected Reload to return an error")
	}

	// Ensure config didn't change.
	checkOptionsEqual(t, golden, server.getOpts())

	if reloaded := server.ConfigTime(); reloaded != loaded {
		t.Fatalf("ConfigTime is incorrect.\nexpected: %s\ngot: %s", loaded, reloaded)
	}
}

// This checks that if we change an option that does not support hot-swapping
// we get an error. Using `listen` for now (test may need to be updated if
// server is changed to support change of listen spec).
func TestConfigReloadUnsupportedHotSwapping(t *testing.T) {
	server, _, config := newServerWithContent(t, []byte("listen: 127.0.0.1:-1"))
	defer os.Remove(config)
	defer server.Shutdown()

	loaded := server.ConfigTime()

	time.Sleep(time.Millisecond)

	// Change config file with unsupported option hot-swap
	changeCurrentConfigContentWithNewContent(t, config, []byte("listen: 127.0.0.1:9999"))

	// This should fail because `listen` host cannot be changed.
	if err := server.Reload(); err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("Expected Reload to return a not supported error, got %v", err)
	}

	if reloaded := server.ConfigTime(); reloaded != loaded {
		t.Fatalf("ConfigTime is incorrect.\nexpected: %s\ngot: %s", loaded, reloaded)
	}
}

// Ensure Reload returns an error when reloading from a bad config file.
func TestConfigReloadInvalidConfig(t *testing.T) {
	server, _, config := newServerWithConfig(t, "./configs/reload/test.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	loaded := server.ConfigTime()

	golden := &Options{
		ConfigFile:     config,
		Host:           "0.0.0.0",
		Port:           2233,
		AuthTimeout:    1.0,
		Debug:          false,
		Trace:          false,
		Logtime:        false,
		MaxControlLine: 4096,
		MaxPayload:     1048576,
		MaxConn:        65536,
		PingInterval:   2 * time.Minute,
		MaxPingsOut:    2,
		WriteDeadline:  2 * time.Second,
		Cluster: ClusterOpts{
			Host: "127.0.0.1",
			Port: -1,
		},
		NoSigs: true,
	}
	processOptions(golden)

	checkOptionsEqual(t, golden, server.getOpts())

	// Change config file to bad config.
	changeCurrentConfigContent(t, config, "./configs/reload/invalid.conf")

	// This should fail because the new config should not parse.
	if err := server.Reload(); err == nil {
		t.Fatal("Expected Reload to return an error")
	}

	// Ensure config didn't change.
	checkOptionsEqual(t, golden, server.getOpts())

	if reloaded := server.ConfigTime(); reloaded != loaded {
		t.Fatalf("ConfigTime is incorrect.\nexpected: %s\ngot: %s", loaded, reloaded)
	}
}

// Ensure Reload returns nil and the config is changed on success.
func TestConfigReload(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/test.conf")
	defer os.Remove(config)
	defer os.Remove("gnatsd.pid")
	defer os.Remove("gnatsd.log")
	defer server.Shutdown()

	dir := filepath.Dir(config)
	var content []byte
	if runtime.GOOS != "windows" {
		content = []byte(`
			remote_syslog: "udp://127.0.0.1:514" # change on reload
			syslog:        true # enable on reload
		`)
	}
	platformConf := filepath.Join(dir, "platform.conf")
	defer os.Remove(platformConf)
	if err := ioutil.WriteFile(platformConf, content, 0666); err != nil {
		t.Fatalf("Unable to write config file: %v", err)
	}

	loaded := server.ConfigTime()

	golden := &Options{
		ConfigFile:     config,
		Host:           "0.0.0.0",
		Port:           2233,
		AuthTimeout:    1.0,
		Debug:          false,
		Trace:          false,
		NoLog:          true,
		Logtime:        false,
		MaxControlLine: 4096,
		MaxPayload:     1048576,
		MaxConn:        65536,
		PingInterval:   2 * time.Minute,
		MaxPingsOut:    2,
		WriteDeadline:  2 * time.Second,
		Cluster: ClusterOpts{
			Host: "127.0.0.1",
			Port: server.ClusterAddr().Port,
		},
		NoSigs: true,
	}
	processOptions(golden)

	checkOptionsEqual(t, golden, opts)

	// Change config file to new config.
	changeCurrentConfigContent(t, config, "./configs/reload/reload.conf")

	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure config changed.
	updated := server.getOpts()
	if !updated.Trace {
		t.Fatal("Expected Trace to be true")
	}
	if !updated.Debug {
		t.Fatal("Expected Debug to be true")
	}
	if !updated.Logtime {
		t.Fatal("Expected Logtime to be true")
	}
	if runtime.GOOS != "windows" {
		if !updated.Syslog {
			t.Fatal("Expected Syslog to be true")
		}
		if updated.RemoteSyslog != "udp://127.0.0.1:514" {
			t.Fatalf("RemoteSyslog is incorrect.\nexpected: udp://127.0.0.1:514\ngot: %s", updated.RemoteSyslog)
		}
	}
	if updated.LogFile != "gnatsd.log" {
		t.Fatalf("LogFile is incorrect.\nexpected: gnatsd.log\ngot: %s", updated.LogFile)
	}
	if updated.TLSConfig == nil {
		t.Fatal("Expected TLSConfig to be non-nil")
	}
	if !server.info.TLSRequired {
		t.Fatal("Expected TLSRequired to be true")
	}
	if !server.info.TLSVerify {
		t.Fatal("Expected TLSVerify to be true")
	}
	if updated.Username != "tyler" {
		t.Fatalf("Username is incorrect.\nexpected: tyler\ngot: %s", updated.Username)
	}
	if updated.Password != "T0pS3cr3t" {
		t.Fatalf("Password is incorrect.\nexpected: T0pS3cr3t\ngot: %s", updated.Password)
	}
	if updated.AuthTimeout != 2 {
		t.Fatalf("AuthTimeout is incorrect.\nexpected: 2\ngot: %f", updated.AuthTimeout)
	}
	if !server.info.AuthRequired {
		t.Fatal("Expected AuthRequired to be true")
	}
	if !updated.Cluster.NoAdvertise {
		t.Fatal("Expected NoAdvertise to be true")
	}
	if updated.PidFile != "gnatsd.pid" {
		t.Fatalf("PidFile is incorrect.\nexpected: gnatsd.pid\ngot: %s", updated.PidFile)
	}
	if updated.MaxControlLine != 512 {
		t.Fatalf("MaxControlLine is incorrect.\nexpected: 512\ngot: %d", updated.MaxControlLine)
	}
	if updated.PingInterval != 5*time.Second {
		t.Fatalf("PingInterval is incorrect.\nexpected 5s\ngot: %s", updated.PingInterval)
	}
	if updated.MaxPingsOut != 1 {
		t.Fatalf("MaxPingsOut is incorrect.\nexpected 1\ngot: %d", updated.MaxPingsOut)
	}
	if updated.WriteDeadline != 3*time.Second {
		t.Fatalf("WriteDeadline is incorrect.\nexpected 3s\ngot: %s", updated.WriteDeadline)
	}
	if updated.MaxPayload != 1024 {
		t.Fatalf("MaxPayload is incorrect.\nexpected 1024\ngot: %d", updated.MaxPayload)
	}

	if reloaded := server.ConfigTime(); !reloaded.After(loaded) {
		t.Fatalf("ConfigTime is incorrect.\nexpected greater than: %s\ngot: %s", loaded, reloaded)
	}
}

// Ensure Reload supports TLS config changes. Test this by starting a server
// with TLS enabled, connect to it to verify, reload config using a different
// key pair and client verification enabled, ensure reconnect fails, then
// ensure reconnect succeeds when the client provides a cert.
func TestConfigReloadRotateTLS(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/tls_test.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Ensure we can connect as a sanity check.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, server.Addr().(*net.TCPAddr).Port)

	nc, err := nats.Connect(addr, nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc.Close()
	sub, err := nc.SubscribeSync("foo")
	if err != nil {
		t.Fatalf("Error subscribing: %v", err)
	}
	defer sub.Unsubscribe()

	// Rotate cert and enable client verification.
	changeCurrentConfigContent(t, config, "./configs/reload/tls_verify_test.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure connecting fails.
	if _, err := nats.Connect(addr, nats.Secure(&tls.Config{InsecureSkipVerify: true})); err == nil {
		t.Fatal("Expected connect to fail")
	}

	// Ensure connecting succeeds when client presents cert.
	cert := nats.ClientCert("./configs/certs/cert.new.pem", "./configs/certs/key.new.pem")
	conn, err := nats.Connect(addr, cert, nats.RootCAs("./configs/certs/cert.new.pem"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	conn.Close()

	// Ensure the original connection can still publish/receive.
	if err := nc.Publish("foo", []byte("hello")); err != nil {
		t.Fatalf("Error publishing: %v", err)
	}
	nc.Flush()
	msg, err := sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving msg: %v", err)
	}
	if string(msg.Data) != "hello" {
		t.Fatalf("Msg is incorrect.\nexpected: %+v\ngot: %+v", []byte("hello"), msg.Data)
	}
}

// Ensure Reload supports enabling TLS. Test this by starting a server without
// TLS enabled, connect to it to verify, reload config with TLS enabled, ensure
// reconnect fails, then ensure reconnect succeeds when using secure.
func TestConfigReloadEnableTLS(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/basic.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Ensure we can connect as a sanity check.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, server.Addr().(*net.TCPAddr).Port)
	nc, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	nc.Close()

	// Enable TLS.
	changeCurrentConfigContent(t, config, "./configs/reload/tls_test.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure connecting is OK (we need to skip server cert verification since
	// the library is not doing that by default now).
	nc, err = nats.Connect(addr, nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	nc.Close()
}

// Ensure Reload supports disabling TLS. Test this by starting a server with
// TLS enabled, connect to it to verify, reload config with TLS disabled,
// ensure reconnect fails, then ensure reconnect succeeds when connecting
// without secure.
func TestConfigReloadDisableTLS(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/tls_test.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Ensure we can connect as a sanity check.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, server.Addr().(*net.TCPAddr).Port)
	nc, err := nats.Connect(addr, nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	nc.Close()

	// Disable TLS.
	changeCurrentConfigContent(t, config, "./configs/reload/basic.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure connecting fails.
	if _, err := nats.Connect(addr, nats.Secure(&tls.Config{InsecureSkipVerify: true})); err == nil {
		t.Fatal("Expected connect to fail")
	}

	// Ensure connecting succeeds when not using secure.
	nc, err = nats.Connect(addr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	nc.Close()
}

// Ensure Reload supports single user authentication config changes. Test this
// by starting a server with authentication enabled, connect to it to verify,
// reload config using a different username/password, ensure reconnect fails,
// then ensure reconnect succeeds when using the correct credentials.
func TestConfigReloadRotateUserAuthentication(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/single_user_authentication_1.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Ensure we can connect as a sanity check.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port)
	nc, err := nats.Connect(addr, nats.UserInfo("tyler", "T0pS3cr3t"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc.Close()
	disconnected := make(chan struct{}, 1)
	asyncErr := make(chan error, 1)
	nc.SetErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
		asyncErr <- err
	})
	nc.SetDisconnectHandler(func(*nats.Conn) {
		disconnected <- struct{}{}
	})

	// Change user credentials.
	changeCurrentConfigContent(t, config, "./configs/reload/single_user_authentication_2.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure connecting fails.
	if _, err := nats.Connect(addr, nats.UserInfo("tyler", "T0pS3cr3t")); err == nil {
		t.Fatal("Expected connect to fail")
	}

	// Ensure connecting succeeds when using new credentials.
	conn, err := nats.Connect(addr, nats.UserInfo("derek", "passw0rd"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	conn.Close()

	// Ensure the previous connection received an authorization error.
	// Note that it is possible that client gets EOF and not able to
	// process async error, so don't fail if we don't get it.
	select {
	case err := <-asyncErr:
		if err != nats.ErrAuthorization {
			t.Fatalf("Expected ErrAuthorization, got %v", err)
		}
	case <-time.After(time.Second):
		// Give it up to 1 sec.
	}

	// Ensure the previous connection was disconnected.
	select {
	case <-disconnected:
	case <-time.After(2 * time.Second):
		t.Fatal("Expected connection to be disconnected")
	}
}

// Ensure Reload supports enabling single user authentication. Test this by
// starting a server with authentication disabled, connect to it to verify,
// reload config using with a username/password, ensure reconnect fails, then
// ensure reconnect succeeds when using the correct credentials.
func TestConfigReloadEnableUserAuthentication(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/basic.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Ensure we can connect as a sanity check.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port)
	nc, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc.Close()
	disconnected := make(chan struct{}, 1)
	asyncErr := make(chan error, 1)
	nc.SetErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
		asyncErr <- err
	})
	nc.SetDisconnectHandler(func(*nats.Conn) {
		disconnected <- struct{}{}
	})

	// Enable authentication.
	changeCurrentConfigContent(t, config, "./configs/reload/single_user_authentication_1.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure connecting fails.
	if _, err := nats.Connect(addr); err == nil {
		t.Fatal("Expected connect to fail")
	}

	// Ensure connecting succeeds when using new credentials.
	conn, err := nats.Connect(addr, nats.UserInfo("tyler", "T0pS3cr3t"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	conn.Close()

	// Ensure the previous connection received an authorization error.
	// Note that it is possible that client gets EOF and not able to
	// process async error, so don't fail if we don't get it.
	select {
	case err := <-asyncErr:
		if err != nats.ErrAuthorization {
			t.Fatalf("Expected ErrAuthorization, got %v", err)
		}
	case <-time.After(time.Second):
	}

	// Ensure the previous connection was disconnected.
	select {
	case <-disconnected:
	case <-time.After(2 * time.Second):
		t.Fatal("Expected connection to be disconnected")
	}
}

// Ensure Reload supports disabling single user authentication. Test this by
// starting a server with authentication enabled, connect to it to verify,
// reload config using with authentication disabled, then ensure connecting
// with no credentials succeeds.
func TestConfigReloadDisableUserAuthentication(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/single_user_authentication_1.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Ensure we can connect as a sanity check.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port)
	nc, err := nats.Connect(addr, nats.UserInfo("tyler", "T0pS3cr3t"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc.Close()
	nc.SetErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
		t.Fatalf("Client received an unexpected error: %v", err)
	})

	// Disable authentication.
	changeCurrentConfigContent(t, config, "./configs/reload/basic.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure connecting succeeds with no credentials.
	conn, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	conn.Close()
}

// Ensure Reload supports token authentication config changes. Test this by
// starting a server with token authentication enabled, connect to it to
// verify, reload config using a different token, ensure reconnect fails, then
// ensure reconnect succeeds when using the correct token.
func TestConfigReloadRotateTokenAuthentication(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/token_authentication_1.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	disconnected := make(chan struct{})
	asyncErr := make(chan error)
	eh := func(nc *nats.Conn, sub *nats.Subscription, err error) { asyncErr <- err }
	dh := func(*nats.Conn) { disconnected <- struct{}{} }

	// Ensure we can connect as a sanity check.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port)
	nc, err := nats.Connect(addr, nats.Token("T0pS3cr3t"), nats.ErrorHandler(eh), nats.DisconnectHandler(dh))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc.Close()

	// Change authentication token.
	changeCurrentConfigContent(t, config, "./configs/reload/token_authentication_2.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure connecting fails.
	if _, err := nats.Connect(addr, nats.Token("T0pS3cr3t")); err == nil {
		t.Fatal("Expected connect to fail")
	}

	// Ensure connecting succeeds when using new credentials.
	conn, err := nats.Connect(addr, nats.Token("passw0rd"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	conn.Close()

	// Ensure the previous connection received an authorization error.
	select {
	case err := <-asyncErr:
		if err != nats.ErrAuthorization {
			t.Fatalf("Expected ErrAuthorization, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Expected authorization error")
	}

	// Ensure the previous connection was disconnected.
	select {
	case <-disconnected:
	case <-time.After(2 * time.Second):
		t.Fatal("Expected connection to be disconnected")
	}
}

// Ensure Reload supports enabling token authentication. Test this by starting
// a server with authentication disabled, connect to it to verify, reload
// config using with a token, ensure reconnect fails, then ensure reconnect
// succeeds when using the correct token.
func TestConfigReloadEnableTokenAuthentication(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/basic.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Ensure we can connect as a sanity check.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port)
	nc, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc.Close()
	disconnected := make(chan struct{}, 1)
	asyncErr := make(chan error, 1)
	nc.SetErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
		asyncErr <- err
	})
	nc.SetDisconnectHandler(func(*nats.Conn) {
		disconnected <- struct{}{}
	})

	// Enable authentication.
	changeCurrentConfigContent(t, config, "./configs/reload/token_authentication_1.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure connecting fails.
	if _, err := nats.Connect(addr); err == nil {
		t.Fatal("Expected connect to fail")
	}

	// Ensure connecting succeeds when using new credentials.
	conn, err := nats.Connect(addr, nats.Token("T0pS3cr3t"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	conn.Close()

	// Ensure the previous connection received an authorization error.
	// Note that it is possible that client gets EOF and not able to
	// process async error, so don't fail if we don't get it.
	select {
	case err := <-asyncErr:
		if err != nats.ErrAuthorization {
			t.Fatalf("Expected ErrAuthorization, got %v", err)
		}
	case <-time.After(time.Second):
	}

	// Ensure the previous connection was disconnected.
	select {
	case <-disconnected:
	case <-time.After(2 * time.Second):
		t.Fatal("Expected connection to be disconnected")
	}
}

// Ensure Reload supports disabling single token authentication. Test this by
// starting a server with authentication enabled, connect to it to verify,
// reload config using with authentication disabled, then ensure connecting
// with no token succeeds.
func TestConfigReloadDisableTokenAuthentication(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/token_authentication_1.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Ensure we can connect as a sanity check.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port)
	nc, err := nats.Connect(addr, nats.Token("T0pS3cr3t"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc.Close()
	nc.SetErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
		t.Fatalf("Client received an unexpected error: %v", err)
	})

	// Disable authentication.
	changeCurrentConfigContent(t, config, "./configs/reload/basic.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure connecting succeeds with no credentials.
	conn, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	conn.Close()
}

// Ensure Reload supports users authentication config changes. Test this by
// starting a server with users authentication enabled, connect to it to
// verify, reload config using a different user, ensure reconnect fails, then
// ensure reconnect succeeds when using the correct credentials.
func TestConfigReloadRotateUsersAuthentication(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/multiple_users_1.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Ensure we can connect as a sanity check.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port)
	nc, err := nats.Connect(addr, nats.UserInfo("alice", "foo"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc.Close()
	disconnected := make(chan struct{}, 1)
	asyncErr := make(chan error, 1)
	nc.SetErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
		asyncErr <- err
	})
	nc.SetDisconnectHandler(func(*nats.Conn) {
		disconnected <- struct{}{}
	})

	// These credentials won't change.
	nc2, err := nats.Connect(addr, nats.UserInfo("bob", "bar"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc2.Close()
	sub, err := nc2.SubscribeSync("foo")
	if err != nil {
		t.Fatalf("Error subscribing: %v", err)
	}
	defer sub.Unsubscribe()

	// Change users credentials.
	changeCurrentConfigContent(t, config, "./configs/reload/multiple_users_2.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure connecting fails.
	if _, err := nats.Connect(addr, nats.UserInfo("alice", "foo")); err == nil {
		t.Fatal("Expected connect to fail")
	}

	// Ensure connecting succeeds when using new credentials.
	conn, err := nats.Connect(addr, nats.UserInfo("alice", "baz"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	conn.Close()

	// Ensure the previous connection received an authorization error.
	// Note that it is possible that client gets EOF and not able to
	// process async error, so don't fail if we don't get it.
	select {
	case err := <-asyncErr:
		if err != nats.ErrAuthorization {
			t.Fatalf("Expected ErrAuthorization, got %v", err)
		}
	case <-time.After(time.Second):
	}

	// Ensure the previous connection was disconnected.
	select {
	case <-disconnected:
	case <-time.After(2 * time.Second):
		t.Fatal("Expected connection to be disconnected")
	}

	// Ensure the connection using unchanged credentials can still
	// publish/receive.
	if err := nc2.Publish("foo", []byte("hello")); err != nil {
		t.Fatalf("Error publishing: %v", err)
	}
	nc2.Flush()
	msg, err := sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving msg: %v", err)
	}
	if string(msg.Data) != "hello" {
		t.Fatalf("Msg is incorrect.\nexpected: %+v\ngot: %+v", []byte("hello"), msg.Data)
	}
}

// Ensure Reload supports enabling users authentication. Test this by starting
// a server with authentication disabled, connect to it to verify, reload
// config using with users, ensure reconnect fails, then ensure reconnect
// succeeds when using the correct credentials.
func TestConfigReloadEnableUsersAuthentication(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/basic.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Ensure we can connect as a sanity check.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port)
	nc, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc.Close()
	disconnected := make(chan struct{}, 1)
	asyncErr := make(chan error, 1)
	nc.SetErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
		asyncErr <- err
	})
	nc.SetDisconnectHandler(func(*nats.Conn) {
		disconnected <- struct{}{}
	})

	// Enable authentication.
	changeCurrentConfigContent(t, config, "./configs/reload/multiple_users_1.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure connecting fails.
	if _, err := nats.Connect(addr); err == nil {
		t.Fatal("Expected connect to fail")
	}

	// Ensure connecting succeeds when using new credentials.
	conn, err := nats.Connect(addr, nats.UserInfo("alice", "foo"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	conn.Close()

	// Ensure the previous connection received an authorization error.
	// Note that it is possible that client gets EOF and not able to
	// process async error, so don't fail if we don't get it.
	select {
	case err := <-asyncErr:
		if err != nats.ErrAuthorization {
			t.Fatalf("Expected ErrAuthorization, got %v", err)
		}
	case <-time.After(time.Second):
	}

	// Ensure the previous connection was disconnected.
	select {
	case <-disconnected:
	case <-time.After(5 * time.Second):
		t.Fatal("Expected connection to be disconnected")
	}
}

// Ensure Reload supports disabling users authentication. Test this by starting
// a server with authentication enabled, connect to it to verify,
// reload config using with authentication disabled, then ensure connecting
// with no credentials succeeds.
func TestConfigReloadDisableUsersAuthentication(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/multiple_users_1.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Ensure we can connect as a sanity check.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port)
	nc, err := nats.Connect(addr, nats.UserInfo("alice", "foo"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc.Close()
	nc.SetErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
		t.Fatalf("Client received an unexpected error: %v", err)
	})

	// Disable authentication.
	changeCurrentConfigContent(t, config, "./configs/reload/basic.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure connecting succeeds with no credentials.
	conn, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	conn.Close()
}

// Ensure Reload supports changing permissions. Test this by starting a server
// with a user configured with certain permissions, test publish and subscribe,
// reload config with new permissions, ensure the previous subscription was
// closed and publishes fail, then ensure the new permissions succeed.
func TestConfigReloadChangePermissions(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/authorization_1.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	addr := fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port)
	nc, err := nats.Connect(addr, nats.UserInfo("bob", "bar"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc.Close()
	asyncErr := make(chan error, 1)
	nc.SetErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
		asyncErr <- err
	})
	// Ensure we can publish and receive messages as a sanity check.
	sub, err := nc.SubscribeSync("_INBOX.>")
	if err != nil {
		t.Fatalf("Error subscribing: %v", err)
	}
	nc.Flush()

	conn, err := nats.Connect(addr, nats.UserInfo("alice", "foo"))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer conn.Close()

	sub2, err := conn.SubscribeSync("req.foo")
	if err != nil {
		t.Fatalf("Error subscribing: %v", err)
	}
	if err := conn.Publish("_INBOX.foo", []byte("hello")); err != nil {
		t.Fatalf("Error publishing message: %v", err)
	}
	conn.Flush()

	msg, err := sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving msg: %v", err)
	}
	if string(msg.Data) != "hello" {
		t.Fatalf("Msg is incorrect.\nexpected: %+v\ngot: %+v", []byte("hello"), msg.Data)
	}

	if err := nc.Publish("req.foo", []byte("world")); err != nil {
		t.Fatalf("Error publishing message: %v", err)
	}
	nc.Flush()

	msg, err = sub2.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving msg: %v", err)
	}
	if string(msg.Data) != "world" {
		t.Fatalf("Msg is incorrect.\nexpected: %+v\ngot: %+v", []byte("world"), msg.Data)
	}

	// Change permissions.
	changeCurrentConfigContent(t, config, "./configs/reload/authorization_2.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure we receive an error for the subscription that is no longer
	// authorized.
	// In this test, since connection is not closed by the server,
	// the client must receive an -ERR
	select {
	case err := <-asyncErr:
		if !strings.Contains(strings.ToLower(err.Error()), "permissions violation for subscription to \"_inbox.>\"") {
			t.Fatalf("Expected permissions violation error, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Expected permissions violation error")
	}

	// Ensure we receive an error when publishing to req.foo and we no longer
	// receive messages on _INBOX.>.
	if err := nc.Publish("req.foo", []byte("hola")); err != nil {
		t.Fatalf("Error publishing message: %v", err)
	}
	nc.Flush()
	if err := conn.Publish("_INBOX.foo", []byte("mundo")); err != nil {
		t.Fatalf("Error publishing message: %v", err)
	}
	conn.Flush()

	select {
	case err := <-asyncErr:
		if !strings.Contains(strings.ToLower(err.Error()), "permissions violation for publish to \"req.foo\"") {
			t.Fatalf("Expected permissions violation error, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Expected permissions violation error")
	}

	queued, _, err := sub2.Pending()
	if err != nil {
		t.Fatalf("Failed to get pending messaged: %v", err)
	}
	if queued != 0 {
		t.Fatalf("Pending is incorrect.\nexpected: 0\ngot: %d", queued)
	}

	queued, _, err = sub.Pending()
	if err != nil {
		t.Fatalf("Failed to get pending messaged: %v", err)
	}
	if queued != 0 {
		t.Fatalf("Pending is incorrect.\nexpected: 0\ngot: %d", queued)
	}

	// Ensure we can publish to _INBOX.foo.bar and subscribe to _INBOX.foo.>.
	sub, err = nc.SubscribeSync("_INBOX.foo.>")
	if err != nil {
		t.Fatalf("Error subscribing: %v", err)
	}
	nc.Flush()
	if err := nc.Publish("_INBOX.foo.bar", []byte("testing")); err != nil {
		t.Fatalf("Error publishing message: %v", err)
	}
	nc.Flush()
	msg, err = sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving msg: %v", err)
	}
	if string(msg.Data) != "testing" {
		t.Fatalf("Msg is incorrect.\nexpected: %+v\ngot: %+v", []byte("testing"), msg.Data)
	}

	select {
	case err := <-asyncErr:
		t.Fatalf("Received unexpected error: %v", err)
	default:
	}
}

// Ensure Reload returns an error when attempting to change cluster address
// host.
func TestConfigReloadClusterHostUnsupported(t *testing.T) {
	server, _, config := runReloadServerWithConfig(t, "./configs/reload/srv_a_1.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Attempt to change cluster listen host.
	changeCurrentConfigContent(t, config, "./configs/reload/srv_c_1.conf")

	// This should fail because cluster address cannot be changed.
	if err := server.Reload(); err == nil {
		t.Fatal("Expected Reload to return an error")
	}
}

// Ensure Reload returns an error when attempting to change cluster address
// port.
func TestConfigReloadClusterPortUnsupported(t *testing.T) {
	server, _, config := runReloadServerWithConfig(t, "./configs/reload/srv_a_1.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Attempt to change cluster listen port.
	changeCurrentConfigContent(t, config, "./configs/reload/srv_b_1.conf")

	// This should fail because cluster address cannot be changed.
	if err := server.Reload(); err == nil {
		t.Fatal("Expected Reload to return an error")
	}
}

// Ensure Reload supports enabling route authorization. Test this by starting
// two servers in a cluster without authorization, ensuring messages flow
// between them, then reloading with authorization and ensuring messages no
// longer flow until reloading with the correct credentials.
func TestConfigReloadEnableClusterAuthorization(t *testing.T) {
	srvb, srvbOpts, srvbConfig := runReloadServerWithConfig(t, "./configs/reload/srv_b_1.conf")
	defer os.Remove(srvbConfig)
	defer srvb.Shutdown()

	srva, srvaOpts, srvaConfig := runReloadServerWithConfig(t, "./configs/reload/srv_a_1.conf")
	defer os.Remove(srvaConfig)
	defer srva.Shutdown()

	checkClusterFormed(t, srva, srvb)

	srvaAddr := fmt.Sprintf("nats://%s:%d", srvaOpts.Host, srvaOpts.Port)
	srvaConn, err := nats.Connect(srvaAddr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer srvaConn.Close()
	sub, err := srvaConn.SubscribeSync("foo")
	if err != nil {
		t.Fatalf("Error subscribing: %v", err)
	}
	defer sub.Unsubscribe()
	if err := srvaConn.Flush(); err != nil {
		t.Fatalf("Error flushing: %v", err)
	}

	srvbAddr := fmt.Sprintf("nats://%s:%d", srvbOpts.Host, srvbOpts.Port)
	srvbConn, err := nats.Connect(srvbAddr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer srvbConn.Close()

	if numRoutes := srvb.NumRoutes(); numRoutes != 1 {
		t.Fatalf("Expected 1 route, got %d", numRoutes)
	}

	// Ensure messages flow through the cluster as a sanity check.
	if err := srvbConn.Publish("foo", []byte("hello")); err != nil {
		t.Fatalf("Error publishing: %v", err)
	}
	srvbConn.Flush()
	msg, err := sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving message: %v", err)
	}
	if string(msg.Data) != "hello" {
		t.Fatalf("Msg is incorrect.\nexpected: %+v\ngot: %+v", []byte("hello"), msg.Data)
	}

	// Enable route authorization.
	changeCurrentConfigContent(t, srvbConfig, "./configs/reload/srv_b_2.conf")
	if err := srvb.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	if numRoutes := srvb.NumRoutes(); numRoutes != 0 {
		t.Fatalf("Expected 0 routes, got %d", numRoutes)
	}

	// Ensure messages no longer flow through the cluster.
	for i := 0; i < 5; i++ {
		if err := srvbConn.Publish("foo", []byte("world")); err != nil {
			t.Fatalf("Error publishing: %v", err)
		}
		srvbConn.Flush()
	}
	if _, err := sub.NextMsg(50 * time.Millisecond); err != nats.ErrTimeout {
		t.Fatalf("Expected ErrTimeout, got %v", err)
	}

	// Reload Server A with correct route credentials.
	changeCurrentConfigContent(t, srvaConfig, "./configs/reload/srv_a_2.conf")
	defer os.Remove(srvaConfig)
	if err := srva.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}
	checkClusterFormed(t, srva, srvb)

	if numRoutes := srvb.NumRoutes(); numRoutes != 1 {
		t.Fatalf("Expected 1 route, got %d", numRoutes)
	}

	// Ensure messages flow through the cluster now.
	if err := srvbConn.Publish("foo", []byte("hola")); err != nil {
		t.Fatalf("Error publishing: %v", err)
	}
	srvbConn.Flush()
	msg, err = sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving message: %v", err)
	}
	if string(msg.Data) != "hola" {
		t.Fatalf("Msg is incorrect.\nexpected: %+v\ngot: %+v", []byte("hola"), msg.Data)
	}
}

// Ensure Reload supports disabling route authorization. Test this by starting
// two servers in a cluster with authorization, ensuring messages flow
// between them, then reloading without authorization and ensuring messages
// still flow.
func TestConfigReloadDisableClusterAuthorization(t *testing.T) {
	srvb, srvbOpts, srvbConfig := runReloadServerWithConfig(t, "./configs/reload/srv_b_2.conf")
	defer os.Remove(srvbConfig)
	defer srvb.Shutdown()

	srva, srvaOpts, srvaConfig := runReloadServerWithConfig(t, "./configs/reload/srv_a_2.conf")
	defer os.Remove(srvaConfig)
	defer srva.Shutdown()

	checkClusterFormed(t, srva, srvb)

	srvaAddr := fmt.Sprintf("nats://%s:%d", srvaOpts.Host, srvaOpts.Port)
	srvaConn, err := nats.Connect(srvaAddr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer srvaConn.Close()

	sub, err := srvaConn.SubscribeSync("foo")
	if err != nil {
		t.Fatalf("Error subscribing: %v", err)
	}
	defer sub.Unsubscribe()
	if err := srvaConn.Flush(); err != nil {
		t.Fatalf("Error flushing: %v", err)
	}

	srvbAddr := fmt.Sprintf("nats://%s:%d", srvbOpts.Host, srvbOpts.Port)
	srvbConn, err := nats.Connect(srvbAddr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer srvbConn.Close()

	if numRoutes := srvb.NumRoutes(); numRoutes != 1 {
		t.Fatalf("Expected 1 route, got %d", numRoutes)
	}

	// Ensure messages flow through the cluster as a sanity check.
	if err := srvbConn.Publish("foo", []byte("hello")); err != nil {
		t.Fatalf("Error publishing: %v", err)
	}
	srvbConn.Flush()
	msg, err := sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving message: %v", err)
	}
	if string(msg.Data) != "hello" {
		t.Fatalf("Msg is incorrect.\nexpected: %+v\ngot: %+v", []byte("hello"), msg.Data)
	}

	// Disable route authorization.
	changeCurrentConfigContent(t, srvbConfig, "./configs/reload/srv_b_1.conf")
	if err := srvb.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	checkClusterFormed(t, srva, srvb)

	if numRoutes := srvb.NumRoutes(); numRoutes != 1 {
		t.Fatalf("Expected 1 route, got %d", numRoutes)
	}

	// Ensure messages still flow through the cluster.
	if err := srvbConn.Publish("foo", []byte("hola")); err != nil {
		t.Fatalf("Error publishing: %v", err)
	}
	srvbConn.Flush()
	msg, err = sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving message: %v", err)
	}
	if string(msg.Data) != "hola" {
		t.Fatalf("Msg is incorrect.\nexpected: %+v\ngot: %+v", []byte("hola"), msg.Data)
	}
}

// Ensure Reload supports changing cluster routes. Test this by starting
// two servers in a cluster, ensuring messages flow between them, then
// reloading with a different route and ensuring messages flow through the new
// cluster.
func TestConfigReloadClusterRoutes(t *testing.T) {
	srvb, srvbOpts, srvbConfig := runReloadServerWithConfig(t, "./configs/reload/srv_b_1.conf")
	defer os.Remove(srvbConfig)
	defer srvb.Shutdown()

	srva, srvaOpts, srvaConfig := runReloadServerWithConfig(t, "./configs/reload/srv_a_1.conf")
	defer os.Remove(srvaConfig)
	defer srva.Shutdown()

	checkClusterFormed(t, srva, srvb)

	srvcOpts, err := ProcessConfigFile("./configs/reload/srv_c_1.conf")
	if err != nil {
		t.Fatalf("Error processing config file: %v", err)
	}
	srvcOpts.NoLog = true
	srvcOpts.NoSigs = true

	srvc := RunServer(srvcOpts)
	defer srvc.Shutdown()

	srvaAddr := fmt.Sprintf("nats://%s:%d", srvaOpts.Host, srvaOpts.Port)
	srvaConn, err := nats.Connect(srvaAddr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer srvaConn.Close()

	sub, err := srvaConn.SubscribeSync("foo")
	if err != nil {
		t.Fatalf("Error subscribing: %v", err)
	}
	defer sub.Unsubscribe()
	if err := srvaConn.Flush(); err != nil {
		t.Fatalf("Error flushing: %v", err)
	}

	srvbAddr := fmt.Sprintf("nats://%s:%d", srvbOpts.Host, srvbOpts.Port)
	srvbConn, err := nats.Connect(srvbAddr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer srvbConn.Close()

	if numRoutes := srvb.NumRoutes(); numRoutes != 1 {
		t.Fatalf("Expected 1 route, got %d", numRoutes)
	}

	// Ensure messages flow through the cluster as a sanity check.
	if err := srvbConn.Publish("foo", []byte("hello")); err != nil {
		t.Fatalf("Error publishing: %v", err)
	}
	srvbConn.Flush()
	msg, err := sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving message: %v", err)
	}
	if string(msg.Data) != "hello" {
		t.Fatalf("Msg is incorrect.\nexpected: %+v\ngot: %+v", []byte("hello"), msg.Data)
	}

	// Reload cluster routes.
	changeCurrentConfigContent(t, srvaConfig, "./configs/reload/srv_a_3.conf")
	if err := srva.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Kill old route server.
	srvbConn.Close()
	srvb.Shutdown()

	checkClusterFormed(t, srva, srvc)

	srvcAddr := fmt.Sprintf("nats://%s:%d", srvcOpts.Host, srvcOpts.Port)
	srvcConn, err := nats.Connect(srvcAddr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer srvcConn.Close()

	// Ensure messages flow through the new cluster.
	for i := 0; i < 5; i++ {
		if err := srvcConn.Publish("foo", []byte("hola")); err != nil {
			t.Fatalf("Error publishing: %v", err)
		}
		srvcConn.Flush()
	}
	msg, err = sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving message: %v", err)
	}
	if string(msg.Data) != "hola" {
		t.Fatalf("Msg is incorrect.\nexpected: %+v\ngot: %+v", []byte("hola"), msg.Data)
	}
}

// Ensure Reload supports removing a solicited route. In this case from A->B
// Test this by starting two servers in a cluster, ensuring messages flow between them.
// Then stop server B, and have server A continue to try to connect. Reload A with a config
// that removes the route and make sure it does not connect to server B when its restarted.
func TestConfigReloadClusterRemoveSolicitedRoutes(t *testing.T) {
	srvb, srvbOpts := RunServerWithConfig("./configs/reload/srv_b_1.conf")
	defer srvb.Shutdown()

	srva, srvaOpts, srvaConfig := runReloadServerWithConfig(t, "./configs/reload/srv_a_1.conf")
	defer os.Remove(srvaConfig)
	defer srva.Shutdown()

	checkClusterFormed(t, srva, srvb)

	srvaAddr := fmt.Sprintf("nats://%s:%d", srvaOpts.Host, srvaOpts.Port)
	srvaConn, err := nats.Connect(srvaAddr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer srvaConn.Close()
	sub, err := srvaConn.SubscribeSync("foo")
	if err != nil {
		t.Fatalf("Error subscribing: %v", err)
	}
	defer sub.Unsubscribe()
	if err := srvaConn.Flush(); err != nil {
		t.Fatalf("Error flushing: %v", err)
	}

	srvbAddr := fmt.Sprintf("nats://%s:%d", srvbOpts.Host, srvbOpts.Port)
	srvbConn, err := nats.Connect(srvbAddr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer srvbConn.Close()

	if err := srvbConn.Publish("foo", []byte("hello")); err != nil {
		t.Fatalf("Error publishing: %v", err)
	}
	srvbConn.Flush()
	msg, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving message: %v", err)
	}
	if string(msg.Data) != "hello" {
		t.Fatalf("Msg is incorrect.\nexpected: %+v\ngot: %+v", []byte("hello"), msg.Data)
	}

	// Now stop server B.
	srvb.Shutdown()

	// Wait til route is dropped.
	checkNumRoutes(t, srva, 0)

	// Now change config for server A to not solicit a route to server B.
	changeCurrentConfigContent(t, srvaConfig, "./configs/reload/srv_a_4.conf")
	defer os.Remove(srvaConfig)
	if err := srva.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Restart server B.
	srvb, _ = RunServerWithConfig("./configs/reload/srv_b_1.conf")
	defer srvb.Shutdown()

	// We should not have a cluster formed here.
	numRoutes := 0
	deadline := time.Now().Add(2 * DEFAULT_ROUTE_RECONNECT)
	for time.Now().Before(deadline) {
		if numRoutes = srva.NumRoutes(); numRoutes != 0 {
			break
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
	if numRoutes != 0 {
		t.Fatalf("Expected 0 routes for server A, got %d", numRoutes)
	}
}

func reloadUpdateConfig(t *testing.T, s *Server, conf, content string) {
	if err := ioutil.WriteFile(conf, []byte(content), 0666); err != nil {
		stackFatalf(t, "Error creating config file: %v", err)
	}
	if err := s.Reload(); err != nil {
		stackFatalf(t, "Error on reload: %v", err)
	}
}

func TestConfigReloadClusterAdvertise(t *testing.T) {
	s, _, conf := runReloadServerWithContent(t, []byte(`
		listen: "0.0.0.0:-1"
		cluster: {
			listen: "0.0.0.0:-1"
		}
	`))
	defer os.Remove(conf)
	defer s.Shutdown()

	orgClusterPort := s.ClusterAddr().Port

	verify := func(expectedHost string, expectedPort int, expectedIP string) {
		s.mu.Lock()
		routeInfo := s.routeInfo
		routeInfoJSON := Info{}
		err := json.Unmarshal(s.routeInfoJSON[5:], &routeInfoJSON) // Skip "INFO "
		s.mu.Unlock()
		if err != nil {
			t.Fatalf("Error on Unmarshal: %v", err)
		}
		if routeInfo.Host != expectedHost || routeInfo.Port != expectedPort || routeInfo.IP != expectedIP {
			t.Fatalf("Expected host/port/IP to be %s:%v, %q, got %s:%d, %q",
				expectedHost, expectedPort, expectedIP, routeInfo.Host, routeInfo.Port, routeInfo.IP)
		}
		// Check that server routeInfoJSON was updated too
		if !reflect.DeepEqual(routeInfo, routeInfoJSON) {
			t.Fatalf("Expected routeInfoJSON to be %+v, got %+v", routeInfo, routeInfoJSON)
		}
	}

	// Update config with cluster_advertise
	reloadUpdateConfig(t, s, conf, `
	listen: "0.0.0.0:-1"
	cluster: {
		listen: "0.0.0.0:-1"
		cluster_advertise: "me:1"
	}
	`)
	verify("me", 1, "nats-route://me:1/")

	// Update config with cluster_advertise (no port specified)
	reloadUpdateConfig(t, s, conf, `
	listen: "0.0.0.0:-1"
	cluster: {
		listen: "0.0.0.0:-1"
		cluster_advertise: "me"
	}
	`)
	verify("me", orgClusterPort, fmt.Sprintf("nats-route://me:%d/", orgClusterPort))

	// Update config with cluster_advertise (-1 port specified)
	reloadUpdateConfig(t, s, conf, `
	listen: "0.0.0.0:-1"
	cluster: {
		listen: "0.0.0.0:-1"
		cluster_advertise: "me:-1"
	}
	`)
	verify("me", orgClusterPort, fmt.Sprintf("nats-route://me:%d/", orgClusterPort))

	// Update to remove cluster_advertise
	reloadUpdateConfig(t, s, conf, `
	listen: "0.0.0.0:-1"
	cluster: {
		listen: "0.0.0.0:-1"
	}
	`)
	verify("0.0.0.0", orgClusterPort, "")
}

func TestConfigReloadClusterNoAdvertise(t *testing.T) {
	s, _, conf := runReloadServerWithContent(t, []byte(`
		listen: "0.0.0.0:-1"
		client_advertise: "me:1"
		cluster: {
			listen: "0.0.0.0:-1"
		}
	`))
	defer os.Remove(conf)
	defer s.Shutdown()

	s.mu.Lock()
	ccurls := s.routeInfo.ClientConnectURLs
	s.mu.Unlock()
	if len(ccurls) != 1 && ccurls[0] != "me:1" {
		t.Fatalf("Unexpected routeInfo.ClientConnectURLS: %v", ccurls)
	}

	// Update config with no_advertise
	reloadUpdateConfig(t, s, conf, `
	listen: "0.0.0.0:-1"
	client_advertise: "me:1"
	cluster: {
		listen: "0.0.0.0:-1"
		no_advertise: true
	}
	`)

	s.mu.Lock()
	ccurls = s.routeInfo.ClientConnectURLs
	s.mu.Unlock()
	if len(ccurls) != 0 {
		t.Fatalf("Unexpected routeInfo.ClientConnectURLS: %v", ccurls)
	}

	// Update config with cluster_advertise (no port specified)
	reloadUpdateConfig(t, s, conf, `
	listen: "0.0.0.0:-1"
	client_advertise: "me:1"
	cluster: {
		listen: "0.0.0.0:-1"
	}
	`)
	s.mu.Lock()
	ccurls = s.routeInfo.ClientConnectURLs
	s.mu.Unlock()
	if len(ccurls) != 1 && ccurls[0] != "me:1" {
		t.Fatalf("Unexpected routeInfo.ClientConnectURLS: %v", ccurls)
	}
}

func TestConfigReloadMaxSubsUnsupported(t *testing.T) {
	s, _, conf := runReloadServerWithContent(t, []byte(`max_subs: 1`))
	defer os.Remove(conf)
	defer s.Shutdown()

	if err := ioutil.WriteFile(conf, []byte(`max_subs: 10`), 0666); err != nil {
		t.Fatalf("Error writing config file: %v", err)
	}
	if err := s.Reload(); err == nil {
		t.Fatal("Expected Reload to return an error")
	}
}

func TestConfigReloadClientAdvertise(t *testing.T) {
	s, _, conf := runReloadServerWithContent(t, []byte(`listen: "0.0.0.0:-1"`))
	defer os.Remove(conf)
	defer s.Shutdown()

	orgPort := s.Addr().(*net.TCPAddr).Port

	verify := func(expectedHost string, expectedPort int) {
		s.mu.Lock()
		info := s.info
		s.mu.Unlock()
		if info.Host != expectedHost || info.Port != expectedPort {
			stackFatalf(t, "Expected host/port to be %s:%d, got %s:%d",
				expectedHost, expectedPort, info.Host, info.Port)
		}
	}

	// Update config with ClientAdvertise (port specified)
	reloadUpdateConfig(t, s, conf, `
	listen: "0.0.0.0:-1"
	client_advertise: "me:1"
	`)
	verify("me", 1)

	// Update config with ClientAdvertise (no port specified)
	reloadUpdateConfig(t, s, conf, `
	listen: "0.0.0.0:-1"
	client_advertise: "me"
	`)
	verify("me", orgPort)

	// Update config with ClientAdvertise (-1 port specified)
	reloadUpdateConfig(t, s, conf, `
	listen: "0.0.0.0:-1"
	client_advertise: "me:-1"
	`)
	verify("me", orgPort)

	// Now remove ClientAdvertise to check that original values
	// are restored.
	reloadUpdateConfig(t, s, conf, `listen: "0.0.0.0:-1"`)
	verify("0.0.0.0", orgPort)
}

// Ensure Reload supports changing the max connections. Test this by starting a
// server with no max connections, connecting two clients, reloading with a
// max connections of one, and ensuring one client is disconnected.
func TestConfigReloadMaxConnections(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/basic.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	// Make two connections.
	addr := fmt.Sprintf("nats://%s:%d", opts.Host, server.Addr().(*net.TCPAddr).Port)
	nc1, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc1.Close()
	closed := make(chan struct{}, 1)
	nc1.SetDisconnectHandler(func(*nats.Conn) {
		closed <- struct{}{}
	})
	nc2, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc2.Close()
	nc2.SetDisconnectHandler(func(*nats.Conn) {
		closed <- struct{}{}
	})

	if numClients := server.NumClients(); numClients != 2 {
		t.Fatalf("Expected 2 clients, got %d", numClients)
	}

	// Set max connections to one.
	changeCurrentConfigContent(t, config, "./configs/reload/max_connections.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure one connection was closed.
	select {
	case <-closed:
	case <-time.After(5 * time.Second):
		t.Fatal("Expected to be disconnected")
	}

	if numClients := server.NumClients(); numClients != 1 {
		t.Fatalf("Expected 1 client, got %d", numClients)
	}

	// Ensure new connections fail.
	_, err = nats.Connect(addr)
	if err == nil {
		t.Fatal("Expected error on connect")
	}
}

// Ensure reload supports changing the max payload size. Test this by starting
// a server with the default size limit, ensuring publishes work, reloading
// with a restrictive limit, and ensuring publishing an oversized message fails
// and disconnects the client.
func TestConfigReloadMaxPayload(t *testing.T) {
	server, opts, config := runReloadServerWithConfig(t, "./configs/reload/basic.conf")
	defer os.Remove(config)
	defer server.Shutdown()

	addr := fmt.Sprintf("nats://%s:%d", opts.Host, server.Addr().(*net.TCPAddr).Port)
	nc, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer nc.Close()
	closed := make(chan struct{})
	nc.SetDisconnectHandler(func(*nats.Conn) {
		closed <- struct{}{}
	})

	conn, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer conn.Close()
	sub, err := conn.SubscribeSync("foo")
	if err != nil {
		t.Fatalf("Error subscribing: %v", err)
	}
	conn.Flush()

	// Ensure we can publish as a sanity check.
	if err := nc.Publish("foo", []byte("hello")); err != nil {
		t.Fatalf("Error publishing: %v", err)
	}
	nc.Flush()
	_, err = sub.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("Error receiving message: %v", err)
	}

	// Set max payload to one.
	changeCurrentConfigContent(t, config, "./configs/reload/max_payload.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Ensure oversized messages don't get delivered and the client is
	// disconnected.
	if err := nc.Publish("foo", []byte("hello")); err != nil {
		t.Fatalf("Error publishing: %v", err)
	}
	nc.Flush()
	_, err = sub.NextMsg(20 * time.Millisecond)
	if err != nats.ErrTimeout {
		t.Fatalf("Expected ErrTimeout, got: %v", err)
	}

	select {
	case <-closed:
	case <-time.After(5 * time.Second):
		t.Fatal("Expected to be disconnected")
	}
}

// Ensure reload supports rotating out files. Test this by starting
// a server with log and pid files, reloading new ones, then check that
// we can rename and delete the old log/pid files.
func TestConfigReloadRotateFiles(t *testing.T) {
	server, _, config := runReloadServerWithConfig(t, "./configs/reload/file_rotate.conf")
	defer func() {
		os.Remove(config)
		os.Remove("log1.txt")
		os.Remove("gnatsd1.pid")
	}()
	defer server.Shutdown()

	// Configure the logger to enable actual logging
	server.ConfigureLogger()

	// Load a config that renames the files.
	changeCurrentConfigContent(t, config, "./configs/reload/file_rotate1.conf")
	if err := server.Reload(); err != nil {
		t.Fatalf("Error reloading config: %v", err)
	}

	// Make sure the new files exist.
	if _, err := os.Stat("log1.txt"); os.IsNotExist(err) {
		t.Fatalf("Error reloading config, no new file: %v", err)
	}
	if _, err := os.Stat("gnatsd1.pid"); os.IsNotExist(err) {
		t.Fatalf("Error reloading config, no new file: %v", err)
	}

	// Check that old file can be renamed.
	if err := os.Rename("log.txt", "log_old.txt"); err != nil {
		t.Fatalf("Error reloading config, cannot rename file: %v", err)
	}
	if err := os.Rename("gnatsd.pid", "gnatsd_old.pid"); err != nil {
		t.Fatalf("Error reloading config, cannot rename file: %v", err)
	}

	// Check that the old files can be removed after rename.
	if err := os.Remove("log_old.txt"); err != nil {
		t.Fatalf("Error reloading config, cannot delete file: %v", err)
	}
	if err := os.Remove("gnatsd_old.pid"); err != nil {
		t.Fatalf("Error reloading config, cannot delete file: %v", err)
	}
}

func TestConfigReloadClusterWorks(t *testing.T) {
	confBTemplate := `
		listen: -1
		cluster: {
			listen: 127.0.0.1:7244
			authorization {
				user: ruser
				password: pwd
				timeout: %d
			}
			routes = [
				nats-route://ruser:pwd@127.0.0.1:7246
			]
		}`
	confB := createConfFile(t, []byte(fmt.Sprintf(confBTemplate, 3)))
	defer os.Remove(confB)

	confATemplate := `
		listen: -1
		cluster: {
			listen: 127.0.0.1:7246
			authorization {
				user: ruser
				password: pwd
				timeout: %d
			}
			routes = [
				nats-route://ruser:pwd@127.0.0.1:7244
			]
		}`
	confA := createConfFile(t, []byte(fmt.Sprintf(confATemplate, 3)))
	defer os.Remove(confA)

	srvb, _ := RunServerWithConfig(confB)
	defer srvb.Shutdown()

	srva, _ := RunServerWithConfig(confA)
	defer srva.Shutdown()

	// Wait for the cluster to form and capture the connection IDs of each route
	checkClusterFormed(t, srva, srvb)

	getCID := func(s *Server) uint64 {
		s.mu.Lock()
		defer s.mu.Unlock()
		for _, r := range s.routes {
			return r.cid
		}
		return 0
	}
	acid := getCID(srva)
	bcid := getCID(srvb)

	// Update auth timeout to force a check of the connected route auth
	reloadUpdateConfig(t, srvb, confB, fmt.Sprintf(confBTemplate, 5))
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, 5))

	// Wait a little bit to ensure that there is no issue with connection
	// breaking at this point (this was an issue before).
	time.Sleep(100 * time.Millisecond)

	// Cluster should still exist
	checkClusterFormed(t, srva, srvb)

	// Check that routes were not re-created
	newacid := getCID(srva)
	newbcid := getCID(srvb)

	if newacid != acid {
		t.Fatalf("Expected server A route ID to be %v, got %v", acid, newacid)
	}
	if newbcid != bcid {
		t.Fatalf("Expected server B route ID to be %v, got %v", bcid, newbcid)
	}
}

func TestConfigReloadClusterPerms(t *testing.T) {
	confATemplate := `
		port: -1
		cluster {
			listen: 127.0.0.1:-1
			permissions {
				import {
					allow: %s
				}
				export {
					allow: %s
				}
			}
		}
	`
	confA := createConfFile(t, []byte(fmt.Sprintf(confATemplate, `"foo"`, `"foo"`)))
	defer os.Remove(confA)
	srva, _ := RunServerWithConfig(confA)
	defer srva.Shutdown()

	confBTemplate := `
		port: -1
		cluster {
			listen: 127.0.0.1:-1
			permissions {
				import {
					allow: %s
				}
				export {
					allow: %s
				}
			}
			routes = [
				"nats://127.0.0.1:%d"
			]
		}
	`
	confB := createConfFile(t, []byte(fmt.Sprintf(confBTemplate, `"foo"`, `"foo"`, srva.ClusterAddr().Port)))
	defer os.Remove(confB)
	srvb, _ := RunServerWithConfig(confB)
	defer srvb.Shutdown()

	checkClusterFormed(t, srva, srvb)

	// Create a connection on A
	nca, err := nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", srva.Addr().(*net.TCPAddr).Port))
	if err != nil {
		t.Fatalf("Error on connect: %v", err)
	}
	defer nca.Close()
	// Create a subscription on "foo" and "bar", only "foo" will be also on server B.
	subFooOnA, err := nca.SubscribeSync("foo")
	if err != nil {
		t.Fatalf("Error on subscribe: %v", err)
	}
	subBarOnA, err := nca.SubscribeSync("bar")
	if err != nil {
		t.Fatalf("Error on subscribe: %v", err)
	}

	// Connect on B and do the same
	ncb, err := nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", srvb.Addr().(*net.TCPAddr).Port))
	if err != nil {
		t.Fatalf("Error on connect: %v", err)
	}
	defer ncb.Close()
	// Create a subscription on "foo" and "bar", only "foo" will be also on server B.
	subFooOnB, err := ncb.SubscribeSync("foo")
	if err != nil {
		t.Fatalf("Error on subscribe: %v", err)
	}
	subBarOnB, err := ncb.SubscribeSync("bar")
	if err != nil {
		t.Fatalf("Error on subscribe: %v", err)
	}

	// Check subscriptions on each server. There should be 3 on each server,
	// foo and bar locally and foo from remote server.
	checkExpectedSubs(t, 3, srva, srvb)

	sendMsg := func(t *testing.T, subj string, nc *nats.Conn) {
		t.Helper()
		if err := nc.Publish(subj, []byte("msg")); err != nil {
			t.Fatalf("Error on publish: %v", err)
		}
	}

	checkSub := func(t *testing.T, sub *nats.Subscription, shouldReceive bool) {
		t.Helper()
		_, err := sub.NextMsg(100 * time.Millisecond)
		if shouldReceive && err != nil {
			t.Fatalf("Expected message on %q, got %v", sub.Subject, err)
		} else if !shouldReceive && err == nil {
			t.Fatalf("Expected no message on %q, got one", sub.Subject)
		}
	}

	// Produce from A and check received on both sides
	sendMsg(t, "foo", nca)
	checkSub(t, subFooOnA, true)
	checkSub(t, subFooOnB, true)
	// Now from B:
	sendMsg(t, "foo", ncb)
	checkSub(t, subFooOnA, true)
	checkSub(t, subFooOnB, true)

	// Publish on bar from A and make sure only local sub receives
	sendMsg(t, "bar", nca)
	checkSub(t, subBarOnA, true)
	checkSub(t, subBarOnB, false)

	// Publish on bar from B and make sure only local sub receives
	sendMsg(t, "bar", ncb)
	checkSub(t, subBarOnA, false)
	checkSub(t, subBarOnB, true)

	// We will now both import/export foo and bar. Start with reloading A.
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, `["foo", "bar"]`, `["foo", "bar"]`))

	// Since B has not been updated yet, the state should remain the same,
	// that is 3 subs on each server.
	checkExpectedSubs(t, 3, srva, srvb)

	// Now update and reload B. Add "baz" for another test down below
	reloadUpdateConfig(t, srvb, confB, fmt.Sprintf(confBTemplate, `["foo", "bar", "baz"]`, `["foo", "bar", "baz"]`, srva.ClusterAddr().Port))

	// Now 4 on each server
	checkExpectedSubs(t, 4, srva, srvb)

	// Make sure that we can receive all messages
	sendMsg(t, "foo", nca)
	checkSub(t, subFooOnA, true)
	checkSub(t, subFooOnB, true)
	sendMsg(t, "foo", ncb)
	checkSub(t, subFooOnA, true)
	checkSub(t, subFooOnB, true)

	sendMsg(t, "bar", nca)
	checkSub(t, subBarOnA, true)
	checkSub(t, subBarOnB, true)
	sendMsg(t, "bar", ncb)
	checkSub(t, subBarOnA, true)
	checkSub(t, subBarOnB, true)

	// Create subscription on baz on server B.
	subBazOnB, err := ncb.SubscribeSync("baz")
	if err != nil {
		t.Fatalf("Error on subscribe: %v", err)
	}
	// Check subscriptions count
	checkExpectedSubs(t, 5, srvb)
	checkExpectedSubs(t, 4, srva)

	sendMsg(t, "baz", nca)
	checkSub(t, subBazOnB, false)
	sendMsg(t, "baz", ncb)
	checkSub(t, subBazOnB, true)

	// Test UNSUB by denying something that was previously imported
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, `"foo"`, `["foo", "bar"]`))
	// Since A no longer imports "bar", we should have one less subscription
	// on B (B will have received an UNSUB for bar)
	checkExpectedSubs(t, 4, srvb)
	// A, however, should still have same number of subs.
	checkExpectedSubs(t, 4, srva)

	// Remove all permissions from A.
	reloadUpdateConfig(t, srva, confA, `
		port: -1
		cluster {
			listen: 127.0.0.1:-1
		}
	`)
	// Server A should now have baz sub
	checkExpectedSubs(t, 5, srvb)
	checkExpectedSubs(t, 5, srva)

	sendMsg(t, "baz", nca)
	checkSub(t, subBazOnB, true)
	sendMsg(t, "baz", ncb)
	checkSub(t, subBazOnB, true)

	// Finally, remove permissions from B
	reloadUpdateConfig(t, srvb, confB, fmt.Sprintf(`
		port: -1
		cluster {
			listen: 127.0.0.1:-1
			routes = [
				"nats://127.0.0.1:%d"
			]
		}
	`, srva.ClusterAddr().Port))
	// Check expected subscriptions count.
	checkExpectedSubs(t, 5, srvb)
	checkExpectedSubs(t, 5, srva)
}

func TestConfigReloadClusterPermsImport(t *testing.T) {
	confATemplate := `
		port: -1
		cluster {
			listen: 127.0.0.1:-1
			permissions {
				import: {
					allow: %s
				}
			}
		}
	`
	confA := createConfFile(t, []byte(fmt.Sprintf(confATemplate, `["foo", "bar"]`)))
	defer os.Remove(confA)
	srva, _ := RunServerWithConfig(confA)
	defer srva.Shutdown()

	confBTemplate := `
		port: -1
		cluster {
			listen: 127.0.0.1:-1
			routes = [
				"nats://127.0.0.1:%d"
			]
		}
	`
	confB := createConfFile(t, []byte(fmt.Sprintf(confBTemplate, srva.ClusterAddr().Port)))
	defer os.Remove(confB)
	srvb, _ := RunServerWithConfig(confB)
	defer srvb.Shutdown()

	checkClusterFormed(t, srva, srvb)

	// Create a connection on A
	nca, err := nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", srva.Addr().(*net.TCPAddr).Port))
	if err != nil {
		t.Fatalf("Error on connect: %v", err)
	}
	defer nca.Close()
	// Create a subscription on "foo" and "bar"
	if _, err := nca.SubscribeSync("foo"); err != nil {
		t.Fatalf("Error on subscribe: %v", err)
	}
	if _, err := nca.SubscribeSync("bar"); err != nil {
		t.Fatalf("Error on subscribe: %v", err)
	}

	checkExpectedSubs(t, 2, srva, srvb)

	// Drop foo
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, `"bar"`))
	checkExpectedSubs(t, 2, srva)
	checkExpectedSubs(t, 1, srvb)

	// Add it back
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, `["foo", "bar"]`))
	checkExpectedSubs(t, 2, srva, srvb)

	// Empty Import means implicit allow
	reloadUpdateConfig(t, srva, confA, `
		port: -1
		cluster {
			listen: 127.0.0.1:-1
			permissions {
				export: ">"
			}
		}
	`)
	checkExpectedSubs(t, 2, srva, srvb)

	confATemplate = `
		port: -1
		cluster {
			listen: 127.0.0.1:-1
			permissions {
				import: {
					allow: ["foo", "bar"]
					deny: %s
				}
			}
		}
	`
	// Now deny all:
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, `["foo", "bar"]`))
	checkExpectedSubs(t, 2, srva)
	checkExpectedSubs(t, 0, srvb)

	// Drop foo from the deny list
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, `"bar"`))
	checkExpectedSubs(t, 2, srva)
	checkExpectedSubs(t, 1, srvb)
}

func TestConfigReloadClusterPermsExport(t *testing.T) {
	confATemplate := `
		port: -1
		cluster {
			listen: 127.0.0.1:-1
			permissions {
				export: {
					allow: %s
				}
			}
		}
	`
	confA := createConfFile(t, []byte(fmt.Sprintf(confATemplate, `["foo", "bar"]`)))
	defer os.Remove(confA)
	srva, _ := RunServerWithConfig(confA)
	defer srva.Shutdown()

	confBTemplate := `
		port: -1
		cluster {
			listen: 127.0.0.1:-1
			routes = [
				"nats://127.0.0.1:%d"
			]
		}
	`
	confB := createConfFile(t, []byte(fmt.Sprintf(confBTemplate, srva.ClusterAddr().Port)))
	defer os.Remove(confB)
	srvb, _ := RunServerWithConfig(confB)
	defer srvb.Shutdown()

	checkClusterFormed(t, srva, srvb)

	// Create a connection on B
	ncb, err := nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", srvb.Addr().(*net.TCPAddr).Port))
	if err != nil {
		t.Fatalf("Error on connect: %v", err)
	}
	defer ncb.Close()
	// Create a subscription on "foo" and "bar"
	if _, err := ncb.SubscribeSync("foo"); err != nil {
		t.Fatalf("Error on subscribe: %v", err)
	}
	if _, err := ncb.SubscribeSync("bar"); err != nil {
		t.Fatalf("Error on subscribe: %v", err)
	}

	checkExpectedSubs(t, 2, srva, srvb)

	// Drop foo
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, `"bar"`))
	checkExpectedSubs(t, 2, srvb)
	checkExpectedSubs(t, 1, srva)

	// Add it back
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, `["foo", "bar"]`))
	checkExpectedSubs(t, 2, srva, srvb)

	// Empty Export means implicit allow
	reloadUpdateConfig(t, srva, confA, `
		port: -1
		cluster {
			listen: 127.0.0.1:-1
			permissions {
				import: ">"
			}
		}
	`)
	checkExpectedSubs(t, 2, srva, srvb)

	confATemplate = `
		port: -1
		cluster {
			listen: 127.0.0.1:-1
			permissions {
				export: {
					allow: ["foo", "bar"]
					deny: %s
				}
			}
		}
	`
	// Now deny all:
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, `["foo", "bar"]`))
	checkExpectedSubs(t, 0, srva)
	checkExpectedSubs(t, 2, srvb)

	// Drop foo from the deny list
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, `"bar"`))
	checkExpectedSubs(t, 1, srva)
	checkExpectedSubs(t, 2, srvb)
}

func TestConfigReloadClusterPermsOldServer(t *testing.T) {
	confATemplate := `
		port: -1
		cluster {
			listen: 127.0.0.1:-1
			permissions {
				export: {
					allow: %s
				}
			}
		}
	`
	confA := createConfFile(t, []byte(fmt.Sprintf(confATemplate, `["foo", "bar"]`)))
	defer os.Remove(confA)
	srva, _ := RunServerWithConfig(confA)
	defer srva.Shutdown()

	optsB := DefaultOptions()
	optsB.Routes = RoutesFromStr(fmt.Sprintf("nats://127.0.0.1:%d", srva.ClusterAddr().Port))
	// Make server B behave like an old server
	testRouteProto = routeProtoZero
	defer func() { testRouteProto = routeProtoInfo }()
	srvb := RunServer(optsB)
	defer srvb.Shutdown()
	testRouteProto = routeProtoInfo

	checkClusterFormed(t, srva, srvb)

	// Get the route's connection ID
	getRouteRID := func() uint64 {
		rid := uint64(0)
		srvb.mu.Lock()
		for _, r := range srvb.routes {
			r.mu.Lock()
			rid = r.cid
			r.mu.Unlock()
			break
		}
		srvb.mu.Unlock()
		return rid
	}
	orgRID := getRouteRID()

	// Cause a config reload on A
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, `"bar"`))

	// Check that new route gets created
	check := func(t *testing.T) {
		t.Helper()
		checkFor(t, 3*time.Second, 15*time.Millisecond, func() error {
			if rid := getRouteRID(); rid == orgRID {
				return fmt.Errorf("Route does not seem to have been recreated")
			}
			return nil
		})
	}
	check(t)

	// Save the current value
	orgRID = getRouteRID()

	// Add another server that supports INFO updates

	optsC := DefaultOptions()
	optsC.Routes = RoutesFromStr(fmt.Sprintf("nats://127.0.0.1:%d", srva.ClusterAddr().Port))
	srvc := RunServer(optsC)
	defer srvc.Shutdown()

	checkClusterFormed(t, srva, srvb, srvc)

	// Cause a config reload on A
	reloadUpdateConfig(t, srva, confA, fmt.Sprintf(confATemplate, `"foo"`))
	// Check that new route gets created
	check(t)
}

func TestConfigReloadBoolFlags(t *testing.T) {
	logfile := "logtime.log"
	defer os.Remove(logfile)
	template := `
		listen: "127.0.0.1:-1"
		logfile: "logtime.log"
		%s
	`

	var opts *Options
	var err error

	for _, test := range []struct {
		name     string
		content  string
		cmdLine  []string
		expected bool
		val      func() bool
	}{
		// Logtime
		{
			"logtime_not_in_config_no_override",
			"",
			nil,
			true,
			func() bool { return opts.Logtime },
		},
		{
			"logtime_not_in_config_override_short_true",
			"",
			[]string{"-T"},
			true,
			func() bool { return opts.Logtime },
		},
		{
			"logtime_not_in_config_override_true",
			"",
			[]string{"-logtime"},
			true,
			func() bool { return opts.Logtime },
		},
		{
			"logtime_false_in_config_no_override",
			"logtime: false",
			nil,
			false,
			func() bool { return opts.Logtime },
		},
		{
			"logtime_false_in_config_override_short_true",
			"logtime: false",
			[]string{"-T"},
			true,
			func() bool { return opts.Logtime },
		},
		{
			"logtime_false_in_config_override_true",
			"logtime: false",
			[]string{"-logtime"},
			true,
			func() bool { return opts.Logtime },
		},
		{
			"logtime_true_in_config_no_override",
			"logtime: true",
			nil,
			true,
			func() bool { return opts.Logtime },
		},
		{
			"logtime_true_in_config_override_short_false",
			"logtime: true",
			[]string{"-T=false"},
			false,
			func() bool { return opts.Logtime },
		},
		{
			"logtime_true_in_config_override_false",
			"logtime: true",
			[]string{"-logtime=false"},
			false,
			func() bool { return opts.Logtime },
		},
		// Debug
		{
			"debug_not_in_config_no_override",
			"",
			nil,
			false,
			func() bool { return opts.Debug },
		},
		{
			"debug_not_in_config_override_short_true",
			"",
			[]string{"-D"},
			true,
			func() bool { return opts.Debug },
		},
		{
			"debug_not_in_config_override_true",
			"",
			[]string{"-debug"},
			true,
			func() bool { return opts.Debug },
		},
		{
			"debug_false_in_config_no_override",
			"debug: false",
			nil,
			false,
			func() bool { return opts.Debug },
		},
		{
			"debug_false_in_config_override_short_true",
			"debug: false",
			[]string{"-D"},
			true,
			func() bool { return opts.Debug },
		},
		{
			"debug_false_in_config_override_true",
			"debug: false",
			[]string{"-debug"},
			true,
			func() bool { return opts.Debug },
		},
		{
			"debug_true_in_config_no_override",
			"debug: true",
			nil,
			true,
			func() bool { return opts.Debug },
		},
		{
			"debug_true_in_config_override_short_false",
			"debug: true",
			[]string{"-D=false"},
			false,
			func() bool { return opts.Debug },
		},
		{
			"debug_true_in_config_override_false",
			"debug: true",
			[]string{"-debug=false"},
			false,
			func() bool { return opts.Debug },
		},
		// Trace
		{
			"trace_not_in_config_no_override",
			"",
			nil,
			false,
			func() bool { return opts.Trace },
		},
		{
			"trace_not_in_config_override_short_true",
			"",
			[]string{"-V"},
			true,
			func() bool { return opts.Trace },
		},
		{
			"trace_not_in_config_override_true",
			"",
			[]string{"-trace"},
			true,
			func() bool { return opts.Trace },
		},
		{
			"trace_false_in_config_no_override",
			"trace: false",
			nil,
			false,
			func() bool { return opts.Trace },
		},
		{
			"trace_false_in_config_override_short_true",
			"trace: false",
			[]string{"-V"},
			true,
			func() bool { return opts.Trace },
		},
		{
			"trace_false_in_config_override_true",
			"trace: false",
			[]string{"-trace"},
			true,
			func() bool { return opts.Trace },
		},
		{
			"trace_true_in_config_no_override",
			"trace: true",
			nil,
			true,
			func() bool { return opts.Trace },
		},
		{
			"trace_true_in_config_override_short_false",
			"trace: true",
			[]string{"-V=false"},
			false,
			func() bool { return opts.Trace },
		},
		{
			"trace_true_in_config_override_false",
			"trace: true",
			[]string{"-trace=false"},
			false,
			func() bool { return opts.Trace },
		},
		// Syslog
		{
			"syslog_not_in_config_no_override",
			"",
			nil,
			false,
			func() bool { return opts.Syslog },
		},
		{
			"syslog_not_in_config_override_short_true",
			"",
			[]string{"-s"},
			true,
			func() bool { return opts.Syslog },
		},
		{
			"syslog_not_in_config_override_true",
			"",
			[]string{"-syslog"},
			true,
			func() bool { return opts.Syslog },
		},
		{
			"syslog_false_in_config_no_override",
			"syslog: false",
			nil,
			false,
			func() bool { return opts.Syslog },
		},
		{
			"syslog_false_in_config_override_short_true",
			"syslog: false",
			[]string{"-s"},
			true,
			func() bool { return opts.Syslog },
		},
		{
			"syslog_false_in_config_override_true",
			"syslog: false",
			[]string{"-syslog"},
			true,
			func() bool { return opts.Syslog },
		},
		{
			"syslog_true_in_config_no_override",
			"syslog: true",
			nil,
			true,
			func() bool { return opts.Syslog },
		},
		{
			"syslog_true_in_config_override_short_false",
			"syslog: true",
			[]string{"-s=false"},
			false,
			func() bool { return opts.Syslog },
		},
		{
			"syslog_true_in_config_override_false",
			"syslog: true",
			[]string{"-syslog=false"},
			false,
			func() bool { return opts.Syslog },
		},
		// Cluster.NoAdvertise
		{
			"cluster_no_advertise_not_in_config_no_override",
			`cluster {
				port: -1
			}`,
			nil,
			false,
			func() bool { return opts.Cluster.NoAdvertise },
		},
		{
			"cluster_no_advertise_not_in_config_override_true",
			`cluster {
				port: -1
			}`,
			[]string{"-no_advertise"},
			true,
			func() bool { return opts.Cluster.NoAdvertise },
		},
		{
			"cluster_no_advertise_false_in_config_no_override",
			`cluster {
				port: -1
				no_advertise: false
			}`,
			nil,
			false,
			func() bool { return opts.Cluster.NoAdvertise },
		},
		{
			"cluster_no_advertise_false_in_config_override_true",
			`cluster {
				port: -1
				no_advertise: false
			}`,
			[]string{"-no_advertise"},
			true,
			func() bool { return opts.Cluster.NoAdvertise },
		},
		{
			"cluster_no_advertise_true_in_config_no_override",
			`cluster {
				port: -1
				no_advertise: true
			}`,
			nil,
			true,
			func() bool { return opts.Cluster.NoAdvertise },
		},
		{
			"cluster_no_advertise_true_in_config_override_false",
			`cluster {
				port: -1
				no_advertise: true
			}`,
			[]string{"-no_advertise=false"},
			false,
			func() bool { return opts.Syslog },
		},
		// -DV override
		{
			"debug_trace_not_in_config_dv_override_true",
			"",
			[]string{"-DV"},
			true,
			func() bool { return opts.Debug && opts.Trace },
		},
		{
			"debug_trace_false_in_config_dv_override_true",
			`debug: false
		     trace: false
			`,
			[]string{"-DV"},
			true,
			func() bool { return opts.Debug && opts.Trace },
		},
		{
			"debug_trace_true_in_config_dv_override_false",
			`debug: true
		     trace: true
			`,
			[]string{"-DV=false"},
			false,
			func() bool { return opts.Debug && opts.Trace },
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := createConfFile(t, []byte(fmt.Sprintf(template, test.content)))
			defer os.Remove(conf)

			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			var args []string
			args = append(args, "-c", conf)
			if test.cmdLine != nil {
				args = append(args, test.cmdLine...)
			}
			opts, err = ConfigureOptions(fs, args, nil, nil, nil)
			if err != nil {
				t.Fatalf("Error processing config: %v", err)
			}
			opts.NoSigs = true
			s := RunServer(opts)
			defer s.Shutdown()

			if test.val() != test.expected {
				t.Fatalf("Expected to be set to %v, got %v", test.expected, test.val())
			}
			if err := s.Reload(); err != nil {
				t.Fatalf("Error on reload: %v", err)
			}
			if test.val() != test.expected {
				t.Fatalf("Expected to be set to %v, got %v", test.expected, test.val())
			}
		})
	}
}
