package org.sebastianruml.miuiapp;

import java.io.IOException;

import org.apache.http.client.ClientProtocolException;
import org.apache.http.client.HttpClient;
import org.apache.http.client.ResponseHandler;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.BasicResponseHandler;
import org.apache.http.impl.client.DefaultHttpClient;

import android.os.AsyncTask;

public class SystemControlTask extends AsyncTask<String, Void, String> {

	@Override
	protected String doInBackground(String... params) {
		String command = params[0];
		String hostUrl = params[1];
		
		if (command != "reboot" && command != "shutdown") {
			return null;
		}
		
		try {
			HttpClient client = new DefaultHttpClient();
			String url = hostUrl + "/chief/system/" + command;
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
		super.onPostExecute(result);
	}

}
