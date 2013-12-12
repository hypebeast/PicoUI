package org.sebastianruml.miuiapp.fragments;

import java.util.ArrayList;
import java.util.HashMap;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;
import org.sebastianruml.miuiapp.R;
import org.sebastianruml.miuiapp.helpers.JSONParser;

import android.app.Activity;
import android.app.Fragment;
import android.app.ProgressDialog;
import android.os.AsyncTask;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ListAdapter;
import android.widget.ListView;
import android.widget.SimpleAdapter;

public class AppListFragment extends Fragment {
	private ListView mAppListView;
	ArrayList<HashMap<String, String>> mApps;
	private String mHostUrl;
	JSONArray mAppsJson = null;
	private OnAppSelectedListener mListener;
	
	public AppListFragment() {}
    
	
    @Override
	public void onAttach(Activity activity) {
		super.onAttach(activity);
		
		try {
			mListener = (OnAppSelectedListener)activity;
		} catch (ClassCastException e) {
			throw new ClassCastException(activity.toString() + " must implement OnAppSelectedListener");
		}
	}


	@Override
	public void onCreate(Bundle savedInstanceState) {
    	Bundle bundle = getArguments();
    	mHostUrl = bundle.getString("hostUrl");
		
		super.onCreate(savedInstanceState);
	}


	@Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
            Bundle savedInstanceState) {
  
    	View rootView = inflater.inflate(R.layout.fragment_applist, container, false);
        
    	mApps = new ArrayList<HashMap<String,String>>();
    	mAppListView = (ListView)rootView.findViewById(R.id.fragment_applist);
    	mAppListView.setOnItemClickListener(new AppItemClickListener());
    	
    	// Get all apps from picoui-chief
    	GetAppsTask getAppsTask = new GetAppsTask();
    	getAppsTask.execute(mHostUrl);
    	
        return rootView;
    }
	
	private class AppItemClickListener implements ListView.OnItemClickListener {
	    @Override
	    public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
	        mListener.onAppSelected(mApps.get(position).get("name"));
	    }
	}
	
	// Container Activity must implement this interface
    public interface OnAppSelectedListener {
        public void onAppSelected(String appName);
    }
    
    private class GetAppsTask extends AsyncTask<String, String, JSONArray> {
    	private ProgressDialog pDialog;
    	
		@Override
		protected void onPreExecute() {
			super.onPreExecute();
			
			pDialog = new ProgressDialog(getActivity());
			pDialog.setMessage("Getting apps list...");
			pDialog.setIndeterminate(false);
			pDialog.setCancelable(false);
			pDialog.show();
		}

		@Override
		protected JSONArray doInBackground(String... arg0) {
			JSONParser parser = new JSONParser();
			
			// Get the app list from picoui-chief
			JSONArray json = parser.getJSONArrayFromUrl(arg0[0] + "/chief/apps");
			
			return json;
		}

		@Override
		protected void onPostExecute(JSONArray result) {
			pDialog.dismiss();
			
			try {
				mAppsJson = result;
				
				for (int i = 0; i < mAppsJson.length(); i++) {
					JSONObject c = mAppsJson.getJSONObject(i);
					
					String appName = c.getString("name");
					
					HashMap<String, String> map = new HashMap<String, String>();
					map.put("name", appName);
					
					mApps.add(map);
					
					ListAdapter adapter = new SimpleAdapter(getActivity(), mApps,
							R.layout.app_list_item,
							new String[] {"name"},
							new int[] {R.id.app_title});
					
					mAppListView.setAdapter(adapter);
				}
			} catch (JSONException e) {
				e.printStackTrace();
			}
			
			super.onPostExecute(result);
		}
    	
    }
}
