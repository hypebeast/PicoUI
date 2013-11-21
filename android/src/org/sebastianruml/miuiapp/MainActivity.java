package org.sebastianruml.miuiapp;

import java.net.InetAddress;

import android.os.Bundle;
import android.preference.PreferenceManager;
import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.Intent;
import android.content.SharedPreferences;
import android.view.Menu;
import android.view.MenuItem;
import android.webkit.WebView;

public class MainActivity extends Activity implements AppListener {
	private WebView webView;
	private AppPollerTask appPoller;
	private String hostUrl;
	
	@SuppressLint("SetJavaScriptEnabled")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        
        webView = (WebView)findViewById(R.id.webView);
        webView.getSettings().setJavaScriptEnabled(true);
        
        webView.loadUrl("file:///android_asset/loading.html");
        
        SharedPreferences pref = PreferenceManager.getDefaultSharedPreferences(this);
		String host = pref.getString("pref_host", "");
		String port = pref.getString("pref_port", "");
		
        appPoller = new AppPollerTask(this);
        appPoller.execute(host, port);
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


	@Override
	public void onServerAddressFound(InetAddress address) {
		SharedPreferences pref = PreferenceManager.getDefaultSharedPreferences(this);
		String port = pref.getString("pref_port", "");
		hostUrl = "http://" + address.getHostAddress() + ":" + port;
		
		webView.loadUrl("file:///android_asset/waitforapp.html");
	}


	@Override
	public void onAppFound() {
		webView.loadUrl(hostUrl);
		
		// TODO start the status thread
	}


	@Override
	public void onStatusUpdate(String status) {
		// TODO Handle the status update
		
	}


	@Override
	protected void onDestroy() {
		if (appPoller != null) {
			appPoller.stop();
		}
		
		super.onDestroy();
	}
    
}
