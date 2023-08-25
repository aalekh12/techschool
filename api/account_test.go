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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/samplebank/db/mock"
	db "github.com/techschool/samplebank/db/sqlc"
	"github.com/techschool/samplebank/util"
)

func TestGetAccountapi(t *testing.T) {
	account := Randomaccount()

	testcases := []struct {
		name          string
		account_id    int64
		buildstubs    func(store *mockdb.MockStore)
		checkresponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "Ok",
			account_id: account.ID,
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Return(account, nil)
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBoydMatch(t, recorder.Body, account)
			},
		},
		{
			name:       "NotFound",
			account_id: account.ID,
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			account_id: account.ID,
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "InvalidID",
			account_id: 0,
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.account_id)
			reqest, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, reqest)
			log.Println("rb", recorder.Body)
			tc.checkresponse(t, recorder)
		})
	}

}

func Randomaccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomString(6),
		Balance:  util.GenerateMoney(),
		Currency: util.GenerateCurrency(),
	}
}

func requireBoydMatch(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotaccount db.Account
	err = json.Unmarshal(data, &gotaccount)
	log.Println("body", string(data))
	log.Println("ac", account)
	require.NoError(t, err)

	require.Equal(t, account, gotaccount)
}
