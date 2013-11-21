package org.sebastianruml.miuiapp;

import android.os.AsyncTask;

public class AppStatusTask extends AsyncTask<String, String, String> {
	private boolean stop;
	private Thread statusThread;
	private String status;
	private AppListener listener;
	
	public AppStatusTask(AppListener listener) {
		this.listener = listener;
	}

	@Override
	protected String doInBackground(String... params) {
		statusThread = new Thread(new Runnable() {

			@Override
			public void run() {
				while (!stop) {
					
				}
			}
			
		});
		
		statusThread.start();
		
		return null;
	}

	@Override
	protected void onPostExecute(String result) {
		// TODO Auto-generated method stub
		super.onPostExecute(result);
	}

	@Override
	protected void onProgressUpdate(String... values) {
		listener.onStatusUpdate(values[0]);
		
		super.onProgressUpdate(values);
	}

	public void stop() {
		stop = true;
		if (statusThread != null) {
			try {
				statusThread.join();
			} catch (InterruptedException e) {
				
			}
		}
	}
}
