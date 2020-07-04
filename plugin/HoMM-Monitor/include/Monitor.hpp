#pragma once
#include "H3API.hpp"
using namespace h3;

#include "Structs.hpp"
using namespace HoMMMonitor;

namespace HoMMMonitor {
	class Monitor {
	public:
		Monitor(int pid) : playerId(pid) {

		}
		void Update();
		State GetState() {
			return prev_state;
		}
		int playerId = 0;

	private:
		std::vector<Hero> iterateHeroes();
		std::vector<Town> iterateTowns();
		Combat getCombat();
		Map getMap();
		Resources getResources();

	private:
		State prev_state;
	};
}