package org.sebastianruml.miuiapp;

import java.io.IOException;

import org.apache.http.client.ClientProtocolException;
import org.apache.http.client.HttpClient;
import org.apache.http.client.ResponseHandler;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.BasicResponseHandler;
import org.apache.http.impl.client.DefaultHttpClient;
import org.sebastianruml.miuiapp.interfaces.AppListener;

import android.os.AsyncTask;

public class StartAppTask extends AsyncTask<String, Void, String> {
	private AppListener mListener;
	private String mAppName;
	private String mHostUrl;
		
	public StartAppTask(AppListener listener) {
		mListener = listener;
	}
	
	@Override
	protected String doInBackground(String... arg0) {
		mAppName = arg0[0];
		mHostUrl = arg0[1];
		
		try {
			HttpClient client = new DefaultHttpClient();
			String url = mHostUrl + "/chief/apps/start?name=" + mAppName;
			HttpGet get = new HttpGet(url);
			ResponseHandler<String> responseHandler = new BasicResponseHandler();
			String response = client.execute(get, responseHandler);
			
			return response;
		} catch (ClientProtocolException e) {
			e.printStackTrace();
		} catch (IOException e) {
			e.printStackTrace();
		}
		
		return null;
	}

	@Override
	protected void onPostExecute(String result) {
		boolean success;
		
		// App was successful started
		if (result.startsWith("ok")) {
			success = true;
		} else { // Some error occurred
			success = false;
		}
		
		mListener.onAppStarted(success);
	}
}
