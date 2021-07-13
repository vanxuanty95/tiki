package usecase

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
	"tiki/internal/api/booking/storages/mock"
	mock2 "tiki/internal/api/user/storages/mock"
	"tiki/internal/api/utils"

	"tiki/internal/api/booking/storages"
	userStorages "tiki/internal/api/user/storages"
)

func TestBooking_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storage := mock.NewMockStore(ctrl)

	type fields struct {
		Store storages.Store
	}
	type args struct {
		ctx         context.Context
		userID      string
		createdDate string
		page        int
		limit       int
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		mockStoreList func()
		want          []*storages.Booking
		wantErr       bool
	}{
		{
			name:   "empty userID",
			fields: fields{Store: storage},
			args: args{
				ctx:         nil,
				userID:      "",
				createdDate: "2021-07-12",
				page:        0,
				limit:       0,
			},
			mockStoreList: func() {
				storage.EXPECT().RetrieveBookings(gomock.Any(), "", "2021-07-12", 0, 0).
					Return(nil, errors.New("storage got error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "empty createDate",
			fields: fields{Store: storage},
			args: args{
				ctx:         nil,
				userID:      "test1",
				createdDate: "",
				page:        0,
				limit:       0,
			},
			mockStoreList: func() {
				storage.EXPECT().RetrieveBookings(gomock.Any(), "test1", "", 0, 0).
					Return(nil, errors.New("storage got error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "empty userID and createDate",
			fields: fields{Store: storage},
			args: args{
				ctx:         nil,
				userID:      "",
				createdDate: "",
				page:        0,
				limit:       0,
			},
			mockStoreList: func() {
				storage.EXPECT().RetrieveBookings(gomock.Any(), "", "", 0, 0).
					Return(nil, errors.New("storage got error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "negative page and limit",
			fields: fields{Store: storage},
			args: args{
				ctx:         nil,
				userID:      "test1",
				createdDate: "2021-07-12",
				page:        -1,
				limit:       -1,
			},
			mockStoreList: func() {
				storage.EXPECT().RetrieveBookings(gomock.Any(), "test1", "2021-07-12", -1, -1).
					Return([]*storages.Booking{{
						ID:          "1",
						Content:     "a",
						UserID:      "test1",
						CreatedDate: "2021-07-12",
					}, {
						ID:          "2",
						Content:     "b",
						UserID:      "test1",
						CreatedDate: "2021-07-12",
					}, {
						ID:          "3",
						Content:     "c",
						UserID:      "test1",
						CreatedDate: "2021-07-12",
					}}, nil)
			},
			want: []*storages.Booking{{
				ID:          "1",
				Content:     "a",
				UserID:      "test1",
				CreatedDate: "2021-07-12",
			}, {
				ID:          "2",
				Content:     "b",
				UserID:      "test1",
				CreatedDate: "2021-07-12",
			}, {
				ID:          "3",
				Content:     "c",
				UserID:      "test1",
				CreatedDate: "2021-07-12",
			}},
			wantErr: false,
		},
		{
			name:   "normal case",
			fields: fields{Store: storage},
			args: args{
				ctx:         nil,
				userID:      "test1",
				createdDate: "2021-07-12",
				page:        0,
				limit:       0,
			},
			mockStoreList: func() {
				storage.EXPECT().RetrieveBookings(gomock.Any(), "test1", "2021-07-12", 0, 0).
					Return([]*storages.Booking{{
						ID:          "1",
						Content:     "a",
						UserID:      "test1",
						CreatedDate: "2021-07-12",
					}, {
						ID:          "2",
						Content:     "b",
						UserID:      "test1",
						CreatedDate: "2021-07-12",
					}, {
						ID:          "3",
						Content:     "c",
						UserID:      "test1",
						CreatedDate: "2021-07-12",
					}}, nil)
			},
			want: []*storages.Booking{{
				ID:          "1",
				Content:     "a",
				UserID:      "test1",
				CreatedDate: "2021-07-12",
			}, {
				ID:          "2",
				Content:     "b",
				UserID:      "test1",
				CreatedDate: "2021-07-12",
			}, {
				ID:          "3",
				Content:     "c",
				UserID:      "test1",
				CreatedDate: "2021-07-12",
			}},
			wantErr: false,
		},
		{
			name:   "set limit and page case",
			fields: fields{Store: storage},
			args: args{
				ctx:         nil,
				userID:      "test1",
				createdDate: "2021-07-12",
				page:        0,
				limit:       1,
			},
			mockStoreList: func() {
				storage.EXPECT().RetrieveBookings(gomock.Any(), "test1", "2021-07-12", 0, 1).
					Return([]*storages.Booking{{
						ID:          "1",
						Content:     "a",
						UserID:      "test1",
						CreatedDate: "2021-07-12",
					}}, nil)
			},
			want: []*storages.Booking{{
				ID:          "1",
				Content:     "a",
				UserID:      "test1",
				CreatedDate: "2021-07-12",
			}},
			wantErr: false,
		},
		{
			name:   "storage got error",
			fields: fields{Store: storage},
			args: args{
				ctx:         nil,
				userID:      "test1",
				createdDate: "2021-07-12",
				page:        0,
				limit:       0,
			},
			mockStoreList: func() {
				storage.EXPECT().RetrieveBookings(gomock.Any(), "test1", "2021-07-12", 0, 0).
					Return(nil, errors.New("storge got error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockStoreList()
			s := &Booking{
				Store: tt.fields.Store,
			}
			got, err := s.List(tt.args.ctx, tt.args.userID, tt.args.createdDate, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Booking.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Booking.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBooking_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookingStorage := mock.NewMockStore(ctrl)
	userStorage := mock2.NewMockStore(ctrl)

	timeNow := utils.GetTimeNowWithDefaultLayoutInString()

	type fields struct {
		Store     storages.Store
		UserStore userStorages.Store
		GenFn     utils.GenerateNewUUIDFn
	}
	type args struct {
		ctx     context.Context
		userID  string
		booking *storages.Booking
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		mockStoreGetUser func()
		mockStoreList    func()
		mockStoreAdd     func()
		want             *storages.Booking
		wantErr          bool
	}{
		{
			name: "user is empty",
			fields: fields{
				Store:     bookingStorage,
				UserStore: userStorage,
				GenFn: func() string {
					return "id_test"
				},
			},
			args: args{
				ctx:     nil,
				userID:  "",
				booking: &storages.Booking{Content: "abc"},
			},
			mockStoreGetUser: func() {
				userStorage.EXPECT().Get(gomock.Any(), "").
					Return(nil, errors.New("storage got error"))
			},
			mockStoreList: nil,
			mockStoreAdd:  nil,
			want:          nil,
			wantErr:       true,
		},
		{
			name: "booking is empty",
			fields: fields{
				Store:     bookingStorage,
				UserStore: userStorage,
				GenFn: func() string {
					return "id_test"
				},
			},
			args: args{
				ctx:     nil,
				userID:  "test1",
				booking: &storages.Booking{},
			},
			mockStoreGetUser: func() {
				userStorage.EXPECT().Get(gomock.Any(), "test1").
					Return(&userStorages.User{
						ID:       "test1",
						Password: "$2a$14$BdgOuNVBU7sdGW9rIDIIv.MWXDdvTVKyTppb3qW03bmvz/6hhA1FO",
						MaxTodo:  5,
					}, nil)
			},
			mockStoreList: func() {
				bookingStorage.EXPECT().RetrieveBookings(gomock.Any(), "test1", timeNow, 0, 6).
					Return([]*storages.Booking{{
						ID:          "1",
						Content:     "a",
						UserID:      "test1",
						CreatedDate: timeNow,
					}, {
						ID:          "2",
						Content:     "b",
						UserID:      "test1",
						CreatedDate: timeNow,
					}, {
						ID:          "3",
						Content:     "c",
						UserID:      "test1",
						CreatedDate: timeNow,
					}}, nil)
			},
			mockStoreAdd: func() {
				bookingStorage.EXPECT().AddBooking(gomock.Any(), &storages.Booking{
					ID:          "id_test",
					Content:     "",
					UserID:      "test1",
					CreatedDate: timeNow,
				}).
					Return(errors.New("storage got error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "booking is empty",
			fields: fields{
				Store:     bookingStorage,
				UserStore: userStorage,
				GenFn: func() string {
					return "id_test"
				},
			},
			args: args{
				ctx:     nil,
				userID:  "test1",
				booking: &storages.Booking{},
			},
			mockStoreGetUser: func() {
				userStorage.EXPECT().Get(gomock.Any(), "test1").
					Return(&userStorages.User{
						ID:       "test1",
						Password: "$2a$14$BdgOuNVBU7sdGW9rIDIIv.MWXDdvTVKyTppb3qW03bmvz/6hhA1FO",
						MaxTodo:  5,
					}, nil)
			},
			mockStoreList: func() {
				bookingStorage.EXPECT().RetrieveBookings(gomock.Any(), "test1", timeNow, 0, 6).
					Return([]*storages.Booking{{
						ID:          "1",
						Content:     "a",
						UserID:      "test1",
						CreatedDate: timeNow,
					}, {
						ID:          "2",
						Content:     "b",
						UserID:      "test1",
						CreatedDate: timeNow,
					}, {
						ID:          "3",
						Content:     "c",
						UserID:      "test1",
						CreatedDate: timeNow,
					}}, nil)
			},
			mockStoreAdd: func() {
				bookingStorage.EXPECT().AddBooking(gomock.Any(), &storages.Booking{
					ID:          "id_test",
					Content:     "",
					UserID:      "test1",
					CreatedDate: timeNow,
				}).
					Return(errors.New("storage got error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "list booking error",
			fields: fields{
				Store:     bookingStorage,
				UserStore: userStorage,
				GenFn: func() string {
					return "id_test"
				},
			},
			args: args{
				ctx:     nil,
				userID:  "test1",
				booking: &storages.Booking{},
			},
			mockStoreGetUser: func() {
				userStorage.EXPECT().Get(gomock.Any(), "test1").
					Return(&userStorages.User{
						ID:       "test1",
						Password: "$2a$14$BdgOuNVBU7sdGW9rIDIIv.MWXDdvTVKyTppb3qW03bmvz/6hhA1FO",
						MaxTodo:  5,
					}, nil)
			},
			mockStoreList: func() {
				bookingStorage.EXPECT().RetrieveBookings(gomock.Any(), "test1", timeNow, 0, 6).
					Return(nil, errors.New("storage got error"))
			},
			mockStoreAdd: nil,
			want:         nil,
			wantErr:      true,
		},
		{
			name: "reach limit",
			fields: fields{
				Store:     bookingStorage,
				UserStore: userStorage,
				GenFn: func() string {
					return "id_test"
				},
			},
			args: args{
				ctx:     nil,
				userID:  "test1",
				booking: &storages.Booking{},
			},
			mockStoreGetUser: func() {
				userStorage.EXPECT().Get(gomock.Any(), "test1").
					Return(&userStorages.User{
						ID:       "test1",
						Password: "$2a$14$BdgOuNVBU7sdGW9rIDIIv.MWXDdvTVKyTppb3qW03bmvz/6hhA1FO",
						MaxTodo:  5,
					}, nil)
			},
			mockStoreList: func() {
				bookingStorage.EXPECT().RetrieveBookings(gomock.Any(), "test1", timeNow, 0, 6).
					Return([]*storages.Booking{{
						ID:          "1",
						Content:     "a",
						UserID:      "test1",
						CreatedDate: timeNow,
					}, {
						ID:          "2",
						Content:     "b",
						UserID:      "test1",
						CreatedDate: timeNow,
					}, {
						ID:          "3",
						Content:     "c",
						UserID:      "test1",
						CreatedDate: timeNow,
					}, {
						ID:          "4",
						Content:     "c",
						UserID:      "test1",
						CreatedDate: timeNow,
					}, {
						ID:          "5",
						Content:     "c",
						UserID:      "test1",
						CreatedDate: timeNow,
					}}, nil)
			},
			mockStoreAdd: nil,
			want:         nil,
			wantErr:      true,
		},
		{
			name: "normal case",
			fields: fields{
				Store:     bookingStorage,
				UserStore: userStorage,
				GenFn: func() string {
					return "id_test"
				},
			},
			args: args{
				ctx:     nil,
				userID:  "test1",
				booking: &storages.Booking{Content: "abc"},
			},
			mockStoreGetUser: func() {
				userStorage.EXPECT().Get(gomock.Any(), "test1").
					Return(&userStorages.User{
						ID:       "test1",
						Password: "$2a$14$BdgOuNVBU7sdGW9rIDIIv.MWXDdvTVKyTppb3qW03bmvz/6hhA1FO",
						MaxTodo:  5,
					}, nil)
			},
			mockStoreList: func() {
				bookingStorage.EXPECT().RetrieveBookings(gomock.Any(), "test1", timeNow, 0, 6).
					Return([]*storages.Booking{{
						ID:          "1",
						Content:     "a",
						UserID:      "test1",
						CreatedDate: timeNow,
					}, {
						ID:          "2",
						Content:     "b",
						UserID:      "test1",
						CreatedDate: timeNow,
					}, {
						ID:          "3",
						Content:     "c",
						UserID:      "test1",
						CreatedDate: timeNow,
					}, {
						ID:          "4",
						Content:     "c",
						UserID:      "test1",
						CreatedDate: timeNow,
					}}, nil)
			},
			mockStoreAdd: func() {
				bookingStorage.EXPECT().AddBooking(gomock.Any(), &storages.Booking{
					ID:          "id_test",
					Content:     "abc",
					UserID:      "test1",
					CreatedDate: timeNow,
				}).
					Return(nil)
			},
			want: &storages.Booking{
				ID:          "id_test",
				Content:     "abc",
				UserID:      "test1",
				CreatedDate: timeNow,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if tt.mockStoreGetUser != nil {
			tt.mockStoreGetUser()
		}
		if tt.mockStoreList != nil {
			tt.mockStoreList()
		}
		if tt.mockStoreAdd != nil {
			tt.mockStoreAdd()
		}

		t.Run(tt.name, func(t *testing.T) {
			s := &Booking{
				Store:           tt.fields.Store,
				UserStore:       tt.fields.UserStore,
				GeneratorUUIDFn: tt.fields.GenFn,
			}
			got, err := s.Add(tt.args.ctx, tt.args.userID, tt.args.booking)
			if (err != nil) != tt.wantErr {
				t.Errorf("Booking.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				if got.Content != tt.want.Content && got.CreatedDate != tt.want.CreatedDate && got.UserID != tt.want.UserID {
					t.Errorf("Booking.Add() = %v, want %v", got, tt.want)
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Booking.Add() = %v, want %v", got, tt.want)
			}
		})

	}
}
