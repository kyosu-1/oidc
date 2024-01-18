package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"slices"
	"time"
)

const (
	// 認可コードの長さ
	authorizeCodeLength = 32
)

// authorizaRequest is the request for the authorization endpoint
// ref: https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
// GETの際はquery parameter、POSTの際はform parameterとして送信される
type authorizeRequest struct {
	ResponseType string `form:"response_type" query:"response_type"`
	ClientID     string `form:"client_id" query:"client_id"`
	RedirectURI  string `form:"redirect_uri" query:"redirect_uri"`
	Scope        string `form:"scope" query:"scope"`
	State        string `form:"state" query:"state"`
}

// authorizeResponse is the response for the authorization endpoint
// ref: https://openid.net/specs/openid-connect-core-1_0.html#AuthResponse
type authorizeResponse struct {
	Code  string
	State string
}

// 認可コードに紐づく情報を雑に定義
// 有効期限や発行時刻など
type authorizeCode struct {
	Code        string
	RedirectURI string
	expiresAt   int64
	issuedAt    int64
}

// 簡易実装として、認可コードを生成し、それはトークンエンドポイントで利用するためにin memoryで保持する
// 本来ならば、認可コードはユーザーに紐づけたDBに保存する
var authorizeCodeMap = map[string]authorizeCode{}

func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// authorization endpoint
// ref: https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
func authorize(w http.ResponseWriter, r *http.Request) {
	// 本来ならば以下の処理を行う必要あり
	// - ユーザー認証
	// - response_typeによるフローの振り分け
	// - client_idの登録状態の確認
	// - redirect_uriとclient_idが示すクライアントとの対応確認
	// - scopeの確認
	// - 属性送出に関する同意画面の表示

	// redirect_uriとclient_idの検証
	// 本来ならば、redirect_uriはclient_idに紐づけたDBに保存されており、比較検証を行う
	// ここではパラメータが存在するかのみを確認する
	if r.FormValue("client_id") == "" {
		respondBadRequest(w, "client_id is required")
		return
	}

	if r.FormValue("redirect_uri") == "" {
		respondBadRequest(w, "redirect_uri is required")
		return
	}

	// scopeの検証
	// SupportedScopesに含まれているかだけ確認する
	if !slices.Contains(SupportedScopes, r.FormValue("scope")) {
		respondBadRequest(w, "unsupported scope")
		return
	}

	// response_typeの検証
	// SupportedResponseTypesに含まれているかだけ確認する
	if !slices.Contains(SupportedResponseTypes, r.FormValue("response_type")) {
		respondBadRequest(w, "unsupported response_type")
		return
	}

	code, err := generateRandomString(authorizeCodeLength)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	authorizeCodeMap[code] = authorizeCode{
		Code:        code,
		RedirectURI: r.FormValue("redirect_uri"),
		expiresAt:   time.Now().Add(10 * time.Minute).Unix(),
		issuedAt:    time.Now().Unix(),
	}

	// https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.2
	// redirect_uriにcodeとstateを付与してリダイレクト
	// stateはCSRF対策のために利用する
	// 本来ならば、redirect_uriはclient_idに紐づけたDBに保存する
	redirectURI := r.FormValue("redirect_uri")
	redirectURI += "?code=" + code
	redirectURI += "&state=" + r.FormValue("state")
	respondFound(w, redirectURI)
}
