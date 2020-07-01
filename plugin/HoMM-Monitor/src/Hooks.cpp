#include "httplib.h"
#include "Hooks.hpp"
#include "H3API.hpp"
using namespace h3;

#include "Monitor.hpp"
using namespace HoMMMonitor;

#include <stdio.h>
#include <chrono>
#include <thread>
using namespace std::chrono_literals;

#include <iconvlite.h>

//!!!!! Only for testing !!!!!
//#include "Cheats.hpp"
//!!!!! Only for testing !!!!!

//#define NOTIFY_SEND

const auto SERVER = "localhost";
const int PORT = 8989;
httplib::Client httpClient(SERVER, PORT);
const auto refreshRate = 2s;

Patcher *_P;
PatcherInstance *_PI;

std::thread reporterThread;
bool started = false;

int clicksCounter = 0;
std::vector<std::chrono::steady_clock::time_point> clicks;
int apm = 0;

std::shared_ptr<Monitor> monitor;
uintptr_t _main = 0;
std::string _tn = "";

void checkAPM() {
	std::vector<std::chrono::steady_clock::time_point> _clicks;
	auto now = std::chrono::steady_clock::now();
	std::copy_if(clicks.begin(), clicks.end(), std::back_inserter(_clicks), 
		[=](std::chrono::steady_clock::time_point c) {
			return std::chrono::duration_cast<std::chrono::seconds>(now - c) <= 60s;
	});
	apm = _clicks.size();
	clicks = _clicks;
}

void resetReport() {
	httpClient.Post("/reset", "", "application/json");
}

void reporterLoop() {
	while (started) {
		if (_tn == "") {
			_tn = std::string(H3Internal::Main()->towns[0].name.String());
		}
		else if (_tn != std::string(H3Internal::Main()->towns[0].name.String())) {
			_tn = "";
			resetReport();
		}

		if (monitor->playerId != H3Internal::Main()->GetPlayerID()) {
			std::this_thread::sleep_for(refreshRate);
			continue;
		}

		auto prev_state = monitor->GetState();
		prev_state.clicks = 0;
		prev_state.apm = 0;
		std::string ps = picojson::convert::to_string(prev_state);
		monitor->Update();
		auto new_state = monitor->GetState();

		std::string ns = picojson::convert::to_string(new_state);
		new_state.clicks = clicksCounter;
		new_state.apm = apm;		

		if (ps !=  ns) {				
			std::string ns = picojson::convert::to_string(new_state);
			auto res = httpClient.Post("/report", cp2utf(ns), "application/json; charset=utf-8");
			
#ifdef NOTIFY_SEND
			F_PrintScreenText(ns.c_str());
			F_PrintScreenText("State sent");
#endif

		}		

		checkAPM();

		std::this_thread::sleep_for(refreshRate);
	}
}

void startReporter() {
	clicksCounter = 0;
	monitor = std::make_shared<Monitor>(P_ActivePlayer->ownerID);
	resetReport();
	started = true;
	reporterThread = std::thread(reporterLoop);
	char buffer[255];
	sprintf(buffer, "HoMM Monitor started\n{Plugin version: %d}", VERSION);
	F_PrintScreenText(buffer);
	return;

}

void stopReporter() {
	started = false;
	reporterThread.join();
	F_PrintScreenText("HoMM Monitor was stopped");
	resetReport();
}

int __stdcall _HH_KeyPressed(HiHook* h, H3AdventureManager* This, H3Msg* msg, int a3, int a4, int a5)
{
	if (msg->KeyPressed() == NH3VKey::H3VK_F6) {
		if (!started) {
			startReporter();
		}
		else {
			stopReporter();
		}
	}

#ifdef CHEATS
	handleKeysForCheats(msg);
#endif // CHEATS

	return THISCALL_5(int, h->GetDefaultFunc(), This, msg, a3, a4, a5);
}

bool keyLock = false;
std::string la; //prevent ghost clicks
int __stdcall _HH_DetectClick(HiHook* h, H3InputManager* This, int* a2)
{
	auto& msg = This->GetCurrentMessage();
	char buffer[50];
	sprintf(buffer, "%p", std::addressof(msg));
	if (started && !keyLock && std::string(buffer) != la
		&& (msg.message == H3InputManager::MT_LBUTTONDOWN || msg.message == H3InputManager::MT_RBUTTONDOWN || msg.message == H3InputManager::MT_KEYDOWN)){
		clicksCounter++;
		clicks.push_back(std::chrono::steady_clock::now());
		la = std::string(buffer);
		keyLock = true;
	}
	if (started && keyLock
		&& (msg.message == H3InputManager::MT_LBUTTONUP || msg.message == H3InputManager::MT_RBUTTONUP || msg.message == H3InputManager::MT_KEYUP)) {
		keyLock = false;
	}

	return THISCALL_2(int, h->GetDefaultFunc(), This, a2);
}

_LHF_(wth)
{
	//F_PrintScreenText("What the hook?");
	return EXEC_DEFAULT;
}

void hooks_init(PatcherInstance* pi)
{
	pi->WriteHiHook(0x408BA0, SPLICE_, THISCALL_, _HH_KeyPressed);
	pi->WriteHiHook(0x4EC660, SPLICE_, THISCALL_, _HH_DetectClick);
	pi->WriteLoHook(0x4F823F, wth);
}