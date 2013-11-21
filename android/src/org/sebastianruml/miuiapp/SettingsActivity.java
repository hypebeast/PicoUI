package org.sebastianruml.miuiapp;

import android.app.Activity;
import android.os.Bundle;

public class SettingsActivity extends Activity {

	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		
		// Display the settings fragment
		getFragmentManager().beginTransaction()
			.replace(android.R.id.content, new SettingsFragement())
			.commit();
	}

}
