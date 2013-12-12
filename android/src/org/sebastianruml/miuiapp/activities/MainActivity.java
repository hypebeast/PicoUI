package org.sebastianruml.miuiapp.activities;

import java.net.InetAddress;

import org.sebastianruml.miuiapp.PicoUiPollerTask;
import org.sebastianruml.miuiapp.PicoUiStatus;
import org.sebastianruml.miuiapp.R;
import org.sebastianruml.miuiapp.StartAppTask;
import org.sebastianruml.miuiapp.StatusTask;
import org.sebastianruml.miuiapp.fragments.AppListFragment;
import org.sebastianruml.miuiapp.fragments.AppListFragment.OnAppSelectedListener;
import org.sebastianruml.miuiapp.fragments.AppViewFragment;
import org.sebastianruml.miuiapp.fragments.HomeFragment;
import org.sebastianruml.miuiapp.interfaces.AppListener;

import android.app.Activity;
import android.app.Fragment;
import android.app.FragmentManager;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.res.Configuration;
import android.os.Bundle;
import android.preference.PreferenceManager;
import android.support.v4.app.ActionBarDrawerToggle;
import android.support.v4.widget.DrawerLayout;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.ListView;

public class MainActivity extends Activity implements AppListener, OnAppSelectedListener {
	private PicoUiPollerTask picouiPoller;
	private StatusTask statusPoller;
	private String hostUrl;
	private String mStatusString;
	private String mSelectedApp;
	
	private DrawerLayout mDrawerLayout;
    private ListView mDrawerList;
    private ActionBarDrawerToggle mDrawerToggle;
    
    // nav drawer title
    private CharSequence mDrawerTitle;
    
    // used to store app title
    private CharSequence mTitle;
    
    // slide menu items
    private String[] navMenuTitles;
    
    // Home Fragment
    HomeFragment mHomeFragment;
	
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        
        mStatusString = "Connecting...";
        mTitle = mDrawerTitle = getTitle();
        
        // load navigation drawer items
        navMenuTitles = getResources().getStringArray(R.array.nav_drawer_items);
        
        // Initialize navigation drawer
        mDrawerLayout = (DrawerLayout) findViewById(R.id.drawer_layout);
        mDrawerList = (ListView) findViewById(R.id.left_drawer);
        
        // Set the adapter for the list view
        mDrawerList.setAdapter(new ArrayAdapter<String>(this,
                R.layout.drawer_list_item, R.id.drawer_title, navMenuTitles));
        
        // Set the list's click listener
        mDrawerList.setOnItemClickListener(new DrawerItemClickListener());
        
        // enabling action bar app icon and behaving it as toggle button
        getActionBar().setDisplayHomeAsUpEnabled(true);
        getActionBar().setHomeButtonEnabled(true);
        
        mDrawerToggle = new ActionBarDrawerToggle(
                this,                  /* host Activity */
                mDrawerLayout,         /* DrawerLayout object */
                R.drawable.ic_drawer,  /* nav drawer icon to replace 'Up' caret */
                R.string.app_name,  /* "open drawer" description */
                R.string.app_name  /* "close drawer" description */
                ) {

            /** Called when a drawer has settled in a completely closed state. */
            public void onDrawerClosed(View view) {
                getActionBar().setTitle(mTitle);
                // calling onPrepareOptionsMenu() to show action bar icons
                invalidateOptionsMenu();
            }

            /** Called when a drawer has settled in a completely open state. */
            public void onDrawerOpened(View drawerView) {
                getActionBar().setTitle(mDrawerTitle);
                // calling onPrepareOptionsMenu() to show action bar icons
                invalidateOptionsMenu();
            }
        };
        
        mDrawerLayout.setDrawerListener(mDrawerToggle);
        
    	SharedPreferences pref = PreferenceManager.getDefaultSharedPreferences(this);
		String host = pref.getString("pref_host", "");
		String port = pref.getString("pref_port", "");
		
		// Listen for preference changes
		PreferenceManager.getDefaultSharedPreferences(this)
				.registerOnSharedPreferenceChangeListener(new PreferenceChangeListener());
		
        picouiPoller = new PicoUiPollerTask(this);
		picouiPoller.execute(host, port);
		
        if (savedInstanceState == null) {
            // on first time display view for first nav item
            displayView(0);
        }
    }
	
	private class DrawerItemClickListener implements ListView.OnItemClickListener {
	    @Override
	    public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
	        displayView(position);
	    }
	}
	
	private void displayView(int position) {
		// update the main content by replacing fragments
        Fragment fragment = null;
        Bundle args = null;
        switch (position) {
        case 0:
        	if (mHomeFragment == null) {
        		mHomeFragment = new HomeFragment();
			}
        	fragment = mHomeFragment;
        	args = new Bundle();
        	args.putString("hostUrl", hostUrl);
        	fragment.setArguments(args);
            break;
        case 1:
        	fragment = new AppListFragment();
        	args = new Bundle();
        	args.putString("hostUrl", hostUrl);
        	fragment.setArguments(args);
            break;
        case 2:
        	fragment = new AppViewFragment();
        	args = new Bundle();
        	args.putString("hostUrl", hostUrl);
        	fragment.setArguments(args);
            break;
 
        default:
            break;
        }
        
        if (fragment != null) {
            FragmentManager fragmentManager = getFragmentManager();
            fragmentManager.beginTransaction()
                    .replace(R.id.frame_container, fragment).commit();
 
            // update selected item and title, then close the drawer
            mDrawerList.setItemChecked(position, true);
            mDrawerList.setSelection(position);
            setTitle(navMenuTitles[position]);
            mDrawerLayout.closeDrawer(mDrawerList);
        } else {
            // error in creating fragment
            Log.e("MainActivity", "Error in creating fragment");
        }
	}
	
	public class PreferenceChangeListener implements SharedPreferences.OnSharedPreferenceChangeListener {
		@Override
		public void onSharedPreferenceChanged(
				SharedPreferences sharedPreferences, String key) {
			onPreferencesChanged();
		}
	}


    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.main, menu);
        return super.onCreateOptionsMenu(menu);
    }
   

	@Override
	public boolean onOptionsItemSelected(MenuItem item) {
		switch (item.getItemId()) {
		case R.id.action_settings:
			startActivity(new Intent(getBaseContext(), SettingsActivity.class));
			return true;

		default:
			return super.onOptionsItemSelected(item);
		}
	}
	
	/***
     * Called when invalidateOptionsMenu() is triggered
     */
    @Override
    public boolean onPrepareOptionsMenu(Menu menu) {
        // if nav drawer is opened, hide the action items
        boolean drawerOpen = mDrawerLayout.isDrawerOpen(mDrawerList);
        menu.findItem(R.id.action_settings).setVisible(!drawerOpen);
        return super.onPrepareOptionsMenu(menu);
    }
 
    @Override
    public void setTitle(CharSequence title) {
        mTitle = title;
        getActionBar().setTitle(mTitle);
    }
 
    
    @Override
    protected void onPostCreate(Bundle savedInstanceState) {
        super.onPostCreate(savedInstanceState);
        // Sync the toggle state after onRestoreInstanceState has occurred.
        mDrawerToggle.syncState();
    }
 
    
    @Override
    public void onConfigurationChanged(Configuration newConfig) {
        super.onConfigurationChanged(newConfig);
        // Pass any configuration change to the drawer toggles
        mDrawerToggle.onConfigurationChanged(newConfig);
    }
    
    
    @Override
   	public void onServerAddressFound(InetAddress address) {
   		SharedPreferences pref = PreferenceManager.getDefaultSharedPreferences(this);
   		String port = pref.getString("pref_port", "");
   		hostUrl = "http://" + address.getHostAddress() + ":" + port;
   	}


   	@Override
   	public void onPicoUiFound() {
   		mStatusString = "Connected";
   		if (mHomeFragment != null) {
			mHomeFragment.SetStatusText(mStatusString);
		}
   		
   		statusPoller = new StatusTask(this, hostUrl);
   		statusPoller.execute("");
   	}


   	@Override
   	public void onStatusUpdate(PicoUiStatus status) {
   		if (mHomeFragment != null) {
   			mHomeFragment.SetStatusText(mStatusString);
   			mHomeFragment.SetStatus(status);	
		}
   	}
   	
   	@Override
	public void onAppStarted(boolean success) {
		if (success) {
			displayView(2);
		}
	}
   	
   	@Override
	public void onAppSelected(String appName) {
		// Start the app on the server
		mSelectedApp = appName;
		StartAppTask task = new StartAppTask(this);
		task.execute(appName, hostUrl);
	}

   	private void onPreferencesChanged() {
   		// Stop the status task
   		if (statusPoller != null) {
			statusPoller.stop();
		}
   		
   		if (picouiPoller != null) {
			picouiPoller.stop();
		}
   		
   		mStatusString = "Connecting...";
   		if (mHomeFragment != null) {
			mHomeFragment.SetStatusText(mStatusString);
		}
   		
   		picouiPoller = new PicoUiPollerTask(this);
   	}
   	
	@Override
	protected void onDestroy() {
		if (picouiPoller != null) {
			picouiPoller.stop();
		}
		
		if (statusPoller != null) {
			statusPoller.stop();
		}
		super.onDestroy();
	}
}
