package org.sebastianruml.miuiapp;

import java.net.InetAddress;
import java.net.URL;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.UnknownHostException;

import android.os.AsyncTask;

public class AppPollerTask extends AsyncTask<String, Void, InetAddress> {
	private AppListener listener;
	private boolean stop;
	private Thread appPinger;
	private String hostPort;
	private static final int POLL_DELAY_MS = 2000;
	
	AppPollerTask(AppListener listener) {
		this.listener = listener;
	}
	
	@Override
	protected InetAddress doInBackground(String ...params) {
		InetAddress address = null;
		
		if (params.length != 2) {
			// TODO throw exception
		}
		
		hostPort = params[1];
		
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
			this.startPingingApp(result);
		}
	}
	
	private void startPingingApp(final InetAddress address) {
		appPinger = new Thread(new Runnable() {

			@Override
			public void run() {
				boolean appFound = false;
				while (!appFound && !stop) {
					URL url;
					try {
						url = new URL("http://" + address.getHostAddress() + ":" + hostPort + "/ping");
						InputStream is = (InputStream)url.getContent();
						String response = (new BufferedReader(new InputStreamReader(is))).readLine();
						
						if (response.equals("pong")) {
							appFound = true;
							listener.onAppFound();
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
		
		appPinger.start();
	}
	
	public void stop() {
		stop = true;
		if (appPinger != null) {
			try {
				appPinger.join();
			} catch (InterruptedException e) {
				
			}
		}
	}
}
