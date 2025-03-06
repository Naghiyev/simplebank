package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	mockdb "simple-banking/db/mock"
	db "simple-banking/db/sqlc"
	"simple-banking/util"
	"testing"
)

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()

	//test cases
	testCases := []struct {
		name          string
		accountId     int64
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:      "Account found is OK ",
			accountId: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				//build stubs
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "Account Not Found  ",
			accountId: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				//build stubs
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "AccountById Bad Request",
			accountId: 0,
			buildStub: func(store *mockdb.MockStore) {
				//build stubs
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "AccountById Internal Server Error ",
			accountId: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				//build stubs
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			//start test http server and send request
			server := NewServer(store)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountId)

			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			//check response
			tc.checkResponse(t, recorder)
		})

	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       int64(util.RandomInt(1, 1000)),
		Owner:    util.RandomOwnerName(),
		Balance:  int64(util.RandomMoney()),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var returnedAccount db.Account
	err = json.Unmarshal(data, &returnedAccount)
	require.NoError(t, err)
	require.Equal(t, returnedAccount, account)

}
