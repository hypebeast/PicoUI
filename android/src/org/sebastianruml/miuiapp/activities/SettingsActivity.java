package org.sebastianruml.miuiapp.activities;

import org.sebastianruml.miuiapp.fragments.SettingsFragement;

import android.app.Activity;
import android.os.Bundle;

public class SettingsActivity extends Activity {

	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		
		getActionBar().setTitle("Settings");
		
		// Display the settings fragment
		getFragmentManager().beginTransaction()
			.replace(android.R.id.content, new SettingsFragement())
			.commit();	
	}

}
