#pragma once
#include <array>
#include "picojson_serializer.h"
#include "picojson_vector_serializer.h"

#define VERSION 5

namespace HoMMMonitor {
	struct Combat {
		int hero = -1;
		int heroArmyValue = -1;
		int enemyArmyValue = -1;
		friend class picojson::convert::access;
		template<class Archive>
		void json(Archive& ar) const
		{
			ar& picojson::convert::member("hero", hero);
			ar& picojson::convert::member("heroArmyValue", heroArmyValue);
			ar& picojson::convert::member("enemyArmyValue", enemyArmyValue);
		}
	};

	struct Town {
		std::string name = "";
		int type = -1;
		bool grail = false;

		int x = 0;
		int y = 0;
		int z = 0;

		bool manaVortextUnused = false;
		//std::array<std::array<int, 6>, 5> spells;
		std::vector<std::vector<int>> spells{};

		bool builtThisTurn = false;
		int guardsValue = 0;

		INT32	garrisonHero;
		INT32	visitingHero;

		bool hasFort = false;

		friend class picojson::convert::access;
		template<class Archive>
		void json(Archive& ar) const
		{
			ar& picojson::convert::member("name", name);
			ar& picojson::convert::member("type", type);
			ar& picojson::convert::member("grail", grail);

			ar& picojson::convert::member("x", x);
			ar& picojson::convert::member("y", y);
			ar& picojson::convert::member("z", z);

			ar& picojson::convert::member("manaVortextUnused", manaVortextUnused);
			ar& picojson::convert::member("spells", spells);
			ar& picojson::convert::member("builtThisTurn", builtThisTurn);

			ar& picojson::convert::member("guardsValue", guardsValue);

			ar& picojson::convert::member("garrisonHero", garrisonHero);
			ar& picojson::convert::member("visitingHero", visitingHero);
			ar& picojson::convert::member("hasFort", hasFort);
		}
	};

	struct Item {
		int id;
		int subtype;
		friend class picojson::convert::access;
		template<class Archive>
		void json(Archive& ar) const
		{
			ar& picojson::convert::member("id", id);
			ar& picojson::convert::member("subtype", subtype);
		}
	};

	struct Hero {
		int id = -1;
		std::string name = "";

		int movement = 2000;
		int maxMovement = 2000;
		int experience = 0;
		int level = 1;
		int gender = 0;

		std::vector<Item> backpack{};
		std::vector<Item> weared{};

		std::vector<int> learned_spells{};
		std::vector<int> available_spells{};

		int attack, defense, spell_power, knowledge = 0;

		std::vector<int> sec_skills{};

		int mana = 0;

		int armyValue = 0;
		int armyPower = 0;

		int x = 0;
		int y = 0;
		int z = 0;

		int moraleBonus = 0;
		int luckBonus = 0;

		bool isVisible = false;

		friend class picojson::convert::access;
		template<class Archive>
		void json(Archive& ar) const
		{
			ar& picojson::convert::member("id", id);
			ar& picojson::convert::member("name", name);

			ar& picojson::convert::member("movement", movement);
			ar& picojson::convert::member("maxMovement", maxMovement);
			ar& picojson::convert::member("experience", experience);
			ar& picojson::convert::member("level", level);
			ar& picojson::convert::member("gender", gender);

			ar& picojson::convert::member("backpack", backpack);
			ar& picojson::convert::member("weared", weared);

			ar& picojson::convert::member("learned_spells", learned_spells);
			ar& picojson::convert::member("available_spells", available_spells);

			ar& picojson::convert::member("attack", attack);
			ar& picojson::convert::member("defense", defense);
			ar& picojson::convert::member("spell_power", spell_power);
			ar& picojson::convert::member("knowledge", knowledge);

			ar& picojson::convert::member("sec_skills", sec_skills);
			ar& picojson::convert::member("mana", mana);
			ar& picojson::convert::member("armyValue", armyValue);
			ar& picojson::convert::member("armyPower", armyPower);

			ar& picojson::convert::member("x", x);
			ar& picojson::convert::member("y", y);
			ar& picojson::convert::member("z", z);

			ar& picojson::convert::member("moraleBonus", moraleBonus);
			ar& picojson::convert::member("luckBonus", luckBonus);

			ar& picojson::convert::member("isVisible", isVisible);
		}
	};

	struct Map {
		std::string mapName = "mapName";
		int size = 0;
		int SubterraneanLevel;
		std::vector<int> visionS;
		std::vector<int> visionU;

		friend class picojson::convert::access;
		template<class Archive>
		void json(Archive& ar) const
		{
			ar& picojson::convert::member("mapName", mapName);
			ar& picojson::convert::member("size", size);

			ar& picojson::convert::member("subterraneanLevel", SubterraneanLevel);
			ar& picojson::convert::member("visionS", visionS);
			ar& picojson::convert::member("visionU", visionU);
		}
	};

	struct Resources {
		int gold = 0;
		int wood = 0;
		int mercury = 0;
		int ore = 0;
		int sulfur = 0;
		int crystal = 0;
		int gems = 0;
		friend class picojson::convert::access;
		template<class Archive>
		void json(Archive& ar) const
		{
			ar& picojson::convert::member("gold", gold);
			ar& picojson::convert::member("wood", wood);
			ar& picojson::convert::member("mercury", mercury);
			ar& picojson::convert::member("ore", ore);
			ar& picojson::convert::member("sulfur", sulfur);
			ar& picojson::convert::member("crystal", crystal);
			ar& picojson::convert::member("gems", gems);
		}
	};

	struct State {
		Resources resources{};
		Map map{};

		int townsCount = 0;
		int visitedObelisks = 0;
		int totalObelisks = 0;

		int color = 0;

		int month = 1;
		int week = 1;
		int day = 1;

		int heroesCount = 1;
		std::vector<Hero> heroes = {};

		std::vector<Town> towns = {};

		int clicks = 0;
		int apm = 0;

		Combat currentCombat{};


		int version = VERSION;


		friend class picojson::convert::access;
		template<class Archive>
		void json(Archive& ar) const
		{
			ar& picojson::convert::member("resources", resources);
			ar& picojson::convert::member("map", map);

			ar& picojson::convert::member("townsCount", townsCount);
			ar& picojson::convert::member("visitedObelisks", visitedObelisks);
			ar& picojson::convert::member("totalObelisks", totalObelisks);

			ar& picojson::convert::member("color", color);

			ar& picojson::convert::member("month", month);
			ar& picojson::convert::member("week", week);
			ar& picojson::convert::member("day", day);

			ar& picojson::convert::member("heroesCount", heroesCount);
			ar& picojson::convert::member("heroes", heroes);

			ar& picojson::convert::member("towns", towns);

			ar& picojson::convert::member("clicks", clicks);
			ar& picojson::convert::member("apm", apm);

			ar& picojson::convert::member("currentCombat", currentCombat);


			ar& picojson::convert::member("version", version);
		}
	};
}