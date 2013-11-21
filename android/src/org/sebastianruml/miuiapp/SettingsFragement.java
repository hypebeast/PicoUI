package org.sebastianruml.miuiapp;

import android.os.Bundle;
import android.preference.PreferenceFragment;

public class SettingsFragement extends PreferenceFragment {

	@Override
	public void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		
		// Load the preferences from preferences.xml
		addPreferencesFromResource(R.xml.preferences);
	}
	
}
