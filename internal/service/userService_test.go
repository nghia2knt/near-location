package service

import (
	"context"
	"errors"
	"fmt"
	"near-location/internal/model"
	"near-location/pkg/config"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	testCase struct {
		name        string
		args        arguments
		mockHandler func(*testing.T, arguments) (UserService, error)
		exp         expect
	}
	arguments struct {
		findNear findUserLocationsNearDatapointRequest
	}
	expect struct {
		err error
	}
)

type findUserLocationsNearDatapointRequest struct {
	datapoint   model.Datapoint
	maxDistance int64
	limit       int64
	skip        int64
}

var testCases []testCase

func init() {
	configPath := "../../configs"
	state := "test"
	if err := config.InitConfig(configPath, state); err != nil {
		fmt.Println("failed to init config: ", err)
		panic(err)
	}
	loadTestCase()
}

func loadTestCase() {
	testCases = []testCase{
		{
			name: "success case",
			args: arguments{
				findNear: findUserLocationsNearDatapointRequest{
					datapoint: model.Datapoint{
						Longitude: 10,
						Latitude:  10,
					},
					maxDistance: 1000000,
					limit:       5,
					skip:        0,
				},
			},
			mockHandler: func(t *testing.T, a arguments) (UserService, error) {
				ctx := context.Background()
				ctl := gomock.NewController(t)
				repo := NewMockUserRepository(ctl)
				uuid, _ := uuid.NewUUID()
				findNearReturn := []model.UserLocation{
					{
						Id:        primitive.NewObjectID(),
						UpdatedAt: time.Now(),
						UserId:    uuid.String(),
						Location: model.GeoJSON{
							Coordinates: []float64{10.000001, 10.000001},
						},
					},
				}
				repo.EXPECT().FindNearUserLocation(ctx, a.findNear.datapoint, a.findNear.maxDistance, a.findNear.limit, a.findNear.skip).Return(findNearReturn, nil)
				return NewUserService(repo), nil
			},
			exp: expect{
				err: nil,
			},
		},

		{
			name: "success case, but reset limit over",
			args: arguments{
				findNear: findUserLocationsNearDatapointRequest{
					datapoint: model.Datapoint{
						Longitude: 10,
						Latitude:  10,
					},
					maxDistance: 1000000,
					limit:       5000,
					skip:        0,
				},
			},
			mockHandler: func(t *testing.T, a arguments) (UserService, error) {
				ctx := context.Background()
				ctl := gomock.NewController(t)
				repo := NewMockUserRepository(ctl)
				uuid, _ := uuid.NewUUID()
				findNearReturn := []model.UserLocation{
					{
						Id:        primitive.NewObjectID(),
						UpdatedAt: time.Now(),
						UserId:    uuid.String(),
						Location: model.GeoJSON{
							Coordinates: []float64{10.000001, 10.000001},
						},
					},
				}
				repo.EXPECT().FindNearUserLocation(ctx, a.findNear.datapoint, a.findNear.maxDistance, int64(config.CV.MongoConfig.QueryMaxLimit), a.findNear.skip).Return(findNearReturn, nil)
				return NewUserService(repo), nil
			},
			exp: expect{
				err: nil,
			},
		},

		{
			name: "failed in database",
			args: arguments{
				findNear: findUserLocationsNearDatapointRequest{
					datapoint: model.Datapoint{
						Longitude: 10,
						Latitude:  10,
					},
					maxDistance: 1000000,
					limit:       5,
					skip:        0,
				},
			},
			mockHandler: func(t *testing.T, a arguments) (UserService, error) {
				ctx := context.Background()
				ctl := gomock.NewController(t)
				repo := NewMockUserRepository(ctl)
				repo.EXPECT().FindNearUserLocation(ctx, a.findNear.datapoint, a.findNear.maxDistance, a.findNear.limit, a.findNear.skip).Return(nil, errors.New("failed in database"))
				return NewUserService(repo), nil
			},
			exp: expect{
				err: errors.New("failed in database"),
			},
		},
	}
}

func TestFindUserLocationsNearDatapoint(t *testing.T) {
	ctx := context.Background()
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			handler, err := test.mockHandler(t, test.args)
			if err != nil && test.exp.err != nil {
				if !strings.Contains(err.Error(), test.exp.err.Error()) {
					t.Fatalf("FindUserLocationsNearDatapoint failed on test case [%s]. Expected error message contains [%s], but got [%s]", test.name, test.exp.err.Error(), err.Error())
				}
				return
			}
			_, err = handler.FindUserLocationsNearDatapoint(
				ctx,
				test.args.findNear.datapoint,
				test.args.findNear.maxDistance,
				test.args.findNear.limit,
				test.args.findNear.skip)
			if err == nil && test.exp.err == nil {
				return
			}
			if err != nil && test.exp.err == nil {
				t.Errorf("FindUserLocationsNearDatapoint failed on test case [%s]. Expected no error, but got [%s]", test.name, err.Error())
				t.Fail()
				return
			}
			if err == nil && test.exp.err != nil {
				t.Errorf("FindUserLocationsNearDatapoint failed on test case [%s]. Expected error [%s], but got no error", test.name, test.exp.err.Error())
				t.Fail()
				return
			}
			if !strings.Contains(err.Error(), test.exp.err.Error()) {
				t.Errorf("FindUserLocationsNearDatapoint failed on test case [%s]. Expected error message contains [%s], but got [%s]", test.name, test.exp.err.Error(), err.Error())
				t.Fail()
			}

		})
	}
}
