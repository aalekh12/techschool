package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/samplebank/db/mock"
	db "github.com/techschool/samplebank/db/sqlc"
	"github.com/techschool/samplebank/util"
)

type eqCreateUserMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util.ComparePassword(e.password, arg.HashedPassword)
	if err != nil {
		log.Println(err)
	}
	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func eqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserMatcher{arg, password}
}

func TestCreateAccountapi(t *testing.T) {
	user, password := RandomUserAccount()
	hashpassword, err := util.HashPassword(password)
	require.NoError(t, err)
	log.Println("pw", password)
	testcases := []struct {
		name          string
		body          gin.H
		buildstubs    func(store *mockdb.MockStore)
		checkresponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"username": user.Username,
				"email":    user.Email,
				"password": password,
				"fullname": user.FullName,
			},
			buildstubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:       user.Username,
					Email:          user.Email,
					FullName:       user.FullName,
					HashedPassword: hashpassword,
				}
				store.EXPECT().CreateUser(gomock.Any(), eqCreateUserParams(arg, password)).Times(1).Return(user, nil)
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireUserBoydMatch(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username,
				"email":    user.Email,
				"password": password,
				"fullname": user.FullName,
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username": user.Username,
				"email":    user.Email,
				"password": password,
				"fullname": user.FullName,
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		//add more cases
	}

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildstubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/user"
			reqest, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, reqest)
			log.Println("rb", recorder.Body)
			tc.checkresponse(t, recorder)
		})
	}

}

func RandomUserAccount() (db.User, string) {
	hashpassword, err := util.HashPassword(util.RandomString(8))
	if err != nil {
		log.Println(err)
	}
	return db.User{
		FullName: util.RandomString(8),
		Username: util.RandomString(6),
		Email:    util.RandomEmail(),
	}, hashpassword
}

func requireUserBoydMatch(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotuser db.User
	err = json.Unmarshal(data, &gotuser)
	log.Println("body", string(data))
	log.Println("ac", user)
	require.NoError(t, err)

	require.Equal(t, user, gotuser)
}
