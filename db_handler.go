package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/suffer-sami/realtor-scraper/internal/database"
)

func (cfg *config) executeTransaction(ctx context.Context, txFunc func(context.Context, *database.Queries) error) error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	tx, err := cfg.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	qtx := cfg.dbQueries.WithTx(tx)

	err = txFunc(ctx, qtx)

	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}
	return tx.Commit()
}

func (cfg *config) storeAgents(agents []Agent) {
	defer cfg.wg.Done()
	for _, agent := range agents {
		cfg.wg.Add(1)
		go func() {
			if err := cfg.storeAgent(agent); err != nil {
				cfg.logger.Errorf("error storing agent (ID: %s): %v", agent.ID, err)
			}
		}()
	}
}

func (cfg *config) storeAgent(agent Agent) error {
	defer cfg.wg.Done()
	return cfg.executeTransaction(context.Background(), func(ctx context.Context, qtx *database.Queries) error {
		cfg.logger.Infof("Agent: %s", agent.PersonName)

		dbAgent, err := qtx.GetAgent(ctx, agent.ID)
		if err != nil {
			if err != sql.ErrNoRows {
				return err
			}

			dbAgent, err = qtx.CreateAgent(ctx, database.CreateAgentParams{
				ID:                   agent.ID,
				FirstName:            toNullString(agent.FirstName),
				LastName:             toNullString(agent.LastName),
				NickName:             toNullString(agent.NickName),
				PersonName:           toNullString(agent.PersonName),
				Title:                toNullString(agent.Title),
				Slogan:               toNullString(agent.Slogan),
				Email:                toNullString(agent.Email),
				AgentRating:          toNullInt(agent.AgentRating),
				Description:          toNullString(agent.Description),
				RecommendationsCount: toNullInt(agent.RecommendationsCount),
				ReviewCount:          toNullInt(agent.ReviewCount),
				LastUpdated:          strToNullTime(agent.LastUpdated, time.RFC1123),
				FirstMonth:           numericToNullInt(agent.FirstMonth),
				FirstYear:            toNullInt(agent.AgentRating),
				Video:                toNullString(agent.Video),
				WebUrl:               toNullString(agent.WebURL),
				Href:                 toNullString(agent.Href),
			})

			if err != nil {
				return err
			}

			if cfg.saveRawAgents {
				cfg.logger.Debugf("- raw agent:")
				if jsonStrAgent, err := toJsonString(agent); err == nil {
					cfg.logger.Debugf("	* %s", jsonStrAgent)
					if err := qtx.CreateRawAgent(ctx, database.CreateRawAgentParams{
						AgentID: toNullString(dbAgent.ID),
						Data:    toNullString(jsonStrAgent),
					}); err != nil {
						cfg.logger.Errorf("error creating raw agent: %v", err)
					}
				}
			}
		}
		agentId := toNullString(dbAgent.ID)

		cfg.logger.Debugf("- sales data: %s", agent.RecentlySold.LastSoldDate)
		if err := qtx.CreateSalesData(ctx, database.CreateSalesDataParams{
			Count:        toNullInt(agent.RecentlySold.Count),
			Min:          toNullInt(agent.RecentlySold.Min),
			Max:          toNullInt(agent.RecentlySold.Max),
			LastSoldDate: strToNullTime(agent.RecentlySold.LastSoldDate, time.DateOnly),
			AgentID:      agentId,
		}); err != nil {
			cfg.logger.Errorf("error creating sales data: %v", err)
		}

		cfg.logger.Debugf("- listing data: %s", agent.ForSalePrice.LastListingDate)
		if err := qtx.CreateListingsData(ctx, database.CreateListingsDataParams{
			Count:           toNullInt(agent.ForSalePrice.Count),
			Min:             toNullInt(agent.ForSalePrice.Min),
			Max:             toNullInt(agent.ForSalePrice.Max),
			LastListingDate: timeToNullTime(agent.ForSalePrice.LastListingDate),
			AgentID:         agentId,
		}); err != nil {
			cfg.logger.Errorf("error creating listing data: %v", err)
		}

		cfg.logger.Debugf("- social medias:")
		for _, socialMedia := range agent.SocialMedias {
			cfg.logger.Debugf("	* %s", socialMedia.Type)
			if err := qtx.CreateSocialMedia(ctx, database.CreateSocialMediaParams{
				Type:    toNullString(socialMedia.Type),
				Href:    toNullString(socialMedia.Href),
				AgentID: agentId,
			}); err != nil {
				cfg.logger.Errorf("error creating social media: %v", err)
			}
		}

		cfg.logger.Debugf("- feed licences:")
		for _, feedLicense := range agent.FeedLicenses {
			cfg.logger.Debugf("	* (%s, %s)", feedLicense.StateCode, feedLicense.Country)
			if err := qtx.CreateFeedLicense(ctx, database.CreateFeedLicenseParams{
				Country:       toNullString(feedLicense.Country),
				LicenseNumber: toNullString(feedLicense.LicenseNumber),
				StateCode:     toNullString(feedLicense.StateCode),
				AgentID:       agentId,
			}); err != nil {
				cfg.logger.Errorf("error creating feed licence: %v", err)
			}
		}

		cfg.logger.Debugf("- mls:")
		for _, mls := range agent.Mls {
			cfg.logger.Debugf("	* %s", mls.Abbreviation)
			dbMls, err := qtx.GetMultipleListingService(ctx, database.GetMultipleListingServiceParams{
				Abbreviation:  toNullString(mls.Abbreviation),
				Type:          toNullString(mls.Type),
				MemberID:      toNullString(mls.MemberID),
				LicenseNumber: toNullString(mls.LicenseNumber),
			})

			if err != nil {
				if err != sql.ErrNoRows {
					return err
				}
				dbMls, err = qtx.CreateMultipleListingService(ctx, database.CreateMultipleListingServiceParams{
					Abbreviation:  toNullString(mls.Abbreviation),
					LicenseNumber: toNullString(mls.LicenseNumber),
					Type:          toNullString(mls.Type),
					MemberID:      toNullString(mls.MemberID),
					IsPrimary:     toNullBool(mls.Primary),
				})

				if err != nil {
					return err
				}
			}

			if err := qtx.CreateAgentMultipleListingService(ctx, database.CreateAgentMultipleListingServiceParams{
				AgentID:                  agentId,
				MultipleListingServiceID: toNullInt64(dbMls.ID),
			}); err != nil {
				cfg.logger.Errorf("error creating agent mls: %v", err)
			}
		}
		cfg.logger.Debugf("- mls history:")
		for _, mls := range agent.MlsHistory {
			cfg.logger.Debugf("	* %s", mls.Abbreviation)
			dbMls, err := qtx.GetMultipleListingService(ctx, database.GetMultipleListingServiceParams{
				Abbreviation:  toNullString(mls.Abbreviation),
				Type:          toNullString(mls.Type),
				MemberID:      toNullString(mls.Member.ID),
				LicenseNumber: toNullString(mls.LicenseNumber),
			})

			if err != nil {
				if err != sql.ErrNoRows {
					return err
				}
				dbMls, err = qtx.CreateMultipleListingService(ctx, database.CreateMultipleListingServiceParams{
					Abbreviation:     toNullString(mls.Abbreviation),
					InactivationDate: timeToNullTime(mls.InactivationDate),
					LicenseNumber:    toNullString(mls.LicenseNumber),
					IsPrimary:        toNullBool(mls.Primary),
					Type:             toNullString(mls.Type),
					MemberID:         toNullString(mls.Member.ID),
				})

				if err != nil {
					return err
				}
			}

			if dbMls.InactivationDate.Time != mls.InactivationDate {
				cfg.logger.Debugf("	* %s (update inactivation_date: %v)", mls.Abbreviation, mls.InactivationDate)
				if err := qtx.UpdateMultipleListingServiceInactivationDate(ctx, database.UpdateMultipleListingServiceInactivationDateParams{
					InactivationDate: dbMls.InactivationDate,
					ID:               dbMls.ID,
				}); err != nil {
					cfg.logger.Errorf("error (update mls inactivation_date: %v)", err)
				}
			}

			if err = qtx.CreateAgentMultipleListingService(ctx, database.CreateAgentMultipleListingServiceParams{
				AgentID:                  agentId,
				MultipleListingServiceID: toNullInt64(dbMls.ID),
			}); err != nil {
				cfg.logger.Errorf("error creating agent mls: %v", err)
			}
		}

		cfg.logger.Debugf("- languages:")
		for _, lang := range agent.Languages {
			cfg.logger.Debugf("	* %s", lang)
			dbLang, err := qtx.GetLanguage(ctx, toNullString(lang))
			if err != nil {
				if err != sql.ErrNoRows {
					return err
				}

				dbLang, err = qtx.CreateLanguage(ctx, toNullString(lang))
				if err != nil {
					return err
				}
			}

			if err := qtx.CreateAgentLanguage(ctx, database.CreateAgentLanguageParams{
				LanguageID: toNullInt64(dbLang.ID),
				AgentID:    agentId,
			}); err != nil {
				cfg.logger.Errorf("error creating agent language: %v", err)
			}
		}
		cfg.logger.Debugf("- user languages:")
		for _, lang := range agent.UserLanguages {
			cfg.logger.Debugf("	* %s", lang)
			dbLang, err := qtx.GetLanguage(ctx, toNullString(lang))
			if err != nil {
				if err != sql.ErrNoRows {
					return err
				}

				dbLang, err = qtx.CreateLanguage(ctx, toNullString(lang))
				if err != nil {
					return err
				}
			}

			if err := qtx.CreateAgentUserLanguage(ctx, database.CreateAgentUserLanguageParams{
				LanguageID: toNullInt64(dbLang.ID),
				AgentID:    agentId,
			}); err != nil {
				cfg.logger.Errorf("error creating agent language: %v", err)
			}
		}

		cfg.logger.Debugf("- zips:")
		for _, zip := range agent.Zips {
			cfg.logger.Debugf("	* %s", zip)
			dbZip, err := qtx.GetZip(ctx, toNullString(zip))
			if err != nil {
				if err != sql.ErrNoRows {
					return err
				}

				dbZip, err = qtx.CreateZip(ctx, toNullString(zip))
				if err != nil {
					return err
				}
			}

			if err := qtx.CreateAgentZip(ctx, database.CreateAgentZipParams{
				ZipID:   toNullInt64(dbZip.ID),
				AgentID: agentId,
			}); err != nil {
				cfg.logger.Errorf("error creating agent zip: %v", err)
			}
		}

		cfg.logger.Debugf("- areas:")
		for _, area := range agent.ServedAreas {
			cfg.logger.Debugf("	* (%s, %s)", area.Name, area.StateCode)
			dbArea, err := qtx.GetArea(ctx, database.GetAreaParams{
				Name:      toNullString(area.Name),
				StateCode: toNullString(area.StateCode),
			})

			if err != nil {
				if err != sql.ErrNoRows {
					return err
				}

				dbArea, err = qtx.CreateArea(ctx, database.CreateAreaParams{
					Name:      toNullString(area.Name),
					StateCode: toNullString(area.StateCode),
				})
				if err != nil {
					return err
				}
			}

			if err := qtx.CreateAgentServedArea(ctx, database.CreateAgentServedAreaParams{
				AreaID:  toNullInt64(dbArea.ID),
				AgentID: agentId,
			}); err != nil {
				cfg.logger.Errorf("error creating agent served area: %v", err)
			}
		}
		cfg.logger.Debugf("- marketing areas:")
		for _, area := range agent.MarketingAreaCities {
			cfg.logger.Debugf("	* (%s, %s)", area.Name, area.StateCode)
			dbArea, err := qtx.GetArea(ctx, database.GetAreaParams{
				Name:      toNullString(area.Name),
				StateCode: toNullString(area.StateCode),
			})

			if err != nil {
				if err != sql.ErrNoRows {
					return err
				}

				dbArea, err = qtx.CreateArea(ctx, database.CreateAreaParams{
					Name:      toNullString(area.Name),
					StateCode: toNullString(area.StateCode),
				})
				if err != nil {
					return err
				}
			}

			if err := qtx.CreateAgentMarketingArea(ctx, database.CreateAgentMarketingAreaParams{
				AreaID:  toNullInt64(dbArea.ID),
				AgentID: agentId,
			}); err != nil {
				cfg.logger.Errorf("error creating agent marketing area: %v", err)
			}
		}

		cfg.logger.Debugf("- designations:")
		for _, designation := range agent.Designations {
			cfg.logger.Debugf("	* %s", designation.Name)
			dbDesignation, err := qtx.GetDesignation(ctx, toNullString(designation.Name))
			if err != nil {
				if err != sql.ErrNoRows {
					return err
				}

				dbDesignation, err = qtx.CreateDesignation(ctx, toNullString(designation.Name))
				if err != nil {
					return err
				}
			}

			if err := qtx.CreateAgentDesignation(ctx, database.CreateAgentDesignationParams{
				DesignationID: toNullInt64(dbDesignation.ID),
				AgentID:       agentId,
			}); err != nil {
				cfg.logger.Errorf("error creating agent designation: %v", err)
			}
		}
		cfg.logger.Debugf("- specializations:")
		for _, specialization := range agent.Specializations {
			cfg.logger.Debugf("	* %s", specialization.Name)
			dbSpecialization, err := qtx.GetSpecialization(ctx, toNullString(specialization.Name))
			if err != nil {
				if err != sql.ErrNoRows {
					return err
				}

				dbSpecialization, err = qtx.CreateSpecialization(ctx, toNullString(specialization.Name))
				if err != nil {
					return err
				}
			}

			if err := qtx.CreateAgentSpecialization(ctx, database.CreateAgentSpecializationParams{
				SpecializationID: toNullInt64(dbSpecialization.ID),
				AgentID:          agentId,
			}); err != nil {
				cfg.logger.Errorf("error creating agent specialization: %v", err)
			}
		}
		return nil
	})
}
