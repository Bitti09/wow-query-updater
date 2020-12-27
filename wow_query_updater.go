package main

import (
	"flag"
	"fmt"
	blizzard_api "github.com/francis-schiavo/blizzard-api-go"
	"wow-query-updater/connections"
	"wow-query-updater/updater"
)

func main() {
	config := &Config{}
	config.LoadFromFile("config.json")

	classic := flag.Bool("classic", false, "Classic mode")
	flag.Parse()

	schema := "public"
	if *classic {
		schema = "classic"
	}

	connections.Connect(config.Username, config.Password, config.Database, 105, schema, config.Hostname, config.Port)
	connections.DatabaseSetup(*classic)
	connections.ReportingMode = false

	cacheProvider := &connections.PostgresCache{Key: "1"}

	connections.WowClient = blizzard_api.NewWoWClient("us", cacheProvider,  nil, *classic)
	connections.WowClient.CreateAccessToken(config.ClientID, config.ClientSecret, "")

	taskManager := updater.NewTaskManager(80, 12, updater.LtInfo)

	//// Common
	//taskManager.AddIndexTask("playable race", "PlayableRaceIndex", "races", "PlayableRace", updater.UpdatePlayableRace)
	//
	//taskManager.AddIndexTask("power type", "PowerTypeIndex", "power_types", "PowerType", updater.UpdatePowerType)
	//taskManager.AddIndexTask("playable class", "PlayableClassIndex", "classes", "PlayableClass", updater.UpdatePlayableClass)
	//taskManager.AddMediaTask("playable class assets", &datasets.PlayableClassMedia{}, "PlayableClassMedia", updater.UpdatePlayableClassMedia)
	//taskManager.AddIndexTask("playable specialization", "PlayableSpecializationIndex", "character_specializations", "PlayableSpecialization", updater.UpdatePlayableSpecialization)
	//taskManager.AddIndexTask("playable pet specialization", "PlayableSpecializationIndex", "pet_specializations", "PlayableSpecialization", updater.UpdatePlayableSpecialization)
	//
	//// Creature
	//taskManager.AddIndexTask("creature family", "CreatureFamilyIndex", "creature_families", "CreatureFamily", updater.UpdateCreatureFamily)
	//taskManager.AddIndexTask("creature type", "CreatureTypeIndex", "creature_types", "CreatureType", updater.UpdateCreatureType)
	//taskManager.AddSearchTask("creature", "CreatureSearch", "Creature", updater.UpdateCreature)
	//
	//taskManager.AddIndexTaskLimited("item class", "ItemClassIndex", "item_classes", "ItemClass", updater.UpdateItemClass, 50)
	//taskManager.AddSimpleTask("add missing classes", updater.InsertMissingItemClasses)
	//taskManager.AddSimpleTask("add missing stats", updater.InsertMissingStats)

	// Item
	if *classic {
		taskManager.AddSearchTask("item", "ItemSearch", "Item", updater.UpdateItem)
	}

	if !*classic {
		//// Preload profession
		//taskManager.AddIndexTask("profession", "ProfessionIndex", "professions", "Profession", updater.UpdateProfession)
		//
		//// Reputation
		//taskManager.AddIndexTask("reputation tier", "ReputationTierIndex", "reputation_tiers", "ReputationTier", updater.UpdateReputationTier)
		//taskManager.AddIndexTask("reputation faction", "ReputationFactionIndex", "root_factions", "ReputationFaction", updater.UpdateReputationFaction)
		//taskManager.AddIndexTask("reputation faction", "ReputationFactionIndex", "factions", "ReputationFaction", updater.UpdateParentReputation)
		//taskManager.AddSimpleTask("add missing reputation tiers", updater.InsertMissingReputationTiers)
		//
		//// Spell
		//taskManager.AddSearchTask("spell", "SpellSearch", "Spell", updater.UpdateSpell)
		//taskManager.AddSimpleTask("add missing spells", updater.InsertMissingSpells)
		//
		//// Items
		//taskManager.AddSearchTask("item", "ItemSearch", "Item", updater.UpdateItem)
		//
		//// Common
		//taskManager.AddIndexTaskLimited("talents", "TalentIndex", "talents", "Talent", updater.UpdateTalent, 20)
		//taskManager.AddMediaTask("playable specialization media", &datasets.PlayableSpecializationMedia{}, "PlayableSpecializationMedia", updater.UpdatePlayableSpecializationMedia)
		//taskManager.AddIndexTaskLimited("pvp talents", "PvPTalentIndex", "pvp_talents", "PvPTalent", updater.UpdatePvpTalent, 20)
		//
		//taskManager.AddIndexTask("title", "TitleIndex", "titles", "Title", updater.UpdateTitle)
		//
		//// Azerite
		//taskManager.AddIndexTaskLimited("azerite essence", "AzeriteEssenceIndex", "azerite_essences", "AzeriteEssence", updater.UpdateAzeriteEssence, 20)
		//taskManager.AddMediaTask("azerite essence media", &datasets.AzeriteEssenceMedia{}, "AzeriteEssenceMedia", updater.UpdateAzeriteEssenceMedia)
		//
		//// Achievement
		//taskManager.AddIndexTask("root achievement category", "AchievementCategoryIndex", "root_categories", "AchievementCategory", updater.UpdateAchievementCategory)
		//taskManager.AddIndexTask("guild achievement category", "AchievementCategoryIndex", "guild_categories", "AchievementCategory", updater.UpdateAchievementCategory)
		//taskManager.AddIndexTask("achievement category", "AchievementCategoryIndex", "categories", "AchievementCategory", updater.UpdateAchievementCategory)
		//
		//taskManager.AddIndexTask("root achievement category update parent", "AchievementCategoryIndex", "root_categories", "AchievementCategory", updater.UpdateParentCategory)
		//taskManager.AddIndexTask("guild achievement category update parent", "AchievementCategoryIndex", "guild_categories", "AchievementCategory", updater.UpdateParentCategory)
		//taskManager.AddIndexTask("achievement category update parent", "AchievementCategoryIndex", "categories", "AchievementCategory", updater.UpdateParentCategory)
		//
		//taskManager.AddIndexTask("achievement", "AchievementIndex", "achievements", "Achievement", updater.UpdateAchievement)
		//taskManager.AddMediaTask("achievement assets", &datasets.AchievementMedia{}, "AchievementMedia", updater.UpdateAchievementMedia)
		//
		//// Quest
		//taskManager.AddIndexTaskLimited("quest category", "QuestCategoryIndex", "categories", "QuestCategory", updater.UpdateQuestCategory, 50)
		//taskManager.AddIndexTaskLimited("quest type", "QuestTypeIndex", "types", "QuestType", updater.UpdateQuestType, 50)
		//taskManager.AddIndexTaskLimited("quest area", "QuestAreaIndex", "areas", "QuestArea", updater.UpdateQuestArea, 50)
		//
		//// Collections
		//taskManager.AddIndexTask("mount", "MountIndex", "mounts", "Mount", updater.UpdateMount)
		//taskManager.AddMediaTask("mount media", &datasets.MountDisplayMedia{}, "CreatureDisplayMedia", updater.UpdateMountDisplayMedia)
		//taskManager.AddIndexTask("pet", "PetIndex", "pets", "Pet", updater.UpdatePet)
		//
		//// Profession
		//taskManager.AddIndexTaskLimited("profession", "ProfessionIndex", "professions", "Profession", updater.UpdateProfessionTiers, 30)
		//taskManager.AddMediaTask("profession media", &datasets.ProfessionMedia{}, "ProfessionMedia", updater.UpdateProfessionMedia)
		//taskManager.AddMediaTask("recipe media", &datasets.RecipeMedia{}, "RecipeMedia", updater.UpdateRecipeMedia)
		//
		////Journal
		//taskManager.AddIndexTask("journal expansion", "JournalExpansionIndex", "tiers", "JournalExpansion", updater.UpdateJournalExpansion)
		//taskManager.AddIndexTask("journal instance", "JournalInstanceIndex", "instances", "JournalInstance", updater.UpdateJournalInstance)
		//taskManager.AddIndexTask("journal encounter", "JournalEncounterIndex", "encounters", "JournalEncounter", updater.UpdateJournalEncounter)
		//taskManager.AddMediaTask("instance media", &datasets.JournalInstanceMedia{}, "JournalInstanceMedia", updater.UpdateInstanceMedia)

		// Tech talent
		taskManager.AddIndexTask("tech talent tree", "TechTalentTreeIndex", "talent_trees", "TechTalentTree", updater.UpdateTechTalentTree)
		//taskManager.AddIndexTask("tech talent", "TechTalentIndex", "talents", "TechTalent", updater.UpdateTechTalent)
		taskManager.AddIndexTask("tech talent workaround", "TechTalentTreeIndex", "talent_trees", "TechTalentTree", updater.UpdateTechTalentUsingTree)

		// Covenant
		taskManager.AddIndexTask("conduit", "ConduitIndex", "conduits", "Conduit", updater.UpdateConduit)
		taskManager.AddIndexTask("covenant", "CovenantIndex", "covenants", "Covenant", updater.UpdateCovenant)
		taskManager.AddIndexTask("soulbind", "SoulbindIndex", "soulbinds", "Soulbind", updater.UpdateSoulbind)
	}

	//// Shared Media
	//taskManager.AddMediaTask("creature media", &datasets.CreatureDisplayMedia{}, "CreatureDisplayMedia", updater.UpdateCreatureDisplayMedia)
	//taskManager.AddMediaTask("item media", &datasets.ItemMedia{}, "ItemMedia", updater.UpdateItemMedia)
	//taskManager.AddMediaTask("creature family media", &datasets.CreatureFamilyMedia{}, "CreatureFamilyMedia", updater.UpdateCreatureFamilyMedia)

	//if !*classic {
	//	taskManager.AddMediaTask("spell media", &datasets.SpellMedia{}, "SpellMedia", updater.UpdateSpellMedia)
	//}

	fmt.Printf("Classic mode: %v\n", *classic)
	taskManager.Run()

	/*if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	defer ui.Close()

	grid := ui.NewGrid()
	termWidth, _ := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, 12)

	taskLabel := NewLabel("Current task", "Items", ui.ColorYellow)
	taskType := NewLabel("Task type", "Search", ui.ColorYellow)
	cachedRequestsLabel := NewLabel("Cached requests", "1000", ui.ColorGreen)
	uncachedRequestsLabel := NewLabel("Uncached requests", "1000", ui.ColorYellow)
	failedRequestsLabel := NewLabel("Failed requests", "1000", ui.ColorRed)
	taskGauge := widgets.NewGauge()
	taskGauge.Title = "Current task progress"
	taskGauge.Percent = 70

	tasksGauge := widgets.NewGauge()
	tasksGauge.Title = "Total progress"
	tasksGauge.Percent = 70

	grid.Set(
		ui.NewRow(3.0/12,
			ui.NewCol(1.5/2, taskLabel),
			ui.NewCol(.5/2, taskType),
		),
		ui.NewRow(3.0/12,
			ui.NewCol(1.0, taskGauge),
		),
		ui.NewRow(3.0/12,
			ui.NewCol(1.0, tasksGauge),
		),
		ui.NewRow(3.0/12,
			ui.NewCol(1.0/3, cachedRequestsLabel),
			ui.NewCol(1.0/3, uncachedRequestsLabel),
			ui.NewCol(1.0/3, failedRequestsLabel),
		),
	)

	ui.Render(grid)

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, 12)
				ui.Clear()
				ui.Render(grid)
			}
		}
	}*/
}
