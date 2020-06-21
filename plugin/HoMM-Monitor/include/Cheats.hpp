#pragma once
#include "H3API.hpp"
using namespace h3;

#define CHEATS 1

void giveItem(int id) {
	H3Artifact art;
	art.id = id;
	art.subtype = -1;
	P_Main->GetHero(P_ActivePlayer->currentHero)->GiveBackpackArtifact(art);
}

void giveScroll(int type) {
	H3Artifact art;
	art.id = H3Artifact::SPELL_SCROLL;
	art.subtype = type;
	P_Main->GetHero(P_ActivePlayer->currentHero)->GiveBackpackArtifact(art);
}

void cheat() {
	giveItem(H3Artifact::ANGEL_WINGS);
	giveItem(H3Artifact::SPELLBINDERS_HAT);
	giveItem(H3Artifact::TOME_OF_AIR_MAGIC);
	giveItem(H3Artifact::TOME_OF_FIRE_MAGIC);
	giveItem(H3Artifact::TOME_OF_EARTH_MAGIC);

	giveItem(H3Artifact::ANGELIC_ALLIANCE);
	giveItem(H3Artifact::ARMOR_OF_THE_DAMNED);
	giveItem(H3Artifact::SHACKLES_OF_WAR);
	//giveItem(H3Artifact::IRONFIST); //art from hota, missed for now

	giveScroll(H3Spell::FLY);
	giveScroll(H3Spell::TOWN_PORTAL);
	giveScroll(H3Spell::ARMAGEDDON);
	giveScroll(H3Spell::RESURRECTION);
	giveScroll(H3Spell::DIMENSION_DOOR);
}

void handleKeysForCheats(H3Msg* msg) {
	if (msg->KeyPressed() == NH3VKey::H3VK_F3) {
		cheat();
	}

	if (msg->KeyPressed() == NH3VKey::H3VK_1) {
		giveScroll(H3Spell::FLY);
	}
	if (msg->KeyPressed() == NH3VKey::H3VK_2) {
		giveScroll(H3Spell::TOWN_PORTAL);
	}
	if (msg->KeyPressed() == NH3VKey::H3VK_3) {
		giveScroll(H3Spell::ARMAGEDDON);
	}
	if (msg->KeyPressed() == NH3VKey::H3VK_4) {
		giveScroll(H3Spell::RESURRECTION);
	}
	if (msg->KeyPressed() == NH3VKey::H3VK_5) {
		giveScroll(H3Spell::DIMENSION_DOOR);
	}
}