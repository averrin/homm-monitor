#pragma once
#include "httplib.h"
#include "iconvlite.h"
#include "Structs.hpp"
using namespace HoMMMonitor;

namespace HoMMMonitor {
	class Reporter {
	public:
		Reporter();
		std::shared_ptr<httplib::Response> Send(State state);
		std::shared_ptr<httplib::Response> SendHeartbeat();
		std::shared_ptr<httplib::Response> Reset();
	private:
		const std::string SERVER = "localhost";
		const int PORT = 8989;
		httplib::Client* httpClient;
	};
}
