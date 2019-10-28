package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/blocktop/mp-common/config"
	"github.com/blocktop/mp-common/server"
	"github.com/dgrijalva/jwt-go"
	"github.com/stellar/go/txnbuild"
	"net/http"
	"net/url"
	"time"
)

func handleGetToken(w http.ResponseWriter, r *http.Request) {
	account := r.URL.Query().Get("account")
	if len(account) == 0 {
		server.ResponseError(w, http.StatusBadRequest, server.BADACCT, fmt.Errorf("invalid account in query"))
		return
	}
	cfg := config.GetConfig()

	txn, err := txnbuild.BuildChallengeTx(cfg.SigningKeySeed, account, cfg.AnchorName, cfg.NetworkPassphrase, 5*time.Minute)
	if err != nil {
		server.ResponseError(w, http.StatusInternalServerError, server.CHLFAIL, err)
		return
	}

	data := map[string]interface{}{
		"transaction":        txn,
		"network_passphrase": cfg.NetworkPassphrase,
	}
	server.ResponseJSONMap(w, data)
}

func handlePostToken(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	var xdr string
	var err error
	switch contentType {
	case "application/x-www-form-urlencoded":
		xdr, err = url.QueryUnescape(r.FormValue("transaction"))
		if err != nil {
			server.ResponseError(w, http.StatusBadRequest, server.BADTXN, fmt.Errorf("could not decode transaction"))
			return
		}
	case "application/json":
		body := make(map[string]string)
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			server.ResponseError(w, http.StatusBadRequest, server.BADTXN, fmt.Errorf("could not decode transaction"))
			return
		}
		xdr = body["transaction"]
	default:
		server.ResponseError(w, http.StatusBadRequest, server.CONTYP, fmt.Errorf("invalid content-type"))
		return
	}

	if len(xdr) == 0 {
		server.ResponseError(w, http.StatusBadRequest, server.BADTXN, fmt.Errorf("no transaction data found"))
		return
	}

	cfg := config.GetConfig()

	valid, err := txnbuild.VerifyChallengeTx(xdr, cfg.SigningKey, cfg.NetworkPassphrase)
	if err != nil {
		server.ResponseError(w, http.StatusBadRequest, server.BADTXN, err)
		return
	}

	if !valid {
		server.ResponseError(w, http.StatusBadRequest, server.BADTXN, fmt.Errorf("challenge transaction invalid"))
		return
	}

	// get subject account ID, no need to check errors because the VerifyChallengeTx above would have caught them
	txn, _ := txnbuild.TransactionFromXDR(xdr)
	op, _ := txn.Operations[0].(*txnbuild.ManageData)
	subject := op.SourceAccount.GetAccountID()

	now := time.Now().Unix()
	var duration time.Duration
	if cfg.IsProduction() {
		// 24 hours for prod
		duration = 24*time.Hour/time.Second
	} else {
		// 5 days for non prod for convenience in testing
		duration = 5*24*time.Hour/time.Second
	}

	tokenClaims := &jwt.StandardClaims{
		ExpiresAt: now + int64(duration),
		Id:        fmt.Sprintf("%x", sha256.Sum256([]byte(xdr))),
		IssuedAt:  now,
		Issuer:    cfg.WebAuthEndpoint,
		Subject:   subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	sToken, err := token.SignedString([]byte(cfg.SigningKeySeed))
	if err != nil {
		server.ResponseError(w, http.StatusInternalServerError, server.JWTSIG, err)
		return
	}
	server.ResponseJSONMap(w, map[string]interface{}{"token": sToken})
}
