package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func (ths *service) GetNearestStores(ctx context.Context, req constant.GetNearestStoresRequest) ([]constant.NearestStore, error) {
	if req.Token.AddressID == nil {
		return nil, errors.New("invalid addressID from user token")
	}

	if req.Token.MerchantID == nil {
		return nil, errors.New("invalid merchantID from user token")
	}

	if req.CategoryType != constant.ALLOWED_CATEGORY_TYPE {
		return constructNearestStores(req, nil, nil, nil), nil
	}

	ac := make(chan *model.Address, 1)
	sc := make(chan []model.StoreWithCategory, 1)
	g, newCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		address, err := ths.repository.GetAddress(newCtx, *req.Token.AddressID)
		if err != nil {
			close(ac)
			return fmt.Errorf("failed at repository.GetAddress: %s", err.Error())
		}

		ac <- address
		return nil
	})

	input := repository.GetNearestStoresInput{
		UserID:     req.Token.UserID,
		MerchantID: *req.Token.MerchantID,
		CategoryID: req.CategoryID,
		Name:       req.Name,
	}

	g.Go(func() error {
		stores, err := ths.repository.GetNearestStores(newCtx, input)
		if err != nil {
			close(sc)
			return fmt.Errorf("failed at repository.GetNearestStores: %s", err.Error())
		}

		sc <- stores
		return nil
	})

	address := <-ac
	stores := <-sc

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if address == nil {
		return nil, errors.New("address not found")
	}

	storeIDs := getStoreIDs(stores)
	if len(storeIDs) == 0 {
		return constructNearestStores(req, nil, nil, nil), nil
	}
	config, err := ths.repository.GetStoreOpenConfig(ctx, storeIDs)
	if err != nil {
		logrus.
			WithField("func", "service.GetNearestStores").
			WithField("input", input).
			WithError(err).
			Warn("failed at repository.GetStoreOpenConfig, fallback to isOpen false")
	}

	storeIsOpen := constructStoreIsOpenMapByStoreID(storeIDs, config)

	return constructNearestStores(req, address, stores, storeIsOpen), nil
}

func getStoreIDs(stores []model.StoreWithCategory) []int64 {
	if len(stores) == 0 {
		return nil
	}

	storeIDs := make([]int64, 0, len(stores))
	for i := range stores {
		storeIDs = append(storeIDs, stores[i].ID)
	}

	return storeIDs
}

func constructStoreIsOpenMapByStoreID(storeIDs []int64, config model.StoreConfig) map[int64]bool {
	storeIsOpen := make(map[int64]bool)

	// Jika config untuk toko-toko terdekat tidak ada, maka seluruh toko tersebut
	// dianggap tutup
	if len(config) == 0 {
		for i := range storeIDs {
			if _, exist := storeIsOpen[storeIDs[i]]; !exist {
				storeIsOpen[storeIDs[i]] = false
			}
		}

		return storeIsOpen
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	currentDay := strconv.Itoa(int(now.Weekday()))

	for i := range storeIDs {
		storeID := strconv.Itoa(int(storeIDs[i]))
		if _, configExist := config[storeID]; configExist {
			// Jika config untuk suatu toko ada dan nilainya "true" untuk key "closed",
			// maka toko tersebut dianggap tutup
			isClosed, isClosedConfigExist := config[storeID]["closed"]
			if isClosedConfigExist && strings.EqualFold(isClosed, "true") {
				storeIsOpen[storeIDs[i]] = false
				continue
			}

			// Jika config untuk suatu toko ada dengan key "closed" bernilai "false",
			// maka lihat config dengan key "service"
			serviceInterval, serviceIntervalConfigExist := config[storeID]["service"]
			if serviceIntervalConfigExist {
				openingHours := make(map[string]struct {
					Start string `json:"start"`
					End   string `json:"end"`
				})

				// Jika config dengan key "service" tidak bisa dibaca, anggap toko
				// tersebut tutup
				err := json.Unmarshal([]byte(serviceInterval), &openingHours)
				if err != nil {
					storeIsOpen[storeIDs[i]] = false
					continue
				}

				startTime, errStart := time.ParseInLocation(constant.DATE_TIME_LAYOUT, fmt.Sprintf("%s %s", now.Format(constant.DATE_ONLY_LAYOUT), openingHours[currentDay].Start), loc)

				// Jika config dengan key "service" ada, tetapi startTime invalid
				// anggap toko tersebut tutup
				if errStart != nil {
					storeIsOpen[storeIDs[i]] = false
					continue
				}

				endTime, errEnd := time.ParseInLocation(constant.DATE_TIME_LAYOUT, fmt.Sprintf("%s %s", now.Format(constant.DATE_ONLY_LAYOUT), openingHours[currentDay].End), loc)

				// Jika config dengan key "service" ada, tetapi endTime invalid
				// anggap toko tersebut tutup
				if errEnd != nil {
					storeIsOpen[storeIDs[i]] = false
					continue
				}

				// Jika waktu sekarang adalah setelah toko buka dan sebelum toko tutup,
				// anggap toko tersebut buka
				if now.After(startTime) && now.Before(endTime) {
					storeIsOpen[storeIDs[i]] = true
					continue
				}

				// Namun, jika waktu sekarang adalah sebelum toko buka atau setelah
				// toko tutup, anggap toko tersebut tutup
				storeIsOpen[storeIDs[i]] = false
				continue
			}
		}

		// Jika
		// 1. config untuk suatu toko tidak ada, ATAU
		// 2. config untuk key closed bernilai "false", tetapi config untuk key "service"
		//    tidak ada,
		// maka toko tersebut dianggap tutup
		storeIsOpen[storeIDs[i]] = false
	}

	return storeIsOpen
}

func constructNearestStores(req constant.GetNearestStoresRequest, address *model.Address, stores []model.StoreWithCategory, storeIsOpen map[int64]bool) []constant.NearestStore {
	nearestStores := make([]constant.NearestStore, 0, req.Limit)
	for i := range stores {
		if stores[i].LocationLat == nil || stores[i].LocationLong == nil {
			continue
		}

		dist := calculateDistanceInKilometer(*stores[i].LocationLat, *stores[i].LocationLong, address.LocationLat, address.LocationLong)

		if stores[i].MaxRadius != nil && *stores[i].MaxRadius <= int64(dist) {
			continue
		}

		nearestStores = append(nearestStores, constant.NearestStore{
			ID:           stores[i].ID,
			Name:         stores[i].Name,
			Picture:      stores[i].Pic,
			Address:      stores[i].Address,
			Distance:     fmt.Sprintf("%.2f km", dist),
			LocationLat:  stores[i].LocationLat,
			LocationLong: stores[i].LocationLong,
			IsOpen:       storeIsOpen[stores[i].ID],
			Dist:         dist,
		})
	}

	if len(nearestStores) > 0 {
		sort.SliceStable(nearestStores, func(i, j int) bool { return nearestStores[i].Dist < nearestStores[j].Dist })
	}

	if len(nearestStores) > req.Limit {
		nearestStores = nearestStores[:req.Limit]
	}

	return nearestStores
}

func calculateDistanceInKilometer(firstLatitude, firstLongitude, secondLatitude, secondLongitude float64) float64 {
	alpha1 := math.Pi * firstLatitude / 180
	alpha2 := math.Pi * secondLatitude / 180

	thetha := math.Pi * (firstLongitude - secondLongitude) / 180

	dist := math.Sin(alpha1)*math.Sin(alpha2) + math.Cos(alpha1)*math.Cos(alpha2)*math.Cos(thetha)
	dist = math.Min(dist, 1)
	dist = math.Max(dist, -1)

	return math.Abs(math.Acos(dist) * 180 * 60 * 1.1515 * 1.609344 / math.Pi)
}
