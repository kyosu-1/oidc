package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// tokenRequest はトークンリクエストを定義します
type tokenRequest struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
	ClientID    string `json:"client_id"`
}

// ref: https://openid.net/specs/openid-connect-core-1_0.html#IDToken
type idTokenPayload struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
}

// トークンレスポンスにIDトークンを追加します
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token"`
}

// tokenEndpoint はトークンエンドポイントのハンドラです
func token(w http.ResponseWriter, r *http.Request) {
	// リクエストボディをパースします
	var req tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// grant_typeがauthorization_codeであることを確認
	if req.GrantType != "authorization_code" {
		http.Error(w, "Unsupported grant type", http.StatusBadRequest)
		return
	}

	// 認可コードを検証
	authCode, ok := authorizeCodeMap[req.Code]
	if !ok {
		http.Error(w, "Invalid code", http.StatusBadRequest)
		return
	}

	// 認可コードに紐づくクライアントIDが一致することを確認
	if authCode.ClientID != req.ClientID {
		http.Error(w, "Invalid client", http.StatusBadRequest)
		return
	}

	// アクセストークンを生成（簡易実装ではランダム文字列を使用）
	accessToken, _ := generateRandomString(32)

	now := time.Now()
	idToken := idTokenPayload{
		Iss: "https://your-issuer.com", // 発行者のURL
		Sub: "1234567890",              // 本来であれば
		Aud: req.ClientID,              // クライアントID
		Exp: now.Add(time.Hour).Unix(), // 有効期限（例: 1時間後）
		Iat: now.Unix(),                // 発行時刻
	}

	// IDトークンをJSON文字列にエンコードします
	// 本来はJWTとして署名する必要があります
	idTokenStr, err := json.Marshal(idToken)
	if err != nil {
		http.Error(w, "Failed to generate ID token", http.StatusInternalServerError)
		return
	}

	// トークンレスポンスにIDトークンを追加します
	resp := tokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		IDToken:     string(idTokenStr),
	}
	respondJSON(w, http.StatusOK, resp)
}
