package usecase

import (
	"context"
	"errors"
	"math/rand"
	"tiki/cmd/api/config"
	"tiki/internal/api/booking/storages"
	"tiki/internal/api/booking/storages/model"
	"tiki/internal/api/dictionary"
	userStorages "tiki/internal/api/user/storages"
	"tiki/internal/api/utils"
	"tiki/internal/pkg/logger"
	"time"
)

type Booking struct {
	Cfg             *config.Config
	Store           storages.Store
	UserStore       userStorages.Store
	GeneratorUUIDFn utils.GenerateNewUUIDFn
}

type Ring struct {
	maxRow    int
	minRow    int
	maxColumn int
	minColumn int
}

type SeatWithoutID struct {
	row    int
	column int
}

func (s *Booking) CreateNewScreen(ctx context.Context, userID string, screen *storages.Screen) (*int, error) {
	now := utils.GetTimeNow()

	screen.UserID = userID
	screen.CreatedDate = now

	if id, err := s.Store.InsertScreen(ctx, screen); err != nil {
		logger.TIKIErrorln(ctx, err)
		return nil, errors.New(dictionary.StoreScreenFailed)
	} else {
		return id, nil
	}
}

func (s *Booking) CheckAvailable(ctx context.Context, userID string, rq *model.CheckAvailableRequest) error {
	screen, err := s.Store.GetScreenByID(ctx, rq.ScreenID)
	if err != nil {
		return err
	}
	if screen == nil {
		return errors.New(dictionary.ScreenIsNotExists)
	}
	currentSeats, err := s.Store.GetAllSeatByScreenID(ctx, rq.ScreenID)
	if err != nil {
		return err
	}
	seat := s.covertLocationToSeat(*rq.Location, screen.ID)
	if s.checkOK(currentSeats, []*storages.Seat{seat}, screen) {
		return nil
	}

	return errors.New("oh")
}

func (s *Booking) BookingNewSeat(ctx context.Context, userID string, bookingRq *model.BookingRequest) ([]*storages.Seat, error) {
	screen, err := s.Store.GetScreenByID(ctx, bookingRq.ScreenID)
	if err != nil {
		return nil, err
	}
	if screen == nil {
		return nil, errors.New(dictionary.ScreenIsNotExists)
	}
	currentSeats, err := s.Store.GetAllSeatByScreenID(ctx, bookingRq.ScreenID)
	if err != nil {
		return nil, err
	}
	var seats []*storages.Seat
	if bookingRq.Number > 0 {
		seats = s.randomSeats(currentSeats, bookingRq.Number, screen)
	} else {
		seats = s.covertLocationsToSeats(*bookingRq.Locations, screen.ID)
		if !s.seatsRequest(seats, currentSeats, screen) {
			seats = nil
		}
	}

	if len(seats) > 0 {
		for _, seat := range seats {
			seat.UserID = userID
			seat.BookedDate = utils.GetTimeNow()
		}
		err = s.Store.InsertSeats(ctx, seats)
		if err != nil {
			return nil, err
		}
	}

	return seats, nil
}

func (s *Booking) randomSeats(currentSeats []*storages.Seat, numberSeat int, screen *storages.Screen) []*storages.Seat {
	startsPointMap := map[storages.Seat]int{}
	shapesDictionary := [][]*storages.Seat{}

	startPoint := s.randomStartPoint(screen)
	startsPointMap[*startPoint] = 1

	shape := s.randomsShapeFromStartPoint(startPoint, screen, numberSeat)
	shapesDictionary = append(shapesDictionary, shape)

	for len(startsPointMap) < screen.NumberSeatRow*screen.NumberSeatColumn {
		if ok := s.checkOK(currentSeats, shape, screen); ok {
			break
		}
		shape = nil
		startPoint = s.randomStartPoint(screen)
		if _, ok := startsPointMap[*startPoint]; ok {
			startsPointMap[*startPoint] += 1
			continue
		}
		startsPointMap[*startPoint] = 1
		shape = s.randomsShapeFromStartPoint(startPoint, screen, numberSeat)
		isSame := false
		for _, shapeDictionary := range shapesDictionary {
			if s.compare2Shape(shapeDictionary, shape) {
				isSame = true
				break
			}
		}
		if !isSame {
			shapesDictionary = append(shapesDictionary, shape)
			continue
		}
	}

	return shape
}

//func separate(numberGroup int) []*storages.Seat {
//	result := []*storages.Seat{}
//
//	for i := 1; i < numberGroup; i++ {
//		result = separate(i)
//		if check(result) {
//			return result
//		}
//	}
//	return result
//}

func (s *Booking) checkOK(currentSeats []*storages.Seat, newSeats []*storages.Seat, screen *storages.Screen) bool {
	ring := s.getRing(newSeats, screen)
	filteredCurrentSeats := s.getRelativeNearest(currentSeats, ring)
	return s.isOk(filteredCurrentSeats, newSeats)
}

func (s *Booking) getRing(seats []*storages.Seat, screen *storages.Screen) *Ring {
	var ring *Ring
	if len(seats) > 0 {
		ring = &Ring{
			minRow:    screen.NumberSeatRow - 1,
			maxRow:    0,
			minColumn: screen.NumberSeatColumn - 1,
			maxColumn: 0,
		}
	} else {
		return nil
	}

	for _, seat := range seats {
		minDistanceRowSeat := utils.Max(seat.Row-s.Cfg.Distance, 0)
		ring.minRow = utils.Min(ring.minRow, minDistanceRowSeat)

		maxDistanceRowSeat := utils.Min(seat.Row+s.Cfg.Distance, screen.NumberSeatRow-1)
		ring.maxRow = utils.Max(ring.maxRow, maxDistanceRowSeat)

		minDistanceColumnSeat := utils.Max(seat.Column-s.Cfg.Distance, 0)
		ring.minColumn = utils.Min(ring.minColumn, minDistanceColumnSeat)

		maxDistanceColumnSeat := utils.Min(seat.Column+s.Cfg.Distance, screen.NumberSeatColumn-1)
		ring.maxColumn = utils.Max(ring.maxColumn, maxDistanceColumnSeat)
	}
	return ring
}

func (s *Booking) getRelativeNearest(currentSeats []*storages.Seat, ring *Ring) []*storages.Seat {
	var newCurrentSeat []*storages.Seat
	if ring != nil && len(currentSeats) > 0 {
		for _, currentSeat := range currentSeats {
			if currentSeat.Row <= ring.maxRow && currentSeat.Row >= ring.minRow &&
				currentSeat.Column <= ring.maxColumn && currentSeat.Column >= ring.minColumn {
				newCurrentSeat = append(newCurrentSeat, currentSeat)
			}
		}
		return newCurrentSeat
	}
	return currentSeats
}

func (s *Booking) isOk(currentSeats []*storages.Seat, newSeats []*storages.Seat) bool {
	for _, newSeat := range newSeats {
		for _, currentSeat := range currentSeats {
			distanceBetweenTwoSeats := utils.Abs(newSeat.Row-currentSeat.Row) + utils.Abs(newSeat.Column-currentSeat.Column)
			if distanceBetweenTwoSeats < s.Cfg.Distance {
				return false
			}
		}
	}
	return true
}

func (s *Booking) randomStartPoint(screen *storages.Screen) *storages.Seat {
	rand.Seed(time.Now().UnixNano())
	row := rand.Intn(screen.NumberSeatRow)
	column := rand.Intn(screen.NumberSeatColumn)
	return &storages.Seat{
		ID:         0,
		Row:        row,
		Column:     column,
		UserID:     "",
		ScreenID:   screen.ID,
		BookedDate: time.Time{},
	}
}

func (s *Booking) randomsShapeFromStartPoint(startPoint *storages.Seat, screen *storages.Screen, numberSeat int) []*storages.Seat {
	if startPoint != nil {
		var aSetSeats []*storages.Seat
		aSetSeats = append(aSetSeats, startPoint)
		for len(aSetSeats) < numberSeat {
			rand.Seed(time.Now().UnixNano())
			newSeat := &storages.Seat{}
			for _, seat := range aSetSeats {
				ok := false
				row := 0
				column := 0
				for !ok {
					minRow := utils.Max(seat.Row-1, 0)
					maxRow := utils.Min(seat.Row+1, screen.NumberSeatRow-1)

					minColumn := utils.Max(seat.Column-1, 0)
					maxColumn := utils.Min(seat.Column+1, screen.NumberSeatColumn-1)

					row = rand.Intn(maxRow-minRow+1) + minRow
					column = rand.Intn(maxColumn-minColumn+1) + minColumn
					if utils.Abs(seat.Row-row)+utils.Abs(seat.Column-column) == 1 {
						ok = true
					}
				}
				newSeat = &storages.Seat{
					ID:         0,
					Row:        row,
					Column:     column,
					UserID:     "",
					ScreenID:   screen.ID,
					BookedDate: time.Time{},
				}
				duplicate := false
				for _, seatIn := range aSetSeats {
					if newSeat.Row == seatIn.Row && newSeat.Column == seatIn.Column {
						duplicate = true
						break
					}
				}
				if !duplicate {
					aSetSeats = append(aSetSeats, newSeat)
					if len(aSetSeats) == numberSeat {
						break
					}
				}
			}
		}
		return aSetSeats
	}
	return nil
}

func (s *Booking) compare2Shape(X, Y []*storages.Seat) bool {
	m := make(map[storages.Seat]int)

	for _, y := range Y {
		m[*y]++
	}

	var ret []storages.Seat
	for _, x := range X {
		if m[*x] > 0 {
			m[*x]--
			continue
		}
		ret = append(ret, *x)
	}
	if len(ret) == 0 {
		return true
	}
	return false
}

func (s *Booking) spiltGroupSeats(seats []*storages.Seat) [][]*storages.Seat {
	groupsSeats := [][]*storages.Seat{}
	for i := 0; i < len(seats); i++ {
		if i == 0 {
			groupsSeats = append(groupsSeats, []*storages.Seat{seats[i]})
			continue
		}
		haveNeighbor := false
		for index, groups := range groupsSeats {
			for _, seatInGroup := range groups {
				if s.isClosed(seatInGroup, seats[i]) {
					groupsSeats[index] = append(groupsSeats[index], seats[i])
					haveNeighbor = true
				}
			}
		}
		if !haveNeighbor {
			groupsSeats = append(groupsSeats, []*storages.Seat{seats[i]})
		}
	}
	return groupsSeats
}

func (s *Booking) isClosed(seatA, seatB *storages.Seat) bool {
	if utils.Abs(seatA.Row-seatB.Row)+utils.Abs(seatA.Column-seatB.Column) == 1 {
		return true
	}
	return false
}

func (s *Booking) seatsRequest(requestSeats, currentSeats []*storages.Seat, screen *storages.Screen) bool {
	groupsSeat := s.spiltGroupSeats(requestSeats)
	sumOk := 0
	for _, group := range groupsSeat {
		if ok := s.checkOK(currentSeats, group, screen); ok {
			sumOk += len(group)
		}
	}
	if sumOk == len(requestSeats) {
		return true
	}
	return false
}

func (s *Booking) covertLocationsToSeats(locations []model.Location, screenID int) []*storages.Seat {
	if len(locations) == 0 {
		return nil
	}
	seats := []*storages.Seat{}
	for _, location := range locations {
		seats = append(seats, s.covertLocationToSeat(location, screenID))
	}
	return seats
}

func (s *Booking) covertLocationToSeat(location model.Location, screenID int) *storages.Seat {
	return &storages.Seat{
		ID:         0,
		Row:        location.Row,
		Column:     location.Column,
		UserID:     "",
		ScreenID:   screenID,
		BookedDate: time.Time{},
	}
}
