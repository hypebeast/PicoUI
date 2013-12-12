package org.sebastianruml.miuiapp;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.HashMap;
import java.util.Map;

import org.sebastianruml.miuiapp.helpers.Utils;
import org.sebastianruml.miuiapp.interfaces.AppListener;

import android.os.AsyncTask;

public class StatusTask extends AsyncTask<String, PicoUiStatus, String> {
	private boolean stop;
	private Thread statusThread;
	private String hostUrl;
	private PicoUiStatus status;
	private AppListener listener;
	private Map<String, URL> statusUrls;
	
	private static final int POLL_DELAY_MS = 2000;
	
	public StatusTask(AppListener listener, String hostUrl) {
		this.listener = listener;
		this.hostUrl = hostUrl;
		status = new PicoUiStatus();
		statusUrls = new HashMap<String, URL>();
		
		// Build all status URLs
		try {
			URL url = new URL(this.hostUrl + "/chief/system/proc_uptime");
			statusUrls.put("uptime", url);
			url = new URL(this.hostUrl + "/chief/system/load");
			statusUrls.put("load", url);
			url = new URL(this.hostUrl + "/chief/system/mem");
			statusUrls.put("mem", url);
			url = new URL(this.hostUrl + "/chief/system/cpu_usage");
			statusUrls.put("cpu", url);
		} catch (MalformedURLException e) {
			
		}
	}

	@Override
	protected String doInBackground(String... params) {
		statusThread = new Thread(new Runnable() {
	
			@Override
			public void run() {
				while (!stop) {
					try {
						for (Map.Entry<String, URL> entry : statusUrls.entrySet()) {
							StringBuilder res = new StringBuilder();
							InputStream is = (InputStream)entry.getValue().getContent();
							BufferedReader reader = new BufferedReader(new InputStreamReader(is));
							String line = reader.readLine();
							while (line != null) {
								res.append(line + "\n");
								line = reader.readLine();
							}
							
							String response = res.toString();
							if (entry.getKey() == "uptime") {
								String seconds = response.split(" ")[0];
								status.SetUptime(Utils.ConvertSecondsToTimeString(seconds));
							} else if (entry.getKey() == "load") {
								status.SetLoad(response.split(" ")[2]);
							} else if (entry.getKey() == "mem") {
								String[] lines = response.split("\n");
								for (String l : lines) {
									if (l.contains("MemFree:")) {
										l = l.replace("MemFree:", "").replace("kB", "").trim();
										status.SetFreeMem(Integer.parseInt(l) / 1024);
										break;
									}
								}
							} else if (entry.getKey() == "cpu") {
								status.SetCpuLoad(Integer.parseInt(response.trim()));
							}
						}
						
						publishProgress(status);
						
						Thread.sleep(POLL_DELAY_MS);
					} catch (IOException e) {
						
					} catch (InterruptedException e) {
						
					}
				}
			}
			
		});
		
		statusThread.start();
		
		return null;
	}

	@Override
	protected void onPostExecute(String result) {
		super.onPostExecute(result);
	}

	@Override
	protected void onProgressUpdate(PicoUiStatus... progress) {
		listener.onStatusUpdate(progress[0]);
		super.onProgressUpdate(progress);
	}

	public void stop() {
		stop = true;
		if (statusThread != null) {
			try {
				statusThread.join();
				statusThread = null;
			} catch (InterruptedException e) {
				
			}
		}
	}
}
