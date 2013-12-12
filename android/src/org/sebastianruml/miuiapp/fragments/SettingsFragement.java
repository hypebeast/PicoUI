package org.sebastianruml.miuiapp.fragments;

import org.sebastianruml.miuiapp.R;

import android.content.SharedPreferences;
import android.content.SharedPreferences.OnSharedPreferenceChangeListener;
import android.os.Bundle;
import android.preference.PreferenceFragment;

public class SettingsFragement extends PreferenceFragment 
								implements OnSharedPreferenceChangeListener {

	private static final String KEY_PREF_HOST = "";
	private static final String KEY_PREF_PORT = "";
	
	@Override
	public void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		
		// Load the preferences from preferences.xml
		addPreferencesFromResource(R.xml.preferences);
	}

	@Override
	public void onSharedPreferenceChanged(SharedPreferences sharedPreferences, String key) {
		if (key == KEY_PREF_HOST || key == KEY_PREF_PORT) {
			// Inform main activity about the change
		}
	}

	@Override
	public void onPause() {
		getPreferenceManager().getSharedPreferences().unregisterOnSharedPreferenceChangeListener(this);
		super.onPause();
	}

	@Override
	public void onResume() {
		super.onResume();
		getPreferenceManager().getSharedPreferences().registerOnSharedPreferenceChangeListener(this);
	}
}
