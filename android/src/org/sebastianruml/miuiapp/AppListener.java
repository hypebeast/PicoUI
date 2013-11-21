package org.sebastianruml.miuiapp;

import java.net.InetAddress;

public interface AppListener {
	void onServerAddressFound(InetAddress address);
	void onAppFound();
	void onStatusUpdate(String status);
}
