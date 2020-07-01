#include "Monitor.hpp"
using namespace HoMMMonitor;

std::vector<Hero> Monitor::iterateHeroes() {
	auto main = H3Internal::Main();
	std::vector<Hero> heroes;
	auto player = main->GetPlayer();

	//for each (H3Hero hero in main->heroes)
	//auto hid = 0;
	//for each (auto _h in main->heroes)
	for (size_t hid = 0; hid < 178; hid++)
	{
		auto hero = *main->GetHero(hid);
		if ( hero.owner == playerId) {
			std::vector<Item> backpack;
			std::vector<Item> weared;
			for (size_t i = 0; i < 19; i++)
			{
				auto item = hero.bodyArtifacts[i];
				if (item.id >= 0) {
					weared.push_back({ item.id, item.subtype });
				}
			}
			for (size_t i = 0; i < 64; i++)
			{
				auto item = hero.backpackArtifacts[i];
				if (item.id >= 0) {
					backpack.push_back({ item.id, item.subtype });
				}
			}
			std::vector<int> learned_spells;
			for (size_t i = 0; i < 70; i++)
			{
				auto sid = hero.learned_spell[i];
				if (sid > 0) {
					learned_spells.push_back(i);
				}
			}

			std::vector<int> available_spells;
			for (size_t i = 0; i < 70; i++)
			{
				auto sid = hero.available_spell[i];
				if (sid > 0) {
					available_spells.push_back(i);
				}
			}

			std::vector<int> secSkills;
			for (size_t i = 0; i < 28; i++)
			{
				secSkills.push_back(hero.secSkill[i]);
			}

			heroes.push_back({ hero.id, hero.name, hero.movement, hero.maxMovement, hero.experience, hero.level, hero.gender, backpack, weared, learned_spells, available_spells,
				hero.primarySkill[0], hero.primarySkill[1], hero.primarySkill[2], hero.primarySkill[3], secSkills, hero.spellPoints, hero.army.GetArmyValue(), hero.GetPower(),
				hero.x, hero.y, hero.z, hero.moraleBonus, hero.luckBonus, bool(hero.isVisible)});
		}
	}
	return heroes;
}

std::vector<Town> Monitor::iterateTowns() {
	std::vector<Town> towns;
	auto main = H3Internal::Main();
	for each (auto t in main->towns)
	{
		if (t.owner == main->GetPlayerID()) {
			auto as = 0;
			if (t.IsBuildingBuilt(H3Town::B_LIBRARY)) {
				as++;
			}
			std::vector<std::vector<int>> spells = {};
			for (int i = 0; i < 5; i++)
			{
				if (t.IsMageGuildBuilt(i)) {
					spells.push_back(std::vector<int>(t.spells[i], t.spells[i] + (5 - i) + as));
				}
			}
			towns.push_back({ t.name.String(), t.type, bool(t.IsBuildingBuilt(H3Town::B_GRAIL) == true),
				t.x, t.y, t.z, bool(t.manaVortextUnused && t.IsBuildingBuilt(H3Town::B_MANA_VORTEX)), spells, bool(t.builtThisTurn), t.Guards.GetArmyValue(), 
				t.garrisonHero, t.visitingHero, bool(t.IsBuildingBuilt(H3Town::B_FORT) == true) });
		}
	}

	return towns;
}

void Monitor::Update() {
		auto heroes = iterateHeroes();
		auto towns = iterateTowns();

		auto main = H3Internal::Main();
		auto player = main->GetPlayer();

		Combat currentCombat{};
		
		auto c = H3Internal::CombatManager();
		if (!c->finished) {
			if (c->hero[0] && c->hero[0]->owner == playerId && c->army[1]) {
				currentCombat.hero = c->hero[0]->id;
				currentCombat.heroArmyValue = c->hero[0]->GetPower();
				currentCombat.enemyArmyValue = c->army[1]->GetArmyValue();
			}
			else if (c->hero[0] && c->hero[1] && c->hero[1]->owner == playerId && c->army[1]) {
				currentCombat.hero = c->hero[1]->id;
				currentCombat.heroArmyValue = c->hero[1]->GetPower();
				currentCombat.enemyArmyValue = c->hero[0]->GetPower();
			}
		}

		
		auto map = H3Internal::AdventureManager()->map;
		std::vector<int> visionS;
		std::vector<int> visionU;
		for (size_t x = 0; x < map->mapSize; x++)
		{
			for (size_t y = 0; y < map->mapSize; y++)
			{
				visionS.push_back(F_CanViewTile(x, y, 0) ? 1 : 0);
				if (map->SubterraneanLevel > 0) {
					visionU.push_back(F_CanViewTile(x, y, 1) ? 1 : 0);
				}
			}
		}
		prev_state = {
			{player->playerResources.gold, player->playerResources.wood, player->playerResources.mercury, player->playerResources.ore, player->playerResources.sulfur,
			player->playerResources.crystal, player->playerResources.gems},
			{main->mapName.String(), map->mapSize, map->SubterraneanLevel, visionS, visionU},
			player->townsCount, player->visitedObelisks, main->obeliskCount, player->ownerID,
			main->date.month, main->date.week, main->date.day, int(heroes.size()), heroes, towns, 
			0, 0, 0, 0,
			currentCombat
		};
}