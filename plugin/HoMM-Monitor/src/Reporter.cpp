#include "Reporter.hpp"
using namespace HoMMMonitor;

Reporter::Reporter() {
	httpClient = new httplib::Client(SERVER, PORT);
}

std::shared_ptr<httplib::Response> Reporter::Send(State state) {
	std::string ns = picojson::convert::to_string(state);
	auto response = httpClient->Post("/report", cp2utf(ns), "application/json; charset=utf-8");
	return response;
}

std::shared_ptr<httplib::Response> Reporter::SendHeartbeat() {
	auto response = httpClient->Post("/heartbeat", "", "application/json; charset=utf-8");
	return response;
}

std::shared_ptr<httplib::Response> Reporter::Reset() {
	auto response = httpClient->Post("/reset", "", "application/json; charset=utf-8");
	return response;
}