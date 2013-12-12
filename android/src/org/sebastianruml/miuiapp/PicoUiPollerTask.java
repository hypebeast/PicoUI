package org.sebastianruml.miuiapp;

import java.net.InetAddress;
import java.net.URL;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.UnknownHostException;

import org.sebastianruml.miuiapp.interfaces.AppListener;

import android.os.AsyncTask;

public class PicoUiPollerTask extends AsyncTask<String, Void, InetAddress> {
	private AppListener listener;
	private boolean stop;
	private Thread chiefPinger;
	private String port;
	private static final int POLL_DELAY_MS = 2000;
	
	public PicoUiPollerTask(AppListener listener) {
		this.listener = listener;
	}
	
	@Override
	protected InetAddress doInBackground(String ...params) {
		InetAddress address = null;
		
		if (params.length != 2) {
			// TODO throw exception
		}
		
		port = params[1];
		
		while(address == null && !stop) {
			try {
				address = InetAddress.getByName(params[0]);
			} catch (UnknownHostException e) {
				try {
					Thread.sleep(POLL_DELAY_MS);
				} catch (InterruptedException e1) {
				
				}
			}
		}
		
		return address;
	}
	
	@Override
	protected void onPostExecute(final InetAddress result) {
		if (result != null) {
			listener.onServerAddressFound(result);
			this.startPingingChief(result);
		}
	}
	
	private void startPingingChief(final InetAddress address) {
		chiefPinger = new Thread(new Runnable() {

			@Override
			public void run() {
				boolean appFound = false;
				while (!appFound && !stop) {
					URL url;
					try {
						url = new URL("http://" + address.getHostAddress() + ":" + port + "/chief/ping");
						InputStream is = (InputStream)url.getContent();
						String response = (new BufferedReader(new InputStreamReader(is))).readLine();
						
						if (response.equals("pong")) {
							appFound = true;
							listener.onPicoUiFound();
						}
					} catch (IOException e) {
						
					}
					
					try {
						Thread.sleep(POLL_DELAY_MS);
					} catch (InterruptedException e) {
					
					}
				}
			}
		});
		
		chiefPinger.start();
	}
	
	public void stop() {
		stop = true;
		if (chiefPinger != null) {
			try {
				chiefPinger.join();
				chiefPinger = null;
			} catch (InterruptedException e) {
				
			}
		}
	}
}
