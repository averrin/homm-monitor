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

HHOOK g_LowLevelMouseHook = NULL;
HHOOK g_keyboardHook = NULL;
HINSTANCE hInstance = NULL;

const auto SERVER = "localhost";
const int PORT = 8989;
httplib::Client httpClient(SERVER, PORT);
const auto refreshRate = 2s;

Patcher *_P;
PatcherInstance *_PI;

std::thread reporterThread;
bool started = false;

int actionsCounter = 0;
std::vector<std::chrono::steady_clock::time_point> actions;
int apm = 0;

int clean_actionsCounter = 0;
std::vector<std::chrono::steady_clock::time_point> clean_actions;
int clean_apm = 0;

std::shared_ptr<Monitor> monitor;
uintptr_t _main = 0;
std::string _tn = "";

LRESULT CALLBACK LowLevelMouseProc(int code, WPARAM wParam, LPARAM lParam)
{
	if (code == HC_ACTION)
	{
		switch (wParam)
		{
		case WM_LBUTTONDOWN:
			clean_actionsCounter++;
			clean_actions.push_back(std::chrono::steady_clock::now());
		case WM_RBUTTONDOWN:
			actionsCounter++;
			actions.push_back(std::chrono::steady_clock::now());
			break;
		}

	}
	return CallNextHookEx(g_LowLevelMouseHook, code, wParam, lParam);
}

LRESULT CALLBACK LowLevelKeyboardProc(int code, WPARAM wParam, LPARAM lParam)
{
	if (wParam == WM_KEYDOWN) {
		actionsCounter++;
		actions.push_back(std::chrono::steady_clock::now());
		clean_actionsCounter++;
		clean_actions.push_back(std::chrono::steady_clock::now());
	}
	return CallNextHookEx(g_keyboardHook, code, wParam, lParam);
}

void setInputHooks() {
	UnhookWindowsHookEx(g_keyboardHook);
	UnhookWindowsHookEx(g_LowLevelMouseHook);
	g_LowLevelMouseHook = SetWindowsHookEx(WH_MOUSE_LL, LowLevelMouseProc, hInstance, 0);
	if (!g_LowLevelMouseHook)
	{
		F_PrintScreenText("failed to set mouse hook");
	}

	g_keyboardHook = SetWindowsHookEx(WH_KEYBOARD_LL, LowLevelKeyboardProc, hInstance, 0);
	if (!g_keyboardHook)
	{
		F_PrintScreenText("failed to set keyboard hook");
	}
}

void checkAPM() {
	std::vector<std::chrono::steady_clock::time_point> _clicks;
	auto now = std::chrono::steady_clock::now();
	std::copy_if(actions.begin(), actions.end(), std::back_inserter(_clicks),
		[=](std::chrono::steady_clock::time_point c) {
			return std::chrono::duration_cast<std::chrono::seconds>(now - c) <= 60s;
	});
	apm = _clicks.size();
	actions = _clicks;

	std::vector<std::chrono::steady_clock::time_point> c_clicks;
	std::copy_if(clean_actions.begin(), clean_actions.end(), std::back_inserter(c_clicks),
		[=](std::chrono::steady_clock::time_point c) {
		return std::chrono::duration_cast<std::chrono::seconds>(now - c) <= 60s;
	});
	clean_apm = c_clicks.size();
	clean_actions = c_clicks;
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
		prev_state.actions = 0;
		prev_state.apm = 0;
		prev_state.cleanActions = 0;
		prev_state.cleanApm = 0;
		std::string ps = picojson::convert::to_string(prev_state);
		monitor->Update();
		auto new_state = monitor->GetState();

		std::string ns = picojson::convert::to_string(new_state);
		new_state.actions = actionsCounter;
		new_state.apm = apm;
		new_state.cleanActions = clean_actionsCounter;
		new_state.cleanApm = clean_apm;

		if (ps !=  ns) {
#ifdef NOTIFY_SEN
			char buffer[255];
			sprintf(buffer, "{APM: %d/%d}", apm, clean_apm);
			F_PrintScreenText(buffer);
			sprintf(buffer, "{Actions: %d/%d}", actionsCounter, clean_actionsCounter);
			F_PrintScreenText(buffer);
#endif

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
	actionsCounter = 0;
	actions.clear();
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
	UnhookWindowsHookEx(g_keyboardHook);
	UnhookWindowsHookEx(g_LowLevelMouseHook);
	F_PrintScreenText("HoMM Monitor was stopped");
	resetReport();
}

int __stdcall _HH_KeyPressed(HiHook* h, H3AdventureManager* This, H3Msg* msg, int a3, int a4, int a5)
{
	if (msg->KeyPressed() == NH3VKey::H3VK_F6) {
		if (!started) {
			startReporter();
			setInputHooks();
		}
		else {
			stopReporter();
		}
	}
	/*
	if (msg->KeyPressed() == NH3VKey::H3VK_F5) {
		setInputHooks();
	}
	*/

#ifdef CHEATS
	handleKeysForCheats(msg);
#endif // CHEATS

	return THISCALL_5(int, h->GetDefaultFunc(), This, msg, a3, a4, a5);
}

void hooks_init(HMODULE hModule, PatcherInstance* pi)
{
	hInstance = hModule;
	pi->WriteHiHook(0x408BA0, SPLICE_, THISCALL_, _HH_KeyPressed);
}