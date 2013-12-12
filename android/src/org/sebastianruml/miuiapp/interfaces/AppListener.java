package org.sebastianruml.miuiapp.interfaces;

import java.net.InetAddress;

import org.sebastianruml.miuiapp.PicoUiStatus;

public interface AppListener {
	void onServerAddressFound(InetAddress address);
	void onPicoUiFound();
	void onStatusUpdate(PicoUiStatus status);
	void onAppStarted(boolean success);
}
