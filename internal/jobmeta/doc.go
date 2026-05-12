// Package jobmeta provides per-job arbitrary metadata storage for cronwrap.
//
// Jobs can carry any number of key-value string pairs (e.g. owner, team,
// environment, cost-centre) that are visible through the HTTP API and can
// be used by alerting or auditing integrations.
//
// # Usage
//
//	store := jobmeta.New()
//	_ = store.Register("nightly-backup")
//	_ = store.Set("nightly-backup", "owner", "ops-team")
//	_ = store.Set("nightly-backup", "env",   "production")
//
//	entry, ok := store.Get("nightly-backup")
//
// # HTTP API
//
// Mount the handler under any prefix:
//
//	mux.Handle("/jobmeta/", jobmeta.Handler(store))
//
// Supported routes:
//
//	GET    /jobmeta            – list all entries
//	GET    /jobmeta/{name}     – get a single job's metadata
//	POST   /jobmeta/{name}     – set a key/value pair
//	DELETE /jobmeta/{name}/{k} – remove a single key
package jobmeta
