package ig

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_buildURL(t *testing.T) {
	type args struct {
		base        string
		token       string
		path        string
		extraParams []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "build correct url",
			want: "https://test.com/me?access_token=TOKEN_TEST",
			args: args{
				base:  "https://test.com",
				token: "TOKEN_TEST",
				path:  "/me",
			},
		},
		{
			name: "build correct url with extra params",
			want: "https://test.com/me?access_token=TOKEN_TEST&name=test&age=12",
			args: args{
				base:        "https://test.com",
				token:       "TOKEN_TEST",
				path:        "/me",
				extraParams: []string{"name=test", "age=12"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildURL(tt.args.base, tt.args.path, tt.args.token, tt.args.extraParams...); got != tt.want {
				t.Errorf("buildURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func fakeServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestClient_GetProfile(t *testing.T) {
	var (
		s1 = fakeServer(func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprint(rw, `{"id": "123", "username": "dan", "account_type": "PERSONAL"}`)
		})
		s2 = fakeServer(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})
	)

	type fields struct {
		accessToken string
		client      *http.Client
		baseURL     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Profile
		wantErr bool
	}{
		{
			name: "calls /me with correct values",
			want: &Profile{ID: "123", Username: "dan", AccountType: "PERSONAL"},
			fields: fields{
				accessToken: "TEST_TOKEN",
				baseURL:     s1.URL,
				client:      s1.Client(),
			},
		},
		{
			name:    "calls /me, respond with err",
			wantErr: true,
			fields: fields{
				accessToken: "TEST_TOKEN",
				baseURL:     s2.URL,
				client:      s2.Client(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				accessToken: tt.fields.accessToken,
				client:      tt.fields.client,
				baseURL:     tt.fields.baseURL,
			}

			got, err := c.GetProfile()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProfile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetMedia(t *testing.T) {
	var (
		c1 = fakeServer(func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprint(rw, `{"data": [
					{
						"caption": "CAPTION_1",
						"id": "TEST_ID_1",
						"thumbnail_url": "TEST_URL_1",
						"username": "USERNAME_1",
						"permalink": "TEST_URL_1",
						"media_url": "TEST_URL_1",
						"media_type": "TYPE_TEST",
						"timestamp": "TIMESTAMP_1"
					},
					{
						"caption": "CAPTION_2",
						"id": "TEST_ID_2",
						"thumbnail_url": "TEST_URL_2",
						"username": "USERNAME_2",
						"permalink": "TEST_URL_2",
						"media_url": "TEST_URL_2",
						"media_type": "TYPE_TEST",
						"timestamp": "TIMESTAMP_2"
					},
					{
						"caption": "CAPTION_3",
						"id": "TEST_ID_3",
						"thumbnail_url": "TEST_URL_3",
						"username": "USERNAME_3",
						"permalink": "TEST_URL_3",
						"media_url": "TEST_URL_3",
						"media_type": "TYPE_TEST",
						"timestamp": "TIMESTAMP_3"
					}
				]}`,
			)
		})
		c2 = fakeServer(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(500)
		})
	)
	type fields struct {
		accessToken string
		client      *http.Client
		baseURL     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Entry
		wantErr bool
	}{
		{
			name:   "fetches correct media",
			fields: fields{accessToken: "TOKEN_TEST", client: c1.Client(), baseURL: c1.URL},
			want: []Entry{
				{
					Caption:      "CAPTION_1",
					ID:           "TEST_ID_1",
					ThumbnailURL: "TEST_URL_1",
					Username:     "USERNAME_1",
					Permalink:    "TEST_URL_1",
					MediaURL:     "TEST_URL_1",
					MediaType:    "TYPE_TEST",
					Timestamp:    "TIMESTAMP_1",
				},
				{
					Caption:      "CAPTION_2",
					ID:           "TEST_ID_2",
					ThumbnailURL: "TEST_URL_2",
					Username:     "USERNAME_2",
					Permalink:    "TEST_URL_2",
					MediaURL:     "TEST_URL_2",
					MediaType:    "TYPE_TEST",
					Timestamp:    "TIMESTAMP_2",
				},
				{
					Caption:      "CAPTION_3",
					ID:           "TEST_ID_3",
					ThumbnailURL: "TEST_URL_3",
					Username:     "USERNAME_3",
					Permalink:    "TEST_URL_3",
					MediaURL:     "TEST_URL_3",
					MediaType:    "TYPE_TEST",
					Timestamp:    "TIMESTAMP_3",
				},
			},
		},
		{
			name:    "handles error",
			fields:  fields{accessToken: "TOKEN_TEST", client: c2.Client(), baseURL: c2.URL},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				accessToken: tt.fields.accessToken,
				client:      tt.fields.client,
				baseURL:     tt.fields.baseURL,
			}
			got, err := c.GetMedia()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMedia() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_Tags(t *testing.T) {

	tests := []struct {
		name    string
		caption string
		want    []string
	}{
		{
			name:    "returns empty if no hash found",
			caption: "there's not hashes here 🦗",
		},
		{
			name:    "splits by space and identifies # chars",
			caption: "Pool time 🧘🏼‍♀️🏊🏽‍♂️ #lizanddan #belize #devacaciones #caribemylove",
			want:    []string{"lizanddan", "belize", "devacaciones", "caribemylove"},
		},
		{
			name:    "splits by space and identifies # chars, no matter if no space between",
			caption: "Pool time 🧘🏼‍♀️🏊🏽‍♂️ #lizanddan#belize#devacaciones#caribemylove",
			want:    []string{"lizanddan", "belize", "devacaciones", "caribemylove"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Entry{Caption: tt.caption}
			if got := e.Tags(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_LongLivedToken(t *testing.T) {

	var (
		s1 = fakeServer(func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprint(rw, `{"access_token":"TEST_TOKEN","token_type":"bearer","expires_in":5183944}`)
		})
		f2 = func(rw http.ResponseWriter, r *http.Request) { rw.WriteHeader(403) }
		s2 = fakeServer(f2)
	)

	type fields struct {
		accessToken string
		client      *http.Client
		baseURL     string
	}
	type args struct {
		secret string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *LongLiveToken
		wantErr bool
	}{
		{
			name: "returns token from request",
			args: args{secret: "SECRET_TEST"},
			fields: fields{
				accessToken: "TEST_TOKEN",
				client:      s1.Client(),
				baseURL:     s1.URL,
			},
			want: &LongLiveToken{AccessToken: "TEST_TOKEN", TokenType: "bearer", ExpiresIn: 5183944},
		},
		{
			name: "handles error",
			args: args{secret: "SECRET_TEST"},
			fields: fields{
				accessToken: "TEST_TOKEN",
				client:      s2.Client(),
				baseURL:     s2.URL,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				accessToken: tt.fields.accessToken,
				client:      tt.fields.client,
				baseURL:     tt.fields.baseURL,
			}
			got, err := c.LongLivedToken(tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("LongLivedToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LongLivedToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
